package main

import "github.com/labstack/echo/v4"

// route sets up the API routes for the Echo instance.
func route(e *echo.Echo) {
	e.GET("/ping", pingHandler) // health check
}
