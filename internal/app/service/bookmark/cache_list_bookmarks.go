package bookmark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
)

type listBookmarksCacheData struct {
	Bookmarks []*model.Bookmark `json:"bookmarks"`
	Total     int64             `json:"total"`
}

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
func (b *bookmarkServiceWithCache) ListBookmarks(ctx context.Context, userID string, opts *common.QueryOptions) ([]*model.Bookmark, error) {
	groupKey := fmt.Sprintf(ListBookmarksCacheGroupKey, userID)
	cacheKey := GenerateCacheKeyFromQueryOptions(opts)

	data := &listBookmarksCacheData{}
	cachedData, err := b.cache.GetCacheData(ctx, groupKey, cacheKey)
	if err == nil && len(cachedData) > 0 {
		err := json.Unmarshal(cachedData, &data)
		if err == nil {
			opts.Total = data.Total
			return data.Bookmarks, nil
		}
	}

	bookmarks, err := b.svc.ListBookmarks(ctx, userID, opts)
	if err != nil {
		return nil, err
	}

	bookmarksBytes, err := json.Marshal(&listBookmarksCacheData{
		Bookmarks: bookmarks,
		Total:     opts.Total,
	})
	if err == nil {
		cacheErr := b.cache.SetCacheData(ctx, groupKey, cacheKey, bookmarksBytes, ListBookmarksCacheTTL)
		if cacheErr != nil {
			log.Error().Err(cacheErr).Msg("failed to set cache for list bookmarks")
		}
	}

	return bookmarks, nil
}

func GenerateCacheKeyFromQueryOptions(opts *common.QueryOptions) string {
	var strBuilder bytes.Buffer
	for _, sort := range opts.Sorting {
		strBuilder.WriteString(sort.Field)
		strBuilder.WriteString("_")
		strBuilder.WriteString(string(sort.Direction))
		strBuilder.WriteString("_")
	}

	return fmt.Sprintf(ListBookmarksCacheKeyFormat, opts.Page, opts.Limit, strBuilder.String())
}
