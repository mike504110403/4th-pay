package payment

import (
	"fmt"
	"log"
	"pay-service/internal/notification"
	"pay-service/internal/payment/payment_types"
	"runtime"
)

// 監聽支付請求
func listenInitRequests() {
	for r := range notification.PaymentRequestChannel {
		fmt.Printf("[Payment] 接收到支付請求，訂單 ID: %s，金額: %.2f", r.Trade, r.Amount)
		// 取得金流商
		provider, err := GetProvider(r.Provider)
		if err != nil {
			continue
		}
		// 交易物件
		trc := payment_types.PaymentTransaction{
			ID:        r.Id,
			Trade:     r.Trade,
			Provider:  r.Provider,
			Channel:   r.Channel,
			Amount:    r.Amount,
			ReturnURL: r.ReturnURL,
		}
		if err := provider.NewTransaction(&trc); err != nil {
			log.Printf("[Payment] 起單失敗: %s", err)
			continue
		}

		// 查找對應的 responseChannel 並回傳起單結果
		if responseChannel, exists := notification.GetResponseChannel(trc.Trade); exists {
			responseChannel <- trc
		}
	}
}

// init : 啟動 Payment 監聽
func init() {
	numWorkers := runtime.NumCPU() // 使用 CPU 核心數來決定 Worker 數量
	log.Printf("啟動 %d 個 Payment Worker", numWorkers)

	for i := 0; i < numWorkers; i++ {
		go listenInitRequests()
	}
}
