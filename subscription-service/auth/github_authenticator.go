package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"subscription-service/data"
	"subscription-service/util"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

const (
	key    = "random string"
	MaxAge = 86400 * 30
	IsProd = false
)

type GitHubAuthenticator struct{}

// NewAuth configures the GitHub authentication mechanism for the application.
//
// This method performs the following steps:
// 1. Retrieves GitHub OAuth credentials (key and secret) from environment variables.
// 2. Sets up a session store using Redis for managing session data.
//   - It attempts to create a new Redis store with a predefined configuration.
//   - If the Redis store creation fails, the application logs the error and terminates.
//
// 3. Configures the session store with specific options such as MaxAge, Path, HttpOnly, and Secure flags.
//   - MaxAge controls the lifetime of the session cookie.
//   - Path sets the URL path where the cookie is valid.
//   - HttpOnly flag, when set to true, prevents client-side scripts from accessing the cookie, enhancing security.
//   - Secure flag is set based on the IsProd variable, which likely indicates whether the application is running in a production environment. This flag ensures cookies are sent over HTTPS only.
//
// 4. Assigns the configured session store to the gothic library, which handles OAuth flows.
// 5. Configures the GitHub authentication provider with the retrieved credentials and the callback URL.
//   - The callback URL is where GitHub redirects the user after authentication. It must match the URL configured in the GitHub OAuth application settings.
//
// This method is essential for initializing the GitHub OAuth authentication process, enabling users to log in with their GitHub accounts. It leverages the goth library to abstract away the complexities of OAuth and session management.
func (g *GitHubAuthenticator) NewAuth() {
	// Get GitHub API credentials from environment variables
	githubKey := os.Getenv("GITHUB_KEY")
	githubSecret := os.Getenv("GITHUB_SECRET")

	// Set up session store
	store, err := data.NewRedisStore(10, "tcp", "redis:6379", "", []byte(key))
	if err != nil {
		log.Fatalf("Failed to create Redis store: %v", err)
	}

	// Configure the store as needed
	store.Options.MaxAge = MaxAge
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store

	// Configure GitHub authentication provider
	goth.UseProviders(
		github.New(githubKey, githubSecret, "http://localhost:80/auth/github/callback"),
	)
}

// CallBack handles the callback from GitHub OAuth after the user has authenticated.
// It uses the Echo framework's context to manage HTTP requests and responses.
//
// Parameters:
// - c: The Echo context, which provides methods for interacting with the HTTP request and response.
//
// Returns:
// - An error if the authentication process fails at any point, otherwise nil.
//
// The function performs the following steps:
// 1. Extracts the provider name from the request parameters.
// 2. Creates a new context from the original HTTP request, adding the provider information to it.
// 3. Attempts to complete the user authentication with GitHub using the gothic library.
//   - If authentication fails, logs the error, sends an HTTP 500 response, and returns an error.
//
// 4. Checks if the authenticated user already exists in the database by their GitHub ID.
//   - If the user does not exist (indicated by pgx.ErrNoRows), a new user record is created with the information
//     obtained from GitHub and inserted into the database.
//   - If the user creation fails, sends an HTTP 500 response and returns an error.
//   - If the user is successfully created, sends an HTTP 200 response indicating success.
//
// 5. For existing users, generates a JWT token for session management.
//   - If token generation fails, logs the error, sends an HTTP 500 response, and returns an error.
//
// 6. Sends an HTTP 200 response with the JWT token and user information if the user exists or is successfully created.
//
// This function is crucial for handling the OAuth callback from GitHub, managing user authentication,
// and ensuring that user records are properly managed in the application's database.
func (g *GitHubAuthenticator) CallBack(c echo.Context) error {
	// Extract the provider parameter from the request URL.
	provider := c.Param("provider")
	// Create a new context with the provider information added to the original request's context.
	req := c.Request().WithContext(context.WithValue(c.Request().Context(), "provider", provider))

	// Attempt to complete the user authentication process using the modified request.
	user, err := gothic.CompleteUserAuth(c.Response().Writer, req)
	if err != nil {
		// Log and return an error response if authentication fails.
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Initialize a variable to hold user data.
	var User data.User
	// Check if the user exists in the database by their GitHub ID.
	if err := User.GetByGitId(user.UserID); err != nil {
		// If the user does not exist, create a new user record.
		if err == pgx.ErrNoRows {
			// Populate the User struct with data from the authenticated GitHub user.
			User.AccessToken = user.AccessToken
			User.Email = user.Email
			User.GithubId = user.UserID
			User.UserName = user.Name
			User.Bio = user.Description
			User.AvatarUrl = user.AvatarURL
			User.FirstName = user.FirstName
			User.LastName = user.LastName
			// Attempt to insert the new user into the database.
			if err := User.InsertUser(User); err != nil {
				// Return an error response if user creation fails.
				return c.JSON(http.StatusInternalServerError, "error while creating user")
			} else {
				// Return a success response if the user is created successfully.
				return c.JSON(http.StatusOK, "Account created successfully!")
			}
		}
	}

	// For existing users, generate a JWT token for session management.
	token, err := util.GenerateJWT(int64(User.ID), User.GithubName)
	if err != nil {
		// Log and return an error response if token generation fails.
		log.Println("failed to generate JWT: ", err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Return a success response with the JWT token and user information.
	return c.JSON(http.StatusOK, map[string]string{
		"token":           token,
		"user_name":       User.UserName,
		"github_username": User.GithubName,
		"message":         "Login successful",
		"id":              fmt.Sprintf("%d", User.ID),
	})
}

// Logout handles the logout request from the authentication provider.
// This method logs the user out of the system by clearing the session data.
// Parameters:
// - c: The echo context containing the request and response objects.
// Returns:
// - An error response if the logout fails.
// - A success response if the logout is successful.
func (g *GitHubAuthenticator) Logout(c echo.Context) error {
	provider := c.Param("provider")
	req := c.Request().WithContext(context.WithValue(c.Request().Context(), "provider", provider))
	res := c.Response().Writer
	gothic.Logout(res, req)
	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusTemporaryRedirect)
	return c.JSON(http.StatusOK, "logout successfully!")
}

// Auth handles the authentication request from the client.
// This method initiates the authentication process with the specified provider.
// Parameters:
// - c: The echo context containing the request and response objects.
// Returns:
// - An error response if the authentication fails.
// - A success response if the authentication is successful.
func (g *GitHubAuthenticator) Auth(c echo.Context) error {
	provider := c.Param("provider")
	req := c.Request().WithContext(context.WithValue(c.Request().Context(), "provider", provider))
	res := c.Response().Writer
	if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
		return c.JSON(http.StatusOK, gothUser)
	} else {
		gothic.BeginAuthHandler(res, req)
	}
	return c.JSON(http.StatusUnauthorized, "Authentication failed")
}

// NewGitHubAuthenticator creates a new GitHubAuthenticator instance.
// Returns a pointer to the instance.
func NewGitHubAuthenticator() *GitHubAuthenticator {
	return &GitHubAuthenticator{}
}
