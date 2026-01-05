package service

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	mockRandomCodeGen "github.com/vukieuhaihoa/bookmark-management/pkg/stringutils/mocks"
)

var urlSafeRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func TestPasswordService_GeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRandomCodeGen func() *mockRandomCodeGen.CodeGenerator

		expectedLength int
		expectedError  error
	}{
		{
			name: "Generate password successfully",

			setupMockRandomCodeGen: func() *mockRandomCodeGen.CodeGenerator {
				codeGenMock := mockRandomCodeGen.NewCodeGenerator(t)
				codeGenMock.On("GenerateCode", 12).Return("AbcD1234EfGh", nil)
				return codeGenMock
			},

			expectedLength: 12,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRandomCodeGen := tc.setupMockRandomCodeGen()
			testSvc := NewPassword(mockRandomCodeGen)
			pass, err := testSvc.GeneratePassword()

			assert.Equal(t, tc.expectedLength, len(pass))
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, urlSafeRegex.MatchString(pass), true)

		})
	}
}
