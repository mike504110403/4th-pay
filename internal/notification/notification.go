package notification

import (
	"pay-service/internal/payment/payment_types"
	"sync"
)

// 訂單請求 channel（Order -> Payment）
var PaymentRequestChannel = make(chan payment_types.PaymentRequest, 100)

// 使用 `map` 來存儲等待支付結果的 channel
var responseChannelMap = make(map[string]chan payment_types.PaymentTransaction)
var mutex sync.Mutex

// 設定 responseChannel
func SetResponseChannel(orderTrade string, responseChannel chan payment_types.PaymentTransaction) {
	mutex.Lock()
	responseChannelMap[orderTrade] = responseChannel
	mutex.Unlock()
}

// 取得 responseChannel
func GetResponseChannel(orderTrade string) (chan payment_types.PaymentTransaction, bool) {
	mutex.Lock()
	channel, exists := responseChannelMap[orderTrade]
	mutex.Unlock()
	return channel, exists
}

// 移除 responseChannel，避免 memory leak
func RemoveResponseChannel(orderTrade string) {
	mutex.Lock()
	delete(responseChannelMap, orderTrade)
	mutex.Unlock()
}

// ----------------------------------------------------------------------------------------
// 支付狀態 channel（Payment -> order)
var PaymentOrderChannel = make(chan payment_types.PaymentOrderRequest, 100)

// 使用 `map` 來存儲等待支付結果的 channel
var orderChannelMap = make(map[string]chan payment_types.PaymentOrderResponse)
var orderMutex sync.Mutex

// 設定 orderChannel
func SetOrderChannel(orderTrade string, orderChannel chan payment_types.PaymentOrderResponse) {
	orderMutex.Lock()
	orderChannelMap[orderTrade] = orderChannel
	orderMutex.Unlock()
}

// 取得 orderChannel
func GetOrderChannel(orderTrade string) (chan payment_types.PaymentOrderResponse, bool) {
	orderMutex.Lock()
	channel, exists := orderChannelMap[orderTrade]
	orderMutex.Unlock()
	return channel, exists
}

// 移除 orderChannel，避免 memory leak
func RemoveOrderChannel(orderTrade string) {
	orderMutex.Lock()
	delete(orderChannelMap, orderTrade)
	orderMutex.Unlock()
}

// ----------------------------------------------------------------------------------------
// 訂單授權驗證通知 channel（Order -> Payment）
var AuthCallBackChannel = make(chan payment_types.AuthCallBackReq, 100)

// 使用 `map` 來存儲等待授權結果的 channel
var authCallBackChannelMap = make(map[string]chan payment_types.AuthCallBackRes)
var authMutex sync.Mutex

// 設定 authCallBackChannel
func SetAuthCallBackChannel(orderTrade string, authCallBackChannel chan payment_types.AuthCallBackRes) {
	authMutex.Lock()
	authCallBackChannelMap[orderTrade] = authCallBackChannel
	authMutex.Unlock()
}

// 取得 authCallBackChannel
func GetAuthCallBackChannel(orderTrade string) (chan payment_types.AuthCallBackRes, bool) {
	authMutex.Lock()
	channel, exists := authCallBackChannelMap[orderTrade]
	authMutex.Unlock()
	return channel, exists
}

// 移除 authCallBackChannel，避免 memory leak
func RemoveAuthCallBackChannel(orderTrade string) {
	authMutex.Lock()
	delete(authCallBackChannelMap, orderTrade)
	authMutex.Unlock()
}
