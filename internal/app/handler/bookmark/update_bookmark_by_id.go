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

// updateBookmarkRequest represents the expected JSON body for updating a bookmark.
type updateBookmarkRequest struct {
	Description string `json:"description" binding:"required"`
	URL         string `json:"url" binding:"required,url"`
}

// UpdateBookmarkByID handles the HTTP request to update a bookmark by its ID for the authenticated user.
// @Summary      Update bookmark by ID
// @Description  Update bookmark by ID for authenticated user
// @Tags         Bookmarks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Bookmark ID"
// @Param        body body      updateBookmarkRequest  true  "Bookmark update info"
// @Success      200  {object}  response.Message
// @Failure      400  {object}  response.Message
// @Failure      401  {object}  response.Message
// @Failure      500  {object}  response.Message
// @Security     BearerAuth
// @Router       /v1/bookmarks/{id} [put]
func (b *bookmarkHandler) UpdateBookmarkByID(c *gin.Context) {
	input := &updateBookmarkRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
	}

	bookmarkID := c.Param("id")
	if bookmarkID == "" {
		c.JSON(http.StatusBadRequest, response.Message{
			Message: "Bookmark ID is required",
		})
		return
	}

	userID, err := utils.GetUserIDFromJWTClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse)
		return
	}

	updatedBookmark := &model.Bookmark{
		Description: input.Description,
		URL:         input.URL,
	}

	err = b.svc.UpdateBookmarkByID(c, bookmarkID, userID, updatedBookmark)
	switch {
	case errors.Is(err, dbutils.ErrRecordNotFoundType):
		c.JSON(http.StatusBadRequest, response.InputErrorResponse)
		return
	case err == nil:
	default:
		log.Error().
			Str("operation", "UpdateBookmarkByID").
			Err(err).
			Msg("service return error when update bookmark by ID")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, response.Message{
		Message: "Success",
	})
}
