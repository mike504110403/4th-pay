package order

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"pay-service/internal/notification"
	"pay-service/internal/payment/payment_types"

	"sync"
	"testing"
	"time"
)

// 定義 OrderRepository 介面
type OrderRepository interface {
	UpdateOrderStatus(tx *sql.Tx, orderID string, status string) error
}

// 真正的 DB 實作
type DBOrderRepository struct{}

func (r *DBOrderRepository) UpdateOrderStatus(tx *sql.Tx, orderID string, status string) error {
	queryStr := `UPDATE OrderRecord SET status = ? WHERE trade_no = ?`
	_, err := tx.Exec(queryStr, status, orderID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("UpdateOrderStatus error: %v", err)
	}
	return tx.Commit()
}

// Mock 版本的 OrderRepository
type MockOrderRepository struct {
	updatedOrders map[string]string
	mu            sync.Mutex
}

func (r *MockOrderRepository) UpdateOrderStatus(tx interface{}, orderID string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.updatedOrders[orderID] = status
	log.Printf("[Mock] 訂單狀態更新成功：%s -> %s", orderID, status)
	return nil
}

func TestListenInitOrderCompleted(t *testing.T) {
	mockRepo := &MockOrderRepository{updatedOrders: make(map[string]string)}

	// 啟動監聽通道
	go listenInitOrderCompleted()

	// 建立 orderChannel 並記錄
	orderChannel := make(chan payment_types.PaymentOrderResponse, 1)
	notification.SetOrderChannel("TEST123", orderChannel)

	// 產生測試交易
	testOrder := payment_types.PaymentOrderRequest{
		Trade:         "TEST123",
		ProviderTrade: "test",
		IsScuess:      true,
	}

	// 發送測試交易到通道
	notification.PaymentOrderChannel <- testOrder

	// 等待 `listenInitOrderCompleted` 處理
	time.Sleep(3 * time.Second)

	// 驗證訂單是否被標記為 `completed`
	mockRepo.mu.Lock()
	if mockRepo.updatedOrders["TEST123"] != "completed" {
		t.Errorf("測試失敗: 訂單狀態未更新為 completed")
	} else {
		t.Logf("測試成功: 訂單已標記為 completed")
	}
	mockRepo.mu.Unlock()

	// 驗證 `orderChannel` 是否收到結果
	select {
	case result := <-orderChannel:
		if !result.IsScuess {
			t.Errorf("測試失敗: OrderChannel 訂單狀態應該為 completed，但收到 %v", result.IsScuess)
		} else {
			t.Logf("測試成功: OrderChannel 訂單狀態為 completed")
		}
	case <-time.After(3 * time.Second):
		t.Errorf("測試失敗: OrderChannel 未收到訂單結果")
	}
}

func TestMain(m *testing.M) {
	// 初始化 DB
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/PayOrder")
	if err != nil {
		log.Fatalf("資料庫連線失敗: %v", err)
	}

	// 設定連線池
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(60 * time.Second)

	if err = db.Ping(); err != nil {
		log.Fatalf("資料庫無法連線: %v", err)
	}

	log.Println("測試資料庫連線成功")

	// 執行測試
	code := m.Run()

	// 測試完成後關閉 DB 連線
	db.Close()

	// 退出測試
	os.Exit(code)
}
