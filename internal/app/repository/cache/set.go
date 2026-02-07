package cache

import (
	"context"
	"time"
)

// SetCacheData sets the cache data for a given group key and key with an expiration time.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - groupKey: The group key under which the cache data is stored.
//   - key: The specific key for the cache data.
//   - value: The value to be cached.
//   - exp: The expiration duration for the cached data.
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
func (db *redisCache) SetCacheData(ctx context.Context, groupKey, key string, value interface{}, exp time.Duration) error {
	err := db.client.HSet(ctx, groupKey, key, value).Err()
	if err != nil {
		return err
	}

	return db.client.Expire(ctx, groupKey, exp).Err()
}
