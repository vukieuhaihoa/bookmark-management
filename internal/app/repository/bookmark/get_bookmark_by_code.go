package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

// GetBookmarkByCode retrieves a bookmark by its unique code.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - code: The unique code of the bookmark to retrieve.
//
// Returns:
//   - A pointer to the Bookmark model if found.
//   - An error if the bookmark is not found or if a database error occurs.
func (b *bookmarkRepository) GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error) {
	res := &model.Bookmark{}
	err := b.db.WithContext(ctx).Where("code = ?", code).First(res).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return res, nil
}
