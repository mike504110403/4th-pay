package provider

import (
	"fmt"

	"pay-service/internal/payment"
	"pay-service/internal/payment/providers/gomypay"
	"pay-service/internal/payment/providers/paypal"
	"pay-service/internal/payment/providers/stripe"
)

func init() {
	// TODO: 這裡可以改成從 config 讀取啟用的金流商
	enabledProviders := []string{"gomypay"}

	// 註冊已啟用的 Provider
	for _, providerName := range enabledProviders {
		switch providerName {
		case "gomypay":
			payment.RegisterProvider("gomypay", &gomypay.GomyPayProvider{})
			fmt.Println("✅ 註冊金流商: GomyPay")
		case "paypal":
			payment.RegisterProvider("paypal", &paypal.PayPalProvider{})
			fmt.Println("✅ 註冊金流商: PayPal")
		case "stripe":
			payment.RegisterProvider("stripe", &stripe.StripeProvider{})
			fmt.Println("✅ 註冊金流商: Stripe")
		default:
			fmt.Printf("⚠️ 未知的金流商 %s，略過註冊", providerName)
		}
	}
}
