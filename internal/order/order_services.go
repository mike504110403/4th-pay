package order

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"pay-service/internal/payment/payment_types"

	mlog "github.com/mike504110403/goutils/log"

	"github.com/valyala/fasthttp"
)

func AppCaller(order payment_types.PaymentOrderResponse) error {
	// 將訂單轉為 JSON 格式
	orderData, err := json.Marshal(order)
	if err != nil {
		mlog.Error(fmt.Sprintf("無法將訂單轉換為 JSON: %v", err))
		return err
	}

	// 建立 HTTP 請求
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(AppURL) // 更換為目標服務的 URL
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBody(orderData)

	resp := fasthttp.AcquireResponse()

	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		mlog.Error(fmt.Sprintf("請求失敗 %s", err.Error()))
		return err
	}

	return nil
}

// Order 訂單
func CreateOrder(tx *sql.DB, order Order) (int, error) {
	id := 0
	// 創建訂單
	queryStr := `
		INSERT INTO OrderRecord (
			app_id, 
			app_trade, 
			trade_no, 
			provider_id, 
			payment_type, 
			amount, 
			currency, 
			status,
			callback_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	if r, err := tx.Exec(
		queryStr,
		order.AppID,
		order.AppTrade,
		order.Trade_No,
		order.ProviderID,
		order.PaymentType,
		order.Amount,
		order.Currency,
		order.Status,
		order.CallbackURL,
	); err != nil {
		return id, err
	} else {
		if i, err := r.LastInsertId(); err != nil {
			return id, err
		} else {
			id = int(i)
		}
	}

	return id, nil
}

// UpdateOrderStatus： 更新訂單狀態
func UpdateOrderStatus(tx *sql.Tx, orderID string, status int) error {
	// 更新訂單狀態
	queryStr := `
		UPDATE OrderRecord SET status = ? WHERE trade_no = ?
	`
	_, err := tx.Exec(queryStr, status, orderID)
	if err != nil {
		tx.Rollback()
		mlog.Error(fmt.Sprintf("UpdateOrderStatus error: %v", err))
		return err
	}
	return tx.Commit()
}

func CallbackUpdate(tx *sql.Tx, orderID string) error {
	queryStr := `
		UPDATE OrderRecord SET is_callback = ? WHERE trade_no = ?
	`
	_, err := tx.Exec(queryStr, "Y", orderID)
	if err != nil {
		tx.Rollback()
		mlog.Error(fmt.Sprintf("CallbackUpdate error: %v", err))
		return err
	}
	return tx.Commit()
}

func CallbackCheck(db *sql.DB, orderID string) (isCallback bool, err error) {
	queryStr := `
		SELECT is_callback FROM OrderRecord WHERE trade_no = ?
	`

	err = db.QueryRow(queryStr, orderID).Scan(&isCallback)
	if err != nil {
		mlog.Error(fmt.Sprintf("CallbackCheck error: %v", err))
		return false, err
	}
	return isCallback, nil
}

func GetOrder(db *sql.DB, orderID string) (Order, error) {
	queryStr := `
		SELECT 
			id,
			app_id, 
			app_trade, 
			trade_no, 
			provider_id, 
			payment_type, 
			amount, 
			currency, 
			status,
			callback_url,
			create_time,
			update_time
		FROM OrderRecord WHERE trade_no = ?
	`

	var order Order
	err := db.QueryRow(queryStr, orderID).Scan(
		&order.Id,
		&order.AppID,
		&order.AppTrade,
		&order.Trade_No,
		&order.ProviderID,
		&order.PaymentType,
		&order.Amount,
		&order.Currency,
		&order.Status,
		&order.CallbackURL,
	)
	if err != nil {
		mlog.Error(fmt.Sprintf("GetOrder error: %v", err))
		return Order{}, err
	}
	return order, nil
}

// 確認App訂單是否重複
func CheckOrderRepeat(db *sql.DB, appTrade string) error {
	queryStr := `
		SELECT COUNT(*) FROM OrderRecord WHERE app_trade = ?
	`

	var count int
	err := db.QueryRow(queryStr, appTrade).Scan(&count)
	if err != nil {
		mlog.Error(fmt.Sprintf("CheckOrderRepeat error: %v", err))
		return err
	}
	if count > 0 {
		mlog.Error("Order already exists")
		return fmt.Errorf("Order already exists")
	} else {
		return nil
	}
}

func OrderSubTypeDescription(db *sql.DB, id int) (string, error) {
	queryStr := `
		SELECT Description FROM SubType WHERE Id = ?
	`

	var subType string
	err := db.QueryRow(queryStr, id).Scan(&subType)
	if err != nil {
		mlog.Error(fmt.Sprintf("OrderSubTypeDescription error: %v", err))
		return "", err
	}
	return subType, nil
}
