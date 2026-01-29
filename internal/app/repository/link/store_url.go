package link

import (
	"context"
	"time"
)

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
func (s *linkRepository) StoreURL(ctx context.Context, code, url string, expireIn int) error {
	if expireIn <= 0 {
		expireIn = int(defaultURLExpireTime)
	}
	return s.redisClient.Set(ctx, code, url, time.Duration(expireIn)*time.Second).Err()
}
