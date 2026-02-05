package link

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/service/link/mocks"
)

func TestHandler_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(ctx *gin.Context) *mocks.Service

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

			setupMockSvc: func(ctx *gin.Context) *mocks.Service {
				svcMock := mocks.NewService(t)
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
			setupMockSvc: func(ctx *gin.Context) *mocks.Service {
				return mocks.NewService(t)
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

			setupMockSvc: func(ctx *gin.Context) *mocks.Service {
				svcMock := mocks.NewService(t)
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

			handler := NewLinkHandler(svcMock)
			handler.ShortenURL(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
