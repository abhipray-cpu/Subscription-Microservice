package clients

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// NewSESClient initializes and returns an SES client or an error.
// It reads the AWS credentials and region from environment variables.
// It returns an error if the session cannot be created.
// It returns an SES client if the session is created successfully.
// It returns an error if the SES client cannot be created.
func NewSESClient() (*ses.SES, error) {
	// Get AWS credentials and region from environment variables
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

	// Create a new session using the loaded environment variables
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Create a new SES client
	svc := ses.New(sess)

	// Return the SES client
	return svc, nil
}
