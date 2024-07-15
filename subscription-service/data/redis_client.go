package data

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8" // Import the go-redis package
)

var ctx = context.Background()

// NewRedisClient creates and returns a new Redis client
func NewRedisClient(address, password string) (*redis.Client, error) {
	// Initialize a new Redis client with the given address and password
	client := redis.NewClient(&redis.Options{
		Addr:     address,  // Redis server address
		Password: password, // No password by default
		DB:       0,        // Default DB
	})

	// Test the connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to Redis")
	return client, nil
}
