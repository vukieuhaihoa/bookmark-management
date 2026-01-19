package repository

import (
	"database/sql"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-management/pkg/sqldb"
	"gorm.io/gorm"
)

func TestHealthCheckRepo_Ping(t *testing.T) {
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
		{
			name: "failed ping",

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

			healthCheckRepo := NewHealthCheck(redisMockClient, nil)

			err := healthCheckRepo.Ping(ctx)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestHealthCheckRepo_DBPing(t *testing.T) {
	t.Parallel()

	// Since DBPing is not implemented, we will just test that it returns nil for now.
	testCases := []struct {
		name string

		setupMockDB func() *gorm.DB

		expectedErrStr string
		expectedError  error
	}{
		{
			name: "DBPing returns nil",

			setupMockDB: func() *gorm.DB {
				db := sqldb.InitMockDB(t)

				return db
			},

			expectedError: nil,
		},
		{
			name: "DBPing on closed DB",

			setupMockDB: func() *gorm.DB {
				db := sqldb.InitMockDB(t)
				sqlDB, _ := db.DB()
				sqlDB.Close()

				return db
			},

			expectedError: sql.ErrConnDone,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupMockDB()

			healthCheckRepo := NewHealthCheck(nil, db)

			err := healthCheckRepo.DBPing(ctx)

			if tc.expectedError != nil {
				assert.Error(t, err)
				if tc.expectedErrStr != "" {
					assert.Contains(t, err.Error(), tc.expectedErrStr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
