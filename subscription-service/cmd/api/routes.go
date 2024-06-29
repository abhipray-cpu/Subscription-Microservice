package main

import "github.com/labstack/echo/v4"

// routes registers the API routes with the provided Echo instance.
func (app *Config) routes(e *echo.Echo) {
	e.GET("/ping", app.pingHandler) // Health check endpoint.
	e.GET("/auth/:provider/callback", app.callBackHandler) // OAuth callback endpoint.
	e.GET("/logout/:provider", app.logoutHandler) // Logout endpoint.
	e.GET("/auth/:provider", app.authHandler) // OAuth authentication endpoint.
}
