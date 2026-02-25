package link

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	"github.com/vukieuhaihoa/bookmark-management/internal/api/middleware"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-management/pkg/utils"
)

func TestGetURLEndpoint_RedirectCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRedis func(ctx context.Context, redisClient *redis.Client) *redis.Client

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode int
		expectedLocation   string
	}{
		{
			name: "successful get original URL",

			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				redisClient.Set(ctx, "abcd1234", "http://example.com", 1000)
				return redisClient
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
		{
			name: "rate limit exceeded",

			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				key := fmt.Sprintf(middleware.RateLimitKeyFormat, "192.0.2.1")
				redisClient.Set(ctx, key, middleware.IPRateLimitMaxCount, middleware.IPRateLimitInterval)
				return redisClient
			},

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/links/redirect/abcd1234", nil) // Body would be added
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusTooManyRequests,
			expectedLocation:   "",
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

			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})

			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark-service",
					InstanceID:  "test_instance_id_1",
				},
				RedisClient:     redisClient,
				SqlDB:           db,
				RandomCodeGen:   utils.NewCodeGenerator(),
				PasswordHashing: nil,
				JWTGenerator:    nil,
				JWTValidator:    nil,
			})

			respRec := tc.setupTestHTTP(apiEngine)

			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			if tc.expectedLocation != "" {
				assert.Equal(t, tc.expectedLocation, respRec.Header().Get("Location"))
			}
		})
	}
}
