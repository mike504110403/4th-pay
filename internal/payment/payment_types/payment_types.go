package payment_types

// PaymentRequest 定義用於發起支付請求的結構
type PaymentRequest struct {
	Id         int     `json:"id"`        // App訂單 ID
	Trade      string  `json:"trade"`     // 訂單編號
	AppTrade   string  `json:"app_trade"` // 來自 app 的訂單編號
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	ProviderId int     `json:"provider_id"`
	Provider   string  `json:"provider"`
	Channel    string  `json:"channel"`
	ReturnURL  string  `json:"return_url"`
}

// PaymentTransaction 定義支付交易記錄
type PaymentTransaction struct {
	ID            int     `json:"id" `            // 交易 ID
	Trade         string  `json:"trade"`          // 訂單編號
	ProviderTrade string  `json:"provider_trade"` // 金流商的訂單編號
	Amount        float64 `json:"amount"`         // 交易金額
	Provider      string  `json:"provider"`       // 金流商名稱
	Channel       string  `json:"channel"`        // 交易通道
	PayUrl        string  `json:"pay_url"`        // 付款連結
	IsScuess      bool    `json:"status"`         // 交易狀態
	ReturnURL     string  `json:"return_url"`     // 交易回調 URL
}

type PaymentOrder string

// PaymentOrderRequest 定義訂單收款請求
type PaymentOrderRequest struct {
	Trade         string `json:"trade"`
	ProviderTrade string `json:"provider_trade"`
	IsScuess      bool   `json:"status"`
}

// PaymentOrderResponse 定義訂單收款回應
type PaymentOrderResponse struct {
	Trade         string `json:"trade"`
	ProviderTrade string `json:"provider_trade"`
	IsScuess      bool   `json:"order_status"`
}

// PayRecord 資料庫結構
type PayRecord struct {
	OrdeId        int    `db:"order_id"`
	ProviderTrade string `db:"provider_trade"`
	Status        int    `db:"status"`
	IsNotify      string `db:"is_notify"`
}

// ResResult : 回應結果
type ResResult string

const (
	ResResultFail    ResResult = "0" // 失敗
	ResResultSuccess ResResult = "1" // 成功
)

// AuthCallBack : Gomypay授權回調
type AuthCallBackReq struct {
	SendType    string    `query:"Send_Type"`              // 傳送型態 (0.信用卡 1.銀聯卡 2.超商條碼 3.WebAtm 4.虛擬帳號 6.超商代碼 7.LinePay)
	Result      ResResult `query:"result"`                 // 回傳結果 (0 失敗 1 成功)
	RetMsg      string    `query:"ret_msg"`                // 回傳訊息
	OrderID     string    `query:"OrderID"`                // 系統訂單編號
	EOrderNo    string    `query:"e_orderno"`              // 自訂訂單編號
	AvCode      string    `query:"avcode,omitempty"`       // 授權碼
	ECur        string    `query:"e_Cur,omitempty"`        // 幣別
	EMoney      string    `query:"e_money,omitempty"`      // 交易金額
	EDate       string    `query:"e_date,omitempty"`       // 交易日期(yyyymmdd)
	ETime       string    `query:"e_time,omitempty"`       // 交易時間(HH:mm:ss)
	EPayaccount string    `query:"e_payaccount,omitempty"` // 繳費帳號
	LimitDate   string    `query:"LimitDate,omitempty"`    // 繳費期限(yyyyMMdd)
	Code1       string    `query:"code1,omitempty"`        // 超商繳費條碼第一段
	Code2       string    `query:"code2,omitempty"`        // 超商繳費條碼第二段
	Code3       string    `query:"code3,omitempty"`        // 超商繳費條碼第三段
	PinCode     string    `query:"PinCode,omitempty"`      // 繳費代碼
	StoreType   string    `query:"StoreType,omitempty"`    // 超商代碼
	CardLastNum string    `query:"CardLastNum,omitempty"`  // 信用卡號後四碼
	StrCheck    string    `query:"str_check,omitempty"`    // 交易驗證密碼
}

type AuthCallBackRes struct {
	Trade_no    string `json:"trade_no"`
	App_id      int    `json:"app_id"`
	Provider_id int    `json:"provider_id"`
	IsScuess    bool   `json:"order_status"`
}
