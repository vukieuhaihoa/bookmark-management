package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetJWTClaimsFromRequest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupCtx func(ctx *gin.Context)

		expectedClaims jwt.MapClaims
		expectedError  error
	}{
		{
			name: "claims not set in context",

			setupCtx: func(ctx *gin.Context) {
				// nothing set
			},

			expectedClaims: nil,
			expectedError:  ErrInvalidToken,
		},
		{
			name: "claims is wrong type in context",

			setupCtx: func(ctx *gin.Context) {
				ctx.Set("claims", "not-a-map-claims")
			},

			expectedClaims: nil,
			expectedError:  ErrInvalidToken,
		},
		{
			name: "valid claims in context",

			setupCtx: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{
					"sub": "de305d54-75b4-431b-adb2-eb6b9e546099",
				})
			},

			expectedClaims: jwt.MapClaims{
				"sub": "de305d54-75b4-431b-adb2-eb6b9e546099",
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			tc.setupCtx(ctx)

			claims, err := GetJWTClaimsFromRequest(ctx)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedClaims, claims)
		})
	}
}

func TestGetUserIDFromJWTClaims(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupCtx func(ctx *gin.Context)

		expectedUserID string
		expectedError  error
	}{
		{
			name: "claims not set in context",

			setupCtx: func(ctx *gin.Context) {
				// nothing set
			},

			expectedUserID: "",
			expectedError:  ErrInvalidToken,
		},
		{
			name: "sub claim missing",

			setupCtx: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{
					"name": "Alice",
				})
			},

			expectedUserID: "",
			expectedError:  ErrEmptyID,
		},
		{
			name: "sub claim is wrong type",

			setupCtx: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{
					"sub": 12345, // int, not string
				})
			},

			expectedUserID: "",
			expectedError:  ErrEmptyID,
		},
		{
			name: "sub claim is empty string",

			setupCtx: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{
					"sub": "",
				})
			},

			expectedUserID: "",
			expectedError:  ErrEmptyID,
		},
		{
			name: "valid sub claim",

			setupCtx: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{
					"sub": "de305d54-75b4-431b-adb2-eb6b9e546099",
				})
			},

			expectedUserID: "de305d54-75b4-431b-adb2-eb6b9e546099",
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			tc.setupCtx(ctx)

			userID, err := GetUserIDFromJWTClaims(ctx)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUserID, userID)
		})
	}
}
