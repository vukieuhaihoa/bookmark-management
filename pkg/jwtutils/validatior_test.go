package jwtutils

import (
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJwtValidator_ValidateToken(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		token         string
		publicKeyPath string

		expectedErrStr string
		expectedClaims jwt.MapClaims
	}{
		{
			name: "valid token",

			token:         "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsIm5hbWUiOiJUZXN0IFVzZXIiLCJzdWIiOiIxMjM0NTY3ODkwIn0.RWWQmqnaAqQzzBlkddcO-OzxmKbu1-yVVbcegzfO_LeftlO7oR0FUptltX30FUkC7Mx4PraToanDkEqYG3aSgIw-JZ-NoanbFbbdXvMsjcqVyRaUeC1VMwUz63J42ZV_f-rLBtPYaRqqbgsE1JuDT_bXddCu8_Xc5fhmxi5lRxuH8uHEgBR_mZn5neZF8BPCcgUZiq6-UEgS_RJWQHoLaqANtKvVIYwQ9I9OXbQzNi-l9MG_U7SmD7aWXSjHYoaL2NNrzAWLpkuqUVb_vqFrHri0X3qESLav8thAwbh0inE7OM4oBk23TBm9vvSPiLRLlhy9Qe5WAw8Yc1NK0NzLVA",
			publicKeyPath: filepath.FromSlash("./public_key.test.pem"),

			expectedClaims: jwt.MapClaims{
				"sub":   "1234567890",
				"name":  "Test User",
				"admin": false,
			},
		},
		{
			name: "invalid token",

			token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
			publicKeyPath: filepath.FromSlash("./public_key.test.pem"),

			expectedErrStr: "invalid token",
		},
		{
			name: "invalid public key path",

			token:         "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsIm5hbWUiOiJUZXN0IFVzZXIiLCJzdWIiOiIxMjM0NTY3ODkwIn0.RWWQmqnaAqQzzBlkddcO-OzxmKbu1-yVVbcegzfO_LeftlO7oR0FUptltX30FUkC7Mx4PraToanDkEqYG3aSgIw-JZ-NoanbFbbdXvMsjcqVyRaUeC1VMwUz63J42ZV_f-rLBtPYaRqqbgsE1JuDT_bXddCu8_Xc5fhmxi5lRxuH8uHEgBR_mZn5neZF8BPCcgUZiq6-UEgS_RJWQHoLaqANtKvVIYwQ9I9OXbQzNi-l9MG_U7SmD7aWXSjHYoaL2NNrzAWLpkuqUVb_vqFrHri0X3qESLav8thAwbh0inE7OM4oBk23TBm9vvSPiLRLlhy9Qe5WAw8Yc1NK0NzLVA",
			publicKeyPath: filepath.FromSlash("./non_existent_public_key.pem"),

			expectedErrStr: "no such file or directory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			validator, err := NewJWTValidator(tc.publicKeyPath)
			if err != nil {
				assert.Contains(t, err.Error(), tc.expectedErrStr)
				return
			}

			claims, err := validator.ValidateToken(tc.token)
			if tc.expectedErrStr != "" {
				assert.Contains(t, err.Error(), tc.expectedErrStr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedClaims, claims)
		})
	}
}
