package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
	"gorm.io/gorm"
)

// Repository defines the interface for bookmark data operations.
//
//go:generate mockery --name=Repository --output=./mocks --filename=repo.go
type Repository interface {
	// CreateBookmark creates a new bookmark in the database.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - bookmark: The bookmark model to be created.
	//
	// Returns:
	//   - *model.Bookmark: The created bookmark model.
	//   - error: An error if the creation fails, otherwise nil.
	CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error)

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

	// UpdateBookmarkByID updates an existing bookmark in the database by its ID and user ID.
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

	// DeleteBookmarkByID deletes a bookmark from the database by its ID and user ID.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - id: The ID of the bookmark to be deleted.
	//   - userID: The ID of the user who owns the bookmark.
	//
	// Returns:
	//   - error: An error if the deletion fails, otherwise nil.
	DeleteBookmarkByID(ctx context.Context, id, userID string) error

	// GetBookmarkByCode retrieves a bookmark from the database by its unique code.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - code: The unique code of the bookmark.
	//
	// Returns:
	//   - *model.Bookmark: The bookmark model if found.
	//   - error: An error if the retrieval fails or the bookmark is not found.
	GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error)

	// GetBookmarkByCodeShortenEncoded retrieves a bookmark from the database by its unique code_shorten_encoded.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - code: The unique code_shorten_encoded of the bookmark.
	//
	// Returns:
	//   - *model.Bookmark: The bookmark model if found.
	//   - error: An error if the retrieval fails or the bookmark is not found.
	GetBookmarkByCodeShortenEncoded(ctx context.Context, code string) (*model.Bookmark, error)
}

// bookmarkRepository implements the Repository interface for bookmark data operations.
type bookmarkRepository struct {
	db *gorm.DB
}

// NewBookmarkRepository creates a new instance of bookmarkRepository.
// It takes a gorm.DB instance as a parameter and returns a Repository interface.
func NewBookmarkRepository(db *gorm.DB) Repository {
	return &bookmarkRepository{
		db: db,
	}
}
