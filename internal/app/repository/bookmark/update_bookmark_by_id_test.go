package bookmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

func TestRepository_UpdateBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputID           string
		inputUserID       string
		setupDB           func(t *testing.T) *gorm.DB
		inputBookmarkData *model.Bookmark

		expectedError error
		verifyFunc    func(t *testing.T, db *gorm.DB, id string, expectedData *model.Bookmark)
	}{
		{
			name: "Update bookmark by ID successfully",

			inputID:     "a1b2c3d4-e5f6-7890-abcd-ef0000000005",
			inputUserID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			inputBookmarkData: &model.Bookmark{
				URL:         "https://updated-example.com",
				Description: "This is an updated description.",
			},

			verifyFunc: func(t *testing.T, db *gorm.DB, id string, expectedData *model.Bookmark) {
				var updatedBookmark model.Bookmark
				err := db.Model(&model.Bookmark{}).First(&updatedBookmark, "id = ?", id).Error
				assert.NoError(t, err)

				assert.Equal(t, expectedData.URL, updatedBookmark.URL)
				assert.Equal(t, expectedData.Description, updatedBookmark.Description)
			},
		},
		{
			name: "Update bookmark by ID failed - bookmark not found",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
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
			name: "Update bookmark by ID failed - invalid user ID",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			inputID:     "a1b2c3d4-e5f6-7890-abcd-ef0000000005",
			inputUserID: "invalid-user-id",

			inputBookmarkData: &model.Bookmark{
				URL:         "https://updated-example.com",
				Description: "This is an updated description.",
			},

			expectedError: dbutils.ErrRecordNotFoundType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewBookmarkRepository(db)

			err := repo.UpdateBookmarkByID(ctx, tc.inputID, tc.inputUserID, tc.inputBookmarkData)
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			tc.verifyFunc(t, db, tc.inputID, tc.inputBookmarkData)
		})
	}
}
