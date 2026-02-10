package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-management/pkg/encoding"
	"gorm.io/gorm"
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

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("CodeShorten").Create(bookmark).Error; err != nil {
			return err
		}

		// Reload code_shorten assigned by BIGSERIAL
		if err := tx.Model(&model.Bookmark{}).Where("id = ?", bookmark.ID).Pluck("code_shorten", &bookmark.CodeShorten).Error; err != nil {
			return err
		}

		encoded, err := encoding.StdEncoding.EncodeInt64ToString(bookmark.CodeShorten)
		if err != nil {
			return err
		}

		bookmark.CodeShortenEncoded = "p_" + encoded

		if err := tx.Model(&model.Bookmark{}).Where("id = ?", bookmark.ID).Update("code_shorten_encoded", bookmark.CodeShortenEncoded).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return bookmark, nil
}
