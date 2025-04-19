package payment

import (
	ipayment "pay-service/internal/payment"

	"github.com/gofiber/fiber/v2"
)

// HandleAuthCallBack 處理授權回調
func HandleAuthCallBack(provider string, c *fiber.Ctx) error {
	// 透過 Registry 取得金流商
	paymentProvider, _ := ipayment.GetProvider(provider)
	if paymentProvider == nil {
		return nil
	}
	return paymentProvider.HandleAuthCallBack(c)
}

// HandleCallback 處理金流回調
func HandleCallback(provider string, c *fiber.Ctx) error {
	// 透過 Registry 取得金流商
	paymentProvider, _ := ipayment.GetProvider(provider)
	if paymentProvider == nil {
		return nil
	}
	return paymentProvider.HandleCallback(c)
}
