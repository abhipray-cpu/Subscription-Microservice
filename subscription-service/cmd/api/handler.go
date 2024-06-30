package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// pingHandler handles the ping request and returns a success message.
func (app *Config) pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "The system is working fine")
}
