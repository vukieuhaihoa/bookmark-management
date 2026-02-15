package link

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
)

// ShortenURL generates a shortened URL code for the given original URL.
// It retries code generation if a collision is detected, up to maxRetryAttempts.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//   - originalURL: The original URL to be shortened
//   - expireIn: The expiration time in seconds for the shortened URL
//
// Returns:
//   - string: The generated shortened URL code
//   - error: An error object if the shortening operation fails, otherwise nil
func (s *linkService) ShortenURL(ctx context.Context, originalURL string, expireIn int) (string, error) {
	var urlCode string
	var err error

	for attempt := 1; attempt <= maxRetryAttempts; attempt++ {
		urlCode, err = s.randomCodeGen.GenerateCode(defaultURLCodeLength)
		if err != nil {
			return "", err
		}

		urlCode = model.RedisShortenPrefix + urlCode

		// Check if the URL code already exists
		_, err = s.repo.GetURL(ctx, urlCode)
		if err != nil {
			// If the error is redis.Nil, the key doesn't exist - this is what we want
			if errors.Is(err, redis.Nil) {
				// URL code is unique, proceed to store it
				break
			}
			// For any other error, return it
			return "", err
		}

		// URL code already exists, retry if we haven't exhausted attempts
		if attempt == maxRetryAttempts {
			return "", ErrMaxRetriesExceeded
		}
	}

	// Store the URL with the unique code
	err = s.repo.StoreURL(ctx, urlCode, originalURL, expireIn)
	if err != nil {
		return "", err
	}

	return urlCode, nil
}
