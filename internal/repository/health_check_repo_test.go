package repository

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
)

func TestHealthCheckRepo_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func() *redis.Client

		expectedError error
	}{
		{
			name: "successful ping",

			setupMock: func() *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				return redisClient
			},

			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMockClient := tc.setupMock()

			healthCheckRepo := NewHealthCheck(redisMockClient)

			err := healthCheckRepo.Ping(ctx)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
