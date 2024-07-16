package worker

import (
	"fmt"

	"go.temporal.io/sdk/client"
)

// StartWorker creates a new Temporal client and returns it or an error.
// It returns an error if the client cannot be created.
func StartWorker() (client.Client, error) {
	c, err := client.Dial(client.Options{
		HostPort: "temporal:7233",
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to Temporal")
	return c, nil
}
