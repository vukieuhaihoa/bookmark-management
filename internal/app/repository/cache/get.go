package cache

import "context"

// GetCacheData retrieves the cache data for a given group key and key.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - groupKey: The group key under which the cache data is stored.
//   - key: The specific key for the cache data.
//
// Returns:
//   - []byte: The cached data as a byte slice.
//   - error: An error if the operation fails, otherwise nil.
func (db *redisCache) GetCacheData(ctx context.Context, groupKey, key string) ([]byte, error) {
	result, err := db.client.HGet(ctx, groupKey, key).Bytes()
	if err != nil {
		return nil, err
	}

	return result, nil
}
