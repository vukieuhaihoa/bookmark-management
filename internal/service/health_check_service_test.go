package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckService_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputServiceName string
		inputInstanceID  string

		expectedMessage     string
		expectedServiceName string
		expectedInstanceID  string
	}{
		{
			name:                "Health check returns OK status",
			inputServiceName:    "TestService",
			inputInstanceID:     "Instance123",
			expectedMessage:     StatusOK,
			expectedServiceName: "TestService",
			expectedInstanceID:  "Instance123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewHealthCheck(tc.inputServiceName, tc.inputInstanceID)
			status, serviceName, instanceID := testSvc.Check()

			assert.Equal(t, tc.expectedMessage, status)
			assert.Equal(t, tc.expectedServiceName, serviceName)
			assert.Equal(t, tc.expectedInstanceID, instanceID)
		})
	}
}
