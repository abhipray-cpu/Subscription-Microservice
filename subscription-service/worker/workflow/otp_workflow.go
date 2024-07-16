package workflow

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type OTPParams struct {
	Name    string // Recipient name
	To      string // Recipient email address
	Contact string // Recipient phone number
	UserID  string // User ID
}

// OTPWorkflow is a Temporal workflow that orchestrates the OTP generation and sending process.
// It takes in a context and OTPParams and returns an error.
// It schedules the GenerateOTP activity with the user ID.
func OTPWorkflow(ctx workflow.Context, params OTPParams) error {
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

	var result string
	err := workflow.ExecuteActivity(ctx, "GenerateOTP", params.UserID).Get(ctx, &result)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, "SendOTPEmail", params.To, result).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, "SendOTPSMS", params.Contact, result).Get(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
