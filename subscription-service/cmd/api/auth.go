// NewAuth initializes the authentication configuration.
package main

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

const (
	key    = "random string"
	MaxAge = 86400 * 30
	IsProd = false
)

// NewAuth initializes the authentication configuration for the API.
// It loads the environment variables from the .env file, sets up the session store,
// and configures the GitHub authentication provider.
func NewAuth() {
	// Get GitHub API credentials from environment variables
	githubKey := os.Getenv("GITHUB_KEY")
	githubSecret := os.Getenv("GITHUB_SECRET")

	// Set up session store
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store

	// Configure GitHub authentication provider
	goth.UseProviders(
		github.New(githubKey, githubSecret, "http://localhost:80/auth/github/callback"),
	)
}
