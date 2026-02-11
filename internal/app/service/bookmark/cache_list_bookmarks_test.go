package bookmark_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	mock_cache "github.com/vukieuhaihoa/bookmark-management/internal/app/repository/cache/mocks"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/service/bookmark"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/service/bookmark/mocks"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
)

func TestBookmarkServiceWithCache_ListBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockSvc   func(ctx context.Context, userID string, opts *common.QueryOptions) *mocks.Service
		setupMockCache func(ctx context.Context, userID string, opts *common.QueryOptions) *mock_cache.DB

		inputUserID string
		inputOpts   *common.QueryOptions

		expectedBookmarks []*model.Bookmark
		expectedError     error
		expectedTotal     int64
	}{
		{
			name: "List bookmarks with cache hit",

			setupMockSvc: func(ctx context.Context, userID string, opts *common.QueryOptions) *mocks.Service {
				svcMock := mocks.NewService(t)
				return svcMock
			},

			setupMockCache: func(ctx context.Context, userID string, opts *common.QueryOptions) *mock_cache.DB {
				cacheMock := mock_cache.NewDB(t)
				groupKey := fmt.Sprintf(bookmark.ListBookmarksCacheGroupKey, userID)
				cacheKey := bookmark.GenerateCacheKeyFromQueryOptions(opts)
				cachedData := `{"bookmarks":[{"id":"a1b2c3d4-e5f6-7890-abcd-ef0000000001","url":"https://example.com/testuser001","code":"XXXX000001","description":"Bookmark for Test User 1 - record 1","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"}],"total":1}`
				cacheMock.On("GetCacheData", ctx, groupKey, cacheKey).Return([]byte(cachedData), nil)
				return cacheMock
			},

			inputUserID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			inputOpts: &common.QueryOptions{
				Paging: common.Paging{
					Page:  1,
					Limit: 10,
				},
				Sorting: []common.SortedField{
					{
						Field:     "created_at",
						Direction: common.SortDesc,
					},
				},
			},

			expectedBookmarks: []*model.Bookmark{
				{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000001",
						CreatedAt: fixture.TestTime,
						UpdatedAt: fixture.TestTime,
					},
					URL:         "https://example.com/testuser001",
					Code:        "XXXX000001",
					Description: "Bookmark for Test User 1 - record 1",
				},
			},
			expectedTotal: 1,
		},
		{
			name: "List bookmarks with cache miss",

			setupMockSvc: func(ctx context.Context, userID string, opts *common.QueryOptions) *mocks.Service {
				svcMock := mocks.NewService(t)
				svcMock.On("ListBookmarks", ctx, userID, opts).Run(
					func(args mock.Arguments) {
						argOpts := args.Get(2).(*common.QueryOptions)
						argOpts.Total = 1
					},
				).Return([]*model.Bookmark{
					{
						Base: model.Base{
							ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000001",
							CreatedAt: fixture.TestTime,
							UpdatedAt: fixture.TestTime,
						},
						URL:         "https://example.com/testuser001",
						Code:        "XXXX000001",
						Description: "Bookmark for Test User 1 - record 1",
					},
				}, nil)
				return svcMock
			},

			setupMockCache: func(ctx context.Context, userID string, opts *common.QueryOptions) *mock_cache.DB {
				cacheMock := mock_cache.NewDB(t)
				groupKey := fmt.Sprintf(bookmark.ListBookmarksCacheGroupKey, userID)
				cacheKey := bookmark.GenerateCacheKeyFromQueryOptions(opts)
				cacheMock.On("GetCacheData", ctx, groupKey, cacheKey).Return([]byte{}, redis.Nil)
				cacheMock.On("SetCacheData", ctx, groupKey, cacheKey, []byte(`{"bookmarks":[{"id":"a1b2c3d4-e5f6-7890-abcd-ef0000000001","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z","description":"Bookmark for Test User 1 - record 1","url":"https://example.com/testuser001","code":"XXXX000001"}],"total":1}`), bookmark.ListBookmarksCacheTTL).Return(nil)
				return cacheMock
			},

			inputUserID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			inputOpts: &common.QueryOptions{
				Paging: common.Paging{
					Page:  1,
					Limit: 10,
				},
				Sorting: []common.SortedField{
					{
						Field:     "created_at",
						Direction: common.SortDesc,
					},
				},
			},

			expectedBookmarks: []*model.Bookmark{
				{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000001",
						CreatedAt: fixture.TestTime,
						UpdatedAt: fixture.TestTime,
					},
					URL:         "https://example.com/testuser001",
					Code:        "XXXX000001",
					Description: "Bookmark for Test User 1 - record 1",
				},
			},
			expectedTotal: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			svcMock := tc.setupMockSvc(ctx, tc.inputUserID, tc.inputOpts)
			cacheMock := tc.setupMockCache(ctx, tc.inputUserID, tc.inputOpts)

			serviceWithCache := bookmark.NewBookmarkServiceWithCache(svcMock, cacheMock)

			bookmarks, err := serviceWithCache.ListBookmarks(ctx, tc.inputUserID, tc.inputOpts)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBookmarks, bookmarks)
				assert.Equal(t, tc.expectedTotal, tc.inputOpts.Total)
			}
		})

	}
}
