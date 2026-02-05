package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

// DeleteBookmarkByID deletes a bookmark from the database by its ID and user ID.
// It takes a context, an ID, and a user ID as input.
// Returns an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - id: The ID of the bookmark to be deleted.
//   - userID: The ID of the user who owns the bookmark.
//
// Returns:
//   - error: An error if the deletion fails, otherwise nil.
func (b *bookmarkRepository) DeleteBookmarkByID(ctx context.Context, id, userID string) error {
	result := b.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Bookmark{})
	if result.Error != nil {
		// log.Error().Err(result.Error).Msg("failed to delete bookmark")
		return dbutils.CatchDBError(result.Error)
	}

	if result.RowsAffected == 0 {
		// log.Warn().Str("bookmark_id", id).Str("user_id", userID).Msg("no bookmark found to delete")
		return dbutils.ErrRecordNotFoundType
	}

	return nil
}
