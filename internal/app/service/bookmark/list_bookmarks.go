package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
)

// ListBookmarks retrieves a list of bookmarks for a specific user based on the provided query options.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - userID: The ID of the user whose bookmarks are to be retrieved.
//   - opts: Query options for filtering, sorting, and pagination.
//
// Returns:
//   - []*model.Bookmark: A slice of bookmark models.
//   - error: An error if the retrieval fails.
func (b *bookmarkService) ListBookmarks(ctx context.Context, userID string, opts *common.QueryOptions) ([]*model.Bookmark, error) {
	return b.repo.ListBookmarks(ctx, userID, opts)
}
