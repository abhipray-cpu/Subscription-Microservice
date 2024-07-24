package test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"subscription-service/data"
	activity "subscription-service/worker/activities"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
)

type ActivitySuite struct {
	activites   activity.Activites
	redisClient *redis.Client
	connection  *pgx.Conn
	userID      []int64
}

func (suite *ActivitySuite) SetupSuite() {
	// db connection
	url := "postgres://root@localhost:26257/defaultdb?sslmode=disable" // Database connection URL.
	conn, err := pgx.Connect(context.Background(), url)                // Attempt to connect to the database.
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	fmt.Println("Connected to the database")
	ensureTableExists(conn) // Ensure the table exists in the database.
	suite.connection = conn

	// initializing new redis client
	redis, err := data.NewRedisClient("localhost:6379", "")
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)

	}

	// sns client
	ses, err := NewSESClient()
	if err != nil {
		log.Fatalf("Failed to create SNS client: %v", err)

	}
	// twilio client
	twilio, err := TwilioClient()
	if err != nil {
		log.Fatalf("Failed to create Twilio client: %v", err)

	}
	suite.redisClient = redis
	suite.activites = activity.NewActivities(ses, twilio, redis, conn)

}

func (suite *ActivitySuite) TeardownSuite() {
	if suite.redisClient != nil {
		suite.redisClient.Close()
	}
	if suite.connection != nil {
		suite.connection.Close(context.Background())
	}
}

// get user activity test
func (suite *ActivitySuite) TestGetUserActitvity(t *testing.T) {
	testCase := []struct {
		name          string
		email         string
		expectedError error
		expectedID    int64
	}{
		{
			name:          "Valid Email, User Found",
			email:         "john.doe@example.com",
			expectedError: nil,
			expectedID:    suite.userID[0],
		},
		{
			name:          "Valid Email, User Not Found",
			email:         "missing@example.com",
			expectedError: errors.New("user not found"),
			expectedID:    0,
		},
		{
			name:          "Invalid Email Format",
			email:         "invalidemail",
			expectedError: errors.New("invalid email format"),
			expectedID:    0,
		},
		{
			name:          "Empty Email",
			email:         "",
			expectedError: errors.New("email cannot be empty"),
			expectedID:    0,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			userResponse, err := suite.activites.GetUser(tc.email)
			if err != nil && tc.expectedError == nil {
				t.Errorf("GetUser() error = %v, expectedError %v", err, tc.expectedError)
				return
			}
			if err == nil && tc.expectedError != nil {
				t.Errorf("GetUser() error = %v, expectedError %v", err, tc.expectedError)
				return
			}
			if userResponse.ID != tc.expectedID {
				t.Errorf("GetUser() ID = %v, expectedID %v", userResponse.ID, tc.expectedID)
			}
		})
	}
}

// updat subscription test
func (suite *ActivitySuite) TestUpdateSubscriptionActitvity(t *testing.T) {
	var testCases = []struct {
		name               string
		id                 int64
		subscriptionStatus string
		subscriptionId     float64
		subscriptionType   string
		expectedError      bool
	}{
		{
			name:               "Successful Update",
			id:                 suite.userID[0],
			subscriptionStatus: "active",
			subscriptionId:     123.45,
			subscriptionType:   "premium",
			expectedError:      false,
		},
		{
			name:               "Non-existent User",
			id:                 999,
			subscriptionStatus: "active",
			subscriptionId:     123.45,
			subscriptionType:   "premium",
			expectedError:      true, // Assuming no error for no rows updated
		},
		{
			name:               "Invalid Subscription ID",
			id:                 2,
			subscriptionStatus: "active",
			subscriptionId:     -1, // Invalid ID
			subscriptionType:   "basic",
			expectedError:      true, // Assuming your validation catches this before making a DB call
		},
		{
			name:               "Database Error",
			id:                 3,
			subscriptionStatus: "inactive",
			subscriptionId:     234.56,
			subscriptionType:   "standard",
			expectedError:      true,
		},
		{
			name:               "Update With Same Data",
			id:                 suite.userID[2],
			subscriptionStatus: "active",
			subscriptionId:     123.45,
			subscriptionType:   "premium",
			expectedError:      true, // Assuming no error for no rows updated
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		u := data.User{}
		err := u.UpdateUserSubscription(suite.connection, tc.id, tc.subscriptionStatus, tc.subscriptionId, tc.subscriptionType)

		if err != nil && !tc.expectedError {
			t.Errorf("UpdateSubscription() error = %v, wantErr %v", err, tc.expectedError)
			return
		}

		if err == nil && tc.expectedError {
			t.Errorf("UpdateSubscription() error = %v, wantErr %v", err, tc.expectedError)
			return
		}
	}
}

// otp activity test
func (suite *ActivitySuite) TestOTPActivity(t *testing.T) {
	var testCases = []struct {
		name          string
		userID        string
		mock          func(userID int64)
		expectedError bool
	}{
		{
			name:   "Successful OTP Generation and Storage",
			userID: strconv.Itoa(int(suite.userID[0])),
			mock: func(userID int64) {
				key := fmt.Sprintf("%d:OTP", userID)
				suite.redisClient.Del(context.Background(), key)                    // Ensure the key does not exist
				suite.redisClient.Set(context.Background(), key, "generatedOTP", 0) // Simulate OTP storage
			},
			expectedError: false,
		},
		{
			name:   "OTP Already Exists",
			userID: "12345",
			mock: func(userID int64) {
				key := fmt.Sprintf("%d:OTP", userID)
				suite.redisClient.Set(aws.BackgroundContext(), key, "existingOTP", 0) // Ensure the OTP already exists
			},
			expectedError: false,
		},
		{
			name:   "Redis Exists Check Error",
			userID: "12345",
			mock: func(userID int64) {
				// Temporarily point the Redis client to an invalid address to simulate an error
				originalOptions := suite.redisClient.Options()
				redisClient := redis.NewClient(&redis.Options{
					Addr: "invalid:6379", // Use an invalid address
				})

				// Defer a function to restore the original Redis client configuration
				// after the test case is done to ensure it does not affect subsequent tests
				defer func() {
					redisClient.Close() // Close the client with the invalid configuration
					// Restore the original Redis client configuration
					redisClient = redis.NewClient(originalOptions)
				}()
			},
			expectedError: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the necessary dependencies
			tc.mock(suite.userID[0])

			// Call the function being tested
			_, err := suite.activites.GenerateOTP(context.Background(), tc.userID)
			// Check the expected error
			if (err != nil) != tc.expectedError {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedError, err)
			}

			if err == nil {
				// Clean up the OTP key in Redis
				key := fmt.Sprintf("%d:OTP", suite.userID[0])
				suite.redisClient.Del(context.Background(), key)
			}
		})
	}
}

func (suite *ActivitySuite) TestWelcomesSMSActivity(t *testing.T) {
	testCases := []struct {
		name      string
		to        string
		nameInput string
		wantErr   bool
	}{
		{"Valid Number", "+917895724996", "John Doe", false},
		{"Invalid Number", "", "John Doe", true},
		{"Invalid Format Number", "12345", "Jane Doe", true},
		{"Valid Number with Special Characters in Name", "+1234567890", "John & Jane Doe", true},
		{"Long Number Beyond Limits", "+12345678901234567890", "John Doe", true},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			err := suite.activites.SendWelcomeSMS(tt.to, tt.nameInput)

			if (err != nil) != tt.wantErr {
				t.Errorf("SendWelcomeSMS() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

// test for creating test users for test activity
func (suite *ActivitySuite) TestInsertUser(t *testing.T) {
	testCases := []struct {
		name    string
		user    TestUser
		wantErr bool
	}{
		{
			name: "ValidUserCompleteInfo",
			user: TestUser{
				ID:          1,
				UserName:    "completeInfo",
				GithubName:  "githubComplete",
				FirstName:   "John",
				SecondName:  "Doe",
				Bio:         "Full Stack Developer",
				Email:       "john.doe@example.com",
				Contact:     "1234567890",
				Password:    "securePassword123",
				AccessToken: "validAccessToken",
				Verified:    true,
			},
			wantErr: false,
		},
		{
			name: "ValidUserTechEnthusiast",
			user: TestUser{
				ID:          2,
				UserName:    "techEnthusiast",
				GithubName:  "githubTech",
				FirstName:   "Alice",
				SecondName:  "Smith",
				Bio:         "Tech enthusiast and blogger",
				Email:       "alice.smith@example.com",
				Contact:     "9876543210",
				Password:    "techLover2023",
				AccessToken: "validTechToken",
				Verified:    true,
			},
			wantErr: false,
		},
		{
			name: "ValidUserDesigner",
			user: TestUser{
				ID:          3,
				UserName:    "creativeDesigner",
				GithubName:  "githubDesign",
				FirstName:   "Bob",
				SecondName:  "Brown",
				Bio:         "Creative designer specializing in UX/UI",
				Email:       "bob.brown@example.com",
				Contact:     "1122334455",
				Password:    "designIsLife",
				AccessToken: "validDesignToken",
				Verified:    true,
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := data.User{
				UserName:    tc.user.UserName,
				GithubName:  tc.user.GithubName,
				FirstName:   tc.user.FirstName,
				LastName:    tc.user.SecondName,
				Bio:         tc.user.Bio,
				Email:       tc.user.Email,
				Contact:     tc.user.Contact,
				Password:    tc.user.Password,
				AccessToken: tc.user.AccessToken,
				Verified:    tc.user.Verified,
			}

			err := u.InsertUser(suite.connection, u)
			if err != nil && !tc.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err == nil && tc.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			suite.userID = append(suite.userID, u.ID)
		})
	}

}
func TestActivitySuite(t *testing.T) {
	activity_suite := ActivitySuite{}
	activity_suite.SetupSuite()
	defer activity_suite.TeardownSuite()

	t.Run("CreateUser", activity_suite.TestInsertUser)
	t.Run("TestGetUserActitvity", activity_suite.TestGetUserActitvity)
	t.Run("TestUpdateSubscriptionActitvity", activity_suite.TestUpdateSubscriptionActitvity)
	t.Run("TestOTPActivity", activity_suite.TestOTPActivity)
	t.Run("TestSMSActivity", activity_suite.TestWelcomesSMSActivity)
}

// setting up test services
func NewSESClient() (*ses.SES, error) {
	// Load environment variables from .env file
	if err := godotenv.Load("../.env"); err != nil {
		// Return nil and the error if .env file is not found or any other error occurs
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
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

// TwilioClient initializes and returns a Twilio REST client or an error.
func TwilioClient() (*twilio.RestClient, error) {
	// Load environment variables from .env file
	if err := godotenv.Load("../.env"); err != nil {
		// Return nil and the error if .env file is not found or any other error occurs
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

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
