package stripe

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

type Service struct {
	*client.API
}

func NewService(api *client.API) Service {
	return Service{api}
}

func (s *Service) CreateStripeCustomer(name, email string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
	}

	newCustomer, err := s.API.Customers.New(params)

	if err != nil {
		return nil, err
	}

	return newCustomer, err
}

func (s *Service) CreatePaymentMethod(customerId string, card CardDetails) (*stripe.PaymentMethod, error) {

	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(card.CardNumber),
			ExpMonth: stripe.String(card.ExpMonth),
			ExpYear:  stripe.String(card.ExpYear),
			CVC:      stripe.String(card.CVC),
		},
		Type: stripe.String("card"),
	}

	method, err := s.API.PaymentMethods.New(params)

	if err != nil {
		return nil, err
	}

	_, attachmentError := s.AttachCustomerPaymentMethod(customerId, method.ID)

	if attachmentError != nil {
		return nil, attachmentError
	}

	return method, nil
}

// AttachCustomerPaymentMethod .
func (s *Service) AttachCustomerPaymentMethod(cId, pId string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(cId),
	}

	method, err := s.API.PaymentMethods.Attach(pId, params)
	if err != nil {
		return nil, err
	}

	return method, nil
}

// CreateStripeSubscription Stripe Subscription
func (s *Service) CreateStripeSubscription(stripeId, priceId, pmId string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(stripeId),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(priceId),
			},
		},
		PaymentBehavior:      stripe.String("default_incomplete"),
		DefaultPaymentMethod: &pmId,
	}
	params.AddExpand("latest_invoice.payment_intent")

	subscription, err := s.API.Subscriptions.New(params)
	if err != nil {
		return nil, err
	}
	return subscription, nil

}

// CreateStripeSubscription Stripe Subscription
func (s *Service) UpdateStripeSubscription(stripeId, priceId, subscriptionId string) (*stripe.Subscription, error) {
	subscription, _ := s.API.Subscriptions.Get(subscriptionId, nil)
	params := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{
				// ID:   stripe.String(subscription.Items.Values[0].ID),
				ID:   stripe.String(subscription.Items.Data[0].ID),
				Plan: stripe.String(priceId),
			},
		},
		CancelAtPeriodEnd: stripe.Bool(false),
		ProrationBehavior: stripe.String(string(stripe.SubscriptionProrationBehaviorCreateProrations)),
		PaymentBehavior:   stripe.String("default_incomplete"),
	}
	params.AddExpand("latest_invoice.payment_intent")
	subscriptions, err := s.API.Subscriptions.Update(subscriptionId, params)
	if err != nil {
		return nil, err
	}
	return subscriptions, nil

}

// ConfirmChargePaymentIntent Confirm Charge
func (s *Service) ConfirmChargePaymentIntent(paymentMethod, pmId string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentConfirmParams{
		PaymentMethod:    stripe.String(pmId),
		SetupFutureUsage: stripe.String("off_session"),
	}

	pi, err := s.API.PaymentIntents.Confirm(
		paymentMethod,
		params,
	)
	return pi, err
}

// CancelSubscription from stripe
func (s *Service) CancelSubscription(subscriptionId string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}
	sub, err := s.API.Subscriptions.Update(subscriptionId, params)
	return sub, err

}

// Delete/Remove card from stripe
func (s *Service) RemovePaymentMethod(paymentMethodId string) (*stripe.PaymentMethod, error) {
	c, err := s.API.PaymentMethods.Detach(
		paymentMethodId,
		nil,
	)
	return c, err

}

// Get all cards for user
func (s *Service) GetPaymentMethod(customerId string) *stripe.PaymentMethod {
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerId),
		Type:     stripe.String("card"),
	}
	i := s.API.PaymentMethods.List(params)
	var pm *stripe.PaymentMethod
	for i.Next() {
		pm = i.PaymentMethod()
	}
	return pm
}
