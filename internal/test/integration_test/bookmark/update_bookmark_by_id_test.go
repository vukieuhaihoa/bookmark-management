package bookmark

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/jwtutils/mocks"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
)

func TestBookmarkEndPoint_UpdateBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			jwtValidator := tc.setupMockJWTValidator(t)
			// Setup API engine
			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark_service",
					InstanceID:  "bookmark_service_instance_01",
				},
				RedisClient:  redisPkg.InitMockRedis(t),
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
