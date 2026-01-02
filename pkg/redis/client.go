package redis

import "github.com/redis/go-redis/v9"

// NewClient creates and returns a new Redis client based on the provided environment prefix.
//
// Parameters:
//   - envPrefix: A string prefix used to load environment variables for Redis configuration
//
// Returns:
//   - *redis.Client: A pointer to the newly created Redis client
//   - error: An error object if the client creation fails, otherwise nil
func NewClient(envPrefix string) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return redisClient, nil
}
