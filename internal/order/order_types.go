package order

type AppOrderRes struct {
	TradeNo string `json:"trade_no"`
	Status  string `json:"status"`
}

const (
	AppURL = "http://localhost:8080"
)

// TODO: 這裡是訂單的相關資料結構 記在中心訂單資料庫的
// Order 訂單結構
type Order struct {
	Id          int     `json:"id" example:"1" validate:"required" db:"id"`                                           // 訂單 ID
	AppID       int     `json:"app_id" example:"1" validate:"required" db:"service_id"`                               // 代理商 ID
	Trade_No    string  `json:"trade_no" example:"123456" validate:"required" db:"trade_no"`                          // 來自 app 的訂單編號
	AppTrade    string  `json:"app_trade" example:"123456" validate:"required" db:"app_trade"`                        // 來自 app 的訂單編號
	Amount      float64 `json:"amount" example:"99.99" validate:"required" db:"amount"`                               // 金額
	Currency    string  `json:"currency" example:"TWD" validate:"required" db:"currency"`                             // 幣別
	Status      int     `json:"status" example:"0" validate:"required" db:"status"`                                   // 訂單狀態
	ProviderID  int     `json:"provider_id" example:"1" validate:"required" db:"provider_id"`                         // 金流商 ID
	PaymentType int     `json:"payment_type" example:"paypal" validate:"required" db:"payment_type"`                  // 付款通道
	CallbackURL string  `json:"callback_url" example:"http://test.com/callback" validate:"require" db:"callback_url"` // 回調 URL
	ReturnURL   string  `json:"return_url" example:"http://test.com/return" `                                         // 前端返回網址                                      // 返回 URL
	CreateTime  string  `json:"create_time" example:"2021-01-01 00:00:00" db:"create_time"`                           // 建立時間
	UpdateTime  string  `json:"update_time" example:"2021-01-01 00:00:00" db:"update_time"`                           // 更新時間
}
