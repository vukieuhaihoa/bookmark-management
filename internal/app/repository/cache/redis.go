package cache

import "github.com/redis/go-redis/v9"

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) DB {
	return &redisCache{
		client: client,
	}
}
