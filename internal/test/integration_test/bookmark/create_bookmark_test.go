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
	"github.com/vukieuhaihoa/bookmark-management/pkg/utils"
)

func TestBookmarkEndpoint_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		setupMockJWTValidator func(t *testing.T) *mocks.JWTValidator

		expectedStatusCode      int
		expectedMessageResponse string
	}{
		{
			name: "successful create bookmark",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/bookmarks", strings.NewReader(`{"url":"http://example.com","description": "A sample bookmark"}`))
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

			expectedStatusCode:      http.StatusCreated,
			expectedMessageResponse: `"message":"Create a bookmark successfully!"`,
		},
		{
			name: "invalid create bookmark payload",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/bookmarks", strings.NewReader(`{"url":"invalid-url","description": ""}`))
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

			expectedStatusCode:      http.StatusBadRequest,
			expectedMessageResponse: `"message":"Invalid input"`,
		},
		{
			name: "create bookmark failed - invalid token",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/bookmarks", strings.NewReader(`{"url":"http://example.com","description": "A sample bookmark"}`))
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

			expectedStatusCode:      http.StatusUnauthorized,
			expectedMessageResponse: `"message":"Invalid token"`,
		},
		{
			name: "create bookmark failed - token does not contain user ID",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/bookmarks", strings.NewReader(`{"url":"http://example.com","description": "A sample bookmark"}`))
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

			expectedStatusCode:      http.StatusUnauthorized,
			expectedMessageResponse: `"message":"Unauthorized"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			jwtValidator := tc.setupMockJWTValidator(t)

			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark_service",
					InstanceID:  "test_instance_id_1",
				},
				RedisClient:     nil,
				SqlDB:           db,
				RandomCodeGen:   utils.NewCodeGenerator(),
				PasswordHashing: nil,
				JWTGenerator:    nil,
				JWTValidator:    jwtValidator,
			})
			respRec := tc.setupTestHTTP(apiEngine)

			// Verify response status code
			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.Contains(t, respRec.Body.String(), tc.expectedMessageResponse)
		})
	}
}
