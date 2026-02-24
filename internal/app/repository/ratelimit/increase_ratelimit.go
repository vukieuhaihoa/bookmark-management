package ratelimit

import (
	"context"
	"time"
)

// IncreaseRateLimit increments the rate limit counter for a given key and sets an expiration time.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - key: The specific key for which the rate limit counter should be incremented.
//   - exp: The expiration duration for the rate limit counter.
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
func (r *redisRepo) IncreaseRateLimit(ctx context.Context, key string, exp time.Duration) error {
	// Increment the rate limit counter for the given key
	err := r.c.Incr(ctx, key).Err()
	if err != nil {
		return err
	}

	// Set the expiration time for the rate limit counter
	err = r.c.ExpireNX(ctx, key, exp).Err()
	if err != nil {
		return err
	}

	return nil
}
