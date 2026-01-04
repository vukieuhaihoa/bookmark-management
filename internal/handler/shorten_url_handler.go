package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/bookmark-management/internal/service"
)

var (
	inValidRequestPayload = errors.New("invalid request payload")
	internalServerError   = errors.New("internal server error")
)

// ShortenURL is the interface for the shorten URL handler.
type ShortenURL interface {
	// ShortenURL handles the URL shortening request.
	// It takes a Gin context as input and processes the request to generate a shortened URL.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	ShortenURL(c *gin.Context)
}

// shortenURLHandler is the concrete implementation of the ShortenURL interface.
type shortenURLHandler struct {
	shortenURLSvc service.ShortenURL
}

func NewShortenURL(shortenURLSvc service.ShortenURL) ShortenURL {
	return &shortenURLHandler{
		shortenURLSvc: shortenURLSvc,
	}
}

type shortenURLRequest struct {
	URL      string `json:"url" binding:"required"`
	ExpireIn int    `json:"exp" binding:"required"`
}

type shortenURLResponse struct {
	Code    string `json:"code"`
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
func (h *shortenURLHandler) ShortenURL(c *gin.Context) {
	// Implementation goes here
	input := &shortenURLRequest{}
	err := c.ShouldBindJSON(input)
	if err != nil || input.URL == "" || input.ExpireIn <= 0 {
		c.JSON(http.StatusBadRequest, shortenURLResponse{
			Code:    "",
			Message: inValidRequestPayload.Error(),
		})
		return
	}

	code, err := h.shortenURLSvc.ShortenURL(c, input.URL, input.ExpireIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shortenURLResponse{
			Code:    "",
			Message: internalServerError.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, shortenURLResponse{
		Code:    code,
		Message: "Shorten URL generated successfully!",
	})
}
