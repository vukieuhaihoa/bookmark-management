package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	mockJWTValidator "github.com/vukieuhaihoa/bookmark-management/pkg/jwtutils/mocks"
)

func TestJWTAuth_JWTAuth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest       func(ctx *gin.Context)
		setupMockValidator func() *mockJWTValidator.JWTValidator

		expectedCode     int
		expectedResponse string
		expectedAborted  bool
		checkClaims      bool
	}{
		{
			name: "valid token - claims set in context",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Authorization", "Bearer valid-token")
			},

			setupMockValidator: func() *mockJWTValidator.JWTValidator {
				mockValidator := &mockJWTValidator.JWTValidator{}
				mockValidator.On("ValidateToken", "valid-token").
					Return(jwt.MapClaims{"sub": "de305d54-75b4-431b-adb2-eb6b9e546099"}, nil)
				return mockValidator
			},

			expectedCode:    http.StatusOK,
			expectedAborted: false,
			checkClaims:     true,
		},
		{
			name: "missing Authorization header",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
			},
			setupMockValidator: func() *mockJWTValidator.JWTValidator {
				return mockJWTValidator.NewJWTValidator(t) // no expectations — not called
			},
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"error":"Authorization header missing"}`,
			expectedAborted:  true,
		},
		{
			name: "invalid format - no space separator",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Authorization", "BearerTokenWithNoSpace")
			},
			setupMockValidator: func() *mockJWTValidator.JWTValidator {
				return mockJWTValidator.NewJWTValidator(t)
			},
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"error":"Invalid Authorization header format"}`,
			expectedAborted:  true,
		},
		{
			name: "invalid format - wrong scheme",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Authorization", "Basic sometoken")
			},
			setupMockValidator: func() *mockJWTValidator.JWTValidator {
				return mockJWTValidator.NewJWTValidator(t)
			},
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"error":"Invalid Authorization header format"}`,
			expectedAborted:  true,
		},
		{
			name: "invalid token - validator returns error",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Authorization", "Bearer invalid-token")
			},
			setupMockValidator: func() *mockJWTValidator.JWTValidator {
				mockValidator := mockJWTValidator.NewJWTValidator(t)
				mockValidator.On("ValidateToken", "invalid-token").
					Return(nil, assert.AnError)
				return mockValidator
			},
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Invalid token"}`,
			expectedAborted:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx)
			mockValidator := tc.setupMockValidator()

			jwtAuthMiddleware := NewJWTAuth(mockValidator)
			handler := jwtAuthMiddleware.JWTAuth()

			handler(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
			assert.Equal(t, tc.expectedAborted, ctx.IsAborted())

			if tc.checkClaims {
				claims, exists := ctx.Get("claims")
				assert.True(t, exists)
				assert.IsType(t, jwt.MapClaims{}, claims)
			}

		})
	}
}
