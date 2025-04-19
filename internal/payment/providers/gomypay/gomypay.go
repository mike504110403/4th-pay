package gomypay

import (
	"database/sql"
	"errors"
	"fmt"
	"pay-service/internal/cachedata"
	"pay-service/internal/payment"
	"pay-service/internal/payment/payment_types"
	apitools "pay-service/utils/apiTools"
	jsonformat "pay-service/utils/jsonFormat"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mike504110403/common-moduals/typeparam"
)

// GomyPayProvider 實現 PaymentProvider 介面
type GomyPayProvider struct{}

// CreateTransaction PayPal 的交易實作
func (p *GomyPayProvider) NewTransaction(order *payment_types.PaymentTransaction) error {
	gomypaySecret := cachedata.GomyPayData()["gomypay"].SecrectInfo
	// 共用參數
	gomypayReq := GomypayRequest{
		CustomerID: gomypaySecret.CustomerID,
		OrderNo:    order.Trade,
		Amount:     strconv.FormatFloat(order.Amount, 'f', -1, 64),
		PayModeNo:  PAY_MODE,
		ReturnURL:  order.ReturnURL,
		// TODO: 背景對帳網址
		// CallbackURL: ,
	}

	// 依據交易通道設定不同的參數
	switch SendType(order.Channel) {
	case SendTypeCreditCard:
		gomypayReq.SendType = SendTypeCreditCard
		gomypayReq.TransCode = TRANS_CODE
		gomypayReq.TransMode = TransModeNormal
		gomypayReq.Installment = "0"
	case SendTypeUnionPay:
		gomypayReq.SendType = SendTypeUnionPay
		gomypayReq.TransCode = TRANS_CODE
	case SendTypeBarcode:
		gomypayReq.SendType = SendTypeBarcode
	case SendTypeWebAtm:
		gomypayReq.SendType = SendTypeWebAtm
	case SendTypeVirtual:
		gomypayReq.SendType = SendTypeVirtual
	case SendTypeStore:
		gomypayReq.SendType = SendTypeStore
	case SendTypeLinePay:
		gomypayReq.SendType = SendTypeLinePay
	default:
		return errors.New("不支援的交易通道")
	}
	// 付款 url API URL
	order.PayUrl = apitools.ProcessUrl(gomypaySecret.OrderUrl, gomypayReq)

	order.IsScuess = true
	if err := payment.CreatePaymentTransaction(order); err != nil {
		return err
	}
	return nil
}

// HandleAuthCallBack PayPal 的授權回調實作
func (p *GomyPayProvider) HandleAuthCallBack(c *fiber.Ctx) error {
	queryData := AuthCallBack{}
	if err := c.QueryParser(&queryData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Query解析失敗")
	}
	fmt.Printf("Response:\n%s\n", jsonformat.PrettyJSON(queryData))
	// 查無此單 || 此單已結束
	payOrderId, err := payment.GetPayOrderId(queryData.EOrderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusBadRequest).SendString("訂單不存在")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("服務錯誤")
	}

	payOrder := payment_types.PayRecord{
		OrdeId: payOrderId,
	}
	// 訂單狀態碼
	payType := typeparam.TypeParam{MainType: "pay_state"}
	switch queryData.Result {
	case ResResultSuccess:
		payType.SubType = "success"
	case ResResultFail:
		payType.SubType = "fail"
	}
	typeInt, _ := payType.Get()
	payOrder.Status = typeInt
	// 更新支付狀態
	if err := payment.UpdatePayOrderStatus(payOrder); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("服務錯誤")
	}

	// 通知訂單
	notify := payment_types.PaymentOrderRequest{
		Trade:         queryData.EOrderNo,
		ProviderTrade: queryData.OrderID,
		IsScuess:      true,
	}
	switch queryData.Result {
	case ResResultSuccess:
		notify.IsScuess = true
	case ResResultFail:
		notify.IsScuess = false
	}

	// 併發通知
	go func(thisNotify payment_types.PaymentOrderRequest, orderId int) {
		if _, err := payment.NotifyOrder(thisNotify); err == nil {
			payment.UpdatePayOrderIsNotify(orderId)
		}
		return
	}(notify, payOrder.OrdeId)

	return nil
}

// TODO: 背景對帳
// HandleCallback PayPal 的回調實作
func (p *GomyPayProvider) HandleCallback(c *fiber.Ctx) error {
	req := BackCallBack{}
	// 將post body 自動解析進 struct
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	notify := payment_types.PaymentOrderRequest{
		Trade:         req.EOrderNo,
		ProviderTrade: req.OrderID,
	}

	notify.IsScuess = true

	payment.NotifyOrder(notify)
	return c.SendStatus(fiber.StatusOK)
}
