package cache

import "context"

// DelCacheData deletes the cache data for a given group key.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - groupKey: The group key under which the cache data is stored.
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
func (db *redisCache) DelCacheData(ctx context.Context, groupKey string) error {
	return db.client.Del(ctx, groupKey).Err()
}
