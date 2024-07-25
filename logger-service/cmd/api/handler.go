package main

import (
	"log"
	"logger-service/data"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (app *Config) pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "The system is working fine")
}

type LogMessage struct {
	Service string `json:"service"`
	Message string `json:"message"`
}

// writeLogHandler handles the HTTP POST request for writing a log entry.
// It reads the request body and binds it to a LogMessage struct.
// If the request payload is invalid, it returns a JSON response with a status code of 400 (Bad Request).
// If the message or service fields are empty, it returns a JSON response with a status code of 400 (Bad Request).
// It creates a LogEntry struct with the log message, service, and current timestamp.
// If inserting the log entry fails, it returns a JSON response with a status code of 500 (Internal Server Error).
// Otherwise, it returns a JSON response with a status code of 201 (Created) and the created log entry.
func (app *Config) writeLogHandler(c echo.Context) error {
	var logmessage LogMessage
	if err := c.Bind(&logmessage); err != nil {
		log.Fatalf("Failed reading the request body: %s", err)
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}
	if logmessage.Message == "" || logmessage.Service == "" {
		return c.JSON(http.StatusBadRequest, "Message and service are required")
	}

	logentry := data.LogEntry{
		Message:   logmessage.Message,
		Service:   logmessage.Service,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := app.Models.LogEntry.Insert(app.client, logentry); err != nil {
		log.Fatalf("Failed inserting log entry: %s", err)
		return c.JSON(http.StatusInternalServerError, "Failed inserting log entry")
	}
	return c.JSON(http.StatusCreated, logentry)
}
