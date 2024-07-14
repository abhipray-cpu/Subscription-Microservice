package activity

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// StoreOTP stores the given OTP for the user in Redis with an expiration time.
func StoreOTP(ctx context.Context, rdb *redis.Client, userID, otp string, expiration time.Duration) error {
	// Construct the key using userID and OTP
	key := fmt.Sprintf("%s:%s", userID, otp)

	// Store the OTP in Redis with the specified expiration time
	err := rdb.Set(ctx, key, otp, expiration).Err()
	if err != nil {
		return fmt.Errorf("error storing OTP in Redis: %w", err)
	}

	return nil
}
