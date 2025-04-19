package cachedata

import (
	"sync"
	"time"
)

type Config struct {
	RefreshDuration time.Duration
	RetryDuration   time.Duration
}

// cachePayTypeData 定義付款方式快取資料
type cacheGomyPayData struct {
	mu              sync.RWMutex
	Data            map[string]GomyType
	NextRefreshTime time.Time
}

// cacheGomyPayIdNameData 定義 Gomypay 快取資料
type cacheGomyPayIdNameData struct {
	mu              sync.RWMutex
	Data            map[int]string
	NextRefreshTime time.Time
}

// cacheGomyType 定義 Gomypay 快取資料
type GomyType struct {
	Id          int           `db:"id"`
	Name        string        `db:"name"`
	SecrectInfo GomypaySecret `db:"secrect_info"`
}

// GomypaySecret 定義 Gomypay 金流的參數
type GomypaySecret struct {
	CustomerID     string `json:"customer_id"`
	StrCheck       string `json:"str_check"`
	OrderUrl       string `json:"order_url"`
	CallOrderUrl   string `json:"call_order_url"`
	GoodReturnUrl  string `json:"good_return_url"`
	TestCardNo     string `json:"test_card_no"`
	TestCardExpire string `json:"test_card_expire"`
	TestCardCvv    string `json:"test_card_cvv"`
}
