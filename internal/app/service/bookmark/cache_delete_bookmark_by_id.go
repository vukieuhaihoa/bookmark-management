package bookmark

import (
	"context"
	"fmt"
)

// DeleteBookmarkByID deletes a bookmark from the database by its ID.
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
func (b *bookmarkServiceWithCache) DeleteBookmarkByID(ctx context.Context, id, userID string) error {
	err := b.cache.DelCacheData(ctx, fmt.Sprintf(ListBookmarksCacheGroupKey, userID))
	if err != nil {
		return err
	}
	return b.svc.DeleteBookmarkByID(ctx, id, userID)
}
