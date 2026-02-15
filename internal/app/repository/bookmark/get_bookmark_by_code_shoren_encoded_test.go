package bookmark

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

func TestRepository_GetBookmarkByCodeShortenEncoded(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputCode string

		setupDB        func(t *testing.T) *gorm.DB
		expectedError  error
		expectedOutput *model.Bookmark
	}{
		{
			name: "Get bookmark by code successfully",

			inputCode: "p_9",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			expectedOutput: &model.Bookmark{
				Base: model.Base{
					ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000009",
					CreatedAt: fixture.TestTime.Add(5 * time.Hour),
					UpdatedAt: fixture.TestTime.Add(5 * time.Hour),
				},
				UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
				URL:                "https://redis.io/docs/manual/data-types/",
				Description:        "Redis data types documentation",
				Code:               "AAAA000009",
				CodeShorten:        9,
				CodeShortenEncoded: "p_9",
			},
		},
		{
			name: "Get bookmark by code failed - bookmark not found",

			inputCode: "nonexistentcode",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			expectedError: dbutils.ErrRecordNotFoundType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := tc.setupDB(t)
			repo := NewBookmarkRepository(db)

			output, err := repo.GetBookmarkByCodeShortenEncoded(context.Background(), tc.inputCode)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}
