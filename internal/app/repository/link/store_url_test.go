package link

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
)

func Test_LinkRepository_StoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func() *redis.Client

		expectedError error

		verifyFunc func(ctx context.Context, redisClient *redis.Client)
	}{
		{
			name: "successful URL storage",

			setupMock: func() *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				return redisClient
			},

			expectedError: nil,

			verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
				url, err := redisClient.Get(ctx, "test").Result()

				assert.Nil(t, err)

				assert.Equal(t, "https://example.com", url)
			},
		},
		{
			name: "failed URL storage due to closed Redis client",

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

			urlStorage := NewLinkRepository(redisMockClient)

			err := urlStorage.StoreURL(ctx, "test", "https://example.com", 0)
			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				tc.verifyFunc(ctx, redisMockClient)
			}

		})
	}
}
