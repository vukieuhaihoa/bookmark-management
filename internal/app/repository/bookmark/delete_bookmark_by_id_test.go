package bookmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

func TestRepository_DeleteBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputID     string
		inputUserID string
		setupDB     func(t *testing.T) *gorm.DB

		expectedError error
		verifyFunc    func(t *testing.T, db *gorm.DB, id, userID string)
	}{
		{
			name: "Delete bookmark by ID successfully",

			inputID:     "a1b2c3d4-e5f6-7890-abcd-ef0000000006",
			inputUserID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			verifyFunc: func(t *testing.T, db *gorm.DB, id, userID string) {
				var bookmark model.Bookmark
				err := db.Model(&model.Bookmark{}).First(&bookmark, "id = ? AND user_id = ?", id, userID).Error
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "Delete bookmark by ID failed - bookmark not found",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			inputID:       "non-existent-id",
			inputUserID:   "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			expectedError: dbutils.ErrRecordNotFoundType,
		},
		{
			name: "Delete bookmark by ID failed - bookmark does not belong to user",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputID:       "a1b2c3d4-e5f6-7890-abcd-ef0000000006",
			inputUserID:   "different-user-id",
			expectedError: dbutils.ErrRecordNotFoundType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewBookmarkRepository(db)

			err := repo.DeleteBookmarkByID(ctx, tc.inputID, tc.inputUserID)
			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}

			assert.NoError(t, err)

			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db, tc.inputID, tc.inputUserID)
			}
		})
	}
}
