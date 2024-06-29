package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

// pingHandler handles the ping request and returns a success message.
func (app *Config) pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "The system is working fine")
}

// callBackHandler handles the callback request from the authentication provider.
// It completes the user authentication process and returns the user information.
func (app *Config) callBackHandler(c echo.Context) error {
    provider := c.Param("provider")
    // Wrap the original request in a new context with the provider information
    req := c.Request().WithContext(context.WithValue(c.Request().Context(), "provider", provider))

    // Attempt to complete the user authentication process using the modified request
    user, err := gothic.CompleteUserAuth(c.Response().Writer, req)
    if err != nil {
        // Directly return an error response using Echo's context method
		fmt.Println(err.Error())
        return c.String(http.StatusInternalServerError, err.Error())
    }
    // Directly return the user information as JSON using Echo's context method
	fmt.Println(user)
    return c.JSON(http.StatusOK, user)
}

// logoutHandler handles the logout request from the authentication provider.
// It logs out the user and redirects to the home page.
func (app *Config) logoutHandler(c echo.Context) error {
	provider := c.Param("provider")
	req := c.Request().WithContext(context.WithValue(c.Request().Context(), "provider", provider))
	res := c.Response().Writer
	gothic.Logout(res, req)
	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusTemporaryRedirect)
	return c.JSON(http.StatusOK, "logout successfully!")
}

// authHandler handles the authentication request from the client.
// If the user is already authenticated, it returns the user information.
// Otherwise, it starts the authentication process.
func (app *Config) authHandler(c echo.Context) error {
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
