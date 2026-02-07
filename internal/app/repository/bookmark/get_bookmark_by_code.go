package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

func (b *bookmarkRepository) GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error) {
	res := &model.Bookmark{}
	err := b.db.WithContext(ctx).Where("code = ?", code).First(res).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return res, nil
}
