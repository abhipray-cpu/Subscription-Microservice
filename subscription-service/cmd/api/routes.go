package main

import "github.com/labstack/echo/v4"

// routes registers the API routes with the provided Echo instance.
func (app *Config) routes(e *echo.Echo) {
	e.GET("/ping", app.pingHandler)                      // Health check endpoint.
	e.GET("/auth/:provider/callback", app.Auth.CallBack) // OAuth callback endpoint.
	e.GET("/logout/:provider", app.Auth.Logout)          // Logout endpoint.
	e.GET("/auth/:provider", app.Auth.Auth)              // OAuth authentication endpoint.
}
