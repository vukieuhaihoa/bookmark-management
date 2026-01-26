package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HandlerError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		inputError  error
		expectPanic bool
	}{
		{
			name:        "nil error",
			inputError:  nil,
			expectPanic: false,
		},
		{
			name:        "non-nil error",
			inputError:  errors.New("test error"),
			expectPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.expectPanic {
				assert.Panics(t, func() {
					HandlerError(tc.inputError)
				})
			} else {
				assert.NotPanics(t, func() {
					HandlerError(tc.inputError)
				})
			}
		})
	}
}
