package bookmark

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

func TestRepository_ListBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockDB func(t *testing.T) *gorm.DB

		inputUserID string
		inputOpts   *common.QueryOptions

		expectedOutput []*model.Bookmark
		expectedError  error
	}{
		{
			name: "list bookmarks successfully - created_at desc sorting, page 1, limit 2",

			setupMockDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			inputUserID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			inputOpts: &common.QueryOptions{
				Paging: common.Paging{
					Page:  1,
					Limit: 4,
				},
				Sorting: []common.SortedField{
					{
						Field:     "created_at",
						Direction: "DESC",
					},
				},
			},

			expectedOutput: []*model.Bookmark{
				{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000009",
						CreatedAt: fixture.TestTime.Add(5 * time.Hour),
						UpdatedAt: fixture.TestTime.Add(5 * time.Hour),
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://redis.io/docs/manual/data-types/",
					Code:        "XXXX000006",
					Description: "Redis data types documentation",
				},
				{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000008",
						CreatedAt: fixture.TestTime.Add(4 * time.Hour),
						UpdatedAt: fixture.TestTime.Add(4 * time.Hour),
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://db-tutorials.dev/postgresql-indexing",
					Code:        "XXXX000005",
					Description: "Learn PostgreSQL indexing basics",
				},
				{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000007",
						CreatedAt: fixture.TestTime.Add(3 * time.Hour),
						UpdatedAt: fixture.TestTime.Add(3 * time.Hour),
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://golang.dev/blog/clean-architecture",
					Code:        "XXXX000004",
					Description: "Go backend best practices article",
				},
				{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000006",
						CreatedAt: fixture.TestTime.Add(2 * time.Hour),
						UpdatedAt: fixture.TestTime.Add(2 * time.Hour),
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://example.com/testuser003",
					Code:        "XXXX000003",
					Description: "Bookmark for Test User 1 - record 3",
				},
			},

			expectedError: nil,
		},
		{
			name: "list bookmarks successfully - created_at desc sorting, page 1, limit 4",

			setupMockDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			inputUserID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			inputOpts: &common.QueryOptions{
				Paging: common.Paging{
					Page:  1,
					Limit: 2,
				},
				Sorting: []common.SortedField{
					{
						Field:     "created_at",
						Direction: "desc",
					},
				},
			},

			expectedOutput: []*model.Bookmark{
				{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000009",
						CreatedAt: fixture.TestTime.Add(5 * time.Hour),
						UpdatedAt: fixture.TestTime.Add(5 * time.Hour),
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://redis.io/docs/manual/data-types/",
					Code:        "XXXX000006",
					Description: "Redis data types documentation",
				},
				{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000008",
						CreatedAt: fixture.TestTime.Add(4 * time.Hour),
						UpdatedAt: fixture.TestTime.Add(4 * time.Hour),
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://db-tutorials.dev/postgresql-indexing",
					Code:        "XXXX000005",
					Description: "Learn PostgreSQL indexing basics",
				},
			},

			expectedError: nil,
		},
		{
			name: "list bookmarks successfully - created_at desc sorting, page 2, limit 2",

			setupMockDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			inputUserID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			inputOpts: &common.QueryOptions{
				Paging: common.Paging{
					Page:  2,
					Limit: 2,
				},
				Sorting: []common.SortedField{
					{
						Field:     "created_at",
						Direction: "desc",
					},
				},
			},

			expectedOutput: []*model.Bookmark{
				{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000007",
						CreatedAt: fixture.TestTime.Add(3 * time.Hour),
						UpdatedAt: fixture.TestTime.Add(3 * time.Hour),
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://golang.dev/blog/clean-architecture",
					Code:        "XXXX000004",
					Description: "Go backend best practices article",
				},
				{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000006",
						CreatedAt: fixture.TestTime.Add(2 * time.Hour),
						UpdatedAt: fixture.TestTime.Add(2 * time.Hour),
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://example.com/testuser003",
					Code:        "XXXX000003",
					Description: "Bookmark for Test User 1 - record 3",
				},
			},

			expectedError: nil,
		},
		{
			name: "list bookmarks successfully - created_at ASC, description DESC sorting, page 1, limit 2",

			setupMockDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			inputUserID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			inputOpts: &common.QueryOptions{
				Paging: common.Paging{
					Page:  1,
					Limit: 2,
				},
				Sorting: []common.SortedField{
					{
						Field:     "created_at",
						Direction: "ASC",
					},
					{
						Field:     "description",
						Direction: "DESC",
					},
				},
			},

			expectedOutput: []*model.Bookmark{

				{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000005",
						CreatedAt: fixture.TestTime,
						UpdatedAt: fixture.TestTime,
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://example.com/testuser001",
					Code:        "XXXX000002",
					Description: "Bookmark for Test User 1 - record 2",
				},
				{
					Base: model.Base{
						ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000004",
						CreatedAt: fixture.TestTime,
						UpdatedAt: fixture.TestTime,
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://example.com/testuser001",
					Code:        "XXXX000001",
					Description: "Bookmark for Test User 1 - record 1",
				},
			},

			expectedError: nil,
		},
		{
			name: "list bookmarks successfully - user not found",
			setupMockDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			inputUserID: "non-existent-user-id",
			inputOpts: &common.QueryOptions{
				Paging: common.Paging{
					Page:  1,
					Limit: 10,
				},
			},
			expectedOutput: []*model.Bookmark{},
			expectedError:  nil,
		},
		{
			name: "sorting by invalid field",

			setupMockDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},

			inputUserID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			inputOpts: &common.QueryOptions{
				Paging: common.Paging{
					Page:  1,
					Limit: 10,
				},
				Sorting: []common.SortedField{
					{
						Field:     "invalid_field",
						Direction: "ASC",
					},
				},
			},

			expectedOutput: nil,
			expectedError:  dbutils.ErrInvalidSortField,
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupMockDB(t)
			repo := NewBookmarkRepository(db)

			output, err := repo.ListBookmarks(ctx, tc.inputUserID, tc.inputOpts)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}
