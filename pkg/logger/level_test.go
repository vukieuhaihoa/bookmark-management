package logger

import "testing"

func TestSetLogLevel(t *testing.T) {
	t.Run("SetLogLevel with valid LOG_LEVEL", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "debug")
		SetLogLevel()
	})

	t.Run("SetLogLevel with invalid LOG_LEVEL", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "invalid_level")
		SetLogLevel()
	})

	t.Run("SetLogLevel with no LOG_LEVEL", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "")
		SetLogLevel()
	})
}
