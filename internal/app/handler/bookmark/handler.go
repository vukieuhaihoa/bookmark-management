package bookmark

import (
	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/service/bookmark"
)

// Handler defines the interface for bookmark-related HTTP handlers.
type Handler interface {

	// CreateBookmark handles the HTTP request to create a new bookmark for the authenticated user.
	CreateBookmark(c *gin.Context)

	// ListBookmarks handles the HTTP request to list bookmarks for the authenticated user.
	ListBookmarks(c *gin.Context)

	// UpdateBookmarkByID handles the HTTP request to update a bookmark by its ID for the authenticated user.
	UpdateBookmarkByID(c *gin.Context)

	// DeleteBookmarkByID handles the HTTP request to delete a bookmark by its ID for the authenticated user.
	DeleteBookmarkByID(c *gin.Context)
}

// bookmarkHandler is the concrete implementation of the Handler interface.
type bookmarkHandler struct {
	svc bookmark.Service
}

// NewBookmarkHandler creates a new instance of bookmarkHandler.
// It takes a bookmark.Service as a parameter and returns a Handler interface.
func NewBookmarkHandler(svc bookmark.Service) Handler {
	return &bookmarkHandler{
		svc: svc,
	}
}
