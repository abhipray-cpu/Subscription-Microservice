package main

import "github.com/labstack/echo/v4"

func (app *Config) routes(e *echo.Echo) {
	s := e.Group("/subscription")                     // Create a new group for subscription-related route
	s.Use(VerifySignatureMiddleware)                  // Add the VerifySignatureMiddleware to the group
	e.GET("/ping", app.pingHandler)                   // Add a ping route to check if the server is running
	s.POST("/created", app.SubscriptionCreated)       // Add a route for handling subscription creation events
	s.POST("/updated", app.SubscriptionUpdated)       // Add a route for handling subscription update events
	s.POST("/cancelled", app.SubscriptionCancelled)   // Add a route for handling subscription cancellation events
	s.POST("/resumed", app.SubscriptionResumed)       // Add a route for handling subscription resumption events
	s.POST("/expired", app.SubscriptionExpired)       // Add a route for handling subscription expiration events
	s.POST("/paused", app.SubscriptionPaused)         // Add a route for handling subscription pause events
	s.POST("/unpaused", app.SubscriptionUnpaused)     // Add a route for handling subscription unpause events
	s.POST("/failed", app.SubscriptionFailedPayment)  // Add a route for handling failed payment events
	s.POST("/success", app.SubscriptionSucessPayment) // Add a route for handling successful payment events
	s.POST("/recovered", app.SubscriptionRecovered)   // Add a route for handling recovered payment events
	s.POST("/refunded", app.SubscriptionRefunded)     // Add a route for handling refunded payment events
	s.POST("/changed", app.SubscriptionChanged)       // Add a route for handling subscription change events
}
