package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// HealthCheck defines the interface for performing health checks on the Redis server.
//
//go:generate mockery --name HealthCheck --filename health_check_repo.go --output ./mocks
type HealthCheck interface {
	// Ping checks the connectivity to the Redis server.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellations
	//
	// Returns:
	//   - error: An error object if the ping operation fails, otherwise nil
	Ping(ctx context.Context) error
}

type healthCheckStorage struct {
	redisClient *redis.Client
}

// NewHealthCheck creates a new instance of HealthCheck using the provided Redis client.
//
// Parameters:
//   - redisClient: The Redis client used for health check operations
//
// Returns:
//   - HealthCheck: A new HealthCheck instance
func NewHealthCheck(redisClient *redis.Client) HealthCheck {
	return &healthCheckStorage{
		redisClient: redisClient,
	}
}

// Ping checks the connectivity to the Redis server.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//
// Returns:
//   - error: An error object if the ping operation fails, otherwise nil
func (h *healthCheckStorage) Ping(ctx context.Context) error {
	return h.redisClient.Ping(ctx).Err()
}
