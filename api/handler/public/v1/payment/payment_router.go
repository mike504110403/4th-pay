package payment

import (
	"github.com/gofiber/fiber/v2"
)

// SetRouter 設定路由
func SetRouter(router fiber.Router) {
	g := router.Group("/payment")
	{
		g.Post("/paypal/callback", paypalCallbackHandler)
		g.Post("/stripe/callback", stripeCallbackHandler)
		g.Get("/gomypay/authCallback", gomypayAuthCallbackHandler)
		g.Post("/gomypay/backCallback", gomypayBackCallbackHandler)
	}
}

// paypalCallbackHandler 處理 Paypal 的回調
func paypalCallbackHandler(c *fiber.Ctx) error {
	HandleCallback("paypal", c)
	return c.SendStatus(fiber.StatusOK)
}

// stripeCallbackHandler 處理 Stripe 的回調
func stripeCallbackHandler(c *fiber.Ctx) error {
	HandleCallback("stripe", c)
	return c.SendStatus(fiber.StatusOK)
}

// gomypayCallbackHandler 處理 Gomypay 的授權回調
func gomypayAuthCallbackHandler(c *fiber.Ctx) error {
	HandleAuthCallBack("gomypay", c)
	// gomypay.HandleAuthCallBack(c)
	return c.SendStatus(fiber.StatusOK)
}

// gomypayCallbackHandler 處理 Gomypay 的背景對帳回調
func gomypayBackCallbackHandler(c *fiber.Ctx) error {
	HandleCallback("gomypay", c)
	return c.SendStatus(fiber.StatusOK)
}
