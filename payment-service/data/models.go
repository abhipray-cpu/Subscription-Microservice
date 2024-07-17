// Package payment provides a set of functions and types to work with payment records in a PostgreSQL database.
package data

import (
	"context" // Used for managing the lifetime of database operations.
	"encoding/json"
	"log"  // Used for logging errors.
	"time" // Used for handling time-related data.

	"github.com/jackc/pgx/v4" // PostgreSQL driver for Go.
)

// Payment represents a single payment transaction.
// It includes details such as the user, subscription, and transaction IDs, the amount, currency, and method of payment,
// the status of the payment, and timestamps for the payment, next billing date, and record creation and update.
type Payment struct {
	ID             int64     `json:"id"`             // Unique identifier for the payment record.
	CustomerID     float64   `json:"userId"`         // Unique identifier for the user making the payment.
	SubscriptionID string    `json:"subscriptionId"` // Identifier for the subscription the payment is for.
	OrderID        float64   `json:"orderId"`        // Unique identifier for the order.
	Status         string    `json:"status"`         // The status of the payment (e.g., completed, pending).
	VariantName    string    `json:"variantName"`    // Name of the variant.
	VariantID      float64   `json:"variantId"`      // Unique identifier for the variant.
	ProductID      float64   `json:"productId"`      // Unique identifier for the product.
	ProductName    string    `json:"productName"`    // Name of the product.
	CardBrand      string    `json:"cardBrand"`      // Brand of the card used for payment.
	CardLastFour   string    `json:"cardLastFour"`   // Last four digits of the card used for payment.
	UserName       string    `json:"userName"`       // Name of the user making the payment.
	UserEmail      string    `json:"userEmail"`      // Email of the user making the payment.
	RenewsAt       time.Time `json:"renewsAt"`       // Timestamp of when the subscription renews.
	CreatedAt      time.Time `json:"createdAt"`      // Timestamp of when the record was created.
	UpdatedAt      time.Time `json:"updatedAt"`      // Timestamp of the last update to the record.
}

// connection holds a global database connection, shared across instances of Models.
// This allows for a single database connection to be reused for multiple operations, improving efficiency.
var connection *pgx.Conn

// Models wraps all the models in the application for easy access.
// Currently, it only contains a Payment model, but it can be expanded to include more models.
type Models struct {
	Payment Payment // Instance of the Payment model.
}

// NewModels initializes a new instance of Models with a database connection.
// It sets the global database connection and ensures the necessary table exists in the database.
func NewModels(conn *pgx.Conn) Models {
	connection = conn       // Set the global connection.
	ensureTableExists(conn) // Ensure the payments table exists in the database.
	return Models{
		Payment: Payment{}, // Initialize the Payment model.
	}
}

// ensureTableExists updated to include new fields
func ensureTableExists(conn *pgx.Conn) {
	query := `
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
        card_last_four CHAR(4),
        user_name VARCHAR(255),
        user_email VARCHAR(255),
        renews_at TIMESTAMP NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`

	if _, err := conn.Exec(context.Background(), query); err != nil {
		log.Fatalf("Failed to ensure payments table exists: %v", err)
	}
}

// CreatePayment updated to include new fields
func (m *Models) CreatePayment(p Payment) (int, error) {
	var id int // Variable to store the ID of the created payment
	query := `
    INSERT INTO payments (customer_id, subscription_id, order_id, status, variant_name, variant_id, product_id, product_name, card_brand, card_last_four, user_name, user_email, renews_at, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
    RETURNING id;`

	err := connection.QueryRow(context.Background(), query,
		p.CustomerID, p.SubscriptionID, p.OrderID, p.Status, p.VariantName, p.VariantID, p.ProductID, p.ProductName, p.CardBrand, p.CardLastFour, p.UserName, p.UserEmail, p.RenewsAt, p.CreatedAt, p.UpdatedAt).Scan(&id)
	if err != nil {
		log.Printf("Failed to create payment: %v", err)
		return 0, err // Return 0 for the ID in case of an error
	}
	return id, nil // Return the ID of the created payment and nil for the error
}

// GetPaymentByID updated to include new fields
func (m *Models) GetPaymentByID(id int) (*Payment, error) {
	query := `
    SELECT id, customer_id, subscription_id, order_id, status, variant_name, variant_id, product_id, product_name, card_brand, card_last_four, user_name, user_email, renews_at, created_at, updated_at
    FROM payments
    WHERE id = $1;`

	var p Payment
	err := connection.QueryRow(context.Background(), query, id).Scan(&p.ID, &p.CustomerID, &p.SubscriptionID, &p.OrderID, &p.Status, &p.VariantName, &p.VariantID, &p.ProductID, &p.ProductName, &p.CardBrand, &p.CardLastFour, &p.UserName, &p.UserEmail, &p.RenewsAt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		log.Printf("Failed to get payment by ID: %v", err)
		return nil, err
	}
	return &p, nil
}

func (m *Models) GetPaymentBySubscriptionID(subscriptionID string) (*Payment, error) {
	query := `
    SELECT id, customer_id, subscription_id, order_id, status, variant_name, variant_id, product_id, product_name, card_brand, card_last_four, user_name, user_email, renews_at, created_at, updated_at
    FROM payments
    WHERE subscription_id = $1;`

	var p Payment
	err := connection.QueryRow(context.Background(), query, subscriptionID).Scan(&p.ID, &p.CustomerID, &p.SubscriptionID, &p.OrderID, &p.Status, &p.VariantName, &p.VariantID, &p.ProductID, &p.ProductName, &p.CardBrand, &p.CardLastFour, &p.UserName, &p.UserEmail, &p.RenewsAt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		log.Printf("Failed to get payment by subscription ID: %v", err)
		return nil, err
	}
	return &p, nil
}

// UpdatePayment updated to include new fields
func (m *Models) UpdatePayment(p Payment) error {
	query := `
    UPDATE payments
    SET customer_id = $2, subscription_id = $3, order_id = $4, status = $5, variant_name = $6, variant_id = $7, product_id = $8, product_name = $9, card_brand = $10, card_last_four = $11, user_name = $12, user_email = $13, renews_at = $14, updated_at = $15
    WHERE id = $1;`

	_, err := connection.Exec(context.Background(), query, p.ID, p.CustomerID, p.SubscriptionID, p.OrderID, p.Status, p.VariantName, p.VariantID, p.ProductID, p.ProductName, p.CardBrand, p.CardLastFour, p.UserName, p.UserEmail, p.RenewsAt, time.Now())
	if err != nil {
		log.Printf("Failed to update payment: %v", err)
		return err
	}
	return nil
}

func GetPayment(body []byte) (*Payment, error) {
	var params map[string]interface{}
	err := json.Unmarshal(body, &params)
	if err != nil {
		return nil, err
	}
	renew_string := params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["renews_at"].(string)
	created_string := params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["created_at"].(string)
	updated_string := params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["updated_at"].(string)

	// Parse the time strings into time.Time objects
	renewsAt, err := time.Parse(time.RFC3339, renew_string)
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, created_string)
	if err != nil {
		return nil, err
	}

	updatedAt, err := time.Parse(time.RFC3339, updated_string)
	if err != nil {
		return nil, err
	}
	payment := Payment{
		CustomerID:     params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["customer_id"].(float64),
		SubscriptionID: params["data"].(map[string]interface{})["id"].(string),
		OrderID:        params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["order_id"].(float64),
		Status:         params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["status"].(string),
		VariantName:    params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["variant_name"].(string),
		VariantID:      params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["variant_id"].(float64),
		ProductID:      params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["product_id"].(float64),
		ProductName:    params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["product_name"].(string),
		CardBrand:      params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["card_brand"].(string),
		CardLastFour:   params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["card_last_four"].(string),
		UserName:       params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["user_name"].(string),
		UserEmail:      params["data"].(map[string]interface{})["attributes"].(map[string]interface{})["user_email"].(string),
		RenewsAt:       renewsAt,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
	return &payment, nil
}
