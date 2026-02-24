package bookmark

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-management/pkg/utils"
)

// DeleteBookmarkByID handles the HTTP request to delete a bookmark by its ID for the authenticated user.
// @Summary      Delete bookmark by ID
// @Description  Delete bookmark by ID for authenticated user
// @Tags         Bookmarks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Bookmark ID"
// @Success      200  {object}  common.Message
// @Failure      400  {object}  common.Message
// @Failure      401  {object}  common.Message
// @Failure      500  {object}  common.Message
// @Security     Bearer
// @Router       /v1/bookmarks/{id} [delete]
func (b *bookmarkHandler) DeleteBookmarkByID(c *gin.Context) {
	bookmarkID := c.Param("id")
	if bookmarkID == "" {
		c.JSON(http.StatusBadRequest, common.Message{
			Message: "Bookmark ID is required",
		})
		return
	}

	userID, err := utils.GetUserIDFromJWTClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	err = b.svc.DeleteBookmarkByID(c, bookmarkID, userID)
	switch {
	case errors.Is(err, dbutils.ErrRecordNotFoundType):
		c.JSON(http.StatusBadRequest, common.InputErrorResponse)
		return
	case err == nil:
	default:
		log.Error().Err(err).Msg("Failed to delete bookmark by ID")
		c.JSON(http.StatusInternalServerError, common.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, common.Message{
		Message: "Success",
	})
}
