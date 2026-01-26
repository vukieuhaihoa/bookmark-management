package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/service"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
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

	// Login is a Gin framework handler that authenticates a user and returns a JWT token.
	// It processes HTTP requests and returns the token or an error.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	Login(c *gin.Context)

	// GetProfile is a Gin framework handler that retrieves the profile of the authenticated user.
	// It processes HTTP requests and returns the user profile or an error.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	GetProfile(c *gin.Context)

	// UpdateProfile is a Gin framework handler that updates the profile of the authenticated user.
	// It processes HTTP requests and returns a success message or an error.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	UpdateProfile(c *gin.Context)
}

// user is the concrete implementation of the User handler interface.
type user struct {
	userSvc service.User
}

// NewUser creates a new instance of the User handler.
// It accepts a user service implementation and returns a handler
// that can process HTTP requests for user creation and login.
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
// @Router       /v1/users/register [post]
func (u *user) CreateUser(c *gin.Context) {
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

type loginRequest struct {
	Username string `json:"username" binding:"required" example:"testuser001"`
	Password string `json:"password" binding:"required,gte=8" example:"my_SECURE_password123@"`
}

type loginResponse struct {
	Data    string `json:"data"`
	Message string `json:"message"`
}

// Login generates a Gin framework handler that authenticates a user and returns a JWT token.
// @Summary      User login
// @Description  Authenticate a user and return a JWT token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        credentials  body      loginRequest  true  "User credentials"
// @Success      200          {object}  loginResponse
// @Failure      400          {object}  response.Message
// @Failure      401          {object}  response.Message
// @Failure      500          {object}  response.Message
// @Router       /v1/users/login [post]
func (u *user) Login(c *gin.Context) {
	// Implementation for user login handler goes here
	input := &loginRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	token, err := u.userSvc.Login(c, input.Username, input.Password)
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusBadRequest, response.Message{
			Message: err.Error(),
		})
		return
	case errors.Is(err, dbutils.ErrRecordNotFoundType):
		c.JSON(http.StatusBadRequest, response.Message{
			Message: "invalid username or password",
		})
		return
	case errors.Is(err, nil):
	default:
		log.Error().
			Str("operation", "Login").
			Err(err).
			Msg("service return error when login user")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return

	}

	c.JSON(http.StatusOK, &loginResponse{
		Data:    token,
		Message: "Logged in successfully!",
	})
}

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
func (u *user) GetProfile(c *gin.Context) {
	// Implementation for user profile retrieval handler goes here
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse)
		return
	}

	userID, ok := userIDValue.(string)
	if !ok {
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
func (u *user) UpdateProfile(c *gin.Context) {
	// Implementation for user profile update handler goes here
	input := &updateProfileRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse)
		return
	}

	userID, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse)
		return
	}

	err := u.userSvc.UpdateUserByID(c, userID, input.DisplayName, input.Email)
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
