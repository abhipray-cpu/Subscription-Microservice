package workflow

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// WelcomeParams defines the parameters for the WelcomeWorkflow.
type WelcomeParams struct {
	To      string // Recipient email address
	Name    string // Recipient name
	Contact string // Recipient phone number
}

// WelcomeWorkflow is a Temporal workflow that orchestrates the welcome email and SMS sending process.
// It takes in a context and WelcomeParams and returns an error.
// It schedules the SendWelcomeEmail and SendWelcomeSMS activities with the recipient email address and phone number.
func WelcomeWorkflow(ctx workflow.Context, params WelcomeParams) error {
	// Define activity options, if needed
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: 10 * time.Second,
		StartToCloseTimeout:    10 * time.Second,
		HeartbeatTimeout:       10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(ctx, "SendWelcomeEmail", params.To, params.Name).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, "SendWelcomeSMS", params.Contact, params.Name).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
