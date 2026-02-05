// Package random_code_gen provides HTTP request handlers for the bookmark management application.
// This file contains the password handler implementation that exposes password generation
// functionality through HTTP endpoints. It supports both Gin and standard library mux frameworks.
package random_code_gen

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/service/random_code_gen"
)

var (
	// ErrPasswordGenerationFailed indicates that password generation has failed.
	ErrPasswordGenerationFailed = errors.New("password generation failed")
)

// Handler defines the interface for password-related HTTP handlers.
// It provides methods for handling password generation requests using different
// web frameworks (Gin and standard library mux).
type Handler interface {
	// GeneratePassword is a Gin framework handler that generates a random password.
	// It processes HTTP requests and returns the generated password or an error.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	GeneratePassword(c *gin.Context)
}

// passwordHandler implements the Handler interface and provides HTTP handlers for password generation.
type passwordHandler struct {
	svc random_code_gen.Service
}

// NewPasswordHandler creates a new instance of the Handler.
// It accepts a password service implementation and returns a handler
// that can process HTTP requests for password generation.
//
// Parameters:
//   - svc: The password service used for generating passwords
//
// Returns:
//   - Handler: A new password handler instance
func NewPasswordHandler(svc random_code_gen.Service) Handler {
	return &passwordHandler{svc: svc}
}
