package util

import (
	"crypto/rand"     // Provides cryptographic random number generation.
	"encoding/base64" // Implements base64 encoding for URL-safe encoding.
	"strconv"         // Provides conversions to and from string representations of basic data types.
	"time"            // Provides functionality for measuring and displaying time.

	"github.com/dgrijalva/jwt-go" // A library for working with JSON Web Tokens (JWT).
	"golang.org/x/crypto/bcrypt"  // Provides password hashing and verification using bcrypt.
)

// GenerateAccessToken generates a secure, random string of the specified length.
// This function is useful for creating tokens or unique identifiers.
//
// Parameters:
// - length: The desired length of the generated string.
//
// Returns:
// - A URL-safe, base64 encoded string of random bytes.
// - An error if the random byte generation fails.
func GenerateAccessToken(length int) (string, error) {
	// Generate a byte slice of the desired length of random bytes.
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		// Return an empty string and the error if there's an issue generating the random bytes.
		return "", err
	}

	// Encode the byte slice to a base64 string to make it URL-safe.
	token := base64.URLEncoding.EncodeToString(b)
	return token, nil
}

// HashPassword generates a bcrypt hash of the input password.
//
// Parameters:
// - password: The plaintext password to hash.
//
// Returns:
// - A hashed version of the input password.
// - An error if the hashing process fails.
func HashPassword(password string) (string, error) {
	// GenerateFromPassword returns a hashed password from the given password string.
	// The cost parameter controls the complexity of the hashing process.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// Return an empty string and the error if hashing fails.
		return "", err
	}
	// Return the hashed password as a string.
	return string(hashedPassword), nil
}

// ComparePasswords checks if a plaintext password matches a bcrypt hashed password.
//
// Parameters:
// - hashedPassword: The bcrypt hashed password.
// - password: The plaintext password to compare.
//
// Returns:
// - nil if the passwords match.
// - An error if the passwords don't match or if there's another error.
func ComparePasswords(hashedPassword, password string) error {
	// CompareHashAndPassword compares a bcrypt hashed password with its possible
	// plaintext equivalent. Returns nil on success, or an error on failure.
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// If the passwords don't match, or if there's another error, return the error.
		return err
	}
	// If the passwords match, return nil.
	return nil
}

// GenerateJWT creates a JWT (JSON Web Token) for a given user.
// This token can be used for authenticating API requests.
//
// Parameters:
// - userID: The user's ID.
// - userName: The user's name.
//
// Returns:
// - A signed JWT string.
// - An error if the JWT signing process fails.
func GenerateJWT(userID int64, userName string) (string, error) {
	var mySigningKey = []byte("secret")      // Use a secret from your environment.
	token := jwt.New(jwt.SigningMethodHS256) // Create a new JWT token using HS256 signing method.
	claims := token.Claims.(jwt.MapClaims)   // Cast the token's claims to a MapClaims object.

	// Set claims for the JWT. These claims include the user's ID, name, and an expiration time.
	claims["authorized"] = true
	claims["user_id"] = strconv.FormatInt(userID, 10) // Convert userID to string.
	claims["user_name"] = userName
	claims["exp"] = time.Now().Add(time.Hour * 720).Unix() // Set expiration to 30 days from now.

	// Sign the token using the specified secret key.
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
