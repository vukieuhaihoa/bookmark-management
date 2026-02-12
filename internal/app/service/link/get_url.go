package link

import (
	"context"
	"errors"
	"strings"

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
	switch {
	case strings.HasPrefix(urlCode, "r"):
		url, err := s.repo.GetURL(ctx, urlCode)
		if errors.Is(err, redis.Nil) {
			return "", ErrCodeNotFound
		}

		return url, err
	case strings.HasPrefix(urlCode, "p"):
		bookmark, err := s.bookmarkRepo.GetBookmarkByCodeShortenEncoded(ctx, urlCode)
		if err != nil {
			return "", err
		}
		return bookmark.URL, nil
	default:
		// For backward compatibility, if the code doesn't have a prefix, check both repositories
		// Determine if the URL code corresponds to a Redis-stored link or a bookmark
		// Deprecated: This block is for backward compatibility and may be removed in future versions
		if len(urlCode) == defaultURLCodeLength {
			url, err := s.repo.GetURL(ctx, urlCode)
			if errors.Is(err, redis.Nil) {
				return "", ErrCodeNotFound
			}

			return url, err
		}

		// Deprecated: This block is for backward compatibility and may be removed in future versions
		bookmark, err := s.bookmarkRepo.GetBookmarkByCode(ctx, urlCode)
		if err != nil {
			return "", err
		}
		return bookmark.URL, nil
	}
}
