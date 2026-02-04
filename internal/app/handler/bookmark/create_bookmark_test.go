package bookmark

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	svcMocks "github.com/vukieuhaihoa/bookmark-management/internal/app/service/bookmark/mocks"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

var jwtClaims = jwt.MapClaims{
	"sub": "user-123",
}

func TestHandler_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputRequest *createBookmarkRequest

		setupRequest func(ctx *gin.Context, inputRequest *createBookmarkRequest)

		setupMockSvc func(ctx *gin.Context, inputRequest *createBookmarkRequest) *svcMocks.Service

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful create bookmark",

			inputRequest: &createBookmarkRequest{
				URL:         "https://example.com",
				Description: "Example Website",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user
				ctx.Set("claims", jwtClaims)
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createBookmarkRequest) *svcMocks.Service {
				svcMock := svcMocks.NewService(t)
				svcMock.On("CreateBookmark", ctx, inputRequest.URL, inputRequest.Description, "user-123").
					Return(&model.Bookmark{
						Base: model.Base{
							ID: "bookmark-456",
						},
						URL:         inputRequest.URL,
						Description: inputRequest.Description,
						UserID:      "user-123",
						Code:        "abcdefghij",
					}, nil)
				return svcMock
			},

			expectedCode:     http.StatusCreated,
			expectedResponse: `{"data":{"id":"bookmark-456","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","description":"Example Website","url":"https://example.com","code":"abcdefghij"},"message":"Create a bookmark successfully!"}`,
		},
		{
			name: "unauthorized user",

			inputRequest: &createBookmarkRequest{
				URL:         "https://example.com",
				Description: "Example Website",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// No authenticated user set
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createBookmarkRequest) *svcMocks.Service {
				return svcMocks.NewService(t)
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "user id invalid - not exist in system",

			inputRequest: &createBookmarkRequest{
				URL:         "https://example.com",
				Description: "Example Website",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user
				ctx.Set("claims", jwtClaims)
			},
			setupMockSvc: func(ctx *gin.Context, inputRequest *createBookmarkRequest) *svcMocks.Service {
				svcMock := svcMocks.NewService(t)
				svcMock.On("CreateBookmark", mock.Anything, inputRequest.URL, inputRequest.Description, "user-123").
					Return(nil, dbutils.ErrForeignKeyType)
				return svcMock
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "invalid request body",

			inputRequest: &createBookmarkRequest{
				URL:         "invalid-url",
				Description: "",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user
				ctx.Set("claims", jwtClaims)
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createBookmarkRequest) *svcMocks.Service {
				return svcMocks.NewService(t)
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid input"}`,
		},
		{
			name: "service internal error",

			inputRequest: &createBookmarkRequest{
				URL:         "https://example.com",
				Description: "Example Website",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user
				ctx.Set("claims", jwtClaims)
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createBookmarkRequest) *svcMocks.Service {
				svcMock := svcMocks.NewService(t)
				svcMock.On("CreateBookmark", mock.Anything, inputRequest.URL, inputRequest.Description, "user-123").
					Return(nil, assert.AnError)
				return svcMock
			},

			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"message":"Internal server error"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx, tc.inputRequest)

			// Setup mock service
			mockSvc := tc.setupMockSvc(ctx, tc.inputRequest)

			// Create handler with mock service
			handler := &bookmarkHandler{
				svc: mockSvc,
			}

			// Invoke the handler
			handler.CreateBookmark(ctx)

			// Assert response code
			assert.Equal(t, tc.expectedCode, rec.Code)

			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}
