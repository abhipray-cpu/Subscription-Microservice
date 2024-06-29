package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

// main is the entry point of the application.
// It initializes the Echo server, sets up routes, and starts the server.
// It also handles server startup retries in case of failure.
func main() {
	e := echo.New()
	route(e)
	initializeConsumer("kafka:9092", "logger")
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
