package payment

import (
	"fmt"
	"log"
	"pay-service/internal/notification"
	"pay-service/internal/payment/payment_types"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotifyOrder(t *testing.T) {
	// **建立 `PaymentOrderChannel`**
	notification.PaymentOrderChannel = make(chan payment_types.PaymentOrderRequest, 1)

	// **啟動監聽訂單支付狀態**
	go listenInitRequest()

	// **建立 `orderChannel`**
	orderResChan := make(chan payment_types.PaymentOrderResponse, 1)
	notification.SetOrderChannel("TEST-123", orderResChan)

	// **發送支付成功通知**
	go func() {
		time.Sleep(1 * time.Second) // 模擬延遲
		notification.PaymentOrderChannel <- payment_types.PaymentOrderRequest{
			Trade:         "TEST-123",
			ProviderTrade: "MOCK-12345",
			IsScuess:      true,
		}
	}()

	// **測試 `NotifyOrder`**
	res, err := NotifyOrder(payment_types.PaymentOrderRequest{
		Trade:         "TEST-123",
		ProviderTrade: "MOCK-12345",
		IsScuess:      true,
	})

	// **驗證結果**
	assert.NoError(t, err)
	assert.Equal(t, "TEST-123", res.Trade)
	assert.Equal(t, "MOCK-12345", res.ProviderTrade)
	log.Println("[測試成功] 訂單支付結果通知成功")

}

func listenInitRequest() {
	for r := range notification.PaymentRequestChannel {
		fmt.Printf("[Payment] 接收到支付請求，訂單 ID: %s，金額: %.2f", r.Trade, r.Amount)
		// 交易物件
		trc := payment_types.PaymentTransaction{
			ID:        r.Id,
			Trade:     r.Trade,
			Provider:  r.Provider,
			Channel:   r.Channel,
			Amount:    r.Amount,
			ReturnURL: r.ReturnURL,
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
		go listenInitRequest()
	}
}
