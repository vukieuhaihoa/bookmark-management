package link

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	mockRepo "github.com/vukieuhaihoa/bookmark-management/internal/app/repository/link/mocks"
)

func TestService_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRepo func(ctx context.Context) *mockRepo.Repository

		inputURLCode string

		expectedOriginalURL string
		expectedError       error
	}{
		{
			name: "Get URL successfully",

			setupMockRepo: func(ctx context.Context) *mockRepo.Repository {
				repoMock := mockRepo.NewRepository(t)
				repoMock.On("GetURL", ctx, "abcd1234").Return("https://example.com", nil)
				return repoMock
			},

			inputURLCode: "abcd1234",

			expectedOriginalURL: "https://example.com",
			expectedError:       nil,
		},
		{
			name: "URL code not found",

			setupMockRepo: func(ctx context.Context) *mockRepo.Repository {
				repoMock := mockRepo.NewRepository(t)
				repoMock.On("GetURL", ctx, "unknowncode").Return("", redis.Nil)
				return repoMock
			},

			inputURLCode: "unknowncode",

			expectedOriginalURL: "",
			expectedError:       ErrCodeNotFound,
		},
		{
			name: "Repository error",

			setupMockRepo: func(ctx context.Context) *mockRepo.Repository {
				repoMock := mockRepo.NewRepository(t)
				repoMock.On("GetURL", ctx, "errorcode").Return("", assert.AnError)
				return repoMock
			},

			inputURLCode: "errorcode",

			expectedOriginalURL: "",
			expectedError:       assert.AnError,
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repoMock := tc.setupMockRepo(ctx)
			service := NewLinkService(repoMock, nil)

			originalURL, err := service.GetURL(ctx, tc.inputURLCode)

			assert.Equal(t, tc.expectedOriginalURL, originalURL)
			assert.Equal(t, tc.expectedError, err)

			repoMock.AssertExpectations(t)
		})
	}
}
