package fixture

import (
	"time"

	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"gorm.io/gorm"
)

var TestTime = time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

type UserCommonTestDB struct {
	base
}

// Migrate migrates the database schema for the UserCommonTestDB fixture.
//
// Returns:
//   - error: An error if migration fails, otherwise nil
func (u *UserCommonTestDB) Migrate() error {
	return u.db.AutoMigrate(&model.User{})
}

// GenerateData populates the test database with common user test data.
//
// Returns:
//   - error: An error if data generation fails, otherwise nil
func (u *UserCommonTestDB) GenerateData() error {
	db := u.db.Session(&gorm.Session{})

	users := []*model.User{
		{
			ID:          "de305d54-75b4-431b-adb2-eb6b9e546000",
			DisplayName: "Alice",
			Username:    "Alice",
			Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
			Email:       "alice@example.com",
			CreatedAt:   TestTime,
			UpdatedAt:   TestTime,
		},
		{
			ID:          "123e4567-e89b-12d3-a456-eb6b9e546001",
			DisplayName: "Bob",
			Username:    "Bob",
			Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
			Email:       "bob@example.com",
			CreatedAt:   TestTime,
			UpdatedAt:   TestTime,
		},
		{
			ID:          "987e6543-e21b-12d3-a456-eb6b9e546002",
			DisplayName: "Charlie",
			Username:    "Charlie",
			Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
			Email:       "charlie@example.com",
			CreatedAt:   TestTime,
			UpdatedAt:   TestTime,
		},
		{
			ID:          "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			DisplayName: "Test User 1",
			Username:    "testuser001",
			Password:    "$2a$10$hhuB9rZrp5ikmRb5yAF9hev6AE2tC404jhtP.bdOjme9lECJClzFu",
			Email:       "testuser001@example.com",
			CreatedAt:   TestTime,
			UpdatedAt:   TestTime,
		},
	}

	return db.CreateInBatches(users, 10).Error
}
