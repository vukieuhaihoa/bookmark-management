package bookmark

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	"github.com/vukieuhaihoa/bookmark-management/internal/api/middleware"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/jwtutils/mocks"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
)

func TestBookmarkEndpoint_DeleteBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRedis func(ctx context.Context, redisClient *redis.Client) *redis.Client

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		setupMockJWTValidator func(t *testing.T) *mocks.JWTValidator

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "delete bookmark by ID successfully",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("DELETE", "/v1/bookmarks/a1b2c3d4-e5f6-7890-abcd-ef0000000006", nil)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"message":"Success"}`,
		},
		{
			name: "delete bookmark by ID failed - bookmark not found",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("DELETE", "/v1/bookmarks/non-existent-id", nil)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"Invalid input"}`,
		},
		{
			name: "delete bookmark failed - deleted another user's bookmark",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("DELETE", "/v1/bookmarks/a1b2c3d4-e5f6-7890-abcd-ef0000000001", nil)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"Invalid input"}`,
		},
		{
			name: "delete bookmark failed - invalid token",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("DELETE", "/v1/bookmarks/a1b2c3d4-e5f6-7890-abcd-ef0000000006", nil)
				req.Header.Set("Authorization", "Bearer invalid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "invalid_jwt_token").Return(nil, assert.AnError)
				return jwtValidator
			},

			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"message":"Invalid token"}`,
		},
		{
			name: "delete bookmark failed - rate limit exceeded",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("DELETE", "/v1/bookmarks/a1b2c3d4-e5f6-7890-abcd-ef0000000006", nil)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				key := fmt.Sprintf(middleware.RateLimitKeyFormat, "4d9326d6-980c-4c62-9709-dbc70a82cbfe")
				redisClient.Set(ctx, key, middleware.UserIDRateLimitMaxCount, middleware.UserIDRateLimitInterval)
				return redisClient
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedStatusCode: http.StatusTooManyRequests,
			expectedResponse:   `{"error":"Too many requests. Please try again later."}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			jwtValidatorMock := tc.setupMockJWTValidator(t)
			redisClient := redisPkg.InitMockRedis(t)
			if tc.setupMockRedis != nil {
				redisClient = tc.setupMockRedis(ctx, redisClient)
			}

			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark_service",
					InstanceID:  "bookmark_service_instance_1",
				},
				RedisClient:  redisClient,
				SqlDB:        db,
				JWTValidator: jwtValidatorMock,
			})

			// setup HTTP request and recorder
			respRec := tc.setupTestHTTP(apiEngine)

			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.Equal(t, tc.expectedResponse, respRec.Body.String())
		})
	}
}
