package bookmark

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	svcMocks "github.com/vukieuhaihoa/bookmark-management/internal/app/service/bookmark/mocks"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
)

var testQueryOpts = &common.QueryOptions{
	Paging: common.Paging{
		Page:  1,
		Limit: 2,
	},
	Sorting: []common.SortedField{
		{
			Field:     "created_at",
			Direction: "DESC",
		},
	},
}

func TestHandler_ListBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRequest func(ctx *gin.Context)

		setupMockSvc func(ctx *gin.Context) *svcMocks.Service

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful list bookmarks",

			setupMockRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=2&sort=-created_at", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user
				ctx.Set("claims", jwt.MapClaims{"sub": "de305d54-75b4-431b-adb2-eb6b9e546099"})
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.Service {
				svcMock := svcMocks.NewService(t)
				svcMock.On("ListBookmarks", ctx, "de305d54-75b4-431b-adb2-eb6b9e546099", mock.MatchedBy(func(opts *common.QueryOptions) bool {
					return opts.Page == testQueryOpts.Page && opts.Limit == testQueryOpts.Limit && opts.Sorting[0].Field == testQueryOpts.Sorting[0].Field && opts.Sorting[0].Direction == testQueryOpts.Sorting[0].Direction
				})).
					Return([]*model.Bookmark{
						{
							Base: model.Base{
								ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000008",
								CreatedAt: fixture.TestTime.Add(4 * time.Hour),
								UpdatedAt: fixture.TestTime.Add(4 * time.Hour),
							},
							UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
							URL:                "https://db-tutorials.dev/postgresql-indexing",
							Description:        "Learn PostgreSQL indexing basics",
							CodeShortenEncoded: "p_8",
						},
						{
							Base: model.Base{
								ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000009",
								CreatedAt: fixture.TestTime.Add(5 * time.Hour),
								UpdatedAt: fixture.TestTime.Add(5 * time.Hour),
							},
							UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
							URL:                "https://redis.io/docs/manual/data-types/",
							Description:        "Redis data types documentation",
							CodeShortenEncoded: "p_9",
						},
					}, nil)
				return svcMock
			},

			expectedCode:     http.StatusOK,
			expectedResponse: `{"data":[{"id":"a1b2c3d4-e5f6-7890-abcd-ef0000000008","created_at":"2023-01-01T04:00:00Z","updated_at":"2023-01-01T04:00:00Z","description":"Learn PostgreSQL indexing basics","url":"https://db-tutorials.dev/postgresql-indexing","code":"p_8"},{"id":"a1b2c3d4-e5f6-7890-abcd-ef0000000009","created_at":"2023-01-01T05:00:00Z","updated_at":"2023-01-01T05:00:00Z","description":"Redis data types documentation","url":"https://redis.io/docs/manual/data-types/","code":"p_9"}],"pagination":{"page":1,"limit":2,"total":0}}`,
		},
		{
			name: "invalid sort field",

			setupMockRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=2&sort=-invalidField", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user
				ctx.Set("claims", jwt.MapClaims{"sub": "de305d54-75b4-431b-adb2-eb6b9e546099"})
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.Service {
				svcMock := svcMocks.NewService(t)
				return svcMock
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid sorted field"}`,
		},
		{
			name: "service error",

			setupMockRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=2&sort=-created_at", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user
				ctx.Set("claims", jwt.MapClaims{"sub": "de305d54-75b4-431b-adb2-eb6b9e546099"})
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.Service {
				svcMock := svcMocks.NewService(t)
				svcMock.On("ListBookmarks", ctx, "de305d54-75b4-431b-adb2-eb6b9e546099", mock.Anything).
					Return(nil, assert.AnError)
				return svcMock
			},
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"message":"Internal server error"}`,
		},
		{
			name: "unauthorized user",
			setupMockRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=2&sort=-created_at", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// No authenticated user set
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.Service {
				svcMock := svcMocks.NewService(t)
				return svcMock
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "invalid paging parameters",
			setupMockRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest("GET", "/v1/bookmarks?page=0&limit=60&sort=-created_at", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user
				ctx.Set("claims", jwt.MapClaims{"sub": "de305d54-75b4-431b-adb2-eb6b9e546099"})
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.Service {
				svcMock := svcMocks.NewService(t)
				return svcMock
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid input fields","details":["Page is invalid (gte)","Limit is invalid (lte)"]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			// Setup mock request
			tc.setupMockRequest(ctx)

			// Setup mock service
			mockBookmarkSvc := tc.setupMockSvc(ctx)

			// Create handler with mock service
			handler := NewBookmarkHandler(mockBookmarkSvc)

			// Call the ListBookmarks handler
			handler.ListBookmarks(ctx)

			// Assert response code
			assert.Equal(t, tc.expectedCode, rec.Code)
			// Assert response body
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}
