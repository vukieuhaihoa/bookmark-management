package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"

	mockRepo "github.com/vukieuhaihoa/bookmark-management/internal/repository/mocks"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	mockJWT "github.com/vukieuhaihoa/bookmark-management/pkg/jwtutils/mocks"
	"github.com/vukieuhaihoa/bookmark-management/pkg/stringutils"
	mockPasswordHashing "github.com/vukieuhaihoa/bookmark-management/pkg/stringutils/mocks"
)

var ErrCannotGenerateToken = errors.New("cannot generate token")

func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockPasswordHashing func(t *testing.T) *mockPasswordHashing.PasswordHashing
		setupMockUserRepo        func(ctx context.Context) *mockRepo.User

		inputUsername    string
		inputPassword    string
		inputDisplayName string
		inputEmail       string

		expectedOutput *model.User
		expectedError  error
	}{
		{
			name: "Create user successfully",

			setupMockPasswordHashing: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				hashingMock := mockPasswordHashing.NewPasswordHashing(t)
				hashingMock.On("Hash", "password123").Return("$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", nil)
				return hashingMock
			},

			setupMockUserRepo: func(ctx context.Context) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("CreateUser", ctx, &model.User{
					Username:    "testuser",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
					DisplayName: "Test User",
					Email:       "testuser@example.com",
				}).Return(&model.User{
					ID:          "de305d54-75b4-431b-adb2-eb6b9e546099",
					Username:    "testuser",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
					DisplayName: "Test User",
					Email:       "testuser@example.com",
				}, nil)
				return repoMock
			},

			inputUsername:    "testuser",
			inputPassword:    "password123",
			inputDisplayName: "Test User",
			inputEmail:       "testuser@example.com",

			expectedOutput: &model.User{
				ID:          "de305d54-75b4-431b-adb2-eb6b9e546099",
				Username:    "testuser",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				DisplayName: "Test User",
				Email:       "testuser@example.com",
			},
		},

		{
			name: "Fail to hash password",

			setupMockPasswordHashing: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				hashingMock := mockPasswordHashing.NewPasswordHashing(t)
				hashingMock.On("Hash", "badpassword").Return("", stringutils.ErrCannotGenerateHash)
				return hashingMock
			},
			setupMockUserRepo: func(ctx context.Context) *mockRepo.User {
				return mockRepo.NewUser(t)
			},

			inputUsername:    "testuser2",
			inputPassword:    "badpassword",
			inputDisplayName: "Test User 2",
			inputEmail:       "testuser2@example.com",

			expectedError: stringutils.ErrCannotGenerateHash,
		},

		{
			name: "Fail to create user in repository",

			setupMockPasswordHashing: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				hashingMock := mockPasswordHashing.NewPasswordHashing(t)
				hashingMock.On("Hash", "password123").Return("$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", nil)
				return hashingMock
			},

			setupMockUserRepo: func(ctx context.Context) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("CreateUser", ctx, &model.User{
					Username:    "testuser3",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
					DisplayName: "Test User 3",
					Email:       "testuser3@example.com",
				}).Return(nil, assert.AnError)
				return repoMock
			},

			inputUsername:    "testuser3",
			inputPassword:    "password123",
			inputDisplayName: "Test User 3",
			inputEmail:       "testuser3@example.com",

			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			passwordHashingMock := tc.setupMockPasswordHashing(t)
			userRepoMock := tc.setupMockUserRepo(ctx)

			userService := NewUser(userRepoMock, passwordHashingMock, nil)

			res, err := userService.CreateUser(ctx, tc.inputUsername, tc.inputPassword, tc.inputDisplayName, tc.inputEmail)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}

func TestUser_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockUserRepo     func(ctx context.Context) *mockRepo.User
		setupMockPasswordHash func(t *testing.T) *mockPasswordHashing.PasswordHashing
		setupMockJWTGen       func(t *testing.T) *mockJWT.JWTGenerator

		inputUsername string
		inputPassword string

		expectedError error

		expectedOutput string
	}{
		{
			name: "Login successfully",

			setupMockUserRepo: func(ctx context.Context) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByUsername", ctx, "testuser").Return(&model.User{
					ID:       "de305d54-75b4-431b-adb2-eb6b9e546099",
					Username: "testuser",
					Password: "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", // hash for "password123"
				}, nil)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				hashingMock := mockPasswordHashing.NewPasswordHashing(t)
				hashingMock.On("CompareHashAndPassword", "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", "password123").Return(true)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mockJWT.JWTGenerator {
				jwtMock := mockJWT.NewJWTGenerator(t)
				jwtMock.On("GenerateToken", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					if claims["sub"] != "de305d54-75b4-431b-adb2-eb6b9e546099" {
						return false
					}

					if _, ok := claims["iat"].(int64); !ok {
						return false
					}

					if _, ok := claims["exp"].(int64); !ok {
						return false
					}

					return true
				})).Return("mocked_jwt_token", nil)
				return jwtMock
			},

			inputUsername: "testuser",
			inputPassword: "password123",

			expectedOutput: "mocked_jwt_token",
		},
		{
			name: "Fail to get user by username",

			setupMockUserRepo: func(ctx context.Context) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByUsername", ctx, "nonexistentuser").Return(nil, ErrInvalidCredentials)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				return mockPasswordHashing.NewPasswordHashing(t)
			},

			setupMockJWTGen: func(t *testing.T) *mockJWT.JWTGenerator {
				return mockJWT.NewJWTGenerator(t)
			},

			inputUsername: "nonexistentuser",
			inputPassword: "somepassword",

			expectedError: ErrInvalidCredentials,
		},
		{
			name: "Invalid password",

			setupMockUserRepo: func(ctx context.Context) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByUsername", ctx, "testuser").Return(&model.User{
					ID:       "de305d54-75b4-431b-adb2-eb6b9e546099",
					Username: "testuser",
					Password: "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", // hash for "password123"
				}, nil)
				return repoMock
			},
			setupMockPasswordHash: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				hashingMock := mockPasswordHashing.NewPasswordHashing(t)
				hashingMock.On("CompareHashAndPassword", "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", "wrongpassword").Return(false)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mockJWT.JWTGenerator {
				return mockJWT.NewJWTGenerator(t)
			},

			inputUsername: "testuser",
			inputPassword: "wrongpassword",

			expectedError: ErrInvalidCredentials,
		},
		{
			name: "Fail to generate JWT token",

			setupMockUserRepo: func(ctx context.Context) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByUsername", ctx, "testuser").Return(&model.User{
					ID:       "de305d54-75b4-431b-adb2-eb6b9e546099",
					Username: "testuser",
					Password: "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", // hash for "password123"
				}, nil)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				hashingMock := mockPasswordHashing.NewPasswordHashing(t)
				hashingMock.On("CompareHashAndPassword", "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", "password123").Return(true)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mockJWT.JWTGenerator {
				jwtMock := mockJWT.NewJWTGenerator(t)
				jwtMock.On("GenerateToken", mock.Anything).Return("", ErrCannotGenerateToken)
				return jwtMock
			},

			inputUsername: "testuser",
			inputPassword: "password123",

			expectedError: ErrCannotGenerateToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			userRepoMock := tc.setupMockUserRepo(ctx)
			passwordHashingMock := tc.setupMockPasswordHash(t)
			jwtGenMock := tc.setupMockJWTGen(t)

			userService := NewUser(userRepoMock, passwordHashingMock, jwtGenMock)

			res, err := userService.Login(ctx, tc.inputUsername, tc.inputPassword)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, res)

		})
	}
}

func TestUser_GetUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockUserRepo func(ctx context.Context) *mockRepo.User

		inputUserID string

		expectedOutput *model.User
		expectedError  error
	}{
		{
			name: "Get user by ID successfully",

			setupMockUserRepo: func(ctx context.Context) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
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

			setupMockUserRepo: func(ctx context.Context) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
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

			userService := NewUser(userRepoMock, nil, nil)

			res, err := userService.GetUserByID(ctx, tc.inputUserID)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}

func TestUser_UpdateUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockUserRepo func(ctx context.Context) *mockRepo.User
		inputUserID       string
		inputDisplayName  string
		inputEmail        string

		expectedError error
	}{
		{
			name: "Update user by ID successfully",

			setupMockUserRepo: func(ctx context.Context) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("UpdateUserByID", ctx, "de305d54-75b4-431b-adb2-eb6b9e546099", &model.User{
					DisplayName: "Updated User",
					Email:       "updateduser@example.com",
				}).Return(nil)
				return repoMock
			},

			inputUserID:      "de305d54-75b4-431b-adb2-eb6b9e546099",
			inputDisplayName: "Updated User",
			inputEmail:       "updateduser@example.com",
		},
		{
			name: "Fail to update user by ID - user not found",

			setupMockUserRepo: func(ctx context.Context) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("UpdateUserByID", ctx, "nonexistentid", &model.User{
					DisplayName: "Updated User",
					Email:       "updateduser@example.com",
				}).Return(dbutils.ErrRecordNotFoundType)
				return repoMock
			},

			inputUserID:      "nonexistentid",
			inputDisplayName: "Updated User",
			inputEmail:       "updateduser@example.com",

			expectedError: dbutils.ErrRecordNotFoundType,
		},
		{
			name: "Fail to update user by ID - duplicate email",

			setupMockUserRepo: func(ctx context.Context) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("UpdateUserByID", ctx, "de305d54-75b4-431b-adb2-eb6b9e546099", &model.User{
					DisplayName: "Updated User",
					Email:       "duplicateemail@example.com",
				}).Return(dbutils.ErrDuplicationType)
				return repoMock
			},

			inputUserID:      "de305d54-75b4-431b-adb2-eb6b9e546099",
			inputDisplayName: "Updated User",
			inputEmail:       "duplicateemail@example.com",

			expectedError: dbutils.ErrDuplicationType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			userRepoMock := tc.setupMockUserRepo(ctx)

			userService := NewUser(userRepoMock, nil, nil)

			err := userService.UpdateUserByID(ctx, tc.inputUserID, tc.inputDisplayName, tc.inputEmail)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
