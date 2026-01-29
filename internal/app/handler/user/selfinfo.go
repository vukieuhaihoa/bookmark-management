package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/handler/utils"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-management/pkg/response"
)

// GetProfile generates a Gin framework handler that retrieves the profile of the authenticated user.
// @Summary      Get user profile
// @Description  Retrieve the profile of the authenticated user
// @Tags         Users
// @Produce      json
// @Success      200  {object}  createUserResponse
// @Failure      401  {object}  response.Message
// @Failure      500  {object}  response.Message
// @Security     Bearer
// @Router       /v1/self/info [get]
func (u *userHandler) GetProfile(c *gin.Context) {
	userID, err := utils.GetUserIDFromJWTClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse)
		return
	}

	user, err := u.userSvc.GetUserByID(c, userID)
	switch {
	case errors.Is(err, dbutils.ErrRecordNotFoundType):
		c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse)
		return
	case errors.Is(err, nil):
	default:
		log.Error().
			Str("operation", "GetProfile").
			Err(err).
			Msg("service return error when get user profile")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, &createUserResponse{
		Data:    user,
		Message: "User profile retrieved successfully!",
	})
}

type updateProfileRequest struct {
	DisplayName string `json:"display_name" binding:"required" example:"Updated User 001"`
	Email       string `json:"email" binding:"required,email" example:"updatedtestuser001@example.com"`
}

// UpdateProfile generates a Gin framework handler that updates the profile of the authenticated user.
// @Summary      Update user profile
// @Description  Update the profile of the authenticated user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        profile  body      updateProfileRequest  true  "Updated user profile"
// @Success      200      {object}  response.Message
// @Failure      400      {object}  response.Message
// @Failure      401      {object}  response.Message
// @Failure      500      {object}  response.Message
// @Security     Bearer
// @Router       /v1/self/info [put]
func (u *userHandler) UpdateProfile(c *gin.Context) {
	// Implementation for user profile update handler goes here
	input := &updateProfileRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	userID, err := utils.GetUserIDFromJWTClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse)
		return
	}

	err = u.userSvc.UpdateUserByID(c, userID, input.DisplayName, input.Email)
	switch {
	case errors.Is(err, dbutils.ErrRecordNotFoundType):
		c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse)
		return
	case errors.Is(err, dbutils.ErrDuplicationType):
		c.JSON(http.StatusBadRequest, response.Message{
			Message: "email already exists",
		})
		return
	case errors.Is(err, nil):
	default:
		log.Error().
			Str("operation", "UpdateProfile").
			Err(err).
			Msg("service return error when update user profile")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, response.Message{
		Message: "Edit current user successfully!",
	})
}
