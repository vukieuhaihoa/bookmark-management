package bookmark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	mockBookmarkRepo "github.com/vukieuhaihoa/bookmark-management/internal/app/repository/bookmark/mocks"
)

func TestService_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRepo func(ctx context.Context) *mockBookmarkRepo.Repository

		inputURL         string
		inputDescription string
		inputUserID      string

		expectedOutput *model.Bookmark
		expectedError  error
	}{
		{
			name: "Create bookmark successfully",

			setupMockRepo: func(ctx context.Context) *mockBookmarkRepo.Repository {
				repoMock := mockBookmarkRepo.NewRepository(t)
				repoMock.On("CreateBookmark", ctx, &model.Bookmark{
					URL:         "https://example.com",
					Description: "Example Website",
					UserID:      "user-123",
				}).Return(&model.Bookmark{
					Base: model.Base{
						ID: "bookmark-456",
					},
					URL:                "https://example.com",
					Description:        "Example Website",
					UserID:             "user-123",
					CodeShortenEncoded: "p_1A",
				}, nil)
				return repoMock
			},

			inputURL:         "https://example.com",
			inputDescription: "Example Website",
			inputUserID:      "user-123",

			expectedOutput: &model.Bookmark{
				Base: model.Base{
					ID: "bookmark-456",
				},
				URL:                "https://example.com",
				Description:        "Example Website",
				UserID:             "user-123",
				CodeShortenEncoded: "p_1A",
			},

			expectedError: nil,
		},
		{
			name: "Fail to create bookmark due to repository error",

			setupMockRepo: func(ctx context.Context) *mockBookmarkRepo.Repository {
				repoMock := mockBookmarkRepo.NewRepository(t)
				repoMock.On("CreateBookmark", ctx, &model.Bookmark{
					URL:         "https://example.com",
					Description: "Example Website",
					UserID:      "user-123",
				}).Return(nil, assert.AnError)
				return repoMock
			},

			inputURL:         "https://example.com",
			inputDescription: "Example Website",
			inputUserID:      "user-123",

			expectedOutput: nil,
			expectedError:  assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			repoMock := tc.setupMockRepo(ctx)

			service := NewBookmarkService(repoMock)

			res, err := service.CreateBookmark(ctx, tc.inputURL, tc.inputDescription, tc.inputUserID)
			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}
