package order

// TODO: 這裡是訂單的相關請求資料結構
// app 過來的訂單請求
type CreateOrderRequest struct {
	AppID       int          `json:"app_id" example:"1" validate:"required"`                              // 代理商 ID
	AppTrade    string       `json:"app_trade" example:"123456" validate:"required"`                      // 來自 app 的訂單編號
	Amount      float64      `json:"amount" example:"99.99" validate:"required"`                          // 金額
	Currency    string       `json:"currency" example:"USD" validate:"required"`                          // 幣別
	Status      int          `json:"status" example:"paid"`                                               // 訂單狀態
	ProviderID  int          `json:"provider_id" example:"1" validate:"required"`                         // 金流商 ID
	PaymentType int          `json:"payment_type" example:"paypal" validate:"required"`                   // 付款通道
	CallbackURL string       `json:"callback_url" example:"http://test.com/callback" validate:"required"` // 回調 URL
	ReturnURL   string       `json:"return_url" example:"http://test.com/return"`                         // 前端返回網址
	Data        *interface{} `json:"data" validate:"require"`                                             // POST資料
}
