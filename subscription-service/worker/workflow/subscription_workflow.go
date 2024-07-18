// Package workflow defines workflows for managing subscriptions using Temporal.
package workflow

import (
	"strconv" // Import strconv for converting strings to other types.
	"time"    // Import time for setting timeouts and intervals.

	"go.temporal.io/sdk/temporal" // Import Temporal's Go SDK for workflow and activity management.
	"go.temporal.io/sdk/workflow" // Import workflow to define and execute workflows.
)

// SubscriptionParams struct holds the parameters required for the subscription workflow.
type SubscriptionParams struct {
	Email       string // Email of the recipient.
	PlanName    string // Name of the subscription plan.
	VariantName string // Name of the subscription variant.
	Status      string // Current status of the subscription.
	Type        string // Type of mail to be sent, used here to demonstrate a custom logic.
}

// SubscriptionWorkflow orchestrates the subscription process using Temporal workflows.
// It schedules and executes several activities based on the provided SubscriptionParams.
func SubscriptionWorkflow(ctx workflow.Context, params SubscriptionParams) error {
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

	// Resp struct is used to capture the response from the GetUser activity.
	type Resp struct {
		ID             int64   // User's ID.
		Contact        string  // User's contact information.
		SubscriptionID float64 // ID of the user's subscription.
	}

	var userResponse Resp // Variable to store the response from GetUser activity.

	// Execute the GetUser activity with the provided email.
	// The result is stored in userResponse.
	err := workflow.ExecuteActivity(ctx, "GetUser", params.Email).Get(ctx, &userResponse)
	if err != nil {
		return err // Return the error if the activity fails.
	}

	// Attempt to convert the Type parameter to a float.
	// If successful, update the SubscriptionID in userResponse and set Type to "Created".
	if floatValue, err := strconv.ParseFloat(params.Type, 64); err == nil {
		userResponse.SubscriptionID = floatValue
		params.Type = "Created"
	}

	// Execute the UpdateSubscription activity with the updated subscription details.
	err = workflow.ExecuteActivity(ctx, "UpdateSubscription", userResponse.ID, params.Status, userResponse.SubscriptionID, params.PlanName+params.VariantName).Get(ctx, nil)
	if err != nil {
		return err // Return the error if the activity fails.
	}

	// Execute the SendSubscriptionStatusEmail activity to send an email to the user.
	err = workflow.ExecuteActivity(ctx, "SendSubscriptionStatusEmail", params.Email, userResponse.SubscriptionID, params.PlanName+params.VariantName, params.Status).Get(ctx, nil)
	if err != nil {
		return nil // Proceed even if sending the email fails.
	}

	// Execute the SendSubscriptionUpdateSMS activity to send an SMS to the user.
	err = workflow.ExecuteActivity(ctx, "SendSubscriptionUpdateSMS", userResponse.Contact, params.PlanName+params.VariantName, params.Status).Get(ctx, nil)
	if err != nil {
		return nil // Proceed even if sending the SMS fails.
	}

	return nil // Return nil to indicate successful completion of the workflow.
}
