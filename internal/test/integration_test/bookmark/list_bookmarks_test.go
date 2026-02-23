package bookmark

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/jwtutils/mocks"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
	"gorm.io/gorm"
)

func TestBookmarkEndpoint_ListBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP         func(api api.Engine) *httptest.ResponseRecorder
		setupDB               func(t *testing.T) *gorm.DB
		setupCache            func(t *testing.T, ctx context.Context) *redis.Client
		setupMockJWTValidator func(t *testing.T) *mocks.JWTValidator

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successful list bookmarks",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=2&sort=-created_at", nil)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			setupCache: func(t *testing.T, ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				return redisClient
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":[{"id":"a1b2c3d4-e5f6-7890-abcd-ef0000000009","created_at":"2023-01-01T05:00:00Z","updated_at":"2023-01-01T05:00:00Z","description":"Redis data types documentation","url":"https://redis.io/docs/manual/data-types/","code":"p_9"},{"id":"a1b2c3d4-e5f6-7890-abcd-ef0000000008","created_at":"2023-01-01T04:00:00Z","updated_at":"2023-01-01T04:00:00Z","description":"Learn PostgreSQL indexing basics","url":"https://db-tutorials.dev/postgresql-indexing","code":"p_8"}],"pagination":{"page":1,"limit":2,"total":6}}`,
		},
		{
			name: "successful list bookmarks with cache hit",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=2&sort=-created_at", nil)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			setupCache: func(t *testing.T, ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				groupKey := "list_bookmarks_4d9326d6-980c-4c62-9709-dbc70a82cbfe"
				cacheKey := "page_1_size_2_sortby_created_at_desc"
				cachedData := `{"bookmarks":[{"id":"a1b2c3d4-e5f6-7890-abcd-ef0000000001","url":"https://example.com/testuser001","code":"p_1","description":"Bookmark for Test User 1 - record 1","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"}],"total":1}`
				redisClient.HSet(ctx, groupKey, cacheKey, []byte(cachedData))
				redisClient.Expire(ctx, groupKey, time.Hour)
				return redisClient
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":[{"id":"a1b2c3d4-e5f6-7890-abcd-ef0000000009","created_at":"2023-01-01T05:00:00Z","updated_at":"2023-01-01T05:00:00Z","description":"Redis data types documentation","url":"https://redis.io/docs/manual/data-types/","code":"p_9"},{"id":"a1b2c3d4-e5f6-7890-abcd-ef0000000008","created_at":"2023-01-01T04:00:00Z","updated_at":"2023-01-01T04:00:00Z","description":"Learn PostgreSQL indexing basics","url":"https://db-tutorials.dev/postgresql-indexing","code":"p_8"}],"pagination":{"page":1,"limit":2,"total":6}}`,
		},
		{
			name: "invalid query parameters",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/bookmarks?page=-1&limit=0&sort=invalidField", nil)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupDB: func(t *testing.T) *gorm.DB {
				return nil
			},

			setupCache: func(t *testing.T, ctx context.Context) *redis.Client {
				return nil
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "{\"message\":\"Invalid input fields\",\"details\":[\"Page is invalid (gte)\",\"Limit is invalid (gte)\"]}",
		},
		{
			name: "list bookmarks failed - invalid token",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=2&sort=-createdAt", nil)
				req.Header.Set("Authorization", "Bearer invalid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupDB: func(t *testing.T) *gorm.DB {
				return nil
			},

			setupCache: func(t *testing.T, ctx context.Context) *redis.Client {
				return nil
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
			name: "list bookmarks failed - token does not contain user ID",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=2&sort=-created_at", nil)
				req.Header.Set("Authorization", "Bearer token_without_user_id")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupDB: func(t *testing.T) *gorm.DB {
				return nil
			},

			setupCache: func(t *testing.T, ctx context.Context) *redis.Client {
				return nil
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "token_without_user_id").Return(jwt.MapClaims{}, nil)
				return jwtValidator
			},

			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"message":"Unauthorized"}`,
		},
		{
			name: "list bookmarks failed - invalid sort field",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=2&sort=invalidField", nil)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupDB: func(t *testing.T) *gorm.DB {
				return nil
			},

			setupCache: func(t *testing.T, ctx context.Context) *redis.Client {
				return nil
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"Invalid sorted field"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			db := tc.setupDB(t)
			cache := tc.setupCache(t, ctx)
			jwtValidator := tc.setupMockJWTValidator(t)

			api := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark_service",
				},
				RedisClient:  cache,
				SqlDB:        db,
				JWTValidator: jwtValidator,
			})

			respRec := tc.setupTestHTTP(api)

			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.Equal(t, tc.expectedResponse, respRec.Body.String())
		})
	}
}
