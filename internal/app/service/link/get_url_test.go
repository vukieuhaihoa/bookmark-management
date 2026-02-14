package link

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	mockBookmarkRepo "github.com/vukieuhaihoa/bookmark-management/internal/app/repository/bookmark/mocks"
	mockLinkRepo "github.com/vukieuhaihoa/bookmark-management/internal/app/repository/link/mocks"
)

func TestService_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockLinkRepo     func(ctx context.Context, code string) *mockLinkRepo.Repository
		setupMockBookmarkRepo func(ctx context.Context, code string) *mockBookmarkRepo.Repository

		inputURLCode string

		expectedOriginalURL string
		expectedError       error
	}{
		{
			name: "Get URL successfully - redis link old format from v1 and v2",

			setupMockLinkRepo: func(ctx context.Context, code string) *mockLinkRepo.Repository {
				repoMock := mockLinkRepo.NewRepository(t)
				repoMock.On("GetURL", ctx, code).Return("https://example.com", nil)
				return repoMock
			},

			setupMockBookmarkRepo: func(ctx context.Context, code string) *mockBookmarkRepo.Repository {
				return mockBookmarkRepo.NewRepository(t)
			},

			inputURLCode: "abcd1234",

			expectedOriginalURL: "https://example.com",
			expectedError:       nil,
		},
		{
			name: "Get URL successfully - redis link new format from v3",

			setupMockLinkRepo: func(ctx context.Context, code string) *mockLinkRepo.Repository {
				repoMock := mockLinkRepo.NewRepository(t)
				repoMock.On("GetURL", ctx, code).Return("https://example.com", nil)
				return repoMock
			},

			setupMockBookmarkRepo: func(ctx context.Context, code string) *mockBookmarkRepo.Repository {
				return mockBookmarkRepo.NewRepository(t)
			},

			inputURLCode: model.RedisShortenPrefix + "abcd1234",

			expectedOriginalURL: "https://example.com",
			expectedError:       nil,
		},
		{
			name: "URL code not found - redis link old format from v1 and v2",

			setupMockLinkRepo: func(ctx context.Context, code string) *mockLinkRepo.Repository {
				repoMock := mockLinkRepo.NewRepository(t)
				repoMock.On("GetURL", ctx, code).Return("", redis.Nil)
				return repoMock
			},

			setupMockBookmarkRepo: func(ctx context.Context, code string) *mockBookmarkRepo.Repository {
				return mockBookmarkRepo.NewRepository(t)
			},

			inputURLCode: "leetcode",

			expectedOriginalURL: "",
			expectedError:       ErrCodeNotFound,
		},
		{
			name: "URL code not found - redis link new format from v3",

			setupMockLinkRepo: func(ctx context.Context, code string) *mockLinkRepo.Repository {
				repoMock := mockLinkRepo.NewRepository(t)
				repoMock.On("GetURL", ctx, code).Return("", redis.Nil)
				return repoMock
			},

			setupMockBookmarkRepo: func(ctx context.Context, code string) *mockBookmarkRepo.Repository {
				return mockBookmarkRepo.NewRepository(t)
			},

			inputURLCode: model.RedisShortenPrefix + "leetcode",

			expectedOriginalURL: "",
			expectedError:       ErrCodeNotFound,
		},
		{
			name: "Get URL Successfully - bookmark link old format from v1 and v2",

			setupMockLinkRepo: func(ctx context.Context, code string) *mockLinkRepo.Repository {
				repoMock := mockLinkRepo.NewRepository(t)
				return repoMock
			},

			setupMockBookmarkRepo: func(ctx context.Context, code string) *mockBookmarkRepo.Repository {
				repoMock := mockBookmarkRepo.NewRepository(t)
				repoMock.On("GetBookmarkByCode", ctx, code).Return(&model.Bookmark{
					URL: "https://example.com/bookmark",
				}, nil)
				return repoMock
			},

			inputURLCode: "mybookmark",

			expectedOriginalURL: "https://example.com/bookmark",
			expectedError:       nil,
		},
		{
			name: "Get URL Successfully - bookmark link new format from v3",

			setupMockLinkRepo: func(ctx context.Context, code string) *mockLinkRepo.Repository {
				repoMock := mockLinkRepo.NewRepository(t)
				return repoMock
			},

			setupMockBookmarkRepo: func(ctx context.Context, code string) *mockBookmarkRepo.Repository {
				repoMock := mockBookmarkRepo.NewRepository(t)
				repoMock.On("GetBookmarkByCodeShortenEncoded", ctx, code).Return(&model.Bookmark{
					URL: "https://example.com/bookmark",
				}, nil)
				return repoMock
			},
			inputURLCode:        model.BookmarkShortenPrefix + "1AE",
			expectedOriginalURL: "https://example.com/bookmark",
			expectedError:       nil,
		},
		{
			name: "Get URL fails - wrong format - not found",

			setupMockLinkRepo: func(ctx context.Context, code string) *mockLinkRepo.Repository {
				repoMock := mockLinkRepo.NewRepository(t)
				return repoMock
			},

			setupMockBookmarkRepo: func(ctx context.Context, code string) *mockBookmarkRepo.Repository {
				repoMock := mockBookmarkRepo.NewRepository(t)
				// repoMock.On("GetBookmarkByCode", ctx, code).Return(nil, dbutils.ErrRecordNotFoundType)
				return repoMock
			},

			inputURLCode: "unknownbookmark",

			expectedOriginalURL: "",
			expectedError:       ErrCodeNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			linkRepoMock := tc.setupMockLinkRepo(ctx, tc.inputURLCode)
			bookmarkRepoMock := tc.setupMockBookmarkRepo(ctx, tc.inputURLCode)
			service := NewLinkService(linkRepoMock, nil, bookmarkRepoMock)

			originalURL, err := service.GetURL(ctx, tc.inputURLCode)

			assert.Equal(t, tc.expectedOriginalURL, originalURL)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
