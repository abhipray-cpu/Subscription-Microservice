// Package workflow defines workflows for managing one-time password (OTP) operations using Temporal.
package workflow

import (
	"time" // Import time for setting timeouts and intervals.

	"go.temporal.io/sdk/temporal" // Import Temporal's Go SDK for defining retry policies and workflow options.
	"go.temporal.io/sdk/workflow" // Import workflow to define and execute workflows.
)

// OTPParams struct holds the parameters required for the OTPWorkflow.
type OTPParams struct {
	Name    string // Recipient name.
	To      string // Recipient email address.
	Contact string // Recipient phone number.
	UserID  string // User ID for whom the OTP is generated.
}

// OTPWorkflow orchestrates the OTP generation and sending process.
// It takes in a context and OTPParams and returns an error if any step in the process fails.
func OTPWorkflow(ctx workflow.Context, params OTPParams) error {
	// ActivityOptions define the execution options for activities, including timeouts and retry policies.
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: 10 * time.Second, // Time allowed to find a worker that can start the activity.
		StartToCloseTimeout:    10 * time.Second, // Time allowed for the activity to complete execution.
		HeartbeatTimeout:       10 * time.Second, // Maximum time between heartbeats. Useful for long-running activities.
		RetryPolicy: &temporal.RetryPolicy{ // Defines the retry policy in case of activity failure.
			InitialInterval:    time.Second, // Initial interval between retries.
			BackoffCoefficient: 2.0,         // Multiplier by which the retry interval increases.
			MaximumInterval:    time.Minute, // Maximum interval between retries.
			MaximumAttempts:    5,           // Maximum number of retry attempts.
		},
	}
	// Apply the defined activity options to the current workflow context.
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string // Variable to store the generated OTP.

	// Execute the GenerateOTP activity with the user ID.
	// The result is stored in the result variable.
	err := workflow.ExecuteActivity(ctx, "GenerateOTP", params.UserID).Get(ctx, &result)
	if err != nil {
		return err // Return the error if the activity fails.
	}

	// Execute the SendOTPEmail activity with the recipient's email address and the generated OTP.
	// If this activity fails, the error is returned and the workflow is terminated.
	err = workflow.ExecuteActivity(ctx, "SendOTPEmail", params.To, result).Get(ctx, nil)
	if err != nil {
		return err // Return the error if the activity fails.
	}

	// Execute the SendOTPSMS activity with the recipient's phone number and the generated OTP.
	// If this activity fails, the error is returned and the workflow is terminated.
	err = workflow.ExecuteActivity(ctx, "SendOTPSMS", params.Contact, result).Get(ctx, nil)
	if err != nil {
		return err // Return the error if the activity fails.
	}

	return nil // Return nil to indicate successful completion of the workflow.
}
