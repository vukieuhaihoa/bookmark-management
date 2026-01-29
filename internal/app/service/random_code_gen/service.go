// Package random_code_gen provides business logic services for the bookmark management application.
// This file contains the password service implementation for generating cryptographically
// secure random passwords.
package random_code_gen

import "github.com/vukieuhaihoa/bookmark-management/pkg/utils"

const (
	defaultPassLength = 12
)

// passwordService implements the Password interface and provides methods for generating secure random passwords.
type passwordService struct {
	randomCodeGen utils.CodeGenerator
}

// Service defines the interface for password generation service.
// It provides a method for generating a secure random password.
//
//go:generate mockery --name Service --filename pass_service.go --output ./mocks
type Service interface {
	GeneratePassword() (string, error)
}

// NewPasswordService creates a new instance of Service.
// It generates random passwords of fixed length using a predefined character set.
func NewPasswordService(randomCodeGen utils.CodeGenerator) Service {
	return &passwordService{
		randomCodeGen: randomCodeGen,
	}
}
