package gomypay

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"pay-service/internal/cachedata"
	"pay-service/internal/database"
	apitools "pay-service/utils/apiTools"
	jsonformat "pay-service/utils/jsonFormat"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	mlog "github.com/mike504110403/goutils/log"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

var gomypaySecret cachedata.GomypaySecret

func Init() {
	os.Setenv("Environment", "dev")
	os.Setenv("Port", "8080")
	os.Setenv("MYSQL_USERNAME", "admin")
	os.Setenv("MYSQL_PASSWORD", "*5reRKAHRw5k*b3r")
	os.Setenv("MYSQL_HOST", "dev-tapio-pay.c9ki2kumaxg5.ap-northeast-3.rds.amazonaws.com")
	// 初始化log
	mlog.Init(mlog.Config{
		EnvMode: mlog.EnvMode("dev"),
		LogType: mlog.LogType("console"),
	})

	// 初始化
	database.Init()
	cachedata.Init(cachedata.Config{
		RefreshDuration: time.Minute * 3,
		RetryDuration:   time.Second * 3,
	})
	gomypaySecret = cachedata.GomyPayData()["gomypay"].SecrectInfo
}

// 測試交易遞交
func TestTransaction(t *testing.T) {
	Init()
	u, _ := uuid.NewV4()
	trade := base64.RawURLEncoding.EncodeToString(u.Bytes())
	fmt.Printf("Trade: %s\n", trade)

	// 模擬交易參數
	reqData := GomypayRequest{
		SendType:    SendTypeCreditCard,
		PayModeNo:   PAY_MODE,
		CustomerID:  gomypaySecret.CustomerID,
		OrderNo:     trade,
		Amount:      "35",
		StoreType:   "1",
		TransCode:   TRANS_CODE,      // 交易類別: 00=授權
		TransMode:   TransModeNormal, // 交易模式: 1=一般; 2=分期
		Installment: "0",             // 分期期數

		BuyerName: "Tapio",
		BuyerTelm: "0912345678",
		BuyerMail: "tapiotest@gmail.com",
		BuyerMemo: "Tapio Test",

		// === 無填寫進預設頁面 ===
		// CardNo:     gomypaySecret.TestCardNo,
		// ExpireDate: gomypaySecret.TestCardExpire,
		// CVV:        gomypaySecret.TestCardCvv,
		// EReturn:     "1",
		ReturnURL:   "https://ba12-220-132-92-168.ngrok-free.app/api/public/v1/payment/gomypay/authCallback",
		CallbackURL: "https://ba12-220-132-92-168.ngrok-free.app/api/public/v1/payment/gomypay/backCallback",
		// === 無填寫進預設頁面 ===

		// StrCheck: gomypaySecret.StrCheck,
	}

	// post 方式送單
	// formValues := apitools.BuildFormValues(reqData)
	// resp, err := apitools.PostForm(gomypaySecret.OrderUrl, formValues)
	// defer fasthttp.ReleaseResponse(resp)

	// fmt.Printf("Response: %s\n", resp.Body())
	// assert.NoError(t, err, "請求失敗")

	// get 方式取得 URL
	url := apitools.ProcessUrl(gomypaySecret.OrderUrl, reqData)
	fmt.Printf("Pay URL: %s\n", url)
}

// 測試多筆交易
func TestMutiTransaction(t *testing.T) {
	for i := 1; i <= 10; i++ {
		t.Run(fmt.Sprintf("第%d筆測試：", i), func(t *testing.T) {
			u, _ := uuid.NewV4()
			trade := base64.RawURLEncoding.EncodeToString(u.Bytes())
			fmt.Printf("Trade: %s\n", trade)

			reqData := GomypayRequest{
				SendType:    SendTypeCreditCard,
				PayModeNo:   PAY_MODE,
				CustomerID:  gomypaySecret.CustomerID,
				OrderNo:     trade,
				Amount:      "35",
				StoreType:   "1",
				TransCode:   TRANS_CODE,
				TransMode:   TransModeNormal,
				Installment: "0",

				BuyerName: "Tapio",
				BuyerTelm: "0912345678",
				BuyerMail: "tapiotest@gmail.com",
				BuyerMemo: fmt.Sprintf("第%d筆測試：", i),

				CardNo:      gomypaySecret.TestCardNo,
				ExpireDate:  gomypaySecret.TestCardExpire,
				CVV:         gomypaySecret.TestCardCvv,
				EReturn:     "1",
				CallbackURL: "https://17f3-220-132-92-168.ngrok-free.app/api/public/v1/payment/gomypay/backCallback",
				StrCheck:    gomypaySecret.StrCheck,
			}

			formValues := apitools.BuildFormValues(reqData)
			resp, err := apitools.PostForm(gomypaySecret.OrderUrl, formValues)
			defer fasthttp.ReleaseResponse(resp)

			fmt.Printf("[第%d筆測試] Response: %s\n", i, resp.Body())
			assert.NoError(t, err, "請求失敗")
		})
	}
}

// TestGetTrade : 取得單筆資料
func TestGetTrade(t *testing.T) {
	req := CallOrderReq{
		OrderNo:    "ZH_ZT_JpS0ufamxlbsCq7w",
		CustomerId: gomypaySecret.CustomerID,
		StrCheck:   gomypaySecret.StrCheck,
	}

	// 動態生成表單資料
	formValues := apitools.BuildFormValues(req)

	// 解析回應
	resp, err := apitools.PostForm(gomypaySecret.CallOrderUrl, formValues)
	// 解析JSON
	resData := CallOrderRes{}
	if err := json.Unmarshal([]byte(resp.Body()), &resData); err != nil {
		fmt.Println("解析失敗:", err)
		return
	}
	assert.NoError(t, err, "請求失敗")

	// 印 res json
	fmt.Printf("Response:\n%s\n", jsonformat.PrettyJSON(resData))

	defer fasthttp.ReleaseResponse(resp)
}

// TestGetTrades : 取得多筆資料
func TestGetTrades(t *testing.T) {
	req := CallOrderReq{
		CustomerId: gomypaySecret.CustomerID,
		StrCheck:   gomypaySecret.StrCheck,
		CreatSdate: "20250218",
		CreatEdate: "20250219",
		CreatStime: "0",
		CreatEtime: "23",
		PaySdate:   "20250218",
		PayEdate:   "20250219",
		PayStime:   "0",
		PayEtime:   "23",
	}

	// 動態生成表單資料
	formValues := apitools.BuildFormValues(req)

	// 解析回應
	resp, err := apitools.PostForm(gomypaySecret.CallOrderUrl, formValues)
	assert.NoError(t, err, "請求失敗")
	fmt.Printf("Response: %s\n", string(resp.Body()))
	// 解析JSON
	resData := CallOrdersRes{}
	if err := json.Unmarshal([]byte(resp.Body()), &resData); err != nil {
		fmt.Println("解析失敗:", err)
		return
	}

	// 印 res json
	prettyJSON, err := json.MarshalIndent(resData, "", "  ")
	if err != nil {
		t.Fatalf("格式化 JSON 失敗: %v", err)
	}
	fmt.Printf("Response:\n%s\n", prettyJSON)

	defer fasthttp.ReleaseResponse(resp)
}

// TestGoodReturn : 測試退貨
func TestGoodReturn(t *testing.T) {
	reqData := GoodsReturnReq{
		OrderNo:      "ZH_ZT_JpS0ufamxlbsCq7w",
		CustomerId:   gomypaySecret.CustomerID,
		StrCheck:     gomypaySecret.StrCheck,
		Goods_Return: "1",
		// Goods_Return_Cancel: "1",
		Goods_Return_Reason: "商品退貨",
	}

	// post 方式送單
	formValues := apitools.BuildFormValues(reqData)
	resp, err := apitools.PostForm(gomypaySecret.GoodReturnUrl, formValues)
	fmt.Printf("Response: %s\n", resp)
	assert.NoError(t, err, "請求失敗")
	defer fasthttp.ReleaseResponse(resp)
}
