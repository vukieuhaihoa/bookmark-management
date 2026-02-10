package link

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	service "github.com/vukieuhaihoa/bookmark-management/internal/app/service/link"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/service/link/mocks"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

func TestHandler_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)

		setupMockSvc func(ctx *gin.Context) *mocks.Service

		expectedCode int
		expectedURL  string
	}{
		{
			name: "successful get URL",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/abcd1234", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "abcd1234"}}
			},

			setupMockSvc: func(ctx *gin.Context) *mocks.Service {
				svcMock := mocks.NewService(t)
				svcMock.On("GetURL", ctx, "abcd1234").Return("http://example.com", nil)
				return svcMock
			},

			expectedCode: http.StatusFound,
			expectedURL:  "http://example.com",
		},
		{
			name: "code not found - redis code",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/leetcode", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "leetcode"}}
			},

			setupMockSvc: func(ctx *gin.Context) *mocks.Service {
				svcMock := mocks.NewService(t)
				svcMock.On("GetURL", ctx, "leetcode").Return("", service.ErrCodeNotFound)
				return svcMock
			},

			expectedCode: http.StatusBadRequest,
			expectedURL:  "",
		},
		{
			name: "code not found - bookmark code",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/leetcodexx", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "leetcodexx"}}
			},

			setupMockSvc: func(ctx *gin.Context) *mocks.Service {
				svcMock := mocks.NewService(t)
				svcMock.On("GetURL", ctx, "leetcodexx").Return("", dbutils.ErrRecordNotFoundType)
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

			setupMockSvc: func(ctx *gin.Context) *mocks.Service {
				svcMock := mocks.NewService(t)
				svcMock.On("GetURL", ctx, "abcd1234").Return("", assert.AnError)
				return svcMock
			},

			expectedCode: http.StatusInternalServerError,
			expectedURL:  "",
		},
		{
			name: "missing code parameter",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect", nil)
				// No code parameter set
			},

			setupMockSvc: func(ctx *gin.Context) *mocks.Service {
				return mocks.NewService(t)
			},

			expectedCode: http.StatusBadRequest,
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

			handler := NewLinkHandler(svcMock)
			handler.GetURL(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedURL, rec.Header().Get("Location"))
		})
	}
}
