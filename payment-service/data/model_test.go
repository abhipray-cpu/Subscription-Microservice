package data

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
)

func TestModels_CreatePayment(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	models := NewModels(conn)
	payment := Payment{
		CustomerID:     123,
		SubscriptionID: "subscription123",
		OrderID:        456,
		Status:         "completed",
		VariantName:    "variant",
		VariantID:      789,
		ProductID:      987,
		ProductName:    "product",
		CardBrand:      "Visa",
		CardLastFour:   "1234",
		UserName:       "john_doe",
		UserEmail:      "john.doe@example.com",
		RenewsAt:       time.Now().Add(time.Hour),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	id, err := models.CreatePayment(payment)
	if err != nil {
		t.Errorf("Failed to create payment: %v", err)
	}

	// Verify that the ID is not zero
	if id == 0 {
		t.Error("Invalid payment ID")
	}
}

func TestModels_GetPaymentByID(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	models := NewModels(conn)
	payment := Payment{
		CustomerID:     123,
		SubscriptionID: "subscription123",
		OrderID:        456,
		Status:         "completed",
		VariantName:    "variant",
		VariantID:      789,
		ProductID:      987,
		ProductName:    "product",
		CardBrand:      "Visa",
		CardLastFour:   "1234",
		UserName:       "john_doe",
		UserEmail:      "john.doe@example.com",
		RenewsAt:       time.Now().Add(time.Hour),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	id, err := models.CreatePayment(payment)
	if err != nil {
		t.Fatalf("Failed to create payment: %v", err)
	}

	p, err := models.GetPaymentByID(id)
	if err != nil {
		t.Errorf("Failed to get payment by ID: %v", err)
	}

	// Verify that the retrieved payment matches the created payment
	if p.CustomerID != payment.CustomerID ||
		p.SubscriptionID != payment.SubscriptionID ||
		p.OrderID != payment.OrderID ||
		p.Status != payment.Status ||
		p.VariantName != payment.VariantName ||
		p.VariantID != payment.VariantID ||
		p.ProductID != payment.ProductID ||
		p.ProductName != payment.ProductName ||
		p.CardBrand != payment.CardBrand ||
		p.CardLastFour != payment.CardLastFour ||
		p.UserName != payment.UserName ||
		p.UserEmail != payment.UserEmail ||
		!p.RenewsAt.Equal(payment.RenewsAt) ||
		!p.CreatedAt.Equal(payment.CreatedAt) ||
		!p.UpdatedAt.Equal(payment.UpdatedAt) {
		t.Error("Retrieved payment does not match the created payment")
	}
}

func TestModels_GetPaymentBySubscriptionID(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	models := NewModels(conn)
	payment := Payment{
		CustomerID:     123,
		SubscriptionID: "subscription123",
		OrderID:        456,
		Status:         "completed",
		VariantName:    "variant",
		VariantID:      789,
		ProductID:      987,
		ProductName:    "product",
		CardBrand:      "Visa",
		CardLastFour:   "1234",
		UserName:       "john_doe",
		UserEmail:      "john.doe@example.com",
		RenewsAt:       time.Now().Add(time.Hour),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err = models.CreatePayment(payment)
	if err != nil {
		t.Fatalf("Failed to create payment: %v", err)
	}

	p, err := models.GetPaymentBySubscriptionID(payment.SubscriptionID)
	if err != nil {
		t.Errorf("Failed to get payment by subscription ID: %v", err)
	}

	// Verify that the retrieved payment matches the created payment
	if p.CustomerID != payment.CustomerID ||
		p.SubscriptionID != payment.SubscriptionID ||
		p.OrderID != payment.OrderID ||
		p.Status != payment.Status ||
		p.VariantName != payment.VariantName ||
		p.VariantID != payment.VariantID ||
		p.ProductID != payment.ProductID ||
		p.ProductName != payment.ProductName ||
		p.CardBrand != payment.CardBrand ||
		p.CardLastFour != payment.CardLastFour ||
		p.UserName != payment.UserName ||
		p.UserEmail != payment.UserEmail ||
		!p.RenewsAt.Equal(payment.RenewsAt) ||
		!p.CreatedAt.Equal(payment.CreatedAt) ||
		!p.UpdatedAt.Equal(payment.UpdatedAt) {
		t.Error("Retrieved payment does not match the created payment")
	}
}

func TestModels_UpdatePayment(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	models := NewModels(conn)
	payment := Payment{
		CustomerID:     123,
		SubscriptionID: "subscription123",
		OrderID:        456,
		Status:         "completed",
		VariantName:    "variant",
		VariantID:      789,
		ProductID:      987,
		ProductName:    "product",
		CardBrand:      "Visa",
		CardLastFour:   "1234",
		UserName:       "john_doe",
		UserEmail:      "john.doe@example.com",
		RenewsAt:       time.Now().Add(time.Hour),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	id, err := models.CreatePayment(payment)
	if err != nil {
		t.Fatalf("Failed to create payment: %v", err)
	}

	// Update the payment
	payment.Status = "pending"
	payment.VariantName = "new_variant"
	payment.UpdatedAt = time.Now()

	err = models.UpdatePayment(payment)
	if err != nil {
		t.Errorf("Failed to update payment: %v", err)
	}

	// Retrieve the updated payment
	p, err := models.GetPaymentByID(id)
	if err != nil {
		t.Errorf("Failed to get payment by ID: %v", err)
	}

	// Verify that the payment has been updated
	if p.Status != payment.Status ||
		p.VariantName != payment.VariantName ||
		!p.UpdatedAt.Equal(payment.UpdatedAt) {
		t.Error("Payment update failed")
	}
}
