package service

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/repository/mocks"
)

func TestHealthCheckService_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputServiceName string
		inputInstanceID  string

		setupMockRepo func(ctx context.Context) *mocks.HealthCheck

		expectedError       error
		expectedMessage     string
		expectedServiceName string
		expectedInstanceID  string
	}{
		{
			name: "Health check returns OK status",

			inputServiceName: "TestService",
			inputInstanceID:  "Instance123",

			setupMockRepo: func(ctx context.Context) *mocks.HealthCheck {
				repoMock := mocks.NewHealthCheck(t)
				repoMock.On("Ping", ctx).Return(nil)
				return repoMock
			},

			expectedError:       nil,
			expectedMessage:     StatusOK,
			expectedServiceName: "TestService",
			expectedInstanceID:  "Instance123",
		},
		{
			name: "Health Check return timeout error",

			inputServiceName: "TestService",
			inputInstanceID:  "Instance123",

			setupMockRepo: func(ctx context.Context) *mocks.HealthCheck {
				repoMock := mocks.NewHealthCheck(t)
				repoMock.On("Ping", ctx).Return(redis.ErrPoolTimeout)
				return repoMock
			},

			expectedError:       redis.ErrPoolTimeout,
			expectedMessage:     RedisPingTimeout,
			expectedServiceName: "TestService",
			expectedInstanceID:  "Instance123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := tc.setupMockRepo(t.Context())

			testSvc := NewHealthCheck(tc.inputServiceName, tc.inputInstanceID, mockRepo)

			status, serviceName, instanceID, err := testSvc.Check(t.Context())

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedMessage, status)
			assert.Equal(t, tc.expectedServiceName, serviceName)
			assert.Equal(t, tc.expectedInstanceID, instanceID)
		})
	}
}
