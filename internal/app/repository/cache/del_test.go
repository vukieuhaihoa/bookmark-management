package cache

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
)

func TestDB_DelCacheData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func() *redis.Client

		expectedError error

		verifyFunc func(ctx context.Context, redisClient *redis.Client)
	}{
		{
			name: "successful cache data deletion",

			setupMock: func() *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				redisClient.Set(t.Context(), "groupKey", "someValue", 10000)
				return redisClient
			},

			expectedError: nil,

			verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
				_, err := redisClient.Get(ctx, "groupKey").Result()

				assert.Equal(t, redis.Nil, err)
			},
		},
		{
			name: "failed cache data deletion due to closed Redis client",

			setupMock: func() *redis.Client {
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

			redisMockClient := tc.setupMock()

			cacheDB := NewRedisCache(redisMockClient)

			err := cacheDB.DelCacheData(ctx, "groupKey")
			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				tc.verifyFunc(ctx, redisMockClient)
			}

		})
	}
}
