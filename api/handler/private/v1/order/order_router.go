package order

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"pay-service/internal/database"
	"pay-service/internal/notification"
	"pay-service/internal/order"
	"pay-service/internal/payment/payment_types"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mike504110403/common-moduals/typeparam"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/order")
	{
		g.Post("/create", createOrderHandler) // 建立訂單
		g.Post("/result", orderResultHandler) // 訂單結果
	}
}

// @Summary 創建訂單
// @Description 用戶請求創建訂單，系統先寫入訂單資料，再通知金流模組
// @Accept json
// @Produce json
// @Param request body CreateOrderRequest true "訂單資料"
// @Success 200 {object} Order
// @Router /api/order/create [post]
func createOrderHandler(c *fiber.Ctx) error {
	var req CreateOrderRequest

	// 解析 JSON 請求體
	if err := c.BodyParser(&req); err != nil {
		log.Printf("BodyParser 解析錯誤: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
	}

	// 開啟 db連線
	db, err := database.PAYORDER.DB()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
	}

	// 確認訂單是否重複
	if err := order.CheckOrderRepeat(db, req.AppTrade); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
	}

	// 設定訂單狀態
	payOrderType := typeparam.TypeParam{
		MainType: "order_status",
		SubType:  "init",
	}

	// 取得訂單狀態
	status, err := payOrderType.Get()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
	}

	u, _ := uuid.NewRandom()
	trade := base64.RawURLEncoding.EncodeToString(u[:])

	orderRes := order.Order{
		AppID:       req.AppID,
		AppTrade:    req.AppTrade,
		Trade_No:    trade,
		ProviderID:  req.ProviderID,
		PaymentType: req.PaymentType,
		Amount:      req.Amount,
		Currency:    req.Currency,
		CallbackURL: req.CallbackURL,
		Status:      status,
		ReturnURL:   req.ReturnURL,
	}

	// 創建訂單
	if i, err := order.CreateOrder(db, orderRes); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "Failed to create order"})
	} else {
		orderRes.Id = i
	}

	log.Printf("[Order] 訂單創建成功，ID: %s，等待支付", trade)

	// 開啟 transaction
	tx, err := database.PAYORDER.TX()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
	}

	res, err := NotifyPayment(orderRes)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
	} else {
		orderStatusType := typeparam.TypeParam{MainType: "order_status"}
		if res.IsScuess {
			orderStatusType.SubType = "inited"
		} else {
			orderStatusType.SubType = "fail"
		}
		statusInt, _ := orderStatusType.Get()
		if err := order.UpdateOrderStatus(tx, trade, statusInt); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
		}
	}
	return c.Status(fiber.StatusOK).JSON(res.PayUrl)
}

// TODO: 收到前端
func orderResultHandler(c *fiber.Ctx) error {
	var req payment_types.AuthCallBackReq
	// 解析 JSON 請求體
	if err := c.BodyParser(&req); err != nil {
		log.Printf("BodyParser 解析錯誤: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
	}

	//result, err := NotifyAuthPayment(req)
	result, err := testNotifyAuthPayment(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
	}
	// 開啟 transaction
	tx, err := database.PAYORDER.TX()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
	} else {
		orderStatusType := typeparam.TypeParam{MainType: "order_status"}
		if result.IsScuess {
			orderStatusType.SubType = "done"
		} else {
			orderStatusType.SubType = "fail"
		}
		statusInt, _ := orderStatusType.Get()
		if err := order.UpdateOrderStatus(tx, result.Trade_no, statusInt); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
		}
		db, err := database.SETTING.DB()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
		}
		status, err := order.OrderSubTypeDescription(db, statusInt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
		}
		// 開啟 transaction
		tx, err = database.PAYORDER.TX()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
		}
		err = order.CallbackUpdate(tx, result.Trade_no)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).SendString("服務錯誤")
		}
		fmt.Printf("[Order] 訂單修改完成，訂單 ID: %v，狀態: %v\n", result.Trade_no, status)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": fiber.Map{
				"message":      "訂單狀態更新成功",
				"trade_no":     result.Trade_no,
				"trade_status": status,
			},
		})
	}
}

func testNotifyAuthPayment(req payment_types.AuthCallBackReq) (payment_types.AuthCallBackRes, error) {
	// 初始化通知通道
	notification.AuthCallBackChannel = make(chan payment_types.AuthCallBackReq, 1)

	// 建立 responseChannel
	orderResChan := make(chan payment_types.AuthCallBackRes, 1)
	notification.SetAuthCallBackChannel("dOeC5pYLQtWFmFfXHvrFkg", orderResChan)

	// 模擬金流模組回傳結果
	go func() {
		time.Sleep(1 * time.Second) // 模擬延遲
		orderResChan <- payment_types.AuthCallBackRes{
			Trade_no:    "a0U4s6ojSGiUXfYku_Mjgw",
			App_id:      1,
			Provider_id: 1,
			IsScuess:    false,
		}
	}()

	// 發送支付請求到 channel
	notification.AuthCallBackChannel <- req

	// 設定超時機制，避免支付結果長時間未返回
	select {
	case result := <-orderResChan:
		notification.RemoveAuthCallBackChannel(req.OrderID)
		return result, nil
	case <-time.After(10 * time.Second): // 超過 10 秒，返回超時
		notification.RemoveAuthCallBackChannel(req.OrderID)
		return payment_types.AuthCallBackRes{}, errors.New("payment timeout")
	}
}
