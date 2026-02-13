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
	case len(urlCode) == defaultURLCodeLength ||
		(strings.HasPrefix(urlCode, "r") && len(urlCode) == defaultURLCodeLength+1):
		// Redis (old 8-char + new r-prefixed 9-char)
		url, err := s.repo.GetURL(ctx, urlCode)
		if errors.Is(err, redis.Nil) {
			return "", ErrCodeNotFound
		}

		return url, err
	case len(urlCode) == 10:
		// backward compatibility for old url code length of 10
		bookmark, err := s.bookmarkRepo.GetBookmarkByCode(ctx, urlCode)
		if err != nil {
			return "", err
		}
		return bookmark.URL, nil
	case strings.HasPrefix(urlCode, "p"):
		bookmark, err := s.bookmarkRepo.GetBookmarkByCodeShortenEncoded(ctx, urlCode)
		if err != nil {
			return "", err
		}
		return bookmark.URL, nil
	default:
		return "", ErrCodeNotFound
	}
}
