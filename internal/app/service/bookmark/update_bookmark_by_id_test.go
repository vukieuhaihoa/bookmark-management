package bookmark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/repository/bookmark/mocks"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

func TestService_UpdateBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputID     string
		inputUserID string

		setupRepo         func(ctx context.Context, id, userID string, bookmarkData *model.Bookmark) *mocks.Repository
		inputBookmarkData *model.Bookmark
		expectedError     error
	}{
		{
			name: "Update bookmark by ID successfully",

			setupRepo: func(ctx context.Context, id, userID string, bookmarkData *model.Bookmark) *mocks.Repository {
				repo := &mocks.Repository{}
				repo.On("UpdateBookmarkByID", ctx, id, userID, bookmarkData).Return(nil)
				return repo
			},

			inputID:     "a1b2c3d4-e5f6-7890-abcd-ef0000000005",
			inputUserID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",

			inputBookmarkData: &model.Bookmark{
				URL:         "https://updated-example.com",
				Description: "This is an updated description.",
			},

			expectedError: nil,
		},
		{
			name: "Update bookmark by ID failed - bookmark not found",

			setupRepo: func(ctx context.Context, id, userID string, bookmarkData *model.Bookmark) *mocks.Repository {
				repo := &mocks.Repository{}
				repo.On("UpdateBookmarkByID", ctx, id, userID, bookmarkData).Return(dbutils.ErrRecordNotFoundType)
				return repo
			},

			inputID:     "non-existent-id",
			inputUserID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",

			inputBookmarkData: &model.Bookmark{
				URL:         "https://nonexistent.com",
				Description: "This bookmark does not exist.",
			},

			expectedError: dbutils.ErrRecordNotFoundType,
		},
		{
			name: "Update bookmark by ID failed - invalid user id",
			setupRepo: func(ctx context.Context, id, userID string, bookmarkData *model.Bookmark) *mocks.Repository {
				repo := &mocks.Repository{}
				repo.On("UpdateBookmarkByID", ctx, id, userID, bookmarkData).Return(dbutils.ErrRecordNotFoundType)
				return repo
			},

			inputID:     "a1b2c3d4-e5f6-7890-abcd-ef0000000005",
			inputUserID: "invalid-user-id",

			inputBookmarkData: &model.Bookmark{
				URL:         "https://example.com",
				Description: "This is a description.",
			},

			expectedError: dbutils.ErrRecordNotFoundType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			repo := tc.setupRepo(ctx, tc.inputID, tc.inputUserID, tc.inputBookmarkData)
			service := NewBookmarkService(repo)

			err := service.UpdateBookmarkByID(ctx, tc.inputID, tc.inputUserID, tc.inputBookmarkData)
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
