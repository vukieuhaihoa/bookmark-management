package sqldb

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// MigrateSQLDB performs database migrations using the provided GORM DB instance.
// Parameters:
//   - db: A GORM DB instance connected to the target database.
//   - migrationPath: The file path to the migration files.
//   - mode: The migration mode, either "up" for full migration or "steps" for a specific number of steps.
//   - steps: The number of migration steps to apply when mode is "steps".
//
// Returns:
//   - error: An error if the migration fails, otherwise nil.
func MigrateSQLDB(db *gorm.DB, migrationPath string, mode string, steps int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance(migrationPath, db.Name(), driver)
	if err != nil {
		return err
	}

	return migrateSchema(migrator, mode, steps)
}

func migrateSchema(m *migrate.Migrate, mode string, steps int) error {
	var migrateErr error
	switch mode {
	case "up":
		migrateErr = m.Up()
	case "steps":
		if steps == 0 {
			return errors.New("[Database migration] Steps must not be 0. Please provide a valid number of steps to migrate.")
		}
		migrateErr = m.Steps(steps)
	default:
		return errors.New("[Database migration] Invalid migration mode. Supported modes are: up, steps.")
	}

	if migrateErr != nil && !errors.Is(migrateErr, migrate.ErrNoChange) {
		return fmt.Errorf("[Database migration] Migration failed: %w", migrateErr)
	}

	return nil
}
