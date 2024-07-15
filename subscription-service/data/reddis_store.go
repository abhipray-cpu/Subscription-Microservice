package data

import (
	"log"

	"github.com/boj/redistore"
)

// NewRedisStore creates a new Redis store with the specified configuration.
// It returns a pointer to the created Redis store and an error, if any.

func NewRedisStore(size int, network, address, password string, keyPairs []byte) (*redistore.RediStore, error) {
	store, err := redistore.NewRediStore(size, network, address, password, keyPairs)
	if err != nil {
		log.Fatalf("Failed to create Redis store:%v", err)
		return nil, err
	}
	log.Println("Successfully created Redis store")
	return store, nil
}
