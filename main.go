package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"pay-service/internal/cachedata"
	"pay-service/internal/config"
	"pay-service/internal/database"

	routers "pay-service/api/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/mike504110403/common-moduals/apiprotocol"
	"github.com/mike504110403/common-moduals/typeparam"
	mlog "github.com/mike504110403/goutils/log"

	// 註冊 swagger
	_ "pay-service/docs"
	// 註冊金流商
	_ "pay-service/internal/payment/providers"
	// 監聽支付請求
	_ "pay-service/internal/payment"
	// 監聽收款請求
	_ "pay-service/internal/order"
)

var (
	Version     = "0.0.0"
	CommitID    = "dev"
	Environment = os.Getenv("Environment")
	Port        = os.Getenv("Port")
)

// ServiceName : 服務名稱 | 不可修改，會影響資料一致性
const ServiceName = "pay-service"

func Init() {
	if Version == "0.0.0" {
		cmd := exec.Command("git", "rev-parse", "HEAD")
		out := bytes.Buffer{}
		cmd.Stdout = &out
		err := cmd.Run()
		if err == nil {
			Version = strings.Replace(out.String(), "\n", "", -1)
		}
	} else {
		CommitID = strings.Replace(CommitID, "\"", "", -1)
		Version = strings.Replace(Version, "\"", "", -1)
	}
	// 環境變數初始化
	config.Init(Version)
	config.EnvInit()
	env := config.GetENV()

	// 初始化log
	mlog.Init(mlog.Config{
		EnvMode: mlog.EnvMode(env.SystemENV.EnvMod),
		LogType: mlog.LogType(env.SystemENV.LogType),
	})

	// 初始化資料庫
	database.Init()
	// 初始化快取
	cachedata.Init(cachedata.Config{
		RefreshDuration: time.Minute * 3,
		RetryDuration:   time.Second * 3,
	})

	// 設定參數初始化
	if err := typeparam.Init(typeparam.Config{
		FuncGetDB: func() (*sql.DB, error) {
			return database.SETTING.DB()
		},
	}); err != nil {
		mlog.Fatal(fmt.Sprintf("參數初始化錯誤: %S", err))
	}
}

func main() {
	Init()
	sysENV := config.GetSystemENV()

	// 設定fiber的config
	app := fiber.New(fiber.Config{
		ReadTimeout:    sysENV.ReadTimeout,
		WriteTimeout:   sysENV.WriteTimeout,
		IdleTimeout:    sysENV.IdleTimeout,
		ReadBufferSize: int(sysENV.FiberHeaderSizeLimitMb) * 1024,
		BodyLimit:      int(sysENV.FiberBodyLimitMb) * 1024 * 1024,
		ErrorHandler:   apiprotocol.ErrorHandler(),
	})

	app.Get("/version", func(c *fiber.Ctx) error {
		statusData := fiber.Map{
			"environment":        Environment,
			"Service":            ServiceName,
			"Version":            Version,
			"commitID":           CommitID,
			"OpenConnections":    c.App().Server().GetOpenConnectionsCount(),
			"CurrentConcurrency": c.App().Server().GetCurrentConcurrency(),
			"goVersion":          runtime.Version(),
			"config":             app.Config(),
		}
		return c.Status(fiber.StatusOK).JSON(statusData)
	})

	// Swagger API
	app.Get("/swagger/*", swagger.HandlerDefault)

	// 設定路由
	if err := routers.Set(app); err != nil {
		mlog.Fatal(fmt.Sprintf("routers設定失敗, err: %v", err))
	}

	// 關閉服務
	go func() {
		if err := app.Listen(":" + Port); err != nil {
			mlog.Error(err.Error())
		}
	}()
	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	app.Shutdown()
	time.Sleep(5 * time.Second)
}
