package utils

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// CodeGenerator defines the interface for generating random codes.
//
//go:generate mockery --name CodeGenerator --filename code.go --output ./mocks
type CodeGenerator interface {
	// GenerateCode generates a random code of the specified length.
	//
	// Parameters:
	//   - length: The desired length of the generated code
	//
	// Returns:
	//   - string: The generated code
	//   - error: An error if random number generation fails, nil otherwise
	GenerateCode(length int) (string, error)
}

type randomCodeGenerator struct {
}

func NewCodeGenerator() CodeGenerator {
	return &randomCodeGenerator{}
}

// GenerateCode generates a cryptographically secure random code of the specified length.
// It uses a predefined character set consisting of alphanumeric characters (a-z, A-Z, 0-9).
//
// Parameters:
//   - length: The desired length of the generated code
//
// Returns:
//   - string: The generated code
//   - error: An error if random number generation fails, nil otherwise
func (r *randomCodeGenerator) GenerateCode(length int) (string, error) {
	var strBuilder bytes.Buffer

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		strBuilder.WriteByte(charset[randomIndex.Int64()])
	}
	return strBuilder.String(), nil
}
