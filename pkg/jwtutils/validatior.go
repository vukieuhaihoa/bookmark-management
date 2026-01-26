package jwtutils

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var errInvalidToken = errors.New("invalid token")

// JWTValidator defines the interface for validating JWT tokens.
//
//go:generate mockery --name=JWTValidator --filename=validatior.go --output=./mocks
type JWTValidator interface {
	// ValidateToken validates the given JWT token string and returns the claims if valid.
	//
	// Parameters:
	//   - tokenString: The JWT token string to be validated
	//
	// Returns:
	//   - jwt.MapClaims: The claims contained in the token if valid
	//   - error: An error if token validation fails, otherwise nil
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

// jwtValidator is the implementation of JWTValidator interface.
type jwtValidator struct {
	publicKey *rsa.PublicKey
}

// NewJWTValidator creates a new instance of the JWTValidator implementation.
//
// Parameters:
//   - publicKeyPath: The file path to the RSA public key in PEM format
//
// Returns:
//   - JWTValidator: A new JWTValidator instance
//   - error: An error if initialization fails, otherwise nil
func NewJWTValidator(publicKeyPath string) (JWTValidator, error) {
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}
	return &jwtValidator{
		publicKey: publicKey,
	}, nil
}

// ValidateToken validates the given JWT token string and returns the claims if valid.
//
// Parameters:
//   - tokenString: The JWT token string to be validated
//
// Returns:
//   - jwt.MapClaims: The claims contained in the token if valid
//   - error: An error if token validation fails, otherwise nil
func (j *jwtValidator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errInvalidToken
	}

	return token.Claims.(jwt.MapClaims), nil
}
