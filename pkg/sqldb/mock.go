package sqldb

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitMockDB initializes an in-memory SQLite database for testing purposes.
//
// Parameters:
//   - t: The testing object used for reporting errors
//
// Returns:
//   - *gorm.DB: A gorm.DB instance connected to the in-memory SQLite database
func InitMockDB(t *testing.T) *gorm.DB {
	cdn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_foreign_keys=on", uuid.New().String())

	db, err := gorm.Open(sqlite.Open(cdn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		t.Fatalf("Fail to create db: %v", err)
	}

	return db
}
