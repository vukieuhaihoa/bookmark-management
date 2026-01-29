package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	service "github.com/vukieuhaihoa/bookmark-management/internal/app/service/user"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-management/pkg/response"
)

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
func (u *userHandler) Login(c *gin.Context) {
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
