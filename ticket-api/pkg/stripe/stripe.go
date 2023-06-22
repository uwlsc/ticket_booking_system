package stripe

import (
	"go.uber.org/fx"
)

// Module Middleware exported
var Module = fx.Module(
	"stripe",
	fx.Provide(
		NewStripeConnection,
		NewService,
	),
)

type CardDetails struct {
	CardNumber     string `json:"number"`
	ExpMonth       string `json:"exp_month"`
	ExpYear        string `json:"exp_year"`
	CVC            string `json:"cvc"`
	CardHolderName string `json:"name"`
	Type           string `json:"type"`
}
