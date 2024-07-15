package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"subscription-service/data"
	"subscription-service/util"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"go.temporal.io/sdk/client"
)

// pingHandler handles the ping request and returns a success message.
func (app *Config) pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "The system is working fine")
}

// signup handles the user registration process.
func (app *Config) signup(c echo.Context) error {
	// Initialize a User struct to store the user's registration details.
	var user data.User

	// Bind the incoming JSON payload to the user struct.
	// This step parses the request body and maps the JSON fields to the struct fields.
	if err := c.Bind(&user); err != nil {
		// If binding fails, publish an error message and return a bad request response.
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to bind user data: "+err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Generate an access token for the user.
	token, err := util.GenerateAccessToken(32)
	if err != nil {
		// If token generation fails, publish an error message and return an internal server error response.
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to generate access token: "+err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Assign the generated token to the user's AccessToken field.
	user.AccessToken = token
	user.Verified = false
	// Hash the user's password for secure storage.
	hash, err := util.HashPassword(user.Password)
	if err != nil {
		// If password hashing fails, publish an error message and return an internal server error response.
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to hash password: "+err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Replace the plain text password with the hashed password.
	user.Password = hash

	// Validate the user's details.
	valid, reason := user.ValidateUser()
	if !valid {
		// If validation fails, return a bad request response with the reason for failure.
		errorMessage := fmt.Sprintf("Failed to create user:%s", reason)
		return c.JSON(http.StatusBadRequest, errorMessage)
	}

	// Insert the user into the database.
	if err := user.InsertUser(user); err != nil {
		// If insertion fails, publish an error message and return an internal server error response.
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to insert user: "+err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	go func() {
		type Param struct {
			To      string // Recipient email address
			Name    string // Recipient name
			Contact string // Recipient phone number
		}

		param := Param{
			To:      user.Email,
			Name:    user.UserName,
			Contact: "+91" + user.Contact,
		}
		// Prepare the workflow options
		workflowOptions := client.StartWorkflowOptions{
			ID:        "WelcomeWorkflow_" + user.Email, // Unique ID for the workflow instance
			TaskQueue: "subscription-service",          // The task queue name should match the one used in worker registration
		}
		_, err := app.Temporal.ExecuteWorkflow(context.Background(), workflowOptions, "WelcomeWorkflow", param)
		if err != nil {
			app.Producer.publishMessage("error", "Subscription-Service", "Failed to start WelcomeWorkflow: "+err.Error())
		}

	}()
	// Return a created response indicating successful account creation.
	return c.JSON(http.StatusCreated, "account created successfully")
}

// login handles user login requests.
func (app *Config) login(c echo.Context) error {
	// Define a struct to hold login details received from the request body.
	var loginDetails struct {
		Credentials string
		Password    string
	}

	// Bind the incoming JSON payload to the loginDetails struct.
	// This step parses the request body and maps the JSON fields to the struct fields.
	if err := c.Bind(&loginDetails); err != nil {
		// If binding fails, publish an error message and return a bad request response.
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to bind login details: "+err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// Extract email and password from the parsed login details.
	credential, password := loginDetails.Credentials, loginDetails.Password

	// switch between mobile and email number
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex := regexp.MustCompile(`^[789]\d{9}$`)
	var credential_type string
	switch {
	case emailRegex.MatchString(credential):
		credential_type = "email"
	case phoneRegex.MatchString(credential):
		credential_type = "phone"
	default:
		credential_type = "email"
	}

	var hashedPassword string
	// Initialize a User struct to hold the fetched user details.
	user := data.User{}

	// switch retreival function based on the credential type
	if credential_type == "email" {
		// Attempt to fetch the user by email from the database.
		if err := user.GetByEmail(credential); err != nil {
			// Check if the error is because the user does not exist in the database.
			if err == pgx.ErrNoRows {
				// If the user does not exist, respond with HTTP 404 Not Found.
				return c.JSON(http.StatusNotFound, "user does not exist")
			}
			// If there is another error, publish an error message and return an internal server error response.
			app.Producer.publishMessage("error", "Subscription-Service", "Failed to get user by email"+err.Error())
			return c.JSON(http.StatusInternalServerError, "Failed to fetch user")
		}
		hashedPassword = user.Password
	} else {
		// Attempt to fetch the user by contact from the database.
		if err := user.GetByContact("+91" + credential); err != nil {
			// Check if the error is because the user does not exist in the database.
			if err == pgx.ErrNoRows {
				// If the user does not exist, respond with HTTP 404 Not Found.
				return c.JSON(http.StatusNotFound, "user does not exist")
			}
			// If there is another error, publish an error message and return an internal server error response.
			app.Producer.publishMessage("error", "Subscription-Service", "Failed to get user by email"+err.Error())
			return c.JSON(http.StatusInternalServerError, "Failed to fetch user")
		}
		hashedPassword = user.Password
	}

	// Compare the provided password with the user's stored password.
	if err := util.ComparePasswords(hashedPassword, password); err != nil {
		// If the password comparison fails, publish an error message and return an unauthorized response.
		app.Producer.publishMessage("error", "Subscription-Service", "Invalid password: "+err.Error())
		return c.JSON(http.StatusUnauthorized, "Wrong password")
	}

	// Generate a JWT token for the authenticated user.
	token, err := util.GenerateJWT(int64(user.ID), user.GithubName)
	if err != nil {
		// If token generation fails, publish an error message and return an internal server error response.
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to generate JWT: "+err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Return a successful login response with the generated token and user details.
	return c.JSON(http.StatusOK, map[string]string{
		"token":           token,                      // The generated JWT token
		"user_name":       user.UserName,              // The user's username
		"github_username": user.GithubName,            // The user's GitHub username
		"message":         "Login successful",         // Success message
		"id":              fmt.Sprintf("%d", user.ID), // The user's ID, converted to a string
	})
}

// deleteAccount handles the deletion of a user's account.
func (app *Config) deleteAccount(c echo.Context) error {
	// Extract the userID from the context. This value is expected to be set by a previous middleware.
	userId := c.Get("userID").(int64)

	// Initialize a User struct. This struct will be used to call the DeleteUser method.
	var user data.User

	// Attempt to delete the user from the database using the DeleteUser method.
	if err := user.DeleteUser(userId); err != nil {
		// Check if the error is because the user does not exist in the database.
		if err == pgx.ErrNoRows {
			// If the user does not exist, respond with HTTP 404 Not Found.
			return c.JSON(http.StatusNotFound, "user does not exist")
		}
		// If there is another error, publish an error message indicating failure to delete the user account.
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to delete user account"+err.Error())
		// Respond with HTTP 500 Internal Server Error indicating failure to delete the account.
		return c.JSON(http.StatusInternalServerError, "Failed to delete account")
	}

	// If the deletion is successful, respond with HTTP 200 OK and a success message.
	return c.JSON(http.StatusOK, "account deleted successfully")
}

// getAccount handles the retrieval of a user's account details.
func (app *Config) getAccount(c echo.Context) error {
	// Extract the userID from the context, which is assumed to be set by a previous middleware.
	userId := c.Get("userID").(int64)

	// Initialize a User struct to hold the user details.
	var user data.User

	// Attempt to fetch the user details from the database using the userID.
	if err := user.GetUser(int64(userId)); err != nil {
		// Check if the error is because the user does not exist in the database.
		if err == pgx.ErrNoRows {
			// Respond with HTTP 404 Not Found if the user does not exist.
			return c.JSON(http.StatusNotFound, "user does not exist")
		}
		// Publish an error message indicating failure to fetch user details.
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to fetch user account"+err.Error())
		// Respond with HTTP 500 Internal Server Error indicating failure to fetch the user.
		return c.JSON(http.StatusInternalServerError, "Failed to fetch user")
	}

	// Respond with HTTP 200 OK and the user details in JSON format if the user is successfully retrieved.
	return c.JSON(http.StatusOK, map[string]string{
		"user_name": user.UserName,                         // User's username
		"github_id": user.GithubId,                         // User's GitHub ID
		"email":     user.Email,                            // User's email address
		"contact":   user.Contact,                          // User's contact information
		"bio":       user.Bio,                              // User's biography
		"avatar":    user.AvatarUrl,                        // URL to the user's avatar image
		"message":   "User details retrieved successfully", // Success message
	})
}

// updateAccount handles account updates.
func (app *Config) updateAccount(c echo.Context) error {
	// Define a struct to hold the new account details received from the request body.
	var newDetails struct {
		FirstName string
		LastName  string
		Email     string
		Contact   string
	}

	// Extract the userID from the context, which is assumed to be set by a previous middleware.
	userId := c.Get("userID").(int64)

	// Initialize a User struct to hold the updated user details.
	var user data.User

	// Bind the incoming JSON payload to the newDetails struct.
	// This step parses the request body and maps the JSON fields to the struct fields.
	if err := c.Bind(&newDetails); err != nil {
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to update user"+err.Error())
		return c.JSON(http.StatusInternalServerError, "Failed to update user details")
	}

	// Update the user struct with the new details received from the request.
	user.FirstName = newDetails.FirstName
	user.LastName = newDetails.LastName
	user.Email = newDetails.Email
	user.Contact = newDetails.Contact
	// Attempt to update the user in the database with the new details.
	if err := user.UpdateUser(userId, user); err != nil {
		// Check if the error is because the user does not exist in the database.
		if errors.Is(err, sql.ErrNoRows) {
			// Respond with HTTP 404 Not Found if the user is not found in the database.
			return c.JSON(http.StatusNotFound, "user not found")
		}
		// Publish an error message indicating failure to update the user in the database.
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to update user: "+err.Error())
		// Respond with HTTP 500 Internal Server Error indicating failure to update the user.
		return c.JSON(http.StatusInternalServerError, "failed to update user account")
	}

	// Respond with HTTP 200 OK on successful update of the user account.
	return c.JSON(http.StatusOK, "account updated successfully")
}

func (app *Config) GenerateOTP(c echo.Context) error {
	userId := c.Get("userID").(int64)
	var user data.User
	if err := user.GetUser(userId); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, "user does not exist")
		}
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to get user by email"+err.Error())
		return c.JSON(http.StatusInternalServerError, "Failed to fetch user")
	}
	go func() {
		type Param struct {
			To      string // Recipient email address
			Name    string // Recipient name
			Contact string // Recipient phone number
			UserID  string // User ID
		}

		param := Param{
			To:      user.Email,
			Name:    user.UserName,
			Contact: user.Contact,
			UserID:  fmt.Sprintf("%d", user.ID),
		}
		// Prepare the workflow options
		workflowOptions := client.StartWorkflowOptions{
			ID:        "OTPWorkflow" + user.Email, // Unique ID for the workflow instance
			TaskQueue: "subscription-service",     // The task queue name should match the one used in worker registration
		}
		_, err := app.Temporal.ExecuteWorkflow(context.Background(), workflowOptions, "OTPWorkflow", param)
		if err != nil {
			app.Producer.publishMessage("error", "Subscription-Service", "Failed to start OTPWorkflow: "+err.Error())
		}
	}()
	return c.JSON(http.StatusOK, "OTP sent successfully please check your email or message")
}

func (app *Config) VerifyOTP(c echo.Context) error {
	type Body struct {
		OTP string
	}
	var body Body
	var user data.User

	if err := c.Bind(&body); err != nil {
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to bind OTP: "+err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user_id := c.Get("userID").(int64)
	if err := user.GetUser(user_id); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, "user does not exist")
		}
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to get user by email"+err.Error())
		return c.JSON(http.StatusInternalServerError, "Failed to fetch user")
	}
	ctx := context.Background()
	key := fmt.Sprintf("%d:%s", user_id, body.OTP)
	exists, err := app.Redis.Exists(ctx, key).Result()
	if err != nil {
		app.Producer.publishMessage("error", "Subscription-Service", "Failed to verify OTP: "+err.Error())
		return c.JSON(http.StatusInternalServerError, "Failed to verify OTP")
	}

	if exists == 1 {
		otp, err := app.Redis.Get(ctx, key).Result()
		if err != nil {
			app.Producer.publishMessage("error", "Subscription-Service", "Failed to get OTP: "+err.Error())
			return c.JSON(http.StatusInternalServerError, "Failed to verify OTP")
		}
		if otp == body.OTP {
			app.Redis.Del(ctx, key)
			user.Verified = true
			if err := user.UpdateUser(user_id, user); err != nil {
				app.Producer.publishMessage("error", "Subscription-Service", "Failed to update user: "+err.Error())
				return c.JSON(http.StatusInternalServerError, "Failed to verify OTP")

			}
			return c.JSON(http.StatusOK, "OTP verified successfully")
		}
		return c.JSON(http.StatusBadRequest, "Invalid OTP")
	}
	return c.JSON(http.StatusBadRequest, "OTP expired")
}
