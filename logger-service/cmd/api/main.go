package main

import (
	"context"
	"log"
	"logger-service/data"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const mongoURL = "mongodb://mongo:27017"

type Config struct {
	Models data.Models
	client *mongo.Client
}

// main is the entry point of the application.
// It connects to MongoDB, initializes the application configuration,
// sets up the routes, and starts the server.
// It also handles server start failures by retrying with an exponential backoff strategy.
func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	e := echo.New()
	app := Config{
		Models: data.New(),
		client: client,
	}
	app.routes(e)
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

// connectToMongo connects to the MongoDB instance and returns the client.
func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	log.Println("Connected to MongoDB")
	return c, nil
}
