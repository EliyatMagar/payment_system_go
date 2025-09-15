package services

import (
	"payment-system/config"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
	"github.com/stripe/stripe-go/v76/webhook"
)

type StripeService struct {
	client        *client.API
	webhookSecret string
}

func NewStripeService(cfg *config.Config) *StripeService {
	stripe.Key = cfg.StripeKey
	sc := &client.API{}
	sc.Init(cfg.StripeKey, nil)

	return &StripeService{
		client:        sc,
		webhookSecret: cfg.WebhookSecret,
	}
}

func (s *StripeService) CreateCustomer(email, name string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}

	customer, err := s.client.Customers.New(params)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *StripeService) CreatePaymentIntent(amount int64, currency, customerID string, metadata map[string]string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		Customer: stripe.String(customerID),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	if metadata != nil {
		params.Metadata = metadata
	}

	pi, err := s.client.PaymentIntents.New(params)
	if err != nil {
		return nil, err
	}

	return pi, nil
}

func (s *StripeService) ConfirmPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	pi, err := s.client.PaymentIntents.Confirm(paymentIntentID, nil)
	if err != nil {
		return nil, err
	}

	return pi, nil
}

func (s *StripeService) HandleWebhook(payload []byte, signature string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, signature, s.webhookSecret)
}

func (s *StripeService) GetPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	pi, err := s.client.PaymentIntents.Get(paymentIntentID, nil)
	if err != nil {
		return nil, err
	}

	return pi, nil
}

func (s *StripeService) RefundPayment(paymentIntentID string) (*stripe.Refund, error) {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentIntentID),
	}

	refund, err := s.client.Refunds.New(params)
	if err != nil {
		return nil, err
	}

	return refund, nil
}
