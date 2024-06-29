package main

import (
	"context"
	"fmt"
	"log"
	"subscription-service/data"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
)

// Config represents the configuration for the application.
type Config struct {
	Models data.Models
	Writer *kafka.Writer
}

func main() {
	// connect to the database
	conn, err := connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	defer conn.Close(context.Background())
	e := echo.New()
	defer e.Close()
	producer, err := createKafkaProducer("kafka:9092", "logger")
	if err != nil {
		log.Fatalf("Failed to create kafka producer: %v", err)
	}

	app := Config{
		Models: data.New(conn),
		Writer: producer,
	}
	app.routes(e)
	// testing the producer
	message := Message{
		Service: "subscription-service",
		Message: "Hello from subscription-service",
	}
	err = app.publishMessage("key", message)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	} else {
		fmt.Println("Message published successfully")
	}
	NewAuth()
	initialWaitTime := 1 * time.Second
	maxRetries := 5
	factor := 2

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := e.Start(":80")
		if err != nil {
			log.Printf("Attempt %d: server failed to start: %v", attempt, err)
			if attempt == maxRetries {
				log.Fatalf("Server failed to start after %d attempts", maxRetries)
			}
			time.Sleep(initialWaitTime)
			initialWaitTime *= time.Duration(factor)
		} else {
			break
		}
	}
}

// connect connects to the CockroachDB database.
func connect() (*pgx.Conn, error) {
	url := "postgres://root@cockroach:26257/defaultdb?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to the database")
	return conn, nil
}
