package activity

import (
	"fmt"
	"os"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func (ac *ActivitiesImpl) SendWelcomeSMS(to, name string) error {
	message := fmt.Sprintf("🌟 Welcome to Our Service, %s! 🌟\nWe're thrilled to have you on board. Stay tuned for updates.", name)
	return sendSMS(ac.twilioClient, to, message)
}

func (ac *ActivitiesImpl) SendOTPSMS(to string, otpCode string) error {
	// Making the OTP stand out by using symbols and spacing
	message := fmt.Sprintf("🔐 Your OTP is:\n\n🌟 %s 🌟\n\nUse this to complete your verification.", otpCode)
	return sendSMS(ac.twilioClient, to, message)
}

func (ac *ActivitiesImpl) SendSubscriptionUpdateSMS(to, subscriptionName, status string) error {
	message := fmt.Sprintf("📢 Update on your subscription: %s 📢\nYour subscription '%s' is now %s.", subscriptionName, subscriptionName, status)
	return sendSMS(ac.twilioClient, to, message)
}

func sendSMS(client *twilio.RestClient, to string, message string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(message)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return err
	}
	fmt.Println("SMS sent successfully to", to)
	return nil
}
