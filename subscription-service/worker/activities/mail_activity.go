package activity

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

// sendEmail encapsulates the functionality to send an email using AWS SES.
func sendEmail(ctx context.Context, svc *ses.SES, to, subject, htmlBody string) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(to)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(htmlBody),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String("dumkaabhipray@gmail.com"), // Replace with your sender email
	}

	// Send the email using the provided SES service with the standard context
	_, err := svc.SendEmailWithContext(ctx, input)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}

// SendWelcomeEmail sends a welcome email to a new user.
func (ac *ActivitiesImpl) SendWelcomeEmail(ctx context.Context, to, name string) error {
	subject := "Welcome to Our Service!"
	htmlBody := fmt.Sprintf(`<html>
<head>
<style>
body {font-family: 'Arial', sans-serif; background-color: #f0f0f0; margin: 0; padding: 20px;}
.container {background-color: #ffffff; padding: 20px; max-width: 600px; margin: auto; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1);}
h1 {color: #333366;}
p {color: #666666;}
.button {background-color: #4CAF50; color: white; padding: 14px 20px; text-align: center; display: inline-block; font-size: 16px; margin: 4px 2px; cursor: pointer; border-radius: 5px; text-decoration: none;}
</style>
</head>
<body>
<div class="container">
<h1>Welcome, %s!</h1>
<p>We're excited to have you on board. Click the button below to get started with our service.</p>
<a href="https://yourwebsite.com/get-started" class="button">Get Started</a>
</div>
</body>
</html>`, name)
	return sendEmail(ctx, ac.sesClient, to, subject, htmlBody)
}

// sendOTPEmail sends an OTP email for account verification.
func (ac *ActivitiesImpl) SendOTPEmail(ctx context.Context, to, otpCode string) error {
	subject := "Your OTP Code"
	htmlBody := fmt.Sprintf(`<html>
<head>
<style>
body {font-family: 'Arial', sans-serif; background-color: #f0f0f0; margin: 0; padding: 20px;}
.container {background-color: #ffffff; padding: 20px; max-width: 600px; margin: auto; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1);}
h1 {color: #333366;}
p {color: #666666;}
.code {font-size: 24px; font-weight: bold; color: #333366;}
</style>
</head>
<body>
<div class="container">
<h1>Your OTP Code</h1>
<p>Please use the following code to verify your account:</p>
<p class="code">%s</p>
<p>If you did not request this code, please ignore this email.</p>
</div>
</body>
</html>`, otpCode)
	return sendEmail(ctx, ac.sesClient, to, subject, htmlBody)
}
