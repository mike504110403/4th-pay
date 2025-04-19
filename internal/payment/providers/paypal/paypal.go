package paypal

import (
	"log"
	"pay-service/internal/payment"
	"pay-service/internal/payment/payment_types"

	"github.com/gofiber/fiber/v2"
)

// PayPalProvider 實現 PaymentProvider 介面
type PayPalProvider struct{}

// CreateTransaction PayPal 的交易實作
func (p *PayPalProvider) NewTransaction(order *payment_types.PaymentTransaction) error {
	return nil
}

// HandleCallback PayPal 的回調實作
func (p *PayPalProvider) HandleCallback(c *fiber.Ctx) error {
	log.Printf("[PayPal] 處理回調: ")
	payment.NotifyOrder(payment_types.PaymentOrderRequest{})

	return nil
}

// HandleAuthCallBack PayPal 的授權回調實作
func (p *PayPalProvider) HandleAuthCallBack(c *fiber.Ctx) error {
	log.Printf("[PayPal] 處理授權回調: ")
	payment.NotifyOrder(payment_types.PaymentOrderRequest{})

	return nil
}
