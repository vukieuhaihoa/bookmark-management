package link

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	"github.com/vukieuhaihoa/bookmark-management/internal/api/middleware"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-management/pkg/sqldb"
	"github.com/vukieuhaihoa/bookmark-management/pkg/utils"
)

func TestShortenURLEndpoint_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRedis func(ctx context.Context, redisClient *redis.Client) *redis.Client
		setupTestHTTP  func(api api.Engine) *httptest.ResponseRecorder

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
			expectedCodeLength: 9,
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
		{
			name: "rate limit exceeded",

			// Key insight: httptest.NewRequest always sets RemoteAddr = "192.0.2.1:1234", so gin's c.ClientIP() will always return 192.0.2.1. The rate limit key becomes rate_limit:192.0.2.1.
			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				key := fmt.Sprintf(middleware.RateLimitKeyFormat, "192.0.2.1")
				redisClient.Set(ctx, key, middleware.IPRateLimitMaxCount, middleware.IPRateLimitInterval)
				return redisClient
			},

			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("POST", "/v1/links/shorten", strings.NewReader(`{"url":"http://example.com","exp":3600}`))
				respRec := httptest.NewRecorder()
				engine.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusTooManyRequests,
			expectedCodeLength: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			redisClient := redisPkg.InitMockRedis(t)

			if tc.setupMockRedis != nil {
				redisClient = tc.setupMockRedis(ctx, redisClient)
			}

			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark-service",
					InstanceID:  "test_instance_id_1",
				},
				RedisClient:     redisClient,
				SqlDB:           sqldb.InitMockDB(t),
				RandomCodeGen:   utils.NewCodeGenerator(),
				PasswordHashing: nil,
				JWTGenerator:    nil,
				JWTValidator:    nil,
			})

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
