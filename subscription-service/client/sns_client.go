package client

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func NewSNSClient() (*sns.SNS, error) {
	// // Load .env file
	// if err := godotenv.Load(); err != nil {
	// 	return nil, fmt.Errorf("error loading .env file: %w", err)
	// }

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

	// Create a new SNS client
	svc := sns.New(sess)

	// Return the SNS client
	return svc, nil
}
