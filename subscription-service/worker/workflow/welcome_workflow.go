// Package workflow defines workflows for managing user interactions, such as sending welcome messages, using Temporal.
package workflow

import (
	"time" // Import time for setting timeouts and intervals.

	"go.temporal.io/sdk/temporal" // Import Temporal's Go SDK for defining retry policies and workflow options.
	"go.temporal.io/sdk/workflow" // Import workflow to define and execute workflows.
)

// WelcomeParams struct holds the parameters required for the WelcomeWorkflow.
type WelcomeParams struct {
	To      string // Email address of the recipient.
	Name    string // Name of the recipient.
	Contact string // Phone number of the recipient.
}

// WelcomeWorkflow orchestrates the process of sending a welcome email and SMS to a new user.
// It takes in a context and WelcomeParams and returns an error if any step in the process fails.
func WelcomeWorkflow(ctx workflow.Context, params WelcomeParams) error {
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

	// Execute the SendWelcomeEmail activity with the recipient's email address and name.
	// If this activity fails, the error is returned and the workflow is terminated.
	err := workflow.ExecuteActivity(ctx, "SendWelcomeEmail", params.To, params.Name).Get(ctx, nil)
	if err != nil {
		return err // Return the error if the activity fails.
	}

	// Execute the SendWelcomeSMS activity with the recipient's phone number and name.
	// If this activity fails, the error is returned and the workflow is terminated.
	err = workflow.ExecuteActivity(ctx, "SendWelcomeSMS", params.Contact, params.Name).Get(ctx, nil)
	if err != nil {
		return err // Return the error if the activity fails.
	}

	return nil // Return nil to indicate successful completion of the workflow.
}
