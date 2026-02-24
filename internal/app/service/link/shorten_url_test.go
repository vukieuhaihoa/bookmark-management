package link

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	mockRepo "github.com/vukieuhaihoa/bookmark-management/internal/app/repository/link/mocks"
	mockRandomCodeGen "github.com/vukieuhaihoa/bookmark-management/pkg/utils/mocks"
)

func TestService_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRandomCodeGen func() *mockRandomCodeGen.CodeGenerator
		setupMockRepo          func(ctx context.Context) *mockRepo.Repository

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

			setupMockRepo: func(ctx context.Context) *mockRepo.Repository {
				repoMock := mockRepo.NewRepository(t)
				repoMock.On("StoreURLIfAbsent", ctx, model.RedisShortenPrefix+"abcd1234", "https://example.com", 3600).Return(true, nil)
				return repoMock
			},

			inputOriginalURL: "https://example.com",
			inputExpireIn:    3600,

			expectedCode:  "rabcd1234",
			expectedError: nil,
		},
		{
			name: "Fail to generate code",

			setupMockRandomCodeGen: func() *mockRandomCodeGen.CodeGenerator {
				codeGenMock := mockRandomCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", 8).Return("", assert.AnError)
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockRepo.Repository {
				return mockRepo.NewRepository(t)
			},

			inputOriginalURL: "https://example.com",
			inputExpireIn:    3600,

			expectedCode:  "",
			expectedError: assert.AnError,
		},
		{
			name: "Fail to store URL due to repository error",

			setupMockRandomCodeGen: func() *mockRandomCodeGen.CodeGenerator {
				codeGenMock := mockRandomCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", 8).Return("abcd1234", nil)
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockRepo.Repository {
				repoMock := mockRepo.NewRepository(t)
				repoMock.On("StoreURLIfAbsent", ctx, model.RedisShortenPrefix+"abcd1234", "https://example.com", 3600).Return(false, assert.AnError)
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

			setupMockRepo: func(ctx context.Context) *mockRepo.Repository {
				repoMock := mockRepo.NewRepository(t)
				repoMock.On("StoreURLIfAbsent", ctx, model.RedisShortenPrefix+"abcd1234", "https://example.com", 3600).Return(false, nil).Once()
				repoMock.On("StoreURLIfAbsent", ctx, model.RedisShortenPrefix+"efgh5678", "https://example.com", 3600).Return(true, nil).Once()
				return repoMock
			},

			inputOriginalURL: "https://example.com",
			inputExpireIn:    3600,

			expectedCode:  "refgh5678",
			expectedError: nil,
		},
		{
			name: "Exceed max retries on code collision",

			setupMockRandomCodeGen: func() *mockRandomCodeGen.CodeGenerator {
				codeGenMock := mockRandomCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", 8).Return("abcd1234", nil)
				return codeGenMock
			},

			setupMockRepo: func(ctx context.Context) *mockRepo.Repository {
				repoMock := mockRepo.NewRepository(t)
				repoMock.On("StoreURLIfAbsent", ctx, model.RedisShortenPrefix+"abcd1234", "https://example.com", 3600).Return(false, nil)
				return repoMock
			},

			inputOriginalURL: "https://example.com",
			inputExpireIn:    3600,

			expectedCode:  "",
			expectedError: ErrMaxRetriesExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRandomCodeGen := tc.setupMockRandomCodeGen()
			mockRepo := tc.setupMockRepo(t.Context())
			testSvc := NewLinkService(mockRepo, mockRandomCodeGen, nil)

			code, err := testSvc.ShortenURL(t.Context(), tc.inputOriginalURL, tc.inputExpireIn)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedCode, code)
		})
	}
}
