package bookmark

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/handler/utils"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-management/pkg/response"
)

// allowedSortFields maps the allowed sorting fields from query parameters to database column names.
var allowedSortFields = map[string]struct{}{
	"created_at":  {},
	"updated_at":  {},
	"url":         {},
	"description": {},
}

// PagingRequest represents the expected query parameters for pagination and sorting when listing bookmarks.
type PagingRequest struct {
	Page  int    `form:"page" binding:"gte=1" default:"1" example:"1"`
	Limit int    `form:"limit" binding:"gte=1,lte=50" default:"5" example:"5"`
	Sort  string `form:"sort" example:"-created_at,updated_at" default:"-created_at"`
}

// ListBookmarksResponse represents the response returned when listing bookmarks.
type ListBookmarksResponse struct {
	Data       []*model.Bookmark `json:"data"`
	Pagination *common.Paging    `json:"pagination"`
}

// ListBookmarks handles the HTTP request to list bookmarks for the authenticated user.
// It supports pagination and sorting based on query parameters.
// @Summary      List bookmarks
// @Description  Retrieve a list of bookmarks for the authenticated user with pagination and sorting
// @Tags         Bookmarks
// @Produce      json
// @Param        page   query     int    false  "Page number"        default(1)    example(1)
// @Param        limit  query     int    false  "Number of items per page"  default(5)    example(5)
// @Param        sort   query     string false  "Sorting criteria("-" for descending), e.g., -created_at,updated_at"  default(-created_at)  example(-created_at,description)
// @Success      200    {object}  ListBookmarksResponse
// @Failure      400    {object}  response.Message
// @Failure      401    {object}  response.Message
// @Failure      500    {object}  response.Message
// @Security		 Bearer
// @Router       /v1/bookmarks [get]
func (b *bookmarkHandler) ListBookmarks(c *gin.Context) {
	input := &PagingRequest{}
	if err := c.ShouldBindQuery(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	userID, err := utils.GetUserIDFromJWTClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse)
		return
	}

	parsedSortFields, err := common.ParseSortParams(input.Sort, allowedSortFields)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.InvalidSortedFieldResponse)
		return
	}

	queryOptions := common.NewQueryOptions(input.Page, input.Limit, parsedSortFields)

	res, err := b.svc.ListBookmarks(c, userID, queryOptions)
	switch {
	case errors.Is(err, dbutils.ErrInvalidSortField):
		c.JSON(http.StatusBadRequest, response.InputErrorResponse)
		return
	case errors.Is(err, nil):
	default:
		log.Error().
			Str("operation", "ListBookmarks").
			Err(err).
			Msg("service return error when list bookmarks")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, &ListBookmarksResponse{
		Data: res,
		Pagination: &common.Paging{
			Page:  queryOptions.Paging.Page,
			Limit: queryOptions.Paging.Limit,
			Total: queryOptions.Total,
		},
	})
}
