package main

import "github.com/labstack/echo/v4"

// routes registers the API routes with the provided Echo instance.
func (app *Config) routes(e *echo.Echo) {
	e.GET("/ping", app.pingHandler) // health check
	e.POST("/write-log", app.writeLogHandler) // write log
}
