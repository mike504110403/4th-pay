package database

import (
	"github.com/mike504110403/goutils/dbconn"
)

// 以下為連線設定範例 需調整
// 連線字符定義在.env
const Envkey string = "MYSQL_URL"

// mysql使用的參數設定
const PAYORDER dbconn.DBName = "PayOrder"
const PROVIDER dbconn.DBName = "Provider"
const SETTING dbconn.DBName = "Setting"

// 組裝用字串
const PayOrder_dsn string = "PayOrder"
const Provider_dsn string = "Provider"
const Setting_dsn string = "Setting"

var DB_Name_Map = map[dbconn.DBName]string{
	PAYORDER: PayOrder_dsn,
	PROVIDER: Provider_dsn,
	SETTING:  Setting_dsn,
}
