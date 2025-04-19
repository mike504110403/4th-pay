package cachedata

import (
	"time"
)

var cfg = Config{
	RefreshDuration: time.Minute * 3,
	RetryDuration:   time.Second * 3,
}

var gomyPayMap = cacheGomyPayData{
	Data:            map[string]GomyType{},
	NextRefreshTime: time.Now(),
}

var gomyPayIdNameMap = cacheGomyPayIdNameData{
	Data:            map[int]string{},
	NextRefreshTime: time.Now(),
}

func Init(initCfg Config) {
	cfg = initCfg
	loadGomyData()
	go refreshDataPeriodically()
}

// 定時更新快取資料
func refreshDataPeriodically() {
	for {
		time.Sleep(cfg.RefreshDuration)
		loadGomyData()
	}
}

// 取得 GomyPay 快取資料
func GomyPayData() map[string]GomyType {
	gomyPayMap.mu.RLock()
	defer gomyPayMap.mu.RUnlock()
	return gomyPayMap.Data
}

// 取得 GomyPayIdName 快取資料
func GomyPayIdNameData() map[int]string {
	gomyPayIdNameMap.mu.RLock()
	defer gomyPayIdNameMap.mu.RUnlock()
	return gomyPayIdNameMap.Data
}
