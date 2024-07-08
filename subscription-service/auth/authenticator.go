package auth

import "github.com/labstack/echo/v4"

// Authenticator is an interface that defines the methods required for an authentication provider.
type Authenticator interface {
	NewAuth()                      // NewAuth initializes the authentication configuration.
	CallBack(c echo.Context) error // CallBack handles the callback request from the authentication provider.
	Logout(c echo.Context) error   // Logout handles the logout request from the authentication provider.
	Auth(c echo.Context) error     // Auth handles the authentication request from the client.
}
