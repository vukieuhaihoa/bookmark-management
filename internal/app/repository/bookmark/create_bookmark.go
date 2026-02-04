package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

// CreateBookmark creates a new bookmark in the database.
// It takes a context and a bookmark model as input.
// Returns the created bookmark and an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - bookmark: The bookmark model to be created.
//
// Returns:
//   - *model.Bookmark: The created bookmark model.
//   - error: An error if the creation fails, otherwise nil.
func (r *bookmarkRepository) CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) {
	if err := r.db.WithContext(ctx).Create(bookmark).Error; err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return bookmark, nil
}
