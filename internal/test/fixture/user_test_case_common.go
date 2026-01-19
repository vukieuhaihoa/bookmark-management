package fixture

import (
	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"gorm.io/gorm"
)

type UserCommonTestDB struct {
	db *gorm.DB
}

// SetupDB sets up the database connection for the UserCommonTestDB fixture.
//
// Parameters:
//   - db: The GORM database connection to be used for the fixture
func (u *UserCommonTestDB) SetupDB(db *gorm.DB) {
	u.db = db
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
		},
		{
			ID:          "123e4567-e89b-12d3-a456-eb6b9e546001",
			DisplayName: "Bob",
			Username:    "Bob",
			Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
			Email:       "bob@example.com",
		},
		{
			ID:          "987e6543-e21b-12d3-a456-eb6b9e546002",
			DisplayName: "Charlie",
			Username:    "Charlie",
			Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
			Email:       "charlie@example.com",
		},
	}

	return db.CreateInBatches(users, 10).Error
}

// DB returns the gorm.DB instance used by the fixture.
//
// Returns:
//   - *gorm.DB: The gorm.DB instance
func (u *UserCommonTestDB) DB() *gorm.DB {
	return u.db
}
