package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/repository/bookmark"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
)

const DEFAULT_CODE_LENGTH = 10

// Service defines the interface for bookmark-related business logic.
//
//go:generate mockery --name=Service --filename=service.go --output=./mocks
type Service interface {
	// CreateBookmark creates a new bookmark with the provided information.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - url: The URL of the bookmark.
	//   - description: A description of the bookmark.
	//   - userID: The ID of the user who owns the bookmark.
	//
	// Returns:
	//   - *model.Bookmark: The created bookmark model.
	//   - error: An error if the creation fails, otherwise nil.
	CreateBookmark(ctx context.Context, url, description, userID string) (*model.Bookmark, error)

	// ListBookmarks retrieves a list of bookmarks for a specific user with optional query options.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - userID: The ID of the user whose bookmarks are to be retrieved.
	//   - opts: QueryOptions containing pagination and sorting details.
	//
	// Returns:
	//   - []*model.Bookmark: A slice of bookmark models.
	//   - error: An error if the retrieval fails, otherwise nil.
	ListBookmarks(ctx context.Context, userID string, opts *common.QueryOptions) ([]*model.Bookmark, error)

	// UpdateBookmarkByID updates a bookmark by its ID and user ID.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - id: The ID of the bookmark to be updated.
	//   - userID: The ID of the user who owns the bookmark.
	//   - updatedBookmark: The bookmark model containing the updated details.
	//
	// Returns:
	//   - error: An error if the update fails, otherwise nil.
	UpdateBookmarkByID(ctx context.Context, id, userID string, updatedBookmark *model.Bookmark) error

	// DeleteBookmarkByID deletes a bookmark by its ID and user ID.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - id: The ID of the bookmark to be deleted.
	//   - userID: The ID of the user who owns the bookmark.
	//
	// Returns:
	//   - error: An error if the deletion fails, otherwise nil.
	DeleteBookmarkByID(ctx context.Context, id, userID string) error
}

// bookmarkService is the concrete implementation of the Service interface.
type bookmarkService struct {
	repo bookmark.Repository
}

// NewBookmarkService creates a new instance of bookmarkService.
// It takes a bookmark.Repository and a utils.CodeGenerator as parameters
// and returns a Service interface.
func NewBookmarkService(repo bookmark.Repository) Service {
	return &bookmarkService{
		repo: repo,
	}
}
