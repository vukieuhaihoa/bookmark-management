package cache

import (
	"context"
	"time"
)

func (db *redisCache) SetCacheData(ctx context.Context, groupKey, key string, value interface{}, exp time.Duration) error {
	err := db.client.HSet(ctx, groupKey, key, value).Err()
	if err != nil {
		return err
	}

	return db.client.Expire(ctx, groupKey, exp).Err()
}
