package workflow

import (
	activity "subscription-service/worker/activities"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/go-redis/redis/v8"
	"github.com/twilio/twilio-go"
	"go.temporal.io/sdk/workflow"
)

type OTPParams struct {
	SVC     *ses.SES           // SES service client
	Twilio  *twilio.RestClient // Twilio service client
	Name    string             // Recipient name
	To      string             // Recipient email address
	Contact string             // Recipient phone number
	UserID  string             // User ID
	Redis   *redis.Client      // Redis client
}

func OTPWorkflow(ctx workflow.Context, params OTPParams) error {
	// Define activity options, if needed
	ao := workflow.ActivityOptions{
		// Activity options here
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	type Result struct {
		OTP string
	}
	var result Result
	err := workflow.ExecuteActivity(ctx, activity.GenerateOTP, params.Redis, params.UserID).Get(ctx, &result)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, activity.SendOTPEmail, params.SVC, params.To, result.OTP).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, activity.SendOTPSMS, params.Twilio, params.Contact, result.OTP).Get(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
