package link

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type shortenURLRequest struct {
	URL      string `json:"url" binding:"required,url" example:"https://www.example.com"`
	ExpireIn int    `json:"exp" binding:"required,lte=604800" example:"10000"`
}

type shortenURLResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

// ShortenURL handles the URL shortening request.
// It takes a Gin context as input and processes the request to generate a shortened URL.
//
// Parameters:
//   - c: The Gin context containing the HTTP request and response
//
// @Summary Shorten a URL
// @Description Shortens a given URL and returns a unique code.
// @Tags URL
// @Accept json
// @Produce json
// @Param shortenURLRequest body shortenURLRequest true "URL to shorten"
// @Success 200 {object} shortenURLResponse
// @Failure 400 {object} shortenURLResponse
// @Failure 500 {object} shortenURLResponse
// @Router /v1/links/shorten [post]
func (h *linkHandler) ShortenURL(c *gin.Context) {
	// Implementation goes here
	input := &shortenURLRequest{}
	err := c.ShouldBindJSON(input)
	if err != nil || input.URL == "" || input.ExpireIn <= 0 {
		c.JSON(http.StatusBadRequest, shortenURLResponse{
			Message: InValidRequestPayload.Error(),
		})
		return
	}

	code, err := h.linkSvc.ShortenURL(c, input.URL, input.ExpireIn)
	if err != nil {
		log.Error().Str("url", input.URL).Err(err).Msg("service return error when shorten url")
		c.JSON(http.StatusInternalServerError, shortenURLResponse{
			Message: InternalServerError.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, shortenURLResponse{
		Code:    code,
		Message: "Shorten URL generated successfully!",
	})
}
