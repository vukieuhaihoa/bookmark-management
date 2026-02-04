package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
)

// CreateBookmark creates a new bookmark with the provided information.
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - url: The URL of the bookmark.
//   - description: A description of the bookmark.
//   - userID: The ID of the user who owns the bookmark.
//
// Returns:
//   - *model.Bookmark: The created bookmark model.
//   - error: An error if the creation fails, otherwise nil.
func (s *bookmarkService) CreateBookmark(ctx context.Context, url, description, userID string) (*model.Bookmark, error) {
	code, err := s.codeGen.GenerateCode(DEFAULT_CODE_LENGTH)
	if err != nil {
		return nil, err
	}

	newBookmark := &model.Bookmark{
		URL:         url,
		Description: description,
		Code:        code,
		UserID:      userID,
	}

	createdBookmark, err := s.repo.CreateBookmark(ctx, newBookmark)
	if err != nil {
		return nil, err
	}

	return createdBookmark, nil
}
