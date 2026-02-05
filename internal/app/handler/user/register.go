package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-management/pkg/response"
)

type createUserRequest struct {
	Username    string `json:"username" binding:"required" example:"testuser001"`
	Password    string `json:"password" binding:"required,min=8,password_strength" example:"my_SECURE_password123@"`
	DisplayName string `json:"display_name" binding:"required" example:"Test User"`
	Email       string `json:"email" binding:"required,email" example:"testuser001@example.com"`
}

type createUserResponse struct {
	Data    *model.User `json:"data"`
	Message string      `json:"message"`
}

// CreateUser generates a Gin framework handler that creates a new user.
// @Summary      Create a new user
// @Description  Create a new user with the provided information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      createUserRequest  true  "User to create"
// @Success      201   {object}  createUserResponse
// @Failure      400   {object}  response.Message
// @Failure      500   {object}  response.Message
// @Router       /v1/users/register [post]
func (u *userHandler) CreateUser(c *gin.Context) {
	input := &createUserRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	createdUser, err := u.userSvc.CreateUser(c, input.Username, input.Password, input.DisplayName, input.Email)
	switch {
	case errors.Is(err, dbutils.ErrDuplicationType):
		c.JSON(http.StatusBadRequest, response.Message{
			Message: "username or email already exists",
		})
		return
	case errors.Is(err, nil):
	default:
		log.Error().
			Str("operation", "CreateUser").
			Err(err).
			Msg("service return error when create user")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusCreated, &createUserResponse{
		Data:    createdUser,
		Message: "Register an user successfully!",
	})
}
