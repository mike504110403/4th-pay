package order

import (
	"errors"
	"pay-service/internal/cachedata"
	"pay-service/internal/notification"
	"pay-service/internal/order"
	"pay-service/internal/payment/payment_types"
	"strconv"
	"time"
)

// NotifyPayment 通知支付
func NotifyPayment(req order.Order) (payment_types.PaymentTransaction, error) {
	// 建立 responseChannel
	responseChannel := make(chan payment_types.PaymentTransaction, 1)

	// 記錄 responseChannel，以便金流模組回傳結果
	notification.SetResponseChannel(req.Trade_No, responseChannel)

	// 取通道 ID
	providerName := cachedata.GomyPayIdNameData()[req.ProviderID]

	// TODO: 確認通知request的通道及支付方式是否傳送int型態，以及訂單ID是否需調整/新增結構
	// 轉換 CreateOrderRequest 為 payment_types.PaymentRequest
	paymentReq := payment_types.PaymentRequest{
		Id:         req.Id,
		Trade:      req.Trade_No, // 確保這些字段匹配
		AppTrade:   req.AppTrade,
		Amount:     req.Amount,
		Currency:   req.Currency,
		ProviderId: req.ProviderID,
		Provider:   providerName,
		Channel:    strconv.Itoa(req.PaymentType),
		ReturnURL:  req.ReturnURL,
	}

	// 發送支付請求到 channel
	notification.PaymentRequestChannel <- paymentReq

	// 設定超時機制，避免支付結果長時間未返回
	select {
	case result := <-responseChannel:
		notification.RemoveResponseChannel(req.Trade_No)
		return result, nil
	case <-time.After(10 * time.Second): // 超過 10 秒，返回超時
		notification.RemoveResponseChannel(req.Trade_No)
		return payment_types.PaymentTransaction{}, errors.New("payment timeout")
	}
}

// 授權結果通知
func NotifyAuthPayment(req payment_types.AuthCallBackReq) (payment_types.AuthCallBackRes, error) {
	// 建立 responseChannel
	PaymentOrderChannel := make(chan payment_types.AuthCallBackRes, 1)

	// 記錄 responseChannel，以便金流模組回傳結果
	notification.SetAuthCallBackChannel(req.OrderID, PaymentOrderChannel)

	// 取通道 ID
	// providerName := cachedata.GomyPayIdNameData()[req.ProviderID]

	// TODO: 確認通知request的通道及支付方式是否傳送int型態，以及訂單ID是否需調整/新增結構
	// 轉換 CreateOrderRequest 為 payment_types.PaymentRequest
	paymentReq := payment_types.AuthCallBackReq{
		SendType:    req.SendType,
		Result:      req.Result,
		RetMsg:      req.RetMsg,
		OrderID:     req.OrderID,
		EOrderNo:    req.EOrderNo,
		AvCode:      req.AvCode,
		ECur:        req.ECur,
		EMoney:      req.EMoney,
		EDate:       req.EDate,
		ETime:       req.ETime,
		EPayaccount: req.EPayaccount,
		LimitDate:   req.LimitDate,
		Code1:       req.Code1,
		Code2:       req.Code2,
		Code3:       req.Code3,
		PinCode:     req.PinCode,
		StoreType:   req.StoreType,
		CardLastNum: req.CardLastNum,
		StrCheck:    req.StrCheck,
	}

	// 發送支付請求到 channel
	notification.AuthCallBackChannel <- paymentReq

	// 設定超時機制，避免支付結果長時間未返回
	select {
	case result := <-PaymentOrderChannel:
		notification.RemoveAuthCallBackChannel(req.OrderID)
		return result, nil
	case <-time.After(10 * time.Second): // 超過 10 秒，返回超時
		notification.RemoveAuthCallBackChannel(req.OrderID)
		return payment_types.AuthCallBackRes{}, errors.New("payment timeout")
	}
}
