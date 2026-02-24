// Package link provides functionality for URL shortening and retrieval.
// It defines the UrlStorage interface and its implementation.
// The UrlStorage interface includes methods for storing and retrieving URLs using a Redis backend.
// This package is essential for managing shortened URLs in the application.
package link

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultURLExpireTime = 24 * time.Hour
)

// Repository defines the interface for URL storage operations.
// It provides methods for retrieving and storing URLs using a Redis backend.
//
//go:generate mockery --name Repository --filename shorten_url_repo.go --output ./mocks
type Repository interface {
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

	// StoreURLIfAbsent stores the given URL with the associated code only if the code does not already exist.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellations
	//   - code: The unique code to associate with the URL
	//   - url: The original URL to be stored
	//   - expireIn: The expiration time in seconds for the stored URL
	//
	// Returns:
	//   - bool: True if the URL was stored successfully, false if the code already exists
	//   - error: An error object if the storage operation fails, otherwise nil
	StoreURLIfAbsent(ctx context.Context, code, url string, expireIn int) (bool, error)
}

// linkRepository is the concrete implementation of Repository interface.
// It uses a Redis client to perform URL storage and retrieval operations.
type linkRepository struct {
	redisClient *redis.Client
}

// NewLinkRepository creates a new instance of linkRepository using the provided Redis client.
//
// Parameters:
//   - redisClient: The Redis client used for URL storage operations
//
// Returns:
//   - Repository: A new Repository instance
func NewLinkRepository(redisClient *redis.Client) Repository {
	return &linkRepository{
		redisClient: redisClient,
	}
}
