// Package handler provides HTTP request handlers for the bookmark management application.
// This file contains the password handler implementation that exposes password generation
// functionality through HTTP endpoints. It supports both Gin and standard library mux frameworks.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/bookmark-management/internal/service"
)

type passwordHandler struct {
	svc service.Password
}

// Password defines the interface for password-related HTTP handlers.
// It provides methods for handling password generation requests using different
// web frameworks (Gin and standard library mux).
type Password interface {
	GeneratePassword(c *gin.Context)
}

// NewPassword creates a new instance of the Password handler.
// It accepts a password service implementation and returns a handler
// that can process HTTP requests for password generation.
//
// Parameters:
//   - svc: The password service used for generating passwords
//
// Returns:
//   - Password: A new password handler instance
func NewPassword(svc service.Password) Password {
	return &passwordHandler{svc: svc}
}

// GeneratePassword is a Gin framework handler that generates a random password.
// It calls the password service to generate a cryptographically secure password
// and returns it as plain text in the HTTP response.
//
// Parameters:
//   - c: The Gin context containing the HTTP request and response
//
// Response:
//   - 200 OK: Returns the generated password as plain text
//   - 500 Internal Server Error: Returns an error message if password generation fails
func (h *passwordHandler) GeneratePassword(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, pass)
}
