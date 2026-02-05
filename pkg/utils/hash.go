package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCannotGenerateHash = errors.New("cannot generate hash from password")
)

// PasswordHashing defines the interface for password hashing operations.
//
//go:generate mockery --name PasswordHashing --filename hash.go --output ./mocks
type PasswordHashing interface {
	// Hash generates a hashed version of the given password.
	//
	// Parameters:
	//   - password: The plain text password to be hashed
	//
	// Returns:
	//   - string: The hashed password
	//   - error: An error if hashing fails, otherwise nil
	Hash(password string) (string, error)

	// CompareHashAndPassword compares a hashed password with a plain text password.
	//
	// Parameters:
	//   - hashedPassword: The hashed password to compare against
	//   - password: The plain text password to verify
	//
	// Returns:
	//   - bool: true if the passwords match, false otherwise
	CompareHashAndPassword(hashedPassword, password string) bool
}

type hashPassword struct {
}

// NewPasswordHashing creates a new instance of the PasswordHashing implementation.
//
// Returns:
//   - PasswordHashing: A new PasswordHashing instance
func NewPasswordHashing() PasswordHashing {
	return &hashPassword{}
}

// Hash generates a hashed version of the given password using bcrypt.
//
// Parameters:
//   - password: The plain text password to be hashed
//
// Returns:
//   - string: The hashed password
//   - error: An error if hashing fails, otherwise nil
func (h *hashPassword) Hash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrCannotGenerateHash
	}

	return string(hashBytes), nil
}

// CompareHashAndPassword compares a hashed password with a plain text password using bcrypt.
//
// Parameters:
//   - hashedPassword: The hashed password to compare against
//   - password: The plain text password to verify
//
// Returns:
//   - bool: true if the passwords match, false otherwise
func (h *hashPassword) CompareHashAndPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
