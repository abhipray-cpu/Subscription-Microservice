package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// writeLog sends a log message to the logger service.
// It takes the service name and log message as parameters.
// The log message is sent as a JSON payload to the logger service's endpoint.
// If there is an error while marshaling the JSON payload or sending the POST request,
// an error message is printed to the console.
// If the response status code is not 200 OK, an error message is printed to the console.
// Otherwise, a success message is printed to the console.
func writeLog(service, message string) {
	data := map[string]string{
		"service": service,
		"message": message,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post("http://logger-service:80/write-log", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected response status:", resp.StatusCode)
		return
	}

	fmt.Println("Log sent successfully")
}

func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "ping")
}
