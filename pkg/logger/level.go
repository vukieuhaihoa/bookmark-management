package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// SetLogLevel sets the global log level based on the LOG_LEVEL environment variable.
// If the variable is not set or invalid, it defaults to Info level.
func SetLogLevel() {
	level, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)
}
