package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"payment-service/data"

	"github.com/labstack/echo/v4"
)

func (app *Config) pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "system is working")
}

func (app *Config) SubscriptionCreated(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to read body")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	id, err := app.Models.CreatePayment(*payment)
	fmt.Println(id)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to create payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil

	}
	app.Producer.publishMessage("key", "Payment Service", "Subscription created successfully")
	return nil
}

func (app *Config) SubscriptionUpdated(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to read body")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to update payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	// run the temporal workflow
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionCancelled(c echo.Context) error {
	fmt.Println("this is working")
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to read body")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to update payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	// run the temporal workflow
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionResumed(c echo.Context) error {
	fmt.Println("this is working")
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to read body")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to update payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	// run the temporal workflow
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionExpired(c echo.Context) error {
	fmt.Println("this is working")
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to read body")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to update payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	// run the temporal workflow
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionPaused(c echo.Context) error {
	fmt.Println("this is working")
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to read body")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to update payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	// run the temporal workflow
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionUnpaused(c echo.Context) error {
	fmt.Println("this is working")
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to read body")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to update payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	// run the temporal workflow
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionFailedPayment(c echo.Context) error {
	fmt.Println("this is working")
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to read body")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to update payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	// run the temporal workflow
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionSucessPayment(c echo.Context) error {
	fmt.Println("this is working")
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to read body")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to update payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	// run the temporal workflow
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionRecovered(c echo.Context) error {
	fmt.Println("this is working")
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to read body")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to update payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	// run the temporal workflow
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionRefunded(c echo.Context) error {
	fmt.Println("this is working")
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to read body")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to update payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	// run the temporal workflow
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}

func (app *Config) SubscriptionChanged(c echo.Context) error {
	fmt.Println("this is working")
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to read body")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment, err := data.GetPayment(body)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	subscription_id := payment.SubscriptionID
	existingpayment, err := app.Models.GetPaymentBySubscriptionID(subscription_id)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to get payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	payment.ID = existingpayment.ID
	err = app.Models.UpdatePayment(*payment)
	if err != nil {
		// will also start a go routine to send a message to the subscription service
		fmt.Println("failed to update payment")
		app.Producer.publishMessage("key", "Payment Service", "Failed to create subscription"+err.Error())
		return nil
	}
	// run the temporal workflow
	app.Producer.publishMessage("key", "Payment Service", "Subscription updated successfully")
	return nil
}
