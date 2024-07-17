package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Publisher holds a Kafka writer instance for publishing messages.
type Publisher struct {
	Writer *kafka.Writer // Writer is a Kafka writer used to publish messages.
}

// Message represents the structure of a message to be published to Kafka.
type Message struct {
	Service string // Service indicates the name of the service sending the message.
	Message string // Message contains the actual message content.
}

// createKafkaProducer initializes a Kafka producer with the specified broker and topic.
// This method configures a Kafka writer for the Publisher instance.
//
// Parameters:
// - broker: A string representing the address of the Kafka broker.
// - topic: A string representing the name of the Kafka topic to which messages will be published.
//
// The method does not return any value. It directly assigns a new Kafka writer to the Publisher's Writer field.
// The Kafka writer is configured with the LeastBytes balancer, which aims to distribute messages evenly across partitions based on their size.
func (publisher *Publisher) createKafkaProducer(broker string, topic string) {
	// Configuration for the Kafka writer, including the broker addresses and the target topic.
	config := kafka.WriterConfig{
		Brokers:  []string{broker},    // List of Kafka broker addresses.
		Topic:    topic,               // Name of the Kafka topic.
		Balancer: &kafka.LeastBytes{}, // Balancer for distributing messages across partitions.
	}

	// Instantiating the Kafka writer with the specified configuration.
	publisher.Writer = kafka.NewWriter(config)
}

// publishMessage sends a message to the Kafka topic configured in the Publisher's writer.
//
// Parameters:
// - key: A string representing the key of the message. Kafka uses this for partitioning.
// - service: The name of the service sending the message.
// - message: The content of the message to be sent.
//
// Returns:
//   - An error if the message could not be marshaled into JSON or if writing the message to Kafka fails.
//     Otherwise, it returns nil indicating the message was published successfully.
//
// This method marshals a Message struct into JSON and publishes it to Kafka using the Writer field of the Publisher.
// It logs a success message upon successful publication.
func (publisher *Publisher) publishMessage(key string, service, message string) error {
	// Constructing the message with the provided service name and message content.
	value := Message{
		Service: service,
		Message: message,
	}
	// Marshaling the message into JSON format.
	valueBytes, err := json.Marshal(value)
	if err != nil {
		// Returning an error if marshaling fails.
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// Creating a Kafka message with the provided key and the marshaled JSON message as the value.
	msg := kafka.Message{
		Key:   []byte(key), // The message key for partitioning.
		Value: valueBytes,  // The JSON marshaled message.
	}
	// Writing the message to Kafka.
	err = publisher.Writer.WriteMessages(context.Background(), msg)
	if err != nil {
		// Returning an error if writing to Kafka fails.
		return fmt.Errorf("failed to write messages: %w", err)
	}

	// Logging a success message.
	fmt.Println("Message published successfully")
	return nil
}

// NewPublisher creates and returns a new instance of Publisher.
// This function is a constructor for the Publisher type.
//
// Returns:
// - A pointer to a new Publisher instance.
//
// The returned Publisher instance is initialized with default values, ready to be configured for message publishing.
func NewPublisher() *Publisher {
	return &Publisher{}
}
