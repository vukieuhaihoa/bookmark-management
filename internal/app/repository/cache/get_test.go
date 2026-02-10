package cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
)

func TestDB_GetCacheData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *redis.Client

		expectedError error

		verifyFunc func(ctx context.Context, redisClient *redis.Client)
	}{
		{
			name: "successful cache get",

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				redisClient.HSet(ctx, "cached_group_key", "cached_key", "cached_value")
				redisClient.Expire(ctx, "cached_group_key", time.Hour)
				return redisClient
			},

			verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
				value, err := redisClient.HGet(ctx, "cached_group_key", "cached_key").Result()

				assert.Nil(t, err)

				assert.Equal(t, "cached_value", value)
			},
		},
		{
			name: "failed cache get due to closed Redis client",

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				redisClient.Close()
				return redisClient
			},

			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMockClient := tc.setupMock(ctx)

			cacheDB := NewRedisCache(redisMockClient)

			_, err := cacheDB.GetCacheData(ctx, "cached_group_key", "cached_key")
			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				tc.verifyFunc(ctx, redisMockClient)
			}

		})
	}
}
