package link

import (
	"context"
	"time"
)

// StoreURLIfAbsent attempts to store the given URL with the associated code in Redis only if the code does not already exist.
// It returns a boolean indicating whether the URL was stored (true if the code was absent and the URL was stored, false if the code already exists).
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//   - code: The unique code to associate with the URL
//   - url: The original URL to be stored
//   - expireIn: The expiration time in seconds for the stored URL
//
// Returns:
//   - bool: True if the URL was stored (code was absent), false if the code already exists
//   - error: An error object if the storage operation fails, otherwise nil
func (s *linkRepository) StoreURLIfAbsent(ctx context.Context, code, url string, expireIn int) (bool, error) {
	if expireIn <= 0 {
		expireIn = int(defaultURLExpireTime.Seconds())
	}
	return s.redisClient.SetNX(ctx, code, url, time.Duration(expireIn)*time.Second).Result()
}
