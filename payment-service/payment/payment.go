package payment

import (
	"os"

	"github.com/NdoleStudio/lemonsqueezy-go"
	"github.com/joho/godotenv"
)

type Payment interface {
	CreateSubscription(planID string, customerEmail string, otherDetails map[string]interface{}) (*lemonsqueezy.Subscription, error)
	CancelSubscription(subscriptionID string) error
	GetAllSubscriptions() ([]*lemonsqueezy.Subscription, error)
	UpdateSubscription(subscriptionID string, updates map[string]interface{}) (*lemonsqueezy.Subscription, error)
	GetSubscription(subscriptionID string) (*lemonsqueezy.Subscription, error)
}

type PaymentService struct {
	Client *lemonsqueezy.Client
}

// NewPayment initializes a new payment client using environment variables.
// Returns a pointer to a lemonsqueezy.Client and any error encountered.
func NewPayment() (Payment, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}
	// Assuming the API key is stored in an environment variable, for example, LEMON_SQUEEZY_API_KEY
	apiKey := os.Getenv("LEMON_SQUEEZY_API_KEY")
	client := lemonsqueezy.New(lemonsqueezy.WithAPIKey(apiKey))
	return &PaymentService{
		Client: client,
	}, nil
}

// these can be implemented as methods of the PaymentService struct

// 1. Create a new subscription
func (p *PaymentService) CreateSubscription(planID string, customerEmail string, otherDetails map[string]interface{}) (*lemonsqueezy.Subscription, error) {
	return nil, nil
}

// 2. Cancel an ongoing subscription
func (p *PaymentService) CancelSubscription(subscriptionID string) error {

	return nil
}

// 3. Get all subscriptions
func (p *PaymentService) GetAllSubscriptions() ([]*lemonsqueezy.Subscription, error) {

	return nil, nil
}

// 4. Update a subscription
func (p *PaymentService) UpdateSubscription(subscriptionID string, updates map[string]interface{}) (*lemonsqueezy.Subscription, error) {

	return nil, nil
}

// 5. Get a specific subscription
func (p *PaymentService) GetSubscription(subscriptionID string) (*lemonsqueezy.Subscription, error) {
	return nil, nil
}
