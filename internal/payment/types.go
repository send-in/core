package payment

import config "core/internal/config"

var razorpay *config.RazorpayConfig

func Configure(cfg *config.RazorpayConfig) {
	razorpay = cfg
}

const (
	MinCredits = 25
	MaxCredits = 200
	CreditStep = 25
	PricePerCredit = 1
)

type OrderRequest struct {
	Amount   int64             `json:"amount"`
	Currency string            `json:"currency"`
	Receipt  string            `json:"receipt"`
	Notes    map[string]string `json:"notes,omitempty"`
}

type OrderResponse struct {
	ID       string `json:"id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
	Receipt  string `json:"receipt"`
}

