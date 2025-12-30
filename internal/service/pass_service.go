// Package service provides business logic services for the bookmark management application.
// This file contains the password service implementation for generating cryptographically
// secure random passwords.
package service

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	passLength = 12
)

// passwordService implements the Password interface and provides methods for generating secure random passwords.
type passwordService struct {
}

// Password defines the interface for password generation service.
// It provides a method for generating a secure random password.
//
//go:generate mockery --name Password --filename pass_service.go --output ./mocks
type Password interface {
	GeneratePassword() (string, error)
}

// NewPassword creates a new instance of Password service.
// It generates random passwords of fixed length using a predefined character set.
func NewPassword() Password {
	return &passwordService{}
}

// GeneratePassword generates a cryptographically secure random password.
// It creates a password of fixed length (12 characters) using alphanumeric characters
// (a-z, A-Z, 0-9). The function uses crypto/rand for secure random number generation.
//
// Returns:
//   - string: The generated password
//   - error: An error if random number generation fails, nil otherwise
func (s *passwordService) GeneratePassword() (string, error) {
	var strBuilder bytes.Buffer
	for i := 0; i < passLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		strBuilder.WriteByte(charset[randomIndex.Int64()])
	}
	return strBuilder.String(), nil
}
