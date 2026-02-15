package bookmark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	mockBookmarkRepo "github.com/vukieuhaihoa/bookmark-management/internal/app/repository/bookmark/mocks"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
)

func TestService_ListBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRepo func(ctx context.Context) *mockBookmarkRepo.Repository

		inputUserID string
		inputOpts   *common.QueryOptions

		expectedOutput []*model.Bookmark
		expectedError  error
	}{
		{
			name: "List bookmarks successfully",

			setupMockRepo: func(ctx context.Context) *mockBookmarkRepo.Repository {
				repoMock := mockBookmarkRepo.NewRepository(t)
				repoMock.On("ListBookmarks", ctx, "4d9326d6-980c-4c62-9709-dbc70a82cbfe", &common.QueryOptions{
					Paging: common.Paging{
						Page:  1,
						Limit: 10,
					},
					Sorting: []common.SortedField{
						{
							Field:     "created_at",
							Direction: "DESC",
						},
					},
				}).Return([]*model.Bookmark{
					{
						Base: model.Base{
							ID: "a1b2c3d4-e5f6-7890-abcd-ef0000000001",
						},
						UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
						URL:                "https://example.com/testuser001",
						Description:        "Bookmark for Test User 1 - record 1",
						Code:               "XXXX000001",
						CodeShortenEncoded: "p_1",
					},
					{
						Base: model.Base{
							ID: "a1b2c3d4-e5f6-7890-abcd-ef0000000002",
						},
						UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
						URL:                "https://example.com/testuser002",
						Description:        "Bookmark for Test User 1 - record 2",
						Code:               "XXXX000002",
						CodeShortenEncoded: "p_2",
					},
				}, nil)
				return repoMock
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
						Direction: "DESC",
					},
				},
			},

			expectedOutput: []*model.Bookmark{
				{
					Base: model.Base{
						ID: "a1b2c3d4-e5f6-7890-abcd-ef0000000001",
					},
					UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:                "https://example.com/testuser001",
					Description:        "Bookmark for Test User 1 - record 1",
					Code:               "XXXX000001",
					CodeShortenEncoded: "p_1",
				},
				{
					Base: model.Base{
						ID: "a1b2c3d4-e5f6-7890-abcd-ef0000000002",
					},
					UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:                "https://example.com/testuser002",
					Description:        "Bookmark for Test User 1 - record 2",
					Code:               "XXXX000002",
					CodeShortenEncoded: "p_2",
				},
			},

			expectedError: nil,
		},
		{
			name: "List bookmarks failed - repository error",

			setupMockRepo: func(ctx context.Context) *mockBookmarkRepo.Repository {
				repoMock := mockBookmarkRepo.NewRepository(t)
				repoMock.On("ListBookmarks", ctx, "4d9326d6-980c-4c62-9709-dbc70a82cbfe", &common.QueryOptions{
					Paging: common.Paging{
						Page:  1,
						Limit: 10,
					},
					Sorting: []common.SortedField{
						{
							Field:     "created_at",
							Direction: "DESC",
						},
					},
				}).Return(nil, assert.AnError)
				return repoMock
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
						Direction: "DESC",
					},
				},
			},

			expectedOutput: nil,
			expectedError:  assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repoMock := tc.setupMockRepo(ctx)
			service := NewBookmarkService(repoMock)
			res, err := service.ListBookmarks(ctx, tc.inputUserID, tc.inputOpts)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}
