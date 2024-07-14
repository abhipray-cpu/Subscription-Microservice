package worker

import (
	"log"
	"subscription-service/worker/workflow"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func StartWorker() {
	c, err := client.Dial(client.Options{
		HostPort: "temporal:7233",
	})
	if err != nil {
		log.Fatalf("failed to create Temporal client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, "subscription-service", worker.Options{})

	// Registering workflows
	w.RegisterWorkflow(workflow.WelcomWorkFlow)
	w.RegisterWorkflow(workflow.OTPWorkflow)

	// start the worker
	// Run the Worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}

}
