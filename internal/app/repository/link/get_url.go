package link

import "context"

// GetURL retrieves the original URL associated with the given ID from Redis.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//   - id: The unique identifier for the URL to be retrieved
//
// Returns:
//   - string: The original URL associated with the given ID
//   - error: An error object if the retrieval fails, otherwise nil
func (s *linkRepository) GetURL(ctx context.Context, id string) (string, error) {
	return s.redisClient.Get(ctx, id).Result()
}
