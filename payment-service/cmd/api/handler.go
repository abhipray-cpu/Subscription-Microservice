package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"payment-service/data"
	"payment-service/grpc/subscription"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func (app *Config) pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "system is working")
}

func (app *Config) SubscriptionCreated(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	id, err := app.Models.CreatePayment(*payment)
	if err != nil {
		go processSubscription("failed create", payment.UserEmail, "failed", payment.ProductName, payment.VariantName)
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil

	}
	app.Producer.publishMessage("key", "Payment Service", "Subscription created successfully")
	go processSubscription(strconv.Itoa(id), payment.UserEmail, payment.Status, payment.ProductName, payment.VariantName)
	return nil
}

func (app *Config) SubscriptionUpdated(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		go processSubscription("failed update", payment.UserEmail, "failed", payment.ProductName, payment.VariantName)
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		go processSubscription("failed update", payment.UserEmail, "failed", payment.ProductName, payment.VariantName)
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	go processSubscription("success update", payment.UserEmail, "failed", payment.ProductName, payment.VariantName)
	return nil
}

func (app *Config) SubscriptionCancelled(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		go processSubscription("failed cancel", payment.UserEmail, "failed", payment.ProductName, payment.VariantName)
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		go processSubscription("failed cancel", payment.UserEmail, "failed", payment.ProductName, payment.VariantName)
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	go processSubscription("success cancel", payment.UserEmail, "failed", payment.ProductName, payment.VariantName)
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionResumed(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		go processSubscription("failed resume", payment.UserEmail, "failed", payment.ProductName, payment.VariantName)

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		go processSubscription("failed resume", payment.UserEmail, "failed", payment.ProductName, payment.VariantName)
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	// run the temporal workflow
	go processSubscription("success resume", payment.UserEmail, payment.Status, payment.ProductName, payment.VariantName)
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionExpired(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	go processSubscription("expired", payment.UserEmail, payment.Status, payment.ProductName, payment.VariantName)
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionPaused(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	go processSubscription("paused", payment.UserEmail, payment.Status, payment.ProductName, payment.VariantName)
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionUnpaused(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	go processSubscription("unpaused", payment.UserEmail, payment.Status, payment.ProductName, payment.VariantName)
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionFailedPayment(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	go processSubscription("payment", payment.UserEmail, payment.Status, payment.ProductName, payment.VariantName)
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionSucessPayment(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	go processSubscription("payment success", payment.UserEmail, payment.Status, payment.ProductName, payment.VariantName)
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionRecovered(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	go processSubscription("recovered", payment.UserEmail, payment.Status, payment.ProductName, payment.VariantName)
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionRefunded(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	go processSubscription("refunded", payment.UserEmail, payment.Status, payment.ProductName, payment.VariantName)
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionChanged(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {

		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	go processSubscription("changed", payment.UserEmail, payment.Status, payment.ProductName, payment.VariantName)
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func processSubscription(mailType, mailId, status, productName, variantName string) {
	req := &subscription.SubscriptionRequest{
		MailType:           mailType,
		EmailId:            mailId,
		SubscriptionStatus: status,
		ProductName:        productName,
		VariantName:        variantName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := app.SubscriptionServiceClient.ProcessSubscription(ctx, req)

	if err != nil {
		fmt.Println(err)
		app.Producer.publishMessage("key", "Payment Service", "Failed to connect to the subscription service"+err.Error())
	}

	log.Printf("Response: %s", r)
}
