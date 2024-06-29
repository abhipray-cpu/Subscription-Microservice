package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	Service string `json:"service"`
	Message string `json:"message"`
}

// initializeConsumer sets up a Kafka consumer to read messages from a specified topic.
// It takes the Kafka broker addresses and the topic name as input parameters.
// The function continuously reads messages from Kafka and processes them.
// It logs any errors encountered during the process.
func initializeConsumer(brokers string, topic string) {
	// Set up Kafka reader configuration
	config := kafka.ReaderConfig{
		Brokers:  []string{brokers}, // Update with your Kafka broker addresses
		Topic:    topic,             // Update with your topic name
		MinBytes: 10e3,              // Minimum number of bytes to fetch from Kafka
		MaxBytes: 10e6,              // Maximum number of bytes to fetch from Kafka
	}

	// Create Kafka reader
	reader := kafka.NewReader(config)

	defer reader.Close()
	// Create a loop to continuously read messages
	for {
		// Read a message from Kafka
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			// Handle error
			log.Printf("Error reading message: %v", err)
			continue
		}

		msgString := string(msg.Value)

		var messageObject Message
		err = json.Unmarshal([]byte(msgString), &messageObject)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			return
		}

		// Now you can access the data in a structured format
		log.Printf("Received message from service %s: %s", messageObject.Service, messageObject.Message)
		writeLog(messageObject.Service, messageObject.Message)
	}
}
