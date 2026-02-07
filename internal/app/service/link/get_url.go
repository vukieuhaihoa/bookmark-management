package link

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

// GetURL retrieves the original URL associated with the given shortened URL code.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//   - urlCode: The shortened URL code
//
// Returns:
//   - string: The original URL associated with the shortened code
//   - error: An error object if the retrieval operation fails, otherwise nil
func (s *linkService) GetURL(ctx context.Context, urlCode string) (string, error) {
	// Determine if the URL code corresponds to a Redis-stored link or a bookmark
	if len(urlCode) == defaultURLCodeLength {
		url, err := s.repo.GetURL(ctx, urlCode)
		if errors.Is(err, redis.Nil) {
			return "", ErrCodeNotFound
		}

		return url, err
	}

	bookmark, err := s.bookmarkRepo.GetBookmarkByCode(ctx, urlCode)
	if err != nil {
		return "", err
	}

	return bookmark.URL, nil
}
