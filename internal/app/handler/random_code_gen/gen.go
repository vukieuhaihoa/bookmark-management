package random_code_gen

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GeneratePassword is a Gin framework handler that generates a random password.
// It calls the password service to generate a cryptographically secure password
// and returns it as plain text in the HTTP response.
//
// Parameters:
//   - c: The Gin context containing the HTTP request and response
//
// Response:
//   - 200 OK: Returns the generated password as plain text
//   - 500 Internal Server Error: Returns an error message if password generation fails
//
// @Summary Generate a random password
// @Description Generates a cryptographically secure random password of fixed length.
// @Tags password
// @Produce plain
// @Success 200 {object} string "w7h3Q9FeXskn"
// @Failure 500 {object} string "Internal Server Error"
// @Router /v1/generate-password [get]
func (h *passwordHandler) GeneratePassword(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		log.Error().Err(err).Msg("service return error when generate password")
		c.String(http.StatusInternalServerError, ErrPasswordGenerationFailed.Error())
		return
	}
	c.String(http.StatusOK, pass)
}
