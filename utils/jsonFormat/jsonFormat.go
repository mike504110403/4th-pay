package jsonformat

import (
	"encoding/json"
	"fmt"
)

// PrettyJSON 格式化 JSON
func PrettyJSON(obj interface{}) string {
	// PrettyJSON
	prettyJSON, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Printf("格式化 JSON 失敗: %v", err)
	}
	return string(prettyJSON)
}
