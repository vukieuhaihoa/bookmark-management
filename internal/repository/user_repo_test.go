package repository

import (
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB   func(t *testing.T) *gorm.DB
		inputUser *model.User

		expectedError  error
		expectedOutput *model.User
		verifyFunc     func(db *gorm.DB, user *model.User)
	}{
		{
			name: "Create user successfully",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputUser: &model.User{
				ID:          "de305d54-75b4-431b-adb2-eb6b9e546099",
				DisplayName: "New User",
				Username:    "New User",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				Email:       "newuser@example.com",
			},

			expectedOutput: &model.User{
				ID:          "de305d54-75b4-431b-adb2-eb6b9e546099",
				DisplayName: "New User",
				Username:    "New User",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				Email:       "newuser@example.com",
			},

			verifyFunc: func(db *gorm.DB, user *model.User) {
				checkUser := &model.User{}
				err := db.Where("id = ?", user.ID).First(checkUser).Error
				assert.Nil(t, err)

				// Verify timestamps are automatically set by GORM
				assert.False(t, user.CreatedAt.IsZero(), "CreatedAt should be automatically set")
				assert.False(t, user.UpdatedAt.IsZero(), "UpdatedAt should be automatically set")

				// On creation, both timestamps should be very close (within a second)
				timeDiff := user.UpdatedAt.Sub(user.CreatedAt).Abs()
				assert.Less(t, timeDiff, time.Second, "CreatedAt and UpdatedAt should be nearly equal on creation")

				// Verify other fields
				assert.Equal(t, user.ID, checkUser.ID)
				assert.Equal(t, user.Username, checkUser.Username)
				assert.Equal(t, user.Email, checkUser.Email)
				assert.Equal(t, user.DisplayName, checkUser.DisplayName)
			},
		},
		{
			name: "Create user failed - duplicate username",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputUser: &model.User{
				ID:          "de305d54-75b4-431b-adb2-eb6b9e546099",
				DisplayName: "Another User",
				Username:    "Alice", // duplicate username
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				Email:       "alice1@example.com",
			},

			expectedError: dbutils.ErrDuplicationType,
		},
		{
			name: "Create user failed - duplicate email",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputUser: &model.User{
				ID:          "de305d54-75b4-431b-adb2-eb6b9e546099",
				DisplayName: "Another User",
				Username:    "AnotherUser",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				Email:       "alice@example.com",
			},

			expectedError: dbutils.ErrDuplicationType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			testUserRepo := NewUser(db)

			res, err := testUserRepo.CreateUser(ctx, tc.inputUser)
			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			tc.verifyFunc(db, res)
		})
	}
}

func TestUser_GetUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB
		inputID string

		expectedError  error
		expectedOutput *model.User
	}{
		{
			name: "Get user by ID successfully",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputID: "de305d54-75b4-431b-adb2-eb6b9e546000",

			expectedOutput: &model.User{
				ID:          "de305d54-75b4-431b-adb2-eb6b9e546000",
				Username:    "Alice",
				DisplayName: "Alice",
				Email:       "alice@example.com",
				CreatedAt:   fixture.TestTime,
				UpdatedAt:   fixture.TestTime,
			},
		},
		{
			name: "Get user by ID failed - user not found",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputID: "non-existent-id",

			expectedError: dbutils.ErrRecordNotFoundType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			testUserRepo := NewUser(db)

			res, err := testUserRepo.GetUserByID(ctx, tc.inputID)
			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			res.Password = "" // omit password field for comparison
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}

func TestUser_GetUserByUsername(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB       func(t *testing.T) *gorm.DB
		inputUsername string

		expectedError  error
		expectedOutput *model.User
	}{
		{
			name: "Get user by username successfully",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputUsername: "Bob",

			expectedOutput: &model.User{
				ID:          "123e4567-e89b-12d3-a456-eb6b9e546001",
				Username:    "Bob",
				DisplayName: "Bob",
				Email:       "bob@example.com",
				CreatedAt:   fixture.TestTime,
				UpdatedAt:   fixture.TestTime,
			},
		},
		{
			name: "Get user by username failed - user not found",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputUsername: "NonExistentUser",

			expectedError: dbutils.ErrRecordNotFoundType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			testUserRepo := NewUser(db)

			res, err := testUserRepo.GetUserByUsername(ctx, tc.inputUsername)
			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			res.Password = "" // omit password field for comparison
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}

func TestUser_UpdateUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB       func(t *testing.T) *gorm.DB
		inputID       string
		inputUserData *model.User

		expectedError error
	}{
		{
			name: "Update user by ID successfully",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputID: "de305d54-75b4-431b-adb2-eb6b9e546000",

			inputUserData: &model.User{
				DisplayName: "Alice Updated",
				Email:       "alice.updated@example.com",
			},
		},
		{
			name: "Update user by ID failed - user not found",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputID: "non-existent-id",

			inputUserData: &model.User{
				DisplayName: "Non Existent User",
				Email:       "",
			},

			expectedError: dbutils.ErrRecordNotFoundType,
		},
		{
			name: "Update user by ID failed - duplicate email",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputID: "de305d54-75b4-431b-adb2-eb6b9e546000",

			inputUserData: &model.User{
				DisplayName: "Alice",
				Email:       "bob@example.com", // duplicate email
			},

			expectedError: dbutils.ErrDuplicationType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			testUserRepo := NewUser(db)

			err := testUserRepo.UpdateUserByID(ctx, tc.inputID, tc.inputUserData)
			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
		})
	}
}
