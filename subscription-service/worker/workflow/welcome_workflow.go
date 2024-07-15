package workflow

import (
	activity "subscription-service/worker/activities"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/twilio/twilio-go"
	"go.temporal.io/sdk/workflow"
)

// WelcomeParams defines the parameters for the WelcomeWorkflow.
type WelcomeParams struct {
	SVC     *ses.SES           // SES service client
	Twilio  *twilio.RestClient // Twilio service client
	To      string             // Recipient email address
	Name    string             // Recipient name
	Contact string             // Recipient phone number
}

func WelcomeWorkflow(ctx workflow.Context, params WelcomeParams) error {
	// Define activity options, if needed
	ao := workflow.ActivityOptions{
		// Activity options here
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(ctx, activity.SendWelcomeEmail, params.SVC, params.To, params.Name).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, activity.SendWelcomeSMS, params.Twilio, params.Contact, params.Name).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
