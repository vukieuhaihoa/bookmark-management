package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/service"
	"github.com/vukieuhaihoa/bookmark-management/internal/service/mocks"
)

func TestShortenURL_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(ctx *gin.Context) *mocks.ShortenURL

		expectedError    error
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "successful URL shortening",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(`{"url":"http://example.com","exp":3600}`))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURL {
				svcMock := mocks.NewShortenURL(t)
				svcMock.On("ShortenURL", ctx, "http://example.com", 3600).Return("abcd1234", nil)
				return svcMock
			},

			expectedError:    nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"code":"abcd1234","message":"Shorten URL generated successfully!"}`,
		},
		{
			name: "invalid request payload",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(`{"url":"", "exp":-1}`))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURL {
				return mocks.NewShortenURL(t)
			},
			expectedError:    nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"invalid request payload"}`,
		},
		{
			name: "service error during URL shortening",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(`{"url":"http://example.com","exp":3600}`))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURL {
				svcMock := mocks.NewShortenURL(t)
				svcMock.On("ShortenURL", ctx, "http://example.com", 3600).Return("", assert.AnError)
				return svcMock
			},

			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"message":"internal server error"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx)
			svcMock := tc.setupMockSvc(ctx)

			handler := NewShortenURL(svcMock)
			handler.ShortenURL(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

func TestShortenURL_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)

		setupMockSvc func(ctx *gin.Context) *mocks.ShortenURL

		expectedCode int
		expectedURL  string
	}{
		{
			name: "successful get URL",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/abcd1234", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "abcd1234"}}
			},

			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURL {
				svcMock := mocks.NewShortenURL(t)
				svcMock.On("GetURL", ctx, "abcd1234").Return("http://example.com", nil)
				return svcMock
			},

			expectedCode: http.StatusFound,
			expectedURL:  "http://example.com",
		},
		{
			name: "code not found",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/unknown", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "unknown"}}
			},

			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURL {
				svcMock := mocks.NewShortenURL(t)
				svcMock.On("GetURL", ctx, "unknown").Return("", service.ErrCodeNotFound)
				return svcMock
			},

			expectedCode: http.StatusBadRequest,
			expectedURL:  "",
		},
		{
			name: "internal server error",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/abcd1234", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "abcd1234"}}
			},

			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURL {
				svcMock := mocks.NewShortenURL(t)
				svcMock.On("GetURL", ctx, "abcd1234").Return("", assert.AnError)
				return svcMock
			},

			expectedCode: http.StatusInternalServerError,
			expectedURL:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx)
			svcMock := tc.setupMockSvc(ctx)

			handler := NewShortenURL(svcMock)
			handler.GetURL(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedURL, rec.Header().Get("Location"))
		})
	}
}
