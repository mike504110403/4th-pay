package cachedata

import (
	"encoding/json"
	"fmt"
	"pay-service/internal/database"
	jsonformat "pay-service/utils/jsonFormat"
	"time"

	mlog "github.com/mike504110403/goutils/log"
)

// 付款方式快取資料
func loadGomyData() {
	if db, err := database.SETTING.DB(); err != nil {
		mlog.Error(err.Error())
	} else {
		queryStr := `
			SELECT id, name, secrect_info
			FROM Providers
			WHERE is_enable = 1
		`

		if rows, err := db.Query(queryStr); err != nil {
			mlog.Error(err.Error())
		} else {
			gomyPayMap.mu.Lock()
			defer gomyPayMap.mu.Unlock()
			gomyPayIdNameMap.mu.Lock()
			defer gomyPayIdNameMap.mu.Unlock()
			gomyPayMap.Data = make(map[string]GomyType)
			gomyPayIdNameMap.Data = make(map[int]string)
			for rows.Next() {
				var data GomyType
				var secrectInfoStr string // 先用字串接收 JSON

				if err := rows.Scan(&data.Id, &data.Name, &secrectInfoStr); err != nil {
					mlog.Error(err.Error())
					continue
				}

				// 解析 JSON
				if err := json.Unmarshal([]byte(secrectInfoStr), &data.SecrectInfo); err != nil {
					mlog.Error(fmt.Sprintf("JSON 解析失敗: %s, err: %v", secrectInfoStr, err))
					continue
				}

				gomyPayMap.Data[data.Name] = data
				gomyPayIdNameMap.Data[data.Id] = data.Name
				fmt.Printf("secrect: %s\n", jsonformat.PrettyJSON(data))
			}
			gomyPayMap.NextRefreshTime = time.Now().Add(cfg.RefreshDuration)
		}
	}
}
