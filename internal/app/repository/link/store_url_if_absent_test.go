package link

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
)

func TestRepository_StoreURLIfAbsent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		code     string
		url      string
		expireIn int

		setupMock func(ctx context.Context) *redis.Client

		expectedResult bool
		expectedError  error
	}{
		{
			name: "store new URL successfully",

			code:     "abc123",
			url:      "https://example.com",
			expireIn: 3600,

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				return redisClient
			},

			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "do not overwrite existing URL",

			code:     "abc123",
			url:      "https://example.com/updated",
			expireIn: 3600,

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				redisClient.Set(ctx, "abc123", "https://example.com", time.Hour)
				return redisClient
			},

			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "failed to store URL due to closed Redis client",

			code:     "abc123",
			url:      "https://example.com",
			expireIn: 3600,

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				redisClient.Close()
				return redisClient
			},

			expectedResult: false,
			expectedError:  redis.ErrClosed,
		},
		{
			name: "invalid expiration time - negative value",

			code:     "abc123",
			url:      "https://example.com",
			expireIn: -1,

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				return redisClient
			},

			expectedResult: true, // should still store with default expiration
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			redisClient := tc.setupMock(ctx)
			repo := NewLinkRepository(redisClient)

			result, err := repo.StoreURLIfAbsent(ctx, tc.code, tc.url, tc.expireIn)
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
