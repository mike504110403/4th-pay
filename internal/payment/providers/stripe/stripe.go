package stripe

import (
	"log"
	"pay-service/internal/payment"
	"pay-service/internal/payment/payment_types"

	"github.com/gofiber/fiber/v2"
)

// StripeProvider 實現 PaymentProvider 介面
type StripeProvider struct{}

// CreateTransaction Stripe 的交易實作
func (s *StripeProvider) NewTransaction(order *payment_types.PaymentTransaction) error {
	return nil
}

// HandleCallback Stripe 的回調實作
func (s *StripeProvider) HandleCallback(c *fiber.Ctx) error {
	log.Printf("[Stripe] 處理回調: ")
	payment.NotifyOrder(payment_types.PaymentOrderRequest{})
	return nil
}

// HandleAuthCallBack Stripe 的授權回調實作
func (s *StripeProvider) HandleAuthCallBack(c *fiber.Ctx) error {
	log.Printf("[Stripe] 處理授權回調: ")
	payment.NotifyOrder(payment_types.PaymentOrderRequest{})
	return nil
}
