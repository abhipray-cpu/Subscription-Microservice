package main

import "github.com/labstack/echo/v4"

func (app *Config) routes(e *echo.Echo) {
	s := e.Group("/subscription")
	s.Use(VerifySignatureMiddleware)
	e.GET("/ping", app.pingHandler)
	s.POST("/created", app.SubscriptionCreated)
	s.POST("/updated", app.SubscriptionUpdated)
	s.POST("/cancelled", app.SubscriptionCancelled)
	s.POST("/resumed", app.SubscriptionResumed)
	s.POST("/expired", app.SubscriptionExpired)
	s.POST("/paused", app.SubscriptionPaused)
	s.POST("/unpaused", app.SubscriptionUnpaused)
	s.POST("/failed", app.SubscriptionFailedPayment)
	s.POST("/success", app.SubscriptionSucessPayment)
	s.POST("/recovered", app.SubscriptionRecovered)
	s.POST("/refunded", app.SubscriptionRefunded)
	s.POST("/changed", app.SubscriptionChanged)
}
