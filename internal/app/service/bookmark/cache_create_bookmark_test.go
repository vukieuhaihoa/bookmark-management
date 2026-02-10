package bookmark_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	mock_cache "github.com/vukieuhaihoa/bookmark-management/internal/app/repository/cache/mocks"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/service/bookmark"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/service/bookmark/mocks"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
)

func TestBookmarkServiceWithCache_CreateBookmark(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string

		url            string
		description    string
		userID         string
		setupMockSvc   func(ctx context.Context, url, description, userID string) *mocks.Service
		setupMockCache func(ctx context.Context, userID string) *mock_cache.DB

		expectedBookmark *model.Bookmark
		expectedError    error
	}{
		{
			name: "Create bookmark with cache invalidation successfully",

			setupMockSvc: func(ctx context.Context, url, description, userID string) *mocks.Service {
				svcMock := mocks.NewService(t)
				res := &model.Bookmark{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000001",
						CreatedAt: fixture.TestTime,
						UpdatedAt: fixture.TestTime,
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://example.com/testuser001",
					Code:        "XXXX000001",
					Description: "Bookmark for Test User 1 - record 1",
				}
				svcMock.On("CreateBookmark", ctx, url, description, userID).Return(res, nil)
				return svcMock
			},

			setupMockCache: func(ctx context.Context, userID string) *mock_cache.DB {
				cacheMock := mock_cache.NewDB(t)
				groupKey := fmt.Sprintf(bookmark.ListBookmarksCacheGroupKey, userID)
				cacheMock.On("DelCacheData", ctx, groupKey).Return(nil)
				return cacheMock
			},

			url:         "https://example.com/testuser001",
			description: "Bookmark for Test User 1 - record 1",
			userID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",

			expectedBookmark: &model.Bookmark{
				Base: model.Base{
					ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000001",
					CreatedAt: fixture.TestTime,
					UpdatedAt: fixture.TestTime,
				},
				UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
				URL:         "https://example.com/testuser001",
				Code:        "XXXX000001",
				Description: "Bookmark for Test User 1 - record 1",
			},
		},
		{
			name: "Create bookmark with cache invalidation failed - cache error",

			setupMockSvc: func(ctx context.Context, url, description, userID string) *mocks.Service {
				svcMock := mocks.NewService(t)
				return svcMock
			},

			setupMockCache: func(ctx context.Context, userID string) *mock_cache.DB {
				cacheMock := mock_cache.NewDB(t)
				groupKey := fmt.Sprintf(bookmark.ListBookmarksCacheGroupKey, userID)
				cacheMock.On("DelCacheData", ctx, groupKey).Return(assert.AnError)
				return cacheMock
			},

			url:         "https://example.com/testuser001",
			description: "Bookmark for Test User 1 - record 1",
			userID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",

			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			svcMock := tc.setupMockSvc(ctx, tc.url, tc.description, tc.userID)
			cacheMock := tc.setupMockCache(ctx, tc.userID)

			service := bookmark.NewBookmarkServiceWithCache(svcMock, cacheMock)
			bookmarkRes, err := service.CreateBookmark(ctx, tc.url, tc.description, tc.userID)
			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBookmark, bookmarkRes)
		})
	}
}
