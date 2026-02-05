package bookmark

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/handler/utils"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-management/pkg/response"
)

// createBookmarkRequest represents the expected payload for creating a new bookmark.
type createBookmarkRequest struct {
	Description string `json:"description" binding:"required" example:"A sample bookmark"`
	URL         string `json:"url" binding:"required,url" example:"https://example.com"`
}

// createBookmarkResponse represents the response returned after creating a new bookmark.
type createBookmarkResponse struct {
	Data    *model.Bookmark `json:"data"`
	Message string          `json:"message"`
}

// CreateBookmark generates a Gin framework handler that creates a new bookmark for the authenticated user.
// @Summary      Create a new bookmark
// @Description  Create a new bookmark for the authenticated user
// @Tags         Bookmarks
// @Accept       json
// @Produce      json
// @Param        request  body      createBookmarkRequest  true  "Create Bookmark Request"
// @Success      201      {object}  createBookmarkResponse
// @Failure      400      {object}  response.Message
// @Failure      401      {object}  response.Message
// @Failure      500      {object}  response.Message
// @Security     Bearer
// @Router       /v1/bookmarks [post]
func (h *bookmarkHandler) CreateBookmark(c *gin.Context) {
	// Implementation of the handler to create a bookmark
	input := &createBookmarkRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputErrorResponse)
		return
	}

	userID, err := utils.GetUserIDFromJWTClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse)
		return
	}

	bookmark, err := h.svc.CreateBookmark(c, input.URL, input.Description, userID)
	switch {
	case errors.Is(err, dbutils.ErrDuplicationType):
		c.JSON(http.StatusBadRequest, response.InputErrorResponse)
		return
	case errors.Is(err, dbutils.ErrForeignKeyType):
		c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse)
		return
	case errors.Is(err, nil):
	default:
		log.Error().
			Str("operation", "CreateBookmark").
			Err(err).
			Msg("service return error when create bookmark")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusCreated, &createBookmarkResponse{
		Data:    bookmark,
		Message: "Create a bookmark successfully!",
	})

}
