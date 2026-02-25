package bookmark

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestBookmarkEndPoint_UpdateBookmarkByID(t *testing.T) {
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
			name: "update bookmark by ID successfully",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("PUT", "/v1/bookmarks/a1b2c3d4-e5f6-7890-abcd-ef0000000005", strings.NewReader(`{"url":"https://updated-example.com","description":"This is an updated description."}`))
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
			name: "update bookmark by ID failed - bookmark not found",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("PUT", "/v1/bookmarks/non-existent-id", strings.NewReader(`{"url":"https://updated-example.com","description":"This is an updated description."}`))
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
			name: "update bookmark of another user - unauthorized",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("PUT", "/v1/bookmarks/a1b2c3d4-e5f6-7890-abcd-ef0000000005", strings.NewReader(`{"url":"https://updated-example.com","description":"This is an updated description."}`))
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "de305d54-75b4-431b-adb2-eb6b9e546000"}, nil)
				return jwtValidator
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"Invalid input"}`,
		},
		{
			name: "missing authorization token",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("PUT", "/v1/bookmarks/a1b2c3d4-e5f6-7890-abcd-ef0000000005", strings.NewReader(`{"url":"https://updated-example.com","description":"This is an updated description."}`))
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				return mocks.NewJWTValidator(t)
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Authorization header missing"}`,
		},
		{
			name: "invalid authorization header format",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("PUT", "/v1/bookmarks/a1b2c3d4-e5f6-7890-abcd-ef0000000005", strings.NewReader(`{"url":"https://updated-example.com","description":"This is an updated description."}`))
				req.Header.Set("Authorization", "InvalidFormatToken")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				return mocks.NewJWTValidator(t)
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Invalid Authorization header format"}`,
		},
		{
			name: "update bookmark failed - rate limit exceeded",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("PUT", "/v1/bookmarks/a1b2c3d4-e5f6-7890-abcd-ef0000000005", strings.NewReader(`{"url":"https://updated-example.com","description":"This is an updated description."}`))
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
		{
			name: "update bookmark failed - invalid token",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("PUT", "/v1/bookmarks/a1b2c3d4-e5f6-7890-abcd-ef0000000005", strings.NewReader(`{"url":"https://updated-example.com","description":"This is an updated description."}`))
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
			name: "update bookmark failed - token does not contain user ID",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("PUT", "/v1/bookmarks/a1b2c3d4-e5f6-7890-abcd-ef0000000005", strings.NewReader(`{"url":"https://updated-example.com","description":"This is an updated description."}`))
				req.Header.Set("Authorization", "Bearer token_without_user_id")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "token_without_user_id").Return(jwt.MapClaims{}, nil)
				return jwtValidator
			},

			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"message":"Unauthorized"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			jwtValidator := tc.setupMockJWTValidator(t)
			redisClient := redisPkg.InitMockRedis(t)
			if tc.setupMockRedis != nil {
				redisClient = tc.setupMockRedis(ctx, redisClient)
			}

			// Setup API engine
			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark_service",
					InstanceID:  "bookmark_service_instance_01",
				},
				RedisClient:  redisClient,
				SqlDB:        db,
				JWTValidator: jwtValidator,
			})
			// Setup HTTP request and recorder
			respRec := tc.setupTestHTTP(apiEngine)

			// Assert response code
			assert.Equal(t, tc.expectedStatusCode, respRec.Code)

			// Assert response body
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(respRec.Body.String()))
		})
	}
}
