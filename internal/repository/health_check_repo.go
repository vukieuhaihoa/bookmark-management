package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
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

	// DBPing checks the connectivity to the database.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellations
	//
	// Returns:
	//   - error: An error object if the ping operation fails, otherwise nil
	DBPing(ctx context.Context) error
}

type healthCheckStorage struct {
	redisClient *redis.Client
	db          *gorm.DB
}

// NewHealthCheck creates a new instance of HealthCheck using the provided Redis client and Gorm DB.
//
// Parameters:
//   - redisClient: The Redis client used for health check operations
//   - db: The Gorm DB used for health check operations
//
// Returns:
//   - HealthCheck: A new HealthCheck instance
func NewHealthCheck(redisClient *redis.Client, db *gorm.DB) HealthCheck {
	return &healthCheckStorage{
		redisClient: redisClient,
		db:          db,
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

// DBPing checks the connectivity to the database.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//
// Returns:
//   - error: An error object if the ping operation fails, otherwise nil
func (h *healthCheckStorage) DBPing(ctx context.Context) error {
	return h.db.WithContext(ctx).Exec("SELECT 1").Error
}
