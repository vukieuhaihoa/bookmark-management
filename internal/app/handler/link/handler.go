// Package link provides HTTP handlers for link-related operations.
// It includes handlers for shortening URLs and retrieving original URLs
// using the Gin web framework.
package link

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/service/link"
)

var (
	ErrCodeNotFound       = errors.New("code not found")
	InValidRequestPayload = errors.New("invalid request payload")
	InternalServerError   = errors.New("internal server error")
)

// Handler is the interface for the shorten URL handler.
type Handler interface {
	// ShortenURL handles the URL shortening request.
	// It takes a Gin context as input and processes the request to generate a shortened URL.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	ShortenURL(c *gin.Context)

	// GetURL handles the request to retrieve the original URL from a shortened code.
	// It takes a Gin context as input and processes the request to fetch the original URL.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	GetURL(c *gin.Context)
}

// linkHandler is the concrete implementation of the Handler interface.
type linkHandler struct {
	linkSvc link.Service
}

// NewLinkHandler creates a new instance of the LinkHandler.
//
// Parameters:
//   - linkSvc: The link service used for link operations
//
// Returns:
//   - Handler: A new shorten URL handler instance
func NewLinkHandler(linkSvc link.Service) Handler {
	return &linkHandler{
		linkSvc: linkSvc,
	}
}
