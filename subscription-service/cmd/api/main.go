package main

import (
	"context"                   // Provides functionality for managing request lifecycles.
	"fmt"                       // Used for formatting and printing output.
	"log"                       // Used for logging error messages.
	"subscription-service/auth" // Custom package for authentication.
	"subscription-service/data" // Custom package for data models.
	"time"                      // Used for time-related operations, such as delays.

	"github.com/jackc/pgx/v4"     // PostgreSQL driver for Go.
	"github.com/labstack/echo/v4" // Echo framework for building web applications.
)

// Config holds the application-wide configurations.
type Config struct {
	Models   data.Models        // Data models for the application.
	Auth     auth.Authenticator // Authentication mechanism.
	Producer *Publisher         // Kafka producer for logging.
}

var app *Config // Global variable to hold the application configuration.

// init is called before the main function. It initializes the application configuration.
func init() {
	Producer := NewPublisher()                           // Create a new Kafka producer.
	Producer.createKafkaProducer("kafka:9092", "logger") // Configure the Kafka producer.
	authenticator := auth.NewGitHubAuthenticator()       // Initialize the GitHub authenticator.
	app = &Config{                                       // Populate the global configuration.
		Producer: Producer,
		Auth:     authenticator,
	}
	// Attempt to publish a startup message to Kafka.
	err := Producer.publishMessage("key", "subscription-service", "Hello from subscription-service")
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err) // Log and exit if the message fails to publish.
	} else {
		fmt.Println("Message published successfully") // Confirm successful message publication.
	}
}

func main() {
	conn, err := connect() // Connect to the database.
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err) // Log and exit if the connection fails.
		app.Producer.publishMessage("key", "Subscription Service", "Failed to connect to the database")
	}

	defer conn.Close(context.Background()) // Ensure the database connection is closed on exit.
	e := echo.New()                        // Create a new Echo instance for the web server.
	defer e.Close()                        // Ensure the Echo server is closed on exit.

	app.Models = data.NewModels(conn) // Initialize the data models.
	app.routes(e)                     // Set up the web routes.

	app.Auth.NewAuth()                 // Initialize the authentication mechanism.
	initialWaitTime := 1 * time.Second // Initial wait time for retrying server start.
	maxRetries := 5                    // Maximum number of retries for starting the server.
	factor := 2                        // Factor by which the wait time increases.

	// Attempt to start the server with exponential backoff.
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := e.Start(":80") // Attempt to start the server.
		if err != nil {
			log.Printf("Attempt %d: server failed to start: %v", attempt, err) // Log the failure.
			if attempt == maxRetries {
				app.Producer.publishMessage("key", "Subscription Service", "Server failed to start") // Log the final failure to Kafka.
				log.Fatalf("Server failed to start after %d attempts", maxRetries)                   // Exit if the server fails to start after max retries.
			}
			time.Sleep(initialWaitTime)              // Wait before retrying.
			initialWaitTime *= time.Duration(factor) // Increase the wait time.
		} else {
			break // Exit the loop if the server starts successfully.
		}
	}
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
