package fixture

import (
	"testing"

	"github.com/vukieuhaihoa/bookmark-management/pkg/sqldb"
	"gorm.io/gorm"
)

// Fixture defines the interface for setting up test fixtures in the database.
type Fixture interface {
	// SetupDB sets up the database connection for the fixture.
	//
	// Parameters:
	//   - db: The GORM database connection to be used for the fixture
	SetupDB(db *gorm.DB)

	// Migrate migrates the database schema for the fixture.
	//
	// Returns:
	//   - error: An error if migration fails, otherwise nil
	Migrate() error

	// GenerateData generates test data for the fixture.
	//
	// Returns:
	//   - error: An error if data generation fails, otherwise nil
	GenerateData() error
	DB() *gorm.DB
}

// NewFixture initializes the test fixture by setting up the database,
// migrating the schema, and generating test data.
//
// Parameters:
//   - t: The testing object used for reporting errors
//   - fix: The fixture to be initialized
//
// Returns:
//   - *gorm.DB: A gorm.DB instance connected to the initialized test database
func NewFixture(t *testing.T, fix Fixture) *gorm.DB {
	// step 1: create test database
	fix.SetupDB(sqldb.InitMockDB(t))

	// step 2: migrate schema
	err := fix.Migrate()
	if err != nil {
		t.Fatalf("Failed to migrate db for testing")
	}
	// step 3: generate test data
	err = fix.GenerateData()
	if err != nil {
		t.Fatalf("Failed to generate test data")
	}

	return fix.DB()
}
