package endpoint

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-management/pkg/sqldb"
)

func TestShortenURLEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode int
		expectedCodeLength int
	}{
		{
			name: "successful shorten URL",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/links/shorten", strings.NewReader(`{"url":"http://example.com","exp":3600}`)) // Body would be added
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusOK,
			expectedCodeLength: 8,
		},
		{
			name: "invalid request payload",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/links/shorten", strings.NewReader(`{"url":"", "exp":-1}`)) // Body would be added
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedCodeLength: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			apiEngine := api.New(&api.Config{
				ServiceName: "bookmark-service",
				InstanceID:  "test_instance_id_1",
			}, redisPkg.InitMockRedis(t), sqldb.InitMockDB(t))

			respRec := tc.setupTestHTTP(apiEngine)

			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			// Check code length
			var respBody struct {
				Code string `json:"code"`
			}
			err := json.Unmarshal(respRec.Body.Bytes(), &respBody)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCodeLength, len(respBody.Code))
		})
	}
}

func TestGetURLEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRedis func(ctx context.Context) *redis.Client

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode int
		expectedLocation   string
	}{
		{
			name: "successful get original URL",

			setupMockRedis: func(ctx context.Context) *redis.Client {
				mockRedis := redisPkg.InitMockRedis(t)
				mockRedis.Set(ctx, "abcd1234", "http://example.com", 1000)
				return mockRedis
			},

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/links/redirect/abcd1234", nil) // Body would be added
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusFound,
			expectedLocation:   "http://example.com",
		},
		{
			name: "code not found",

			setupMockRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/links/redirect/unknown", nil) // Body would be added
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedLocation:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			mockRedis := tc.setupMockRedis(ctx)

			apiEngine := api.New(&api.Config{
				ServiceName: "bookmark-service",
				InstanceID:  "test_instance_id_1",
			}, mockRedis, nil)

			respRec := tc.setupTestHTTP(apiEngine)

			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			if tc.expectedLocation != "" {
				assert.Equal(t, tc.expectedLocation, respRec.Header().Get("Location"))
			}
		})
	}
}
