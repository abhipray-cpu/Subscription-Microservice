package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8" // Import the go-redis package
)

var ctx = context.Background()

var redisClient *redis.Client

// NewRedisClient creates and returns a new Redis client
func NewRedisClient(address, password string) error {
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
		return err
	}

	log.Println("Successfully connected to Redis")
	redisClient = client
	return nil
}

// Create a new key-value pair with an expiration of 60 minutes
func CreateKey(key, value string) error {
	if len(value) != 6 {
		return fmt.Errorf("value must be 6 characters long")
	}
	status := redisClient.Set(ctx, key, value, 60*time.Minute)
	if err := status.Err(); err != nil {
		log.Printf("Failed to create key: %v", err)
		return err
	}
	log.Printf("Key %s created successfully", key)
	return nil
}

// Read the value of a given key
func ReadKey(key string) (string, error) {
	result, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Key %s does not exist", key)
		return "", fmt.Errorf("key %s does not exist", key)
	} else if err != nil {
		log.Printf("Failed to read key: %v", err)
		return "", err
	}
	log.Printf("Value for key %s is %s", key, result)
	return result, nil
}

// Update the value of an existing key and reset its expiration to 60 minutes
func UpdateKey(key, newValue string) error {
	if len(newValue) != 6 {
		return fmt.Errorf("value must be 6 characters long")
	}
	status := redisClient.Set(ctx, key, newValue, 60*time.Minute)
	if err := status.Err(); err != nil {
		log.Printf("Failed to update key: %v", err)
		return err
	}
	log.Printf("Key %s updated successfully", key)
	return nil
}

// Delete a key from the store
func DeleteKey(key string) error {
	result, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		log.Printf("Failed to delete key: %v", err)
		return err
	}
	if result == 0 {
		log.Printf("Key %s does not exist", key)
		return fmt.Errorf("key %s does not exist", key)
	}
	log.Printf("Key %s deleted successfully", key)
	return nil
}
