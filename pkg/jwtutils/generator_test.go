package jwtutils

import (
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJwtGenerator_GenerateToken(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		privateKeyPath string
		inputClaims    jwt.MapClaims
		expectedErrStr string
		expectErr      error
		expectedOutput string
	}{
		{
			name: "valid case",

			inputClaims: jwt.MapClaims{
				"sub":   "1234567890",
				"name":  "Test User",
				"admin": false,
			},

			privateKeyPath: filepath.FromSlash("./private_key.test.pem"),

			expectedOutput: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsIm5hbWUiOiJUZXN0IFVzZXIiLCJzdWIiOiIxMjM0NTY3ODkwIn0.RWWQmqnaAqQzzBlkddcO-OzxmKbu1-yVVbcegzfO_LeftlO7oR0FUptltX30FUkC7Mx4PraToanDkEqYG3aSgIw-JZ-NoanbFbbdXvMsjcqVyRaUeC1VMwUz63J42ZV_f-rLBtPYaRqqbgsE1JuDT_bXddCu8_Xc5fhmxi5lRxuH8uHEgBR_mZn5neZF8BPCcgUZiq6-UEgS_RJWQHoLaqANtKvVIYwQ9I9OXbQzNi-l9MG_U7SmD7aWXSjHYoaL2NNrzAWLpkuqUVb_vqFrHri0X3qESLav8thAwbh0inE7OM4oBk23TBm9vvSPiLRLlhy9Qe5WAw8Yc1NK0NzLVA",
		},
		{
			name: "invalid private key path",

			inputClaims: jwt.MapClaims{
				"sub":   "1234567890",
				"name":  "Test User",
				"admin": false,
			},

			privateKeyPath: filepath.FromSlash("./non_existent_private_key.pem"),

			expectedErrStr: "no such file or directory",
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testGenerator, err := NewJWTGenerator(tc.privateKeyPath)
			if err != nil {
				assert.Contains(t, err.Error(), tc.expectedErrStr)
				return
			}
			token, err := testGenerator.GenerateToken(tc.inputClaims)
			assert.Equal(t, tc.expectedOutput, token)
		})
	}
}
