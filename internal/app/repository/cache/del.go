package cache

import "context"

func (db *redisCache) DelCacheData(ctx context.Context, groupKey string) error {
	return db.client.Del(ctx, groupKey).Err()
}
