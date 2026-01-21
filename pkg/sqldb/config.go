package sqldb

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host     string `default:"localhost" envconfig:"DB_HOST"`
	User     string `default:"admin" envconfig:"DB_USER"`
	Password string `default:"admin" envconfig:"DB_PASSWORD"`
	DBName   string `default:"bookmark_service" envconfig:"DB_NAME"`
	Port     string `default:"5432" envconfig:"DB_PORT"`
}

// NewConfig creates a new database configuration by reading environment variables.
//
// Parameters:
//   - envPrefix: The prefix for environment variables containing database configuration
//
// Returns:
//   - *config: A pointer to the config struct containing database configuration
//   - error: An error if the configuration cannot be processed, otherwise nil
func NewConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetDSN constructs the Data Source Name (DSN) for connecting to the PostgreSQL database.
//
// Returns:
//   - string: The DSN string for database connection
func (cfg *config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)
}
