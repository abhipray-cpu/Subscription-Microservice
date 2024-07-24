package activity

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	twilio "github.com/twilio/twilio-go"
)

// Activites is an interface that defines the activities that can be performed by the worker.
// It defines methods for sending welcome emails and SMS messages, sending OTP emails and SMS messages,
type Activites interface {
	SendWelcomeEmail(ctx context.Context, to, name string) error
	SendWelcomeSMS(to, name string) error
	SendOTPSMS(to string, otpCode string) error
	SendOTPEmail(ctx context.Context, to, otpCode string) error
	GenerateOTP(ctx context.Context, userID string) (string, error)
	GetUser(email string) (UserResponse, error)
	UpdateSubscription(id int64, subscriptionStatus string, subscriptionId float64, subscriptionType string) error
	SendSubscriptionUpdateSMS(to, subscriptionName, status string) error
	SendSubscriptionStatusEmail(ctx context.Context, to string, subscriptionID float64, subscriptionName, status string) error
}

// ActivitiesImpl is an implementation of the Activites interface.
// It contains the necessary clients and services to perform the activities.
type ActivitiesImpl struct {
	sesClient    *ses.SES
	twilioClient *twilio.RestClient
	redis        *redis.Client
	connection   *pgx.Conn
}

// NewActivities creates a new ActivitiesImpl instance with the given clients and services.
// It returns an Activites interface.
func NewActivities(sesClient *ses.SES, twilioClient *twilio.RestClient, redis *redis.Client, connection *pgx.Conn) Activites {
	return &ActivitiesImpl{
		sesClient:    sesClient,
		twilioClient: twilioClient,
		redis:        redis,
		connection:   connection,
	}
}
