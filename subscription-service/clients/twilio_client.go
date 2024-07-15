package clients

import (
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
)

// TwilioClient initializes and returns a Twilio REST client or an error.
func TwilioClient() (*twilio.RestClient, error) {
	// Load environment variables from .env file
	// if err := godotenv.Load(".env"); err != nil {
	// 	// Return nil and the error if .env file is not found or any other error occurs
	// 	return nil, fmt.Errorf("error loading .env file: %w", err)
	// }

	// Retrieve Twilio account SID and auth token from environment variables
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	// Check if the environment variables are set
	if accountSID == "" || authToken == "" {
		return nil, fmt.Errorf("TWILIO_ACCOUNT_SID or TWILIO_AUTH_TOKEN is not set")
	}

	// Initialize the Twilio client with the account SID and auth token
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	// Return the Twilio client and nil as the error
	return client, nil
}
