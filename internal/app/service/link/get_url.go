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
	url, err := s.repo.GetURL(ctx, urlCode)
	if errors.Is(err, redis.Nil) {
		return "", ErrCodeNotFound
	}

	return url, err
}
