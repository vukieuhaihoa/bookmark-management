// Package link provides functionality for URL shortening and retrieval services.
// It defines the ShortenURL interface and its implementation.
// The ShortenURL service uses a repository for storing and retrieving URLs,
// and a code generator for creating unique shortened URL codes.
// This package is essential for managing URL shortening operations in the application.
package link

import (
	"context"
	"errors"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/repository/bookmark"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/repository/link"
	"github.com/vukieuhaihoa/bookmark-management/pkg/utils"
)

const (
	defaultURLCodeLength = 8
	maxRetryAttempts     = 10
)

var (
	ErrCodeNotFound       = errors.New("shortened URL not found")
	ErrMaxRetriesExceeded = errors.New("maximum retry attempts exceeded for generating unique URL code")
)

// Service defines the interface for shortening URLs.
//
//go:generate mockery --name Service --filename shorten_url_service.go --output ./mocks
type Service interface {
	// ShortenURL generates a shortened URL code for the given original URL.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellations
	//   - originalURL: The original URL to be shortened
	//   - expireIn: The expiration time in seconds for the shortened URL
	//
	// Returns:
	//   - string: The generated shortened URL code
	//   - error: An error object if the shortening operation fails, otherwise nil
	ShortenURL(ctx context.Context, originalURL string, expireIn int) (string, error)

	// GetURL retrieves the original URL associated with the given shortened URL code.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellations
	//   - urlCode: The shortened URL code
	//
	// Returns:
	//   - string: The original URL associated with the shortened code
	//   - error: An error object if the retrieval operation fails, otherwise nil
	GetURL(ctx context.Context, urlCode string) (string, error)
}

// linkService implements the Service interface and provides methods for shortening URLs.
type linkService struct {
	repo          link.Repository
	randomCodeGen utils.CodeGenerator
	bookmarkRepo  bookmark.Repository
}

// NewLinkService creates a new instance of Service.
//
// Parameters:
//   - repo: The repository used for URL storage operations
//   - randomCodeGen: The code generator used for generating random codes
//   - bookmarkRepo: The bookmark repository for additional bookmark-related operations
//
// Returns:
//   - Service: The initialized Service instance
func NewLinkService(repo link.Repository, randomCodeGen utils.CodeGenerator, bookmarkRepo bookmark.Repository) Service {
	return &linkService{
		repo:          repo,
		randomCodeGen: randomCodeGen,
		bookmarkRepo:  bookmarkRepo,
	}
}
