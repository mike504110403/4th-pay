package payment

import (
	"database/sql"
	"errors"
	"pay-service/internal/database"
	"pay-service/internal/notification"
	"pay-service/internal/payment/payment_types"
	"time"
)

// CreatePaymentTransaction 將支付交易存入資料庫
func CreatePaymentTransaction(transaction *payment_types.PaymentTransaction) error {
	db, err := database.PAYORDER.DB()
	if err != nil {
		return err
	}
	sqlStr := `
		INSERT INTO PayRecord (
			order_id
		) VALUES (?)
	`
	if _, err := db.Exec(sqlStr, transaction.ID); err != nil {
		return err
	}
	return nil
}

// GetPaymentTransaction 根據 OrderID 獲取支付交易
func GetPaymentTransaction(tx *sql.Tx, orderID uint64) (*payment_types.PaymentTransaction, error) {
	var transaction payment_types.PaymentTransaction
	return &transaction, nil
}

// NotifyOrder 通知支付
func NotifyOrder(req payment_types.PaymentOrderRequest) (payment_types.PaymentOrderResponse, error) {
	// 建立 responseChannel
	resChan := make(chan payment_types.PaymentOrderResponse)

	// 記錄 responseChannel
	notification.SetOrderChannel(req.Trade, resChan)

	// TODO: 完整收款物件
	paymentOrderReq := payment_types.PaymentOrderRequest{
		Trade:         req.Trade,
		ProviderTrade: req.ProviderTrade,
		IsScuess:      req.IsScuess,
	}

	// 發送收款請求到 channel
	notification.PaymentOrderChannel <- paymentOrderReq

	// 設定超時機制，避免支付結果長時間未返回
	select {
	case result := <-resChan:
		notification.RemoveOrderChannel(req.Trade)
		return result, nil
	case <-time.After(5 * time.Second):
		notification.RemoveOrderChannel(req.Trade)
		return payment_types.PaymentOrderResponse{}, errors.New("通知超時")
	}
}

// GetPayOrder 獲取支付訂單
func GetPayOrderId(eOrderNo string) (int, error) {
	var id int
	db, err := database.PAYORDER.DB()
	if err != nil {
		return id, err
	}
	sqlStr := `
		SELECT ord.id FROM PayRecord AS prd
		JOIN OrderRecord AS ord
		ON prd.order_id = ord.id
		WHERE ord.trade_no = ?
		AND prd.is_notify = 'N'
	`
	if err := db.QueryRow(sqlStr, eOrderNo).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

// UpdatePayOrderStatus 更新支付訂單狀態
func UpdatePayOrderStatus(pr payment_types.PayRecord) error {
	db, err := database.PAYORDER.DB()
	if err != nil {
		return err
	}
	sqlStr := `
		UPDATE PayRecord SET provider_trade = ?, status = ? WHERE order_id = ?
	`
	if _, err := db.Exec(sqlStr, pr.ProviderTrade, pr.Status, pr.OrdeId); err != nil {
		return err
	}
	return nil
}

// UpdatePayOrderIsNotify 更新支付訂單已通知
func UpdatePayOrderIsNotify(orderId int) {
	db, err := database.PAYORDER.DB()
	if err != nil {
		return
	}
	db.Exec(`UPDATE PayRecord SET is_notify = 'Y' WHERE order_id = ?`, orderId)
	return
}
