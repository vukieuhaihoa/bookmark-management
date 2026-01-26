package jwtutils

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// JWTGenerator defines the interface for generating JWT tokens.
//
//go:generate mockery --name=JWTGenerator --filename=generator.go --output=./mocks
type JWTGenerator interface {
	// GenerateToken generates a JWT token with the given claims.
	//
	// Parameters:
	//   - claims: The JWT claims to be included in the token
	//
	// Returns:
	//   - string: The generated JWT token
	//   - error: An error if token generation fails, otherwise nil
	GenerateToken(claims jwt.Claims) (string, error)
}

type jwtGenerator struct {
	privateKey *rsa.PrivateKey
}

// NewJWTGenerator creates a new instance of the JWTGenerator implementation.
//
// Parameters:
//   - privateKeyPath: The file path to the RSA private key in PEM format
//
// Returns:
//   - JWTGenerator: A new JWTGenerator instance
//   - error: An error if initialization fails, otherwise nil
func NewJWTGenerator(privateKeyPath string) (JWTGenerator, error) {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}

	return &jwtGenerator{
		privateKey: privateKey,
	}, nil
}

// GenerateToken generates a JWT token with the given claims.
//
// Parameters:
//   - claims: The JWT claims to be included in the token
//
// Returns:
//   - string: The generated JWT token
//   - error: An error if token generation fails, otherwise nil
func (j *jwtGenerator) GenerateToken(claims jwt.Claims) (string, error) {
	// Implementation for generating JWT token using j.privateKey and claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
