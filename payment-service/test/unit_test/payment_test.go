package test

import (
	"context"
	"fmt"
	"log"
	"payment-service/data"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
)

type PaymentSuite struct {
	connection     *pgx.Conn
	paymentId      []int
	subscriptionId []string
	model          data.Models
}

func (suite *PaymentSuite) SetupSuite() {
	url := "postgres://root@localhost:26257/defaultdb?sslmode=disable" // Database connection URL.
	conn, err := pgx.Connect(context.Background(), url)                // Attempt to connect to the database.
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err) // Log and exit if the connection fails.
	}
	fmt.Println("Connected to the database") // Confirm successful connection.
	ensureTableExists(conn)                  // Ensure the table exists.
	suite.connection = conn
	suite.model = data.NewModels(conn)
}

// ensureTableExists updated to include new fields
func ensureTableExists(conn *pgx.Conn) {
	query := `
	DROP TABLE IF EXISTS payments;
    CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    customer_id FLOAT NOT NULL,
    subscription_id VARCHAR(255) NOT NULL,
    order_id FLOAT UNIQUE NOT NULL,
    status VARCHAR(50) NOT NULL,
    variant_name VARCHAR(255),
    variant_id FLOAT NOT NULL,
    product_id FLOAT NOT NULL,
    product_name VARCHAR(255),
    card_brand VARCHAR(50),
    card_last_four CHAR(4) NOT NULL CHECK (LENGTH(card_last_four) = 4),
    user_name VARCHAR(255) NOT NULL CHECK (user_name <> '' AND user_name ~ '^[A-Za-z ]+$'),
    user_email VARCHAR(255) NOT NULL CHECK (user_email <> '' AND user_email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    renews_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (renews_at >= created_at)
);`

	if _, err := conn.Exec(context.Background(), query); err != nil {
		log.Fatalf("Failed to ensure payments table exists: %v", err)
	}
}

func (suite *PaymentSuite) TearDown() {
	suite.connection.Close(context.Background())
	fmt.Println("Connection to the database closed")
}

func (suite *PaymentSuite) TestCreatePayment(t *testing.T) {
	testCases := []struct {
		name    string
		payment data.Payment
		wantErr bool
	}{
		{
			name: "ValidCompletePayment",
			payment: data.Payment{
				ID:             1,
				CustomerID:     1001,
				SubscriptionID: "sub_001",
				OrderID:        5001,
				Status:         "completed",
				VariantName:    "Premium",
				VariantID:      101,
				ProductID:      201,
				ProductName:    "Product A",
				CardBrand:      "Visa",
				CardLastFour:   "1234",
				UserName:       "John Doe",
				UserEmail:      "john.doe@example.com",
				RenewsAt:       time.Now().AddDate(0, 1, 0),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			wantErr: false,
		},
		{
			name: "ValidCompletePayment2",
			payment: data.Payment{
				ID:             12,
				CustomerID:     1002,
				SubscriptionID: "sub_011",
				OrderID:        5021,
				Status:         "completed",
				VariantName:    "Premium",
				VariantID:      101,
				ProductID:      201,
				ProductName:    "Product A",
				CardBrand:      "Visa",
				CardLastFour:   "1234",
				UserName:       "John Doe",
				UserEmail:      "john.doe@example.com",
				RenewsAt:       time.Now().AddDate(0, 1, 0),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			wantErr: false,
		},
		{
			name: "InvalidEmailFormat",
			payment: data.Payment{
				ID:             2,
				CustomerID:     1002,
				SubscriptionID: "sub_002",
				OrderID:        5002,
				Status:         "pending",
				VariantName:    "Basic",
				VariantID:      102,
				ProductID:      202,
				ProductName:    "Product B",
				CardBrand:      "MasterCard",
				CardLastFour:   "5678",
				UserName:       "Jane Doe",
				UserEmail:      "jane.doe",
				RenewsAt:       time.Now().AddDate(0, 1, 0),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			wantErr: true,
		},
		{
			name: "EmptyUserName",
			payment: data.Payment{
				ID:             3,
				CustomerID:     1003,
				SubscriptionID: "sub_003",
				OrderID:        5003,
				Status:         "completed",
				VariantName:    "Standard",
				VariantID:      103,
				ProductID:      203,
				ProductName:    "Product C",
				CardBrand:      "Amex",
				CardLastFour:   "9012",
				UserName:       "",
				UserEmail:      "empty.user@example.com",
				RenewsAt:       time.Now().AddDate(0, 1, 0),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			wantErr: true,
		},
		{
			name: "FutureCreatedAt",
			payment: data.Payment{
				ID:             4,
				CustomerID:     1004,
				SubscriptionID: "sub_004",
				OrderID:        5004,
				Status:         "pending",
				VariantName:    "Exclusive",
				VariantID:      104,
				ProductID:      204,
				ProductName:    "Product D",
				CardBrand:      "Visa",
				CardLastFour:   "3456",
				UserName:       "Future User",
				UserEmail:      "future.user@example.com",
				RenewsAt:       time.Now().AddDate(0, 1, 0),
				CreatedAt:      time.Now().AddDate(0, 0, 1), // Future date
				UpdatedAt:      time.Now(),
			},
			wantErr: true,
		},
		{
			name: "RenewsAtBeforeCreatedAt",
			payment: data.Payment{
				ID:             5,
				CustomerID:     1005,
				SubscriptionID: "sub_005",
				OrderID:        5005,
				Status:         "completed",
				VariantName:    "Limited Edition",
				VariantID:      105,
				ProductID:      205,
				ProductName:    "Product E",
				CardBrand:      "MasterCard",
				CardLastFour:   "6789",
				UserName:       "Retro User",
				UserEmail:      "retro.user@example.com",
				RenewsAt:       time.Now().AddDate(0, -1, 0), // Past date
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			wantErr: true,
		},
		{
			name: "ExpiredRenewal",
			payment: data.Payment{
				ID:             6,
				CustomerID:     1006,
				SubscriptionID: "sub_006",
				OrderID:        5006,
				Status:         "completed",
				VariantName:    "Premium",
				VariantID:      101,
				ProductID:      201,
				ProductName:    "Product A",
				CardBrand:      "Visa",
				CardLastFour:   "1234",
				UserName:       "Expired User",
				UserEmail:      "expired.user@example.com",
				RenewsAt:       time.Now().AddDate(0, -1, 0), // Past date
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			wantErr: true,
		},
		{
			name: "InvalidCardBrand",
			payment: data.Payment{
				ID:             7,
				CustomerID:     1007,
				SubscriptionID: "sub_007",
				OrderID:        5007,
				Status:         "completed",
				VariantName:    "Premium",
				VariantID:      101,
				ProductID:      201,
				ProductName:    "Product A",
				CardBrand:      "InvalidCardBrand",
				CardLastFour:   "1234",
				UserName:       "Invalid Card User",
				UserEmail:      "invalid.card.user@example.com",
				RenewsAt:       time.Now().AddDate(0, 1, 0),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			wantErr: true,
		},
		{
			name: "InvalidStatus",
			payment: data.Payment{
				ID:             8,
				CustomerID:     1008,
				SubscriptionID: "sub_008",
				OrderID:        5008,
				Status:         "invalid",
				VariantName:    "Premium",
				VariantID:      101,
				ProductID:      201,
				ProductName:    "Product A",
				CardBrand:      "Visa",
				CardLastFour:   "1234",
				UserName:       "Invalid Status User",
				UserEmail:      "invalid.status.user@example.com",
				RenewsAt:       time.Now().AddDate(0, 1, 0),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			wantErr: true,
		},
		{
			name: "ValidPendingPayment",
			payment: data.Payment{
				ID:             9,
				CustomerID:     1009,
				SubscriptionID: "sub_009",
				OrderID:        5009,
				Status:         "pending",
				VariantName:    "Basic",
				VariantID:      102,
				ProductID:      202,
				ProductName:    "Product B",
				CardBrand:      "MasterCard",
				CardLastFour:   "5678",
				UserName:       "Pending User",
				UserEmail:      "pending.user@example.com",
				RenewsAt:       time.Now().AddDate(0, 1, 0),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			wantErr: false,
		},
		{
			name: "ValidEmptyProductName",
			payment: data.Payment{
				ID:             10,
				CustomerID:     1010,
				SubscriptionID: "sub_010",
				OrderID:        5010,
				Status:         "completed",
				VariantName:    "Standard",
				VariantID:      103,
				ProductID:      203,
				ProductName:    "",
				CardBrand:      "Amex",
				CardLastFour:   "9012",
				UserName:       "Empty Product Name User",
				UserEmail:      "empty.product.name.user@example.com",
				RenewsAt:       time.Now().AddDate(0, 1, 0),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := suite.model.CreatePayment(suite.connection, tc.payment)
			if err != nil && !tc.wantErr {
				t.Errorf("CreatePayment() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			suite.paymentId = append(suite.paymentId, id)
			suite.subscriptionId = append(suite.subscriptionId, tc.payment.SubscriptionID)
		})
	}
}

func (suite *PaymentSuite) TestGetPaymentByID(t *testing.T) {
	testCases := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{"Valid ID - Existing Payment", suite.paymentId[0], false},
		{"Valid ID - Another Existing Payment", suite.paymentId[1], false},
		{"Non-Existing ID", 999, true},
		{"Boundary Condition - Lowest Valid ID", 0, true},
		{"Negative ID", -1, true},
		{"Very Large ID", 1000000000, true},
		{"ID Leading to DB Timeout", 12345, true},       // Simulate a scenario where the database operation times out.
		{"ID Causing Unexpected DB Error", 54321, true}, // Simulate a scenario where an unexpected error occurs in the database layer.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payment, err := suite.model.GetPaymentByID(suite.connection, tc.id)
			if err != nil {
				if !tc.wantErr {
					t.Errorf("GetPaymentByID() error = %v, wantErr %v", err, tc.wantErr)
				}
				// If an error was expected, there's no need to proceed with further checks.
				return
			}
			if tc.wantErr {
				t.Errorf("Expected error but got none")
				return
			}
			// Check for the scenario where payment is nil but no error was returned.
			if payment == nil {
				t.Errorf("GetPaymentByID() returned nil payment, but error was also nil")
				return
			}
			if tc.id != int(payment.ID) {
				t.Errorf("GetPaymentByID() got ID = %v, want %v", payment.ID, tc.id)
			}
		})
	}
}

func (suite *PaymentSuite) TestGetSubscriptionByID(t *testing.T) {
	testCases := []struct {
		name    string
		id      string // Changed type to string for subscription ID
		wantErr bool
	}{
		{"Valid ID - Existing Subscription", suite.subscriptionId[0], false},
		{"Valid ID - Another Existing Subscription", suite.subscriptionId[1], false},
		{"Non-Existing ID", "nonexistent", true},
		{"Empty ID", "", true},
		{"Very Long ID", "00000000-0000-0000-0000-00000000000000000", true},
		{"ID Leading to DB Timeout", "timeoutID", true},               // Simulate a scenario where the database operation times out.
		{"ID Causing Unexpected DB Error", "unexpectedErrorID", true}, // Simulate a scenario where an unexpected error occurs in the database layer.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			subscription, err := suite.model.GetPaymentBySubscriptionID(suite.connection, tc.id)
			if err != nil {
				if !tc.wantErr {
					t.Errorf("GetSubscriptionByID() error = %v, wantErr %v", err, tc.wantErr)
				}
				// If an error was expected, there's no need to proceed with further checks.
				return
			}
			if tc.wantErr {
				t.Errorf("Expected error but got none")
				return
			}
			// Check for the scenario where subscription is nil but no error was returned.
			if subscription == nil {
				t.Errorf("GetSubscriptionByID() returned nil subscription, but error was also nil")
				return
			}
			if tc.id != subscription.SubscriptionID {
				t.Errorf("GetSubscriptionByID() got ID = %v, want %v", subscription.ID, tc.id)
			}
		})
	}
}

func TestPaymentSuite(t *testing.T) {

	paymentSuite := PaymentSuite{}

	paymentSuite.SetupSuite()

	defer paymentSuite.TearDown()

	t.Run("TestCreatePayment", paymentSuite.TestCreatePayment)
	t.Run("TestGetPaymentByID", paymentSuite.TestGetPaymentByID)
	t.Run("TestGetPaymentBySubscriptionID", paymentSuite.TestGetSubscriptionByID)

}
