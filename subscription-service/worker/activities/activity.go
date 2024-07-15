package activity

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/go-redis/redis/v8"
	twilio "github.com/twilio/twilio-go"
)

type Activites interface {
	SendWelcomeEmail(ctx context.Context, to, name string) error
	SendWelcomeSMS(to, name string) error
	SendOTPSMS(to string, otpCode string) error
	SendOTPEmail(ctx context.Context, to, otpCode string) error
	GenerateOTP(ctx context.Context, userID string) (string, error)
}

type ActivitiesImpl struct {
	sesClient    *ses.SES
	twilioClient *twilio.RestClient
	redis        *redis.Client
}

func NewActivities(sesClient *ses.SES, twilioClient *twilio.RestClient, redis *redis.Client) Activites {
	return &ActivitiesImpl{
		sesClient:    sesClient,
		twilioClient: twilioClient,
		redis:        redis,
	}
}
