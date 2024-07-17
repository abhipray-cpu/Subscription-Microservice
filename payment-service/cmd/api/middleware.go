package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func VerifySignatureMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Read the request body
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read request body")
		}
		// Since ioutil.ReadAll closes the body, we need to repopulate it for future handlers
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// The secret used for signing the requests
		secret := os.Getenv("WEBHOOK_SECRET")

		// Compute HMAC with SHA256 using the secret and the request body
		h := hmac.New(sha256.New, []byte(secret))
		h.Write(body)
		computedSignature := hex.EncodeToString(h.Sum(nil))

		// Retrieve the X-Signature header from the request
		signatureHeader := c.Request().Header.Get("X-Signature")

		// Convert the computed signature and the X-Signature header to byte slices for comparison
		computedSignatureBytes, err := hex.DecodeString(computedSignature)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode computed signature")
		}
		signatureBytes, err := hex.DecodeString(signatureHeader)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode X-Signature header")
		}

		// Use hmac.Equal for a timing-safe comparison
		if !hmac.Equal(computedSignatureBytes, signatureBytes) {
			// If not, log the attempt and return an unauthorized error
			fmt.Println("Unauthorized webhook request from IP: ", c.RealIP())
			// Assuming app.Producer.publishMessage is a method to log or notify about the unauthorized attempt
			app.Producer.publishMessage("key", "Payment Service", "Unauthorized webhook request from IP: "+c.RealIP())
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid signature")
		}

		// If the signature is valid, call the next handler in the chain
		return next(c)
	}
}
