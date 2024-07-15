package worker

import (
	"fmt"

	"go.temporal.io/sdk/client"
)

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
