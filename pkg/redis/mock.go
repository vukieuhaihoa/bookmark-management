package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// InitMockRedis initializes and returns a mock Redis client for testing purposes.
//
// Parameters:
//   - t: The testing object used for managing test lifecycle
//
// Returns:
//   - *redis.Client: A mock Redis client connected to an in-memory Redis server
func InitMockRedis(t *testing.T) *redis.Client {
	mock := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{
		Addr: mock.Addr(),
	})
}
