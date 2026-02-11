package link

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-management/pkg/utils"
)

func TestGetURLEndpoint_RedirectCode(t *testing.T) {
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

			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})

			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark-service",
					InstanceID:  "test_instance_id_1",
				},
				RedisClient:     mockRedis,
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
