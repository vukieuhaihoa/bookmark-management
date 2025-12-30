package service

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var urlSafeRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func TestPasswordService_GeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		expectedLength int
		expectedError  error
	}{
		{
			name:           "Generate password successfully",
			expectedLength: 12,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewPassword()
			pass, err := testSvc.GeneratePassword()

			assert.Equal(t, tc.expectedLength, len(pass))
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, urlSafeRegex.MatchString(pass), true)

		})
	}
}
