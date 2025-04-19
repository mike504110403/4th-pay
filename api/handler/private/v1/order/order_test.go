package order

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pay-service/internal/database"
	"pay-service/internal/notification"
	"pay-service/internal/payment/payment_types"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	database.Init()
	reqData := CreateOrderRequest{
		AppID:       0,
		AppTrade:    uuid.New().String(),
		Amount:      500.0,
		Currency:    "TWD",
		PaymentType: 0,
		ProviderID:  1,
		CallbackURL: "http://google.com",
	}

	reqBody, _ := json.Marshal(reqData)

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/v1/v1/order/create", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("🚨 創建請求失敗: ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err.Error() == "payment timeout" {

	}
	if err != nil {
		fmt.Println("🚨 發送請求失敗: ", err)
	}
	defer resp.Body.Close()
}

func TestNotifyAuthPayment(t *testing.T) {
	// 初始化通知通道
	notification.AuthCallBackChannel = make(chan payment_types.AuthCallBackReq, 1)

	// 建立 responseChannel
	orderResChan := make(chan payment_types.AuthCallBackRes, 1)
	notification.SetAuthCallBackChannel("dOeC5pYLQtWFmFfXHvrFkg", orderResChan)

	// 模擬金流模組回傳結果
	go func() {
		time.Sleep(1 * time.Second) // 模擬延遲
		orderResChan <- payment_types.AuthCallBackRes{
			Trade_no:    "dOeC5pYLQtWFmFfXHvrFkg",
			App_id:      1,
			Provider_id: 1,
			IsScuess:    true,
		}
	}()

	// 測試 NotifyAuthPayment
	req := payment_types.AuthCallBackReq{
		OrderID: "dOeC5pYLQtWFmFfXHvrFkg",
	}
	res, err := NotifyAuthPayment(req)

	// 驗證結果
	assert.NoError(t, err)
	assert.Equal(t, "dOeC5pYLQtWFmFfXHvrFkg", res.Trade_no)
	assert.Equal(t, 1, res.App_id)
	assert.Equal(t, 1, res.Provider_id)
	assert.True(t, res.IsScuess)
}
