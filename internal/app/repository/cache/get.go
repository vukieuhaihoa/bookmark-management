package cache

import "context"

func (db *redisCache) GetCacheData(ctx context.Context, groupKey, key string) ([]byte, error) {
	result, err := db.client.HGet(ctx, groupKey, key).Bytes()
	if err != nil {
		return nil, err
	}

	return result, nil
}
