package bookmark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	mockBookmarkRepo "github.com/vukieuhaihoa/bookmark-management/internal/app/repository/bookmark/mocks"
	mockCodeGen "github.com/vukieuhaihoa/bookmark-management/pkg/utils/mocks"
)

func TestService_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockCodeGen func(t *testing.T) *mockCodeGen.CodeGenerator
		setupMockRepo    func(ctx context.Context) *mockBookmarkRepo.Repository

		inputURL         string
		inputDescription string
		inputUserID      string

		expectedOutput *model.Bookmark
		expectedError  error
	}{
		{
			name: "Create bookmark successfully",

			setupMockCodeGen: func(t *testing.T) *mockCodeGen.CodeGenerator {
				codeGenMock := mockCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", DEFAULT_CODE_LENGTH).Return("abcdefghij", nil)
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockBookmarkRepo.Repository {
				repoMock := mockBookmarkRepo.NewRepository(t)
				repoMock.On("CreateBookmark", ctx, &model.Bookmark{
					URL:         "https://example.com",
					Description: "Example Website",
					UserID:      "user-123",
					Code:        "abcdefghij",
				}).Return(&model.Bookmark{
					Base: model.Base{
						ID: "bookmark-456",
					},
					URL:         "https://example.com",
					Description: "Example Website",
					UserID:      "user-123",
					Code:        "abcdefghij",
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
				URL:         "https://example.com",
				Description: "Example Website",
				UserID:      "user-123",
				Code:        "abcdefghij",
			},

			expectedError: nil,
		},
		{
			name: "Fail to create bookmark due to code generation error",

			setupMockCodeGen: func(t *testing.T) *mockCodeGen.CodeGenerator {
				codeGenMock := mockCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", DEFAULT_CODE_LENGTH).Return("", assert.AnError)
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockBookmarkRepo.Repository {
				return mockBookmarkRepo.NewRepository(t)
			},

			inputURL:         "https://example.com",
			inputDescription: "Example Website",
			inputUserID:      "user-123",

			expectedOutput: nil,
			expectedError:  assert.AnError,
		},
		{
			name: "Fail to create bookmark due to repository error",

			setupMockCodeGen: func(t *testing.T) *mockCodeGen.CodeGenerator {
				codeGenMock := mockCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", DEFAULT_CODE_LENGTH).Return("abcdefghij", nil)
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockBookmarkRepo.Repository {
				repoMock := mockBookmarkRepo.NewRepository(t)
				repoMock.On("CreateBookmark", ctx, &model.Bookmark{
					URL:         "https://example.com",
					Description: "Example Website",
					UserID:      "user-123",
					Code:        "abcdefghij",
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
			codeGenMock := tc.setupMockCodeGen(t)
			repoMock := tc.setupMockRepo(ctx)

			service := NewBookmarkService(repoMock, codeGenMock)

			res, err := service.CreateBookmark(ctx, tc.inputURL, tc.inputDescription, tc.inputUserID)
			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}
