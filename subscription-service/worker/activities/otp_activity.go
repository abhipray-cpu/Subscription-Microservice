package activity

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// GenerateOTP generates a unique 6-digit OTP, checks Redis to ensure it's unique, and returns it.
// The Redis client is passed as a parameter.
func (ac *ActivitiesImpl) GenerateOTP(ctx context.Context, userID string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000)) // Generate a 6-digit OTP

	// Construct the key with userID and OTP for Redis check
	key := fmt.Sprintf("%s:%s", userID, otp)

	// Check if the OTP already exists in Redis
	exists, err := ac.redis.Exists(ctx, key).Result()
	if err != nil {
		return "", err // Return error if any during Redis check
	}

	if exists > 0 {
		// If OTP already exists, recursively call GenerateOTP to generate a new one
		return ac.GenerateOTP(ctx, userID)
	} else {
		// Store the OTP in Redis with an expiration time, e.g., 5 minutes
		err := ac.redis.Set(ctx, key, otp, 60*time.Minute).Err()
		if err != nil {
			return "", err // Return error if any during Redis set
		}
		// If OTP is unique and stored successfully, return it
		return otp, nil
	}
}
