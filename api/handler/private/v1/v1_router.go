package v1

import (
	"pay-service/api/handler/private/v1/order"

	"github.com/gofiber/fiber/v2"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/v1")
	{
		order.SetRouter(g)
	}
}
