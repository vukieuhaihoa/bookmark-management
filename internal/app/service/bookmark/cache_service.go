package bookmark

import (
	"time"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/repository/cache"
)

const (
	ListBookmarksCacheGroupKey  = "list_bookmarks_%s"         // list_bookmarks_{userID}
	ListBookmarksCacheKeyFormat = "page_%d_size_%d_sortby_%s" // page_{page}_size_{size}_sort_{sort}
	ListBookmarksCacheTTL       = 24 * time.Hour              // 24 hours TTL for cached bookmark lists
)

type bookmarkServiceWithCache struct {
	svc   Service
	cache cache.DB
}

func NewBookmarkServiceWithCache(svc Service, cache cache.DB) Service {
	return &bookmarkServiceWithCache{
		svc:   svc,
		cache: cache,
	}
}
