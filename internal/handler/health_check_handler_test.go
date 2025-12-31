package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/service"
	"github.com/vukieuhaihoa/bookmark-management/internal/service/mocks"
)

func TestHealthCheckService_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func() *mocks.HealthCheck

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "successful health check",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockSvc: func() *mocks.HealthCheck {
				svcMock := mocks.NewHealthCheck(t)
				svcMock.On("Check").Return(service.StatusOK, "TestService", "Instance123")
				return svcMock
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"OK","service_name":"TestService","instance_id":"Instance123"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)

			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc()
			testHandler := &healthCheckHandler{svc: mockSvc}

			testHandler.Check(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
