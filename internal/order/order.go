package order

import (
	"log"
	"pay-service/internal/database"
	"pay-service/internal/notification"
	"pay-service/internal/payment/payment_types"
	"runtime"

	"github.com/mike504110403/common-moduals/typeparam"
)

// 監聽訂單支付狀態
func listenInitOrderCompleted() {
	for order := range notification.PaymentOrderChannel {

		tx, err := database.PAYORDER.TX()
		if err != nil {
			log.Printf("[Payment] Failed to open transaction")
			return
		}

		orderStatusType := typeparam.TypeParam{MainType: "order_status"}
		if order.IsScuess {
			orderStatusType.SubType = "success"
		} else {
			orderStatusType.SubType = "fail"
		}
		statusInt, _ := orderStatusType.Get()
		err = UpdateOrderStatus(tx, order.Trade, statusInt)
		if err != nil {
			log.Printf("[Payment] Failed to update order status")
			tx.Rollback()
			continue
		}

		// TODO: 新狀態的tx 要跟app端的回應綁在一起
		// TODO: 這邊要做通知app端的動作
		// 成功才commit 加上內部通知payment模組
		// db, err := database.PAYORDER.DB()
		// if err != nil {
		// 	log.Printf("[Payment] Failed to get db")
		// 	return
		// }
		// appOrder, err := orders.GetOrder(db, order.Trade)
		// if err != nil {
		// 	log.Printf("[Payment] Failed to get order")
		// 	return
		// }

		tx, err = database.PAYORDER.TX()
		if err != nil {
			log.Printf("[Payment] Failed to open transaction")
			return
		}

		err = CallbackUpdate(tx, order.Trade)
		if err != nil {
			log.Printf("[Payment] Failed to update order status")
			tx.Rollback()
			continue
		}

		res := payment_types.PaymentOrderResponse(order)

		log.Printf("[Payment] 支付成功，訂單 ID: %s", order.Trade)

		err = AppCaller(res)
		if err != nil {
			log.Printf("[Payment] Failed to call app")
			return
		}

		// 查找對應的 responseChannel 並回傳起單結果
		orderResChan, exist := notification.GetOrderChannel(order.Trade)
		if !exist {
			continue
		}
		orderResChan <- res
		tx.Commit()
	}
}

// init : 啟動 Payment 監聽
func init() {
	numWorkers := runtime.NumCPU() // 使用 CPU 核心數來決定 Worker 數量
	log.Printf("啟動 %d 個 Order Worker", numWorkers)

	for i := 0; i < numWorkers; i++ {
		go listenInitOrderCompleted()
	}
}
