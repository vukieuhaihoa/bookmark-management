package sqldb

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewClient creates a new gorm.DB client connected to a PostgreSQL database.
//
// Parameters:
//   - envPrefix: The prefix for environment variables containing database configuration
//
// Returns:
//   - *gorm.DB: A gorm.DB instance connected to the PostgreSQL database
//   - error: An error if the connection fails, otherwise nil
func NewClient(envPrefix string) (*gorm.DB, error) {
	cfg, err := NewConfig(envPrefix)
	if err != nil {
		return nil, err
	}
	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
