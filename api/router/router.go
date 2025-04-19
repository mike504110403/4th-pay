package router

import (
	"pay-service/api/router/middleware/reshandler"

	"pay-service/api/handler/private"
	"pay-service/api/handler/public"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
	"github.com/mike504110403/common-moduals/fiber/handler/recoverhandler"
	"github.com/mike504110403/common-moduals/fiber/handler/tracehandler"
	"github.com/mike504110403/common-moduals/ilog"
)

// Set : 設定全部的路由
func Set(r *fiber.App) error {

	r.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		// AllowCredentials: true,
	}))

	r.Use(recoverhandler.New())

	r.Use(tracehandler.New())

	r.Use(func(c *fiber.Ctx) error { // 設定每次request建立一個log物件，並在最後處理或印出log
		if c.Method() == fiber.MethodOptions {
			c.Set("Access-Control-Allow-Origin", "*") // 或具體的域名
			c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			return c.SendStatus(fiber.StatusOK)
		}

		reqUrl := c.Request().URI().RequestURI()
		logData := ilog.Basic().Trace(tracehandler.Get(c))

		defer func() {
			if c.Response().StatusCode() != fiber.StatusOK {
				if len(c.Request().Body()) > 200 {
					logData.Log(`[req][path->%v][body]%v`, c.Path(), string(c.Request().Body()[:200]))
				} else {
					logData.Log(`[req][path->%v][body]%v`, c.Path(), string(c.Request().Body()))
				}
				logData.Log(`[res->%v][path->%v] %v`, c.Response().StatusCode(), c.Path(), string(reqUrl))
			}
		}()

		return c.Next()
	})
	r.Use(helmet.New(helmet.Config{
		XFrameOptions: "SAMEORIGIN",
	}))
	r.Use(reshandler.ResponseMiddleware)
	api := r.Group("/api")
	{
		private.SetRouter(api)
		public.SetRouter(api)
	}

	return nil
}
