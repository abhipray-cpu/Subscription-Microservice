package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	// Other imports...
)

// JWTAuthMiddleware creates a middleware for JWT authentication.
// This middleware function is designed to be used with the Echo framework to secure endpoints by validating JWT tokens.
//
// Parameters:
// - next: A function of type echo.HandlerFunc. This represents the next handler in the middleware chain.
//
// Returns:
//   - A function of type echo.HandlerFunc. This function takes an Echo context and returns an error.
//     It performs the JWT validation and, upon success, calls the next handler in the chain.
//
// The middleware performs the following steps:
// 1. Extracts the Authorization header from the incoming HTTP request.
//   - If the Authorization header is missing, it returns an HTTP 401 Unauthorized error.
//
// 2. Strips the "Bearer " prefix from the Authorization header to isolate the JWT token.
// 3. Parses the JWT token using the jwt-go library.
//   - The parsing function requires a callback to provide the signing key. Here, a hardcoded "secret" is used.
//   - The callback also validates that the token's signing method matches the expected HMAC signing method.
//     If it doesn't, an HTTP 401 Unauthorized error is returned.
//
// 4. If the token parsing fails (due to being invalid or expired), an HTTP 401 Unauthorized error is returned.
// 5. If the token is successfully parsed, it extracts the "user_id" claim from the token's payload.
//   - If the "user_id" claim is missing or not a string, an HTTP 401 Unauthorized error is returned.
//   - If the "user_id" claim is present but its format is invalid (not an integer), an HTTP 401 Unauthorized error is returned.
//     6. If the "user_id" claim is valid, it is added to the Echo context using c.Set("userID", userID).
//     This allows downstream handlers to access the authenticated user's ID.
//     7. Finally, if the JWT is valid and the "user_id" claim is processed successfully, the next handler in the middleware chain is called.
//
// This middleware is crucial for securing routes that require user authentication. It ensures that only requests with a valid JWT,
// which signifies an authenticated user, can access certain endpoints.
func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the Authorization header from the request.
		authHeader := c.Request().Header.Get("Authorization")
		// Check if the Authorization header is missing.
		if authHeader == "" {
			// Return an HTTP 401 Unauthorized error if the header is missing.
			return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
		}

		// Remove the "Bearer " prefix from the Authorization header to get the JWT token.
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// Parse the JWT token.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the token's signing method.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// Return an error if the signing method is not HMAC.
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}
			// Provide the signing key ("secret") for token verification.
			return []byte("secret"), nil
		})

		// Handle parsing errors (invalid or expired token).
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
		}

		// Check if the token's claims are valid.
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract the "user_id" claim as a string.
			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				// Return an error if the "user_id" claim is missing or not a string.
				return echo.NewHTTPError(http.StatusUnauthorized, "user_id claim must be a string")
			}

			// Convert the "user_id" string to an integer.
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				// Return an error if the "user_id" format is invalid.
				return echo.NewHTTPError(http.StatusUnauthorized, "user_id format is invalid")
			}

			// Add the user ID to the Echo context for use in downstream handlers.
			c.Set("userID", userID)
			// Call the next handler in the middleware chain.
			return next(c)
		} else {
			// Return an error if the token is invalid.
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}
	}
}
