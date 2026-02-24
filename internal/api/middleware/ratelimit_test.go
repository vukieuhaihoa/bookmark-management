package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mockRateLimit "github.com/vukieuhaihoa/bookmark-management/internal/app/repository/ratelimit/mocks"
)

func TestRateLimit_RateLimit(t *testing.T) {
	t.Parallel()

	// httptest.NewRequest sets RemoteAddr to "192.0.2.1:1234" by default,
	// so Gin's ClientIP() always returns "192.0.2.1" in tests.
	testKey := fmt.Sprintf(RateLimitKeyFormat, "192.0.2.1")

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockRepo    func(ctx context.Context) *mockRateLimit.Repository
		expectedCode     int
		expectedResponse string
		expectedAborted  bool
	}{
		{
			name: "first request - counter not exists (returns -1)",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/api/resource", nil)
			},
			setupMockRepo: func(ctx context.Context) *mockRateLimit.Repository {
				repoMock := mockRateLimit.NewRepository(t)
				repoMock.On("GetCurrentRateLimit", mock.Anything, testKey).
					Return(-1, nil)
				repoMock.On("IncreaseRateLimit", mock.Anything, testKey, IP_RateLimitInterval).
					Return(nil)
				return repoMock
			},
			expectedCode:     http.StatusOK,
			expectedResponse: ``,
			expectedAborted:  false,
		},
		{
			name: "under the limit - counter below max",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/api/resource", nil)
			},
			setupMockRepo: func(ctx context.Context) *mockRateLimit.Repository {
				repoMock := mockRateLimit.NewRepository(t)
				repoMock.On("GetCurrentRateLimit", ctx, testKey).
					Return(IP_RateLimitMaxCount-1, nil)
				repoMock.On("IncreaseRateLimit", ctx, testKey, IP_RateLimitInterval).
					Return(nil)
				return repoMock
			},
			expectedCode:     http.StatusOK,
			expectedResponse: ``,
			expectedAborted:  false,
		},
		{
			name: "rate limit exceeded - counter at max",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/api/resource", nil)
			},
			setupMockRepo: func(ctx context.Context) *mockRateLimit.Repository {
				repoMock := mockRateLimit.NewRepository(t)
				repoMock.On("GetCurrentRateLimit", ctx, testKey).
					Return(IP_RateLimitMaxCount, nil)
				// IncreaseRateLimit must NOT be called when rate limit is exceeded
				return repoMock
			},
			expectedCode:     http.StatusTooManyRequests,
			expectedResponse: `{"error":"Too many requests. Please try again later."}`,
			expectedAborted:  true,
		},
		{
			name: "GetCurrentRateLimit error - logs and falls through",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/api/resource", nil)
			},
			setupMockRepo: func(ctx context.Context) *mockRateLimit.Repository {
				repoMock := mockRateLimit.NewRepository(t)
				repoMock.On("GetCurrentRateLimit", ctx, testKey).
					Return(-1, assert.AnError)
				// currRate == -1, so rate limit check is skipped; Increase is still called
				repoMock.On("IncreaseRateLimit", ctx, testKey, IP_RateLimitInterval).
					Return(nil)
				return repoMock
			},
			expectedCode:     http.StatusOK,
			expectedResponse: ``,
			expectedAborted:  false,
		},
		{
			name: "IncreaseRateLimit error - logs and proceeds",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/api/resource", nil)
			},
			setupMockRepo: func(ctx context.Context) *mockRateLimit.Repository {
				repoMock := mockRateLimit.NewRepository(t)
				repoMock.On("GetCurrentRateLimit", ctx, testKey).
					Return(5, nil)
				repoMock.On("IncreaseRateLimit", ctx, testKey, IP_RateLimitInterval).
					Return(assert.AnError)
				return repoMock
			},
			expectedCode:     http.StatusOK,
			expectedResponse: ``,
			expectedAborted:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx)
			repoMock := tc.setupMockRepo(ctx)

			rateLimitMiddleware := NewRateLimit(repoMock)
			rateLimitMiddleware.RateLimit(RateLimitIPKey)(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
			assert.Equal(t, tc.expectedAborted, ctx.IsAborted())
		})
	}
}
