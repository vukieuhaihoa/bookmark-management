package repository

import (
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
)

func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB   func(t *testing.T) *gorm.DB
		inputUser *model.User

		expectedErrorString string
		expectedOutput      *model.User
		verifyFunc          func(db *gorm.DB, user *model.User)
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

			expectedErrorString: "UNIQUE constraint failed: users.username",
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

			expectedErrorString: "UNIQUE constraint failed: users.email",
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
				assert.Contains(t, err.Error(), tc.expectedErrorString)
				return
			}
			tc.verifyFunc(db, res)
		})
	}
}
