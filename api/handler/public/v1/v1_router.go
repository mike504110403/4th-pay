package v1

import (
	"pay-service/api/handler/public/v1/payment"

	"github.com/gofiber/fiber/v2"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/v1")
	{
		payment.SetRouter(g)
	}
}
