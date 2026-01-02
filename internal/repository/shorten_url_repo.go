package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultURLExpireTime = 24 * time.Hour
)

// UrlStorage defines the interface for URL storage operations.
// It provides methods for retrieving and storing URLs using a Redis backend.
type UrlStorage interface {
	// GetURL retrieves the original URL associated with the given ID.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellations
	//   - id: The unique identifier for the URL to be retrieved
	//
	// Returns:
	//   - string: The original URL associated with the given ID
	//   - error: An error object if the retrieval fails, otherwise nil
	GetURL(ctx context.Context, id string) (string, error)

	// StoreURL stores the given URL with the associated code.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellations
	//   - code: The unique code to associate with the URL
	//   - url: The original URL to be stored
	//   - expireIn: The expiration time in seconds for the stored URL
	//
	// Returns:
	//   - error: An error object if the storage operation fails, otherwise nil
	StoreURL(ctx context.Context, code, url string, expireIn int) error
}

// urlStorage is the concrete implementation of UrlStorage interface.
// It uses a Redis client to perform URL storage and retrieval operations.
type urlStorage struct {
	redisClient *redis.Client
}

// NewUrlStorage creates a new instance of UrlStorage using the provided Redis client.
//
// Parameters:
//   - redisClient: The Redis client used for URL storage operations
//
// Returns:
//   - UrlStorage: A new UrlStorage instance
func NewUrlStorage(redisClient *redis.Client) UrlStorage {
	return &urlStorage{
		redisClient: redisClient,
	}
}

// GetURL retrieves the original URL associated with the given ID from Redis.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//   - id: The unique identifier for the URL to be retrieved
//
// Returns:
//   - string: The original URL associated with the given ID
//   - error: An error object if the retrieval fails, otherwise nil
func (s *urlStorage) GetURL(ctx context.Context, id string) (string, error) {
	return s.redisClient.Get(ctx, id).Result()
}

// StoreURL stores the given URL with the associated code in Redis with a default expiration time.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//   - code: The unique code to associate with the URL
//   - url: The original URL to be stored
//   - expireIn: The expiration time in seconds for the stored URL
//
// Returns:
//   - error: An error object if the storage operation fails, otherwise nil
func (s *urlStorage) StoreURL(ctx context.Context, code, url string, expireIn int) error {
	expiration := defaultURLExpireTime
	if expireIn > 0 {
		expiration = time.Duration(expireIn) * time.Second
	}
	return s.redisClient.Set(ctx, code, url, expiration).Err()
}
