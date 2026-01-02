package main

import (
	"context"
	"time"

	"github.com/vukieuhaihoa/bookmark-management/pkg/redis"
)

func main() {
	ctx := context.Background()

	redisClient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}

	redisClient.Set(ctx, "key", "123", time.Hour)

	cacheDB, err := redis.NewClient("CACHE")
	if err != nil {
		panic(err)
	}

	cacheDB.Set(ctx, "key", "456", time.Hour)
}
