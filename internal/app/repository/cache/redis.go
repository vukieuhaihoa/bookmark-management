package cache

import "github.com/redis/go-redis/v9"

type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new instance of Redis cache implementing the DB interface.
//
// Parameters:
//   - client: The Redis client to be used for cache operations.
//
// Returns:
//   - DB: An instance of the cache database interface.
func NewRedisCache(client *redis.Client) DB {
	return &redisCache{
		client: client,
	}
}
