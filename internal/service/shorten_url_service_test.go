package service

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	mockRepo "github.com/vukieuhaihoa/bookmark-management/internal/repository/mocks"
	mockRandomCodeGen "github.com/vukieuhaihoa/bookmark-management/pkg/stringutils/mocks"
)

func TestShortenURL_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRandomCodeGen func() *mockRandomCodeGen.CodeGenerator
		setupMockRepo          func(ctx context.Context) *mockRepo.UrlStorage

		inputOriginalURL string
		inputExpireIn    int

		expectedCode  string
		expectedError error
	}{
		{
			name: "Shorten URL successfully",

			setupMockRandomCodeGen: func() *mockRandomCodeGen.CodeGenerator {
				codeGenMock := mockRandomCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", 8).Return("abcd1234", nil)
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockRepo.UrlStorage {
				repoMock := mockRepo.NewUrlStorage(t)
				repoMock.On("GetURL", ctx, "abcd1234").Return("", redis.Nil)
				repoMock.On("StoreURL", ctx, "abcd1234", "https://example.com", 3600).Return(nil)
				return repoMock
			},

			inputOriginalURL: "https://example.com",
			inputExpireIn:    3600,

			expectedCode:  "abcd1234",
			expectedError: nil,
		},
		{
			name: "Fail to generate code",

			setupMockRandomCodeGen: func() *mockRandomCodeGen.CodeGenerator {
				codeGenMock := mockRandomCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", 8).Return("", assert.AnError)
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockRepo.UrlStorage {
				return mockRepo.NewUrlStorage(t)
			},

			inputOriginalURL: "https://example.com",
			inputExpireIn:    3600,

			expectedCode:  "",
			expectedError: assert.AnError,
		},
		{
			name: "Fail to store URL",

			setupMockRandomCodeGen: func() *mockRandomCodeGen.CodeGenerator {
				codeGenMock := mockRandomCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", 8).Return("abcd1234", nil)
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockRepo.UrlStorage {
				repoMock := mockRepo.NewUrlStorage(t)
				repoMock.On("GetURL", ctx, "abcd1234").Return("", redis.Nil)
				repoMock.On("StoreURL", ctx, "abcd1234", "https://example.com", 3600).Return(assert.AnError)
				return repoMock
			},

			inputOriginalURL: "https://example.com",
			inputExpireIn:    3600,

			expectedCode:  "",
			expectedError: assert.AnError,
		},
		{
			name: "Retry on code collision and succeed",

			setupMockRandomCodeGen: func() *mockRandomCodeGen.CodeGenerator {
				codeGenMock := mockRandomCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", 8).Return("abcd1234", nil).Once()
				codeGenMock.On("GenerateCode", 8).Return("efgh5678", nil).Once()
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockRepo.UrlStorage {
				repoMock := mockRepo.NewUrlStorage(t)
				repoMock.On("GetURL", ctx, "abcd1234").Return("https://collision.com", nil).Once()
				repoMock.On("GetURL", ctx, "efgh5678").Return("", redis.Nil).Once()
				repoMock.On("StoreURL", ctx, "efgh5678", "https://example.com", 3600).Return(nil).Once()
				return repoMock
			},

			inputOriginalURL: "https://example.com",
			inputExpireIn:    3600,

			expectedCode:  "efgh5678",
			expectedError: nil,
		},
		{
			name: "Exceed max retries on code collision",

			setupMockRandomCodeGen: func() *mockRandomCodeGen.CodeGenerator {
				codeGenMock := mockRandomCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", 8).Return("abcd1234", nil)
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockRepo.UrlStorage {
				repoMock := mockRepo.NewUrlStorage(t)
				repoMock.On("GetURL", ctx, "abcd1234").Return("https://collision.com", nil)
				return repoMock
			},

			inputOriginalURL: "https://example.com",
			inputExpireIn:    3600,

			expectedCode:  "",
			expectedError: ErrMaxRetriesExceeded,
		},
		{
			name: "Fail to check existing URL due to repo error",

			setupMockRandomCodeGen: func() *mockRandomCodeGen.CodeGenerator {
				codeGenMock := mockRandomCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", 8).Return("abcd1234", nil)
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockRepo.UrlStorage {
				repoMock := mockRepo.NewUrlStorage(t)
				repoMock.On("GetURL", ctx, "abcd1234").Return("", assert.AnError)
				return repoMock
			},

			inputOriginalURL: "https://example.com",
			inputExpireIn:    3600,

			expectedCode:  "",
			expectedError: assert.AnError,
		},
		{
			name: "Fail to store URL after retries",
			setupMockRandomCodeGen: func() *mockRandomCodeGen.CodeGenerator {
				codeGenMock := mockRandomCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", 8).Return("abcd1234", nil)
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockRepo.UrlStorage {
				repoMock := mockRepo.NewUrlStorage(t)
				repoMock.On("GetURL", ctx, "abcd1234").Return("", redis.Nil)
				repoMock.On("StoreURL", ctx, "abcd1234", "https://example.com", 3600).Return(assert.AnError)
				return repoMock
			},

			inputOriginalURL: "https://example.com",
			inputExpireIn:    3600,

			expectedCode:  "",
			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRandomCodeGen := tc.setupMockRandomCodeGen()
			mockRepo := tc.setupMockRepo(t.Context())
			testSvc := NewShortenURL(mockRepo, mockRandomCodeGen)

			code, err := testSvc.ShortenURL(t.Context(), tc.inputOriginalURL, tc.inputExpireIn)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedCode, code)
		})
	}
}

func TestShortenURL_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRepo func(ctx context.Context) *mockRepo.UrlStorage

		inputURLCode string

		expectedOriginalURL string
		expectedError       error
	}{
		{
			name: "Get URL successfully",

			setupMockRepo: func(ctx context.Context) *mockRepo.UrlStorage {
				repoMock := mockRepo.NewUrlStorage(t)
				repoMock.On("GetURL", ctx, "abcd1234").Return("https://example.com", nil)
				return repoMock
			},

			inputURLCode: "abcd1234",

			expectedOriginalURL: "https://example.com",
			expectedError:       nil,
		},
		{
			name: "URL code not found",

			setupMockRepo: func(ctx context.Context) *mockRepo.UrlStorage {
				repoMock := mockRepo.NewUrlStorage(t)
				repoMock.On("GetURL", ctx, "unknowncode").Return("", redis.Nil)
				return repoMock
			},

			inputURLCode: "unknowncode",

			expectedOriginalURL: "",
			expectedError:       ErrCodeNotFound,
		},
		{
			name: "Repository error",

			setupMockRepo: func(ctx context.Context) *mockRepo.UrlStorage {
				repoMock := mockRepo.NewUrlStorage(t)
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
			service := NewShortenURL(repoMock, nil)

			originalURL, err := service.GetURL(ctx, tc.inputURLCode)

			assert.Equal(t, tc.expectedOriginalURL, originalURL)
			assert.Equal(t, tc.expectedError, err)

			repoMock.AssertExpectations(t)
		})
	}
}
