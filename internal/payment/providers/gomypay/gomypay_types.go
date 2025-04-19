package gomypay

// SendType : 傳送型態
type SendType string

const (
	SendTypeCreditCard SendType = "0" // 信用卡
	SendTypeUnionPay   SendType = "1" // 銀聯卡
	SendTypeBarcode    SendType = "2" // 超商條碼
	SendTypeWebAtm     SendType = "3" // WebAtm
	SendTypeVirtual    SendType = "4" // 虛擬帳號
	SendTypeStore      SendType = "6" // 超商代碼
	SendTypeLinePay    SendType = "7" // LinePay
)

// TransMode : 交易模式
type TransMode string

const (
	TransModeNormal      TransMode = "1" // 一般交易
	TransModeInstallment TransMode = "2" // 分期交易
)

const (
	PAY_MODE   = "2"  // 付款模式
	TRANS_CODE = "00" // 交易類別
)

// ResResult : 回應結果
type ResResult string

const (
	ResResultFail    ResResult = "0" // 失敗
	ResResultSuccess ResResult = "1" // 成功
)

// GomypayRequest : Gomypay信用卡交易請求參數
type GomypayRequest struct {
	SendType    SendType  `form:"Send_Type" url:"Send_Type,omitempty"`       // 傳送型態 (0:信用卡, 1:銀聯卡, 2:超商條碼, 3:WebAtm, 4:虛擬帳號, 6:超商代碼)
	PayModeNo   string    `form:"Pay_Mode_No" url:"Pay_Mode_No,omitempty"`   // 付款模式，信用卡固定為 "2"
	CustomerID  string    `form:"CustomerId" url:"CustomerId,omitempty"`     // 商店代號，由 Gomypay 提供
	OrderNo     string    `form:"Order_No" url:"Order_No,omitempty"`         // 交易單號，商店自行產生須唯一
	Amount      string    `form:"Amount" url:"Amount,omitempty"`             // 交易金額，單位: 元 (新台幣)
	TransCode   string    `form:"TransCode" url:"TransCode,omitempty"`       // 交易類別，信用卡固定為 "00"
	BuyerName   string    `form:"Buyer_Name" url:"Buyer_Name,omitempty"`     // 消費者姓名
	BuyerTelm   string    `form:"Buyer_Telm" url:"Buyer_Telm,omitempty"`     // 消費者電話 (手機號碼)
	BuyerMail   string    `form:"Buyer_Mail" url:"Buyer_Mail,omitempty"`     // 消費者 Email
	BuyerMemo   string    `form:"Buyer_Memo" url:"Buyer_Memo,omitempty"`     // 商品或備註資訊
	CardNo      string    `form:"CardNo" url:"CardNo,omitempty"`             // 信用卡號 (信用卡交易必填)
	ExpireDate  string    `form:"ExpireDate" url:"ExpireDate,omitempty"`     // 信用卡有效日期(YYMM)
	CVV         string    `form:"CVV" url:"CVV,omitempty"`                   // 信用卡背面3碼
	TransMode   TransMode `form:"TransMode" url:"TransMode,omitempty"`       // 交易模式，1:一般交易，2:分期交易
	Installment string    `form:"Installment" url:"Installment,omitempty"`   // 分期數，0:不分期，3:3期，6:6期，12:12期
	ReturnURL   string    `form:"Return_url" url:"Return_url,omitempty"`     // 授權結果回傳網址
	CallbackURL string    `form:"Callback_Url" url:"Callback_Url,omitempty"` // 背景對帳網址
	EReturn     string    `form:"e_return" url:"e_return,omitempty"`         // 是否使用 JSON 回傳，1: 使用 JSON，0: 使用純文字
	StrCheck    string    `form:"Str_Check" url:"Str_Check,omitempty"`       // 驗證碼，由 Gomypay 提供 (確保交易安全)
	StoreType   string    `form:"StoreType" url:"StoreType,omitempty"`       // 超商代碼類別 (超商交易專用)
}

// AuthCallBack : Gomypay授權回調
type AuthCallBack struct {
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

type (
	// CallOrderReq : Gomypay訂單查詢請求
	CallOrderReq struct {
		// 必填
		CustomerId string `form:"CustomerId" url:"CustomerId,omitempty"` // 商店代號
		StrCheck   string `form:"Str_Check" url:"Str_Check,omitempty"`   // 驗證碼
		// 以下為選填 (單筆查詢)
		OrderNo string `form:"Order_No" url:"Order_No,omitempty"` // 交易單號(我方自訂編號)
		// 以下為選填 (區間查詢)
		CreatSdate string `form:"CreatSdate" url:"CreatSdate,omitempty"` // 訂單建立日期(yyyyMMdd)
		CreatEdate string `form:"CreatEdate" url:"CreatEdate,omitempty"` // 訂單建立日期(yyyyMMdd)
		CreatStime string `form:"CreatStime" url:"CreatStime,omitempty"` // 訂單建立時間 (24小時制 0~23)
		CreatEtime string `form:"CreatEtime" url:"CreatEtime,omitempty"` // 訂單建立時間 (24小時制 0~23)
		PaySdate   string `form:"PaySdate" url:"PaySdate,omitempty"`     // 繳費日期(yyyyMMdd)
		PayEdate   string `form:"PayEdate" url:"PayEdate,omitempty"`     // 繳費日期(yyyyMMdd)
		PayStime   string `form:"PayStime" url:"PayStime,omitempty"`     // 繳費時間 (24小時制 0~23)
		PayEtime   string `form:"PayEtime" url:"PayEtime,omitempty"`     // 繳費時間 (24小時制 0~23)
	}

	// CallOrderRes : Gomypay訂單查詢回應
	CallOrdersRes struct {
		Check    string         `json:"check"`     // 查詢結果(0:失敗/1:成功)
		CheckMsg string         `json:"check_msg"` // 查詢結果訊息
		Order    []CallOrderRes `json:"Order"`     // 訂單資訊
	}
	CallOrderRes struct {
		Result           string `json:"result"`             // 長度(1) 交易結果(0:失敗/1:成功/2:待付款/3:交易中斷)
		RetMsg           string `json:"ret_msg"`            // 最大長度(100) 交易結果訊息
		OrderID          string `json:"OrderID"`            // 長度(19) 訂單編號
		ECur             string `json:"e_Cur"`              // 長度(2) 幣別
		EMoney           string `json:"e_money"`            // 最大長度(10) 交易金額
		PayAmount        string `json:"PayAmount"`          // 最大長度(10) 實際繳費金額
		EDate            string `json:"e_date"`             // 長度(8) 訂單建立日期(yyyyMMdd)
		ETime            string `json:"e_time"`             // 長度(8) 訂單建立時間(HH:mm:ss)
		PDate            string `json:"p_date"`             // 長度(8) 繳費日期(yyyyMMdd)
		PTime            string `json:"p_time"`             // 長度(8) 繳費時間(HH:mm:ss)
		EOrderNo         string `json:"e_orderno"`          // 最大長度(25) 自訂訂單編號
		ENo              string `json:"e_no"`               // 最大長度(11) 商店代號(等同於 CustomerId )
		EOutlay          string `json:"e_outlay"`           // 最大長度(10) 交易總手續費
		BankName         string `json:"bankname"`           // 最大長度(50) 閘道銀行
		AvCode           string `json:"avcode"`             // 最大長度(10) 授權碼
		BuyerName        string `json:"Buyer_Name"`         // 最大長度(20) 消費者名稱
		BuyerMail        string `json:"Buyer_Mail"`         // 最大長度(50) 消費者聯絡信箱
		BuyerTelm        string `json:"Buyer_Telm"`         // 最大長度(20) 消費者聯絡電話
		BuyerMemo        string `json:"Buyer_Memo"`         // 最大長度(500) 消費備註(交易內容)
		CreditcardNo     string `json:"Creditcard_No"`      // 最大長度(20) 信用卡號碼(前六後四中間打*號)
		Installment      string `json:"Installment"`        // 最大長度(3) 期數(交易類型為信用卡時，無分期則填 0)
		ShopPaymentCode  string `json:"Shop_PaymentCode"`   // 最大長度(20) 超商代碼
		VirtualAccount   string `json:"Virtual_Account"`    // 最大長度(20) 虛擬帳號、超商條碼(第二段條碼)、WebAtm(第二段條碼)
		EPayInfo         string `json:"e_PayInfo"`          // 長度(9) 帳號繳費資訊(銀行代號三碼+,+繳費帳號後五碼)
		SendType         string `json:"Send_Type"`          // 長度(1) 傳送型態(0.信用卡 1.銀聯卡 2.超商條碼 3.WebAtm 4.虛擬帳號 5.定期扣款 6.超商代碼 7.LinePay)
		GoodsReturn      string `json:"Goods_Return"`       // 長度(1) 退貨、取消交易註記(0:無退貨狀態/1:申請退貨)
		GoodsReturnStatu string `json:"Goods_Return_Statu"` // 最大長度(100) 退貨處理訊息
		PayResult        string `json:"pay_result"`         // 長度(1) 回傳付款情況 (0 未付款 1 己付款)
		LimitDate        string `json:"LimitDate"`          // 長度(8) 繳費期限(yyyyMMdd)
		MarketID         string `json:"Market_ID"`          // 長度(2) FM：全家, OK：OK, HL：萊爾富, SE：7-11
		ShopStoreName    string `json:"Shop_Store_Name"`    // 最大長度(100) 繳費門市+(門市地址)
	}
)

// GoodsReturnReq : Gomypay退貨請求
type GoodsReturnReq struct {
	OrderNo             string `form:"Order_No" url:"Order_No,omitempty"`                       // 交易單號(我方自訂編號)
	CustomerId          string `form:"CustomerId" url:"CustomerId,omitempty"`                   // 商店代號
	StrCheck            string `form:"Str_Check" url:"Str_Check,omitempty"`                     // 驗證碼
	Goods_Return        string `form:"Goods_Return" url:"Goods_Return,omitempty"`               // 退貨、取消交易註記(1:申請退貨)
	Goods_Return_Reason string `form:"Goods_Return_Reason" url:"Goods_Return_Reason,omitempty"` // 退貨原因
	Goods_Return_Cancel string `form:"Goods_Return_Cancel" url:"Goods_Return_Cancel,omitempty"` // 退貨取消註記(1:取消退貨)
}

// 背景對帳回調
type BackCallBack struct {
	SendType    SendType  `form:"Send_Type"`
	Result      ResResult `form:"result"`      // 回傳結果 (0:失敗, 1:成功)
	RetMsg      string    `form:"ret_msg"`     // 回傳訊息
	OrderID     string    `form:"OrderID"`     // 系統訂單編號
	ECur        string    `form:"e_Cur"`       // 幣別
	EMoney      string    `form:"e_money"`     // 交易金額
	PayAmount   string    `form:"PayAmount"`   // 實際繳費金額
	EDate       string    `form:"e_date"`      // 交易日期 (yyyymmdd)
	ETime       string    `form:"e_time"`      // 交易時間 (HH:mm:ss)
	EOrderNo    string    `form:"e_orderno"`   // 自訂訂單編號
	ENo         string    `form:"e_no"`        // 商店代號
	EOutlay     string    `form:"e_outlay"`    // 交易總手續費
	AvCode      string    `form:"avcode"`      // 授權碼
	StrCheck    string    `form:"str_check"`   // 商店驗證碼
	CardLastNum string    `form:"CardLastNum"` // 信用卡號後四碼
	// 超商繳費相關 虛擬帳號相關
	EPayAccount string `form:"e_payaccount"` // 繳費帳號
	// 超商代碼相關
	PinCode       string `form:"PinCode"`         // 繳費代碼
	Barcode2      string `form:"Barcode2"`        // 第二段序號
	MarketID      string `form:"Market_ID"`       // FM：全家, OK：OK, HL：萊爾富, SE：7-11
	ShopStoreName string `form:"Shop_Store_Name"` // 繳費門市 + (門市地址)
	Invoice_No    string `form:"Invoice_No"`      // 發票號碼
}
