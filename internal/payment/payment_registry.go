package payment

import (
	"errors"
	"pay-service/internal/payment/payment_types"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// PaymentProvider 定義所有金流商的行為 -> 金流核心介面
type PaymentProvider interface {
	// 建立交易
	NewTransaction(order *payment_types.PaymentTransaction) error
	// 接收授權回調
	HandleAuthCallBack(c *fiber.Ctx) error
	// 接收回調
	HandleCallback(c *fiber.Ctx) error
}

// providerRegistry 儲存已註冊的金流商 -> 金流商註冊池
var providerRegistry = struct {
	sync.RWMutex
	providers map[string]PaymentProvider
}{
	providers: make(map[string]PaymentProvider),
}

// RegisterProvider 註冊金流商
func RegisterProvider(name string, provider PaymentProvider) {
	providerRegistry.Lock()
	defer providerRegistry.Unlock()
	providerRegistry.providers[name] = provider
}

// GetProvider 取得對應的金流商
func GetProvider(name string) (PaymentProvider, error) {
	providerRegistry.RLock()
	defer providerRegistry.RUnlock()

	provider, exists := providerRegistry.providers[name]
	if !exists {
		return nil, errors.New("provider not found")
	}

	return provider, nil
}
