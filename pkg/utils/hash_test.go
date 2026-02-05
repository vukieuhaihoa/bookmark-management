package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing_Hash(t *testing.T) {
	t.Parallel()

	passwordHashing := NewPasswordHashing()

	testCases := []struct {
		name string

		inputPassword string

		expectedHashedPassword string
		expectedError          error
	}{
		{
			name: "Hash password successfully",

			inputPassword: "my_secure_password",

			expectedError: nil,
		},
		{
			name: "Hash too long password",

			inputPassword: "this_password_is_way_too_long_and_should_trigger_an_error_because_it_exceeds_the_maximum_length_allowed_by_bcrypt_which_is_72_bytes",

			expectedError: ErrCannotGenerateHash,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			hashedPassword, err := passwordHashing.Hash(tc.inputPassword)
			assert.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.NotEmpty(t, hashedPassword)
			}
		})
	}
}

func TestPasswordHashing_CompareHashAndPassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputPassword       string
		inputHashedPassword string

		expectedMatch bool
	}{
		{
			name: "Password matches hash",

			inputPassword:       "my_secure_password",
			inputHashedPassword: "$2a$10$yIIizEHMEKSm.OARDrSjHe4otTolPuCjjEy6IQ3RtRny3ZB7ToN.e", // Hash for "my_secure_password"

			expectedMatch: true,
		},
		{
			name: "Password does not match hash",

			inputPassword:       "wrong_password",
			inputHashedPassword: "$2a$10$yIIizEHMEKSm.OARDrSjHe4otTolPuCjjEy6IQ3RtRny3ZB7ToN.e", // Hash for "my_secure_password"

			expectedMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			passwordHashing := NewPasswordHashing()

			match := passwordHashing.CompareHashAndPassword(tc.inputHashedPassword, tc.inputPassword)

			assert.Equal(t, tc.expectedMatch, match)
		})
	}
}
