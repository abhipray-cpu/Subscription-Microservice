package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Message represents a Kafka message.
type Message struct {
	Service string // Service represents the name of the service.
	Message string // Message represents the content of the message.
}

// createKafkaProducer creates a Kafka producer and returns it.
func createKafkaProducer(broker string, topic string) (*kafka.Writer, error) {

	// Create Kafka writer configuration
	config := kafka.WriterConfig{
		Brokers:  []string{broker}, // Brokers is a list of Kafka broker addresses.
		Topic:    topic,            // Topic is the name of the Kafka topic.
		Balancer: &kafka.LeastBytes{},
	}

	// Create Kafka writer
	writer := kafka.NewWriter(config)

	return writer, nil
}

// publishMessage publishes a message to Kafka.
func (app *Config) publishMessage(key string, value Message) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	msg := kafka.Message{
		Key:   []byte(key),     // Key is the message key.
		Value: valueBytes,      // Value is the message value.
	}
	err = app.Writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("failed to write messages: %w", err)
	}

	fmt.Println("Message published successfully")
	return nil
}
