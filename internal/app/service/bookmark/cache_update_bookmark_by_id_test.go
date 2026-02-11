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
)

func TestBookmarkServiceWithCache_UpdateBookmarkByID(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string

		id              string
		userID          string
		updatedBookmark *model.Bookmark
		setupMockSvc    func(ctx context.Context, id, userID string, updatedBookmark *model.Bookmark) *mocks.Service
		setupMockCache  func(ctx context.Context, userID string) *mock_cache.DB

		expectedError error
	}{
		{
			name: "Update bookmark by ID with cache invalidation successfully",

			id:     "a1b2c3d4-e5f6-7890-abcd-ef0000000001",
			userID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			updatedBookmark: &model.Bookmark{
				Description: "Updated description",
			},

			setupMockSvc: func(ctx context.Context, id, userID string, updatedBookmark *model.Bookmark) *mocks.Service {
				svcMock := mocks.NewService(t)
				svcMock.On("UpdateBookmarkByID", ctx, id, userID, updatedBookmark).Return(nil)
				return svcMock
			},

			setupMockCache: func(ctx context.Context, userID string) *mock_cache.DB {
				cacheMock := mock_cache.NewDB(t)
				groupKey := fmt.Sprintf(bookmark.ListBookmarksCacheGroupKey, userID)
				cacheMock.On("DelCacheData", ctx, groupKey).Return(nil)
				return cacheMock
			},
		},
		{
			name: "Update bookmark by ID with cache invalidation failed - cache error",

			id:     "a1b2c3d4-e5f6-7890-abcd-ef0000000001",
			userID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			updatedBookmark: &model.Bookmark{
				Description: "Updated description",
			},

			setupMockSvc: func(ctx context.Context, id, userID string, updatedBookmark *model.Bookmark) *mocks.Service {
				svcMock := mocks.NewService(t)
				return svcMock
			},

			setupMockCache: func(ctx context.Context, userID string) *mock_cache.DB {
				cacheMock := mock_cache.NewDB(t)
				groupKey := fmt.Sprintf(bookmark.ListBookmarksCacheGroupKey, userID)
				cacheMock.On("DelCacheData", ctx, groupKey).Return(assert.AnError)
				return cacheMock
			},

			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			svcMock := tc.setupMockSvc(ctx, tc.id, tc.userID, tc.updatedBookmark)
			cacheMock := tc.setupMockCache(ctx, tc.userID)
			serviceWithCache := bookmark.NewBookmarkServiceWithCache(svcMock, cacheMock)

			err := serviceWithCache.UpdateBookmarkByID(ctx, tc.id, tc.userID, tc.updatedBookmark)
			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
