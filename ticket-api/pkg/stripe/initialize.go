package stripe

import (
	"github.com/stripe/stripe-go/client"

	"ticket-api/config"
)

func NewStripeConnection(config config.Config) *client.API {
	return client.New(config.StripeKey, nil)
}
