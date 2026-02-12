package bookmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

func TestRepository_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB       func(t *testing.T) *gorm.DB
		inputBookmark *model.Bookmark

		expectedOutput *model.Bookmark
		expectedError  error
	}{
		{
			name: "Create bookmark successfully",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			inputBookmark: &model.Bookmark{
				Base: model.Base{
					ID: "a1b2c3d4-e5f6-7890-abcd-ef0000000077",
				},
				UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
				URL:         "https://example.com/newbookmark",
				Description: "New bookmark for Test User 1",
			},

			expectedOutput: &model.Bookmark{
				Base: model.Base{
					ID: "a1b2c3d4-e5f6-7890-abcd-ef0000000077",
				},
				UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
				URL:         "https://example.com/newbookmark",
				Description: "New bookmark for Test User 1",
			},
		},

		{
			name: "Fail to create bookmark with non-existing user",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			inputBookmark: &model.Bookmark{
				Base: model.Base{
					ID: "c1d2e3f4-g5h6-7890-abcd-ef0000000099",
				},
				UserID:      "non-existing-user-id",
				URL:         "https://example.com/invaliduserbookmark",
				Description: "Bookmark with non-existing user",
			},

			expectedError: dbutils.ErrForeignKeyType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewBookmarkRepository(db)

			output, err := repo.CreateBookmark(ctx, tc.inputBookmark)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}

			helperCompareBookmarks(t, db, output)
		})
	}
}

func helperCompareBookmarks(t *testing.T, db *gorm.DB, actual *model.Bookmark) {
	expected := &model.Bookmark{}
	err := db.First(expected, "id = ?", actual.ID).Error
	assert.Nil(t, err)

	// Verify timestamps are automatically set by GORM
	assert.False(t, actual.CreatedAt.IsZero(), "CreatedAt should be automatically set")
	assert.False(t, actual.UpdatedAt.IsZero(), "UpdatedAt should be automatically set")

	assert.NoError(t, err)
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.URL, actual.URL)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.CodeShorten, actual.CodeShorten)
	assert.Equal(t, expected.CodeShortenEncoded, actual.CodeShortenEncoded)
}
