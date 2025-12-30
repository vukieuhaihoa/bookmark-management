package api

import "github.com/kelseyhightower/envconfig"

type Config struct {
	AppPort string `envconfig:"APP_PORT" default:"8080"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
