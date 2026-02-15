package bookmark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/repository/bookmark/mocks"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

func TestService_DeleteBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRepo func(ctx context.Context) *mocks.Repository

		inputID     string
		inputUserID string

		expectedError error
	}{
		{
			name: "Delete bookmark by ID successfully",

			setupMockRepo: func(ctx context.Context) *mocks.Repository {
				repoMock := mocks.NewRepository(t)
				repoMock.On("DeleteBookmarkByID", ctx, "a1b2c3d4-e5f6-7890-abcd-ef0000000006", "4d9326d6-980c-4c62-9709-dbc70a82cbfe").Return(nil)
				return repoMock
			},

			inputID:     "a1b2c3d4-e5f6-7890-abcd-ef0000000006",
			inputUserID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
		},
		{
			name: "Delete bookmark by ID failed - bookmark not found",

			setupMockRepo: func(ctx context.Context) *mocks.Repository {
				repoMock := mocks.NewRepository(t)
				repoMock.On("DeleteBookmarkByID", ctx, "non-existent-id", "4d9326d6-980c-4c62-9709-dbc70a82cbfe").Return(dbutils.ErrRecordNotFoundType)
				return repoMock
			},

			inputID:       "non-existent-id",
			inputUserID:   "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			expectedError: dbutils.ErrRecordNotFoundType,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			repoMock := tc.setupMockRepo(ctx)
			service := NewBookmarkService(repoMock)
			err := service.DeleteBookmarkByID(ctx, tc.inputID, tc.inputUserID)
			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
