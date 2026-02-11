package cache

import (
	"context"
	"time"
)

// DB defines the interface for cache database operations.
// It includes methods for setting, getting, and deleting cache data.
//
//go:generate mockery --name=DB --output=./mocks --filename=cache.go --outpkg=mock_cache
type DB interface {
	// SetCacheData sets the cache data for a given group key and key with an expiration time.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - groupKey: The group key under which the cache data is stored.
	//   - key: The specific key for the cache data.
	//   - value: The value to be cached.
	//   - exp: The expiration duration for the cached data.
	//
	// Returns:
	//   - error: An error if the operation fails, otherwise nil.
	SetCacheData(ctx context.Context, groupKey, key string, value interface{}, exp time.Duration) error

	// GetCacheData retrieves the cache data for a given group key and key.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - groupKey: The group key under which the cache data is stored.
	//   - key: The specific key for the cache data.
	//
	// Returns:
	//   - []byte: The cached data as a byte slice.
	//   - error: An error if the operation fails, otherwise nil.
	GetCacheData(ctx context.Context, groupKey, key string) ([]byte, error)

	// DelCacheData deletes the cache data for a given group key.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - groupKey: The group key under which the cache data is stored.
	//
	// Returns:
	//   - error: An error if the operation fails, otherwise nil.
	DelCacheData(ctx context.Context, groupKey string) error
}
