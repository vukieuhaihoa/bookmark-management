package common

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func Test_InputFieldError(t *testing.T) {
	t.Parallel()

	validate := validator.New()

	testCases := []struct {
		name          string
		inputError    error
		expectedError Message
	}{
		{
			name:       "Non-validation error",
			inputError: errors.New("some random error"),
			expectedError: Message{
				Message: "Invalid input",
			},
		},
		{
			name: "Validation error with multiple fields",
			inputError: func() error {

				type TestStruct struct {
					Name  string `validate:"required"`
					Email string `validate:"required,email"`
				}
				testObj := TestStruct{
					Name:  "",
					Email: "invalid-email",
				}
				return validate.Struct(testObj)
			}(),
			expectedError: Message{
				Message: "Invalid input fields",
				Details: []string{
					"Name is invalid (required)",
					"Email is invalid (email)",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := InputFieldError(tc.inputError)
			assert.Equal(t, tc.expectedError, res)
			if res.Details != nil {
				assert.ElementsMatch(t, tc.expectedError.Details.([]string), res.Details.([]string))
			}
		})
	}
}
