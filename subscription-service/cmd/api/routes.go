package main

import (
	"github.com/labstack/echo/v4"
)

// routes registers the API routes with the provided Echo instance.
func (app *Config) routes(e *echo.Echo) {
	g := e.Group("/account")
	g.Use(JWTAuthMiddleware)
	e.GET("/ping", app.pingHandler)                      // Health check endpoint.
	e.GET("/auth/:provider/callback", app.Auth.CallBack) // OAuth callback endpoint.
	e.GET("/logout/:provider", app.Auth.Logout)          // Logout endpoint.
	e.GET("/auth/:provider", app.Auth.Auth)              // OAuth authentication endpoint.
	e.POST("/signup", app.signup)                        // Signup endpoint.
	e.POST("/login", app.login)                          // Login endpoint.                           // Logout endpoint.
	g.DELETE("/", app.deleteAccount)                     // Delete account endpoint.
	g.GET("/", app.getAccount)                           // Get account endpoint.
	g.PUT("/", app.updateAccount)                        // Update account endpoint.

}
