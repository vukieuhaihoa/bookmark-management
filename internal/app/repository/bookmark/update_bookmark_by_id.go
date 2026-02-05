package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

// UpdateBookmarkByID updates an existing bookmark in the database by its ID.
// It takes a context, an ID, and a bookmark model with updated details as input.
// Returns an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - id: The ID of the bookmark to be updated.
//   - updatedBookmark: The bookmark model containing the updated details.
//
// Returns:
//   - error: An error if the update fails, otherwise nil.
func (b *bookmarkRepository) UpdateBookmarkByID(ctx context.Context, id, userID string, updatedBookmark *model.Bookmark) error {
	result := b.db.WithContext(ctx).Model(&model.Bookmark{}).Where("id = ? and user_id = ?", id, userID).Updates(updatedBookmark)
	if result.Error != nil {
		return dbutils.CatchDBError(result.Error)
	}

	if result.RowsAffected == 0 {
		return dbutils.ErrRecordNotFoundType
	}

	return nil
}
