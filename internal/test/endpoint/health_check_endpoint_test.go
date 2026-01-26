package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-management/pkg/sqldb"
)

func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Health check returns OK status",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("GET", "/health-check", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"OK","service_name":"bookmark-service","instance_id":"test_instance_id_1"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			apiEngine := api.New(&api.Config{
				ServiceName: "bookmark-service",
				InstanceID:  "test_instance_id_1",
			}, redisPkg.InitMockRedis(t), sqldb.InitMockDB(t), nil, nil)

			respRec := tc.setupTestHTTP(apiEngine)

			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.JSONEq(t, tc.expectedResponseBody, respRec.Body.String())
		})
	}
}
