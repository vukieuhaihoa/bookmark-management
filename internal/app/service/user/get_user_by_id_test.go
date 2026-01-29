package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"

	mockUserRepo "github.com/vukieuhaihoa/bookmark-management/internal/app/repository/user/mocks"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

func TestService_GetUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockUserRepo func(ctx context.Context) *mockUserRepo.Repository

		inputUserID string

		expectedOutput *model.User
		expectedError  error
	}{
		{
			name: "Get user by ID successfully",

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("GetUserByID", ctx, "de305d54-75b4-431b-adb2-eb6b9e546099").Return(&model.User{
					ID:          "de305d54-75b4-431b-adb2-eb6b9e546099",
					Username:    "testuser",
					DisplayName: "Test User",
					Email:       "testuser@example.com",
					CreatedAt:   fixture.TestTime,
					UpdatedAt:   fixture.TestTime,
				}, nil)
				return repoMock
			},

			inputUserID: "de305d54-75b4-431b-adb2-eb6b9e546099",

			expectedOutput: &model.User{
				ID:          "de305d54-75b4-431b-adb2-eb6b9e546099",
				Username:    "testuser",
				DisplayName: "Test User",
				Email:       "testuser@example.com",
				CreatedAt:   fixture.TestTime,
				UpdatedAt:   fixture.TestTime,
			},
		},
		{
			name: "Fail to get user by ID",

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("GetUserByID", ctx, "nonexistentid").Return(nil, dbutils.ErrRecordNotFoundType)
				return repoMock
			},

			inputUserID:   "nonexistentid",
			expectedError: dbutils.ErrRecordNotFoundType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			userRepoMock := tc.setupMockUserRepo(ctx)

			userService := NewUserService(userRepoMock, nil, nil)

			res, err := userService.GetUserByID(ctx, tc.inputUserID)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}
