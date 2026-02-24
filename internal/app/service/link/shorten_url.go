package link

import (
	"context"

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
	for attempt := 1; attempt <= maxRetryAttempts; attempt++ {
		urlCode, err := s.randomCodeGen.GenerateCode(defaultURLCodeLength)
		if err != nil {
			return "", err
		}

		urlCode = model.RedisShortenPrefix + urlCode

		stored, err := s.repo.StoreURLIfAbsent(ctx, urlCode, originalURL, expireIn)
		if err != nil {
			return "", err
		}
		if stored {
			return urlCode, nil // atomically claimed
		}
		// collision — another request already holds this key, retry
	}

	return "", ErrMaxRetriesExceeded
}
