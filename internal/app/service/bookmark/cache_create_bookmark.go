package bookmark

import (
	"context"
	"fmt"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
)

// CreateBookmark creates a new bookmark in the database.
// It takes a context, URL, description, and userID as input.
// Returns the created bookmark model and an error if the operation fails.
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
func (b *bookmarkServiceWithCache) CreateBookmark(ctx context.Context, url, description, userID string) (*model.Bookmark, error) {
	err := b.cache.DelCacheData(ctx, fmt.Sprintf(ListBookmarksCacheGroupKey, userID))
	if err != nil {
		return nil, err
	}

	return b.svc.CreateBookmark(ctx, url, description, userID)
}
