package main

import (
	"context"
	"fmt"
	"log"
	"payment-service/data"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

type Config struct {
	Models   data.Models // Data models for the application.
	Producer *Publisher  // Kafka producer for logging.
}

var app *Config

func init() {
	Producer := NewPublisher()                           // Create a new Kafka producer.
	Producer.createKafkaProducer("kafka:9092", "logger") // Configure the Kafka producer.      // Initialize the GitHub authenticator.
	app = &Config{                                       // Populate the global configuration.
		Producer: Producer,
	}
}

func main() {
	var wg sync.WaitGroup
	conn, err := connect() // Connect to the database.
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err) // Log and exit if the connection fails.
		app.Producer.publishMessage("key", "Payment Service", "Failed to connect to the database")
	}

	defer conn.Close(context.Background())

	e := echo.New()
	defer e.Close()
	app.Models = data.NewModels(conn)
	app.routes(e)
	wg.Add(1)
	go func() {
		initialWaitTime := 1 * time.Second // Initial wait time for retrying server start.
		maxRetries := 5                    // Maximum number of retries for starting the server.
		factor := 2                        // Factor by which the wait time increases.

		// Attempt to start the server with exponential backoff.
		for attempt := 1; attempt <= maxRetries; attempt++ {
			err := e.Start(":80") // Attempt to start the server.
			if err != nil {
				log.Printf("Attempt %d: server failed to start: %v", attempt, err) // Log the failure.
				if attempt == maxRetries {
					app.Producer.publishMessage("key", "Payment Service", "Server failed to start") // Log the final failure to Kafka.
					log.Fatalf("Server failed to start after %d attempts", maxRetries)              // Exit if the server fails to start after max retries.
				}
				time.Sleep(initialWaitTime)              // Wait before retrying.
				initialWaitTime *= time.Duration(factor) // Increase the wait time.
			} else {
				break // Exit the loop if the server starts successfully.
			}
		}
	}()
	wg.Wait()
}

// connect establishes a connection to the CockroachDB database.
// Returns a pointer to the database connection and an error, if any.
func connect() (*pgx.Conn, error) {
	url := "postgres://root@cockroach:26257/defaultdb?sslmode=disable" // Database connection URL.
	conn, err := pgx.Connect(context.Background(), url)                // Attempt to connect to the database.
	if err != nil {
		return nil, err // Return the error if the connection fails.
	}
	fmt.Println("Connected to the database") // Confirm successful connection.
	return conn, nil                         // Return the database connection.
}
