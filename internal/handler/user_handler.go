package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/service"
	"github.com/vukieuhaihoa/bookmark-management/pkg/response"
)

// User defines the interface for user-related HTTP handlers.
// It provides methods for handling user creation requests using the Gin web framework.
type User interface {
	// CreateUser is a Gin framework handler that creates a new user.
	// It processes HTTP requests and returns the created user or an error.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	CreateUser(c *gin.Context)
}

type user struct {
	userSvc service.User
}

// NewUser creates a new instance of the User handler.
// It accepts a user service implementation and returns a handler
// that can process HTTP requests for user creation.
//
// Parameters:
//   - userSvc: The user service used for user-related operations
//
// Returns:
//   - User: A new user handler instance
func NewUser(userSvc service.User) User {
	return &user{userSvc: userSvc}
}

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
// @Router       /v1/user/register [post]
func (u *user) CreateUser(c *gin.Context) {
	input := &createUserRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	createdUser, err := u.userSvc.CreateUser(c, input.Username, input.Password, input.DisplayName, input.Email)
	if err != nil {
		log.Error().Err(err).Msg("service return error when create user")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(201, &createUserResponse{
		Data:    createdUser,
		Message: "Register an user successfully!",
	})
}
