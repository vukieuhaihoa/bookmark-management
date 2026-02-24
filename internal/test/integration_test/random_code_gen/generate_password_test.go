package random_code_gen

import (
	"context"
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
	"github.com/vukieuhaihoa/bookmark-management/pkg/utils"
)

func TestPasswordEndpoint_GeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRedis    func(ctx context.Context, redisClient *redis.Client) *redis.Client
		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode      int
		expectedResponseLength  int
		expectedMessageResponse string
	}{
		{
			name: "Generate password successfully",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("GET", "/v1/generate-password", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:     http.StatusOK,
			expectedResponseLength: 12,
		},
		{
			name: "rate limit exceeded",

			setupRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				key := fmt.Sprintf(middleware.RateLimitKeyFormat, "192.0.2.1")
				redisClient.Set(ctx, key, middleware.IPRateLimitMaxCount, middleware.IPRateLimitInterval)
				return redisClient
			},

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("GET", "/v1/generate-password", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusTooManyRequests,
			expectedResponseLength:  0,
			expectedMessageResponse: `{"error":"Too many requests. Please try again later."}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			redisClient := redisPkg.InitMockRedis(t)

			if tc.setupRedis != nil {
				redisClient = tc.setupRedis(ctx, redisClient)
			}

			apiEngine := api.New(&api.EngineOpts{
				Engine:          gin.New(),
				Cfg:             &api.Config{},
				RedisClient:     redisClient,
				SqlDB:           nil,
				RandomCodeGen:   utils.NewCodeGenerator(),
				PasswordHashing: nil,
				JWTGenerator:    nil,
				JWTValidator:    nil,
			})

			respRec := tc.setupTestHTTP(apiEngine)

			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			if tc.expectedResponseLength > 0 {
				assert.Equal(t, tc.expectedResponseLength, respRec.Body.Len())
			}
			if tc.expectedMessageResponse != "" {
				assert.Equal(t, tc.expectedMessageResponse, strings.TrimSpace(respRec.Body.String()))
			}
		})
	}
}
