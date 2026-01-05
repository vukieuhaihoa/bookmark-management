package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-management/internal/service"
)

var (
	ErrCodeNotFound       = errors.New("code not found")
	InValidRequestPayload = errors.New("invalid request payload")
	InternalServerError   = errors.New("internal server error")
)

// ShortenURL is the interface for the shorten URL handler.
type ShortenURL interface {
	// ShortenURL handles the URL shortening request.
	// It takes a Gin context as input and processes the request to generate a shortened URL.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	ShortenURL(c *gin.Context)

	// GetURL handles the request to retrieve the original URL from a shortened code.
	// It takes a Gin context as input and processes the request to fetch the original URL.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	GetURL(c *gin.Context)
}

// shortenURLHandler is the concrete implementation of the ShortenURL interface.
type shortenURLHandler struct {
	shortenURLSvc service.ShortenURL
}

// NewShortenURL creates a new instance of the ShortenURL handler.
//
// Parameters:
//   - shortenURLSvc: The shorten URL service used for URL shortening operations
//
// Returns:
//   - ShortenURL: A new shorten URL handler instance
func NewShortenURL(shortenURLSvc service.ShortenURL) ShortenURL {
	return &shortenURLHandler{
		shortenURLSvc: shortenURLSvc,
	}
}

type shortenURLRequest struct {
	URL      string `json:"url" binding:"required,url"`
	ExpireIn int    `json:"exp" binding:"required,lte=604800"`
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
func (h *shortenURLHandler) ShortenURL(c *gin.Context) {
	// Implementation goes here
	input := &shortenURLRequest{}
	err := c.ShouldBindJSON(input)
	if err != nil || input.URL == "" || input.ExpireIn <= 0 {
		c.JSON(http.StatusBadRequest, shortenURLResponse{
			Message: InValidRequestPayload.Error(),
		})
		return
	}

	code, err := h.shortenURLSvc.ShortenURL(c, input.URL, input.ExpireIn)
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

// GetURL handles the request to retrieve the original URL from a shortened code.
// It takes a Gin context as input and processes the request to fetch the original URL.
//
// Parameters:
//   - c: The Gin context containing the HTTP request and response
//
// @Summary Get original URL
// @Description Retrieves the original URL from a given shortened code.
// @Tags URL
// @Accept json
// @Produce json
// @Param code path string true "Shortened code"
// @Success 200 {object} shortenURLResponse
// @Failure 400 {object} shortenURLResponse
// @Failure 500 {object} shortenURLResponse
// @Router /v1/links/redirect/{code} [get]
func (h *shortenURLHandler) GetURL(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, shortenURLResponse{
			Message: InValidRequestPayload.Error(),
		})
		return
	}

	url, err := h.shortenURLSvc.GetURL(c, code)
	// Handle different error cases with switch statement
	switch {
	case errors.Is(err, service.ErrCodeNotFound):
		c.JSON(http.StatusBadRequest, shortenURLResponse{
			Message: ErrCodeNotFound.Error(),
		})
		return
	case err == nil: // MUST: to redirect when no error
	default:
		log.Error().Str("code", code).Err(err).Msg("service return error when get original url")
		c.JSON(http.StatusInternalServerError, shortenURLResponse{
			Message: InternalServerError.Error(),
		})
		return
	}

	// Handle errors with if else statement
	// if err != nil {
	// 	if errors.Is(err, service.ErrCodeNotFound) {
	// 		c.JSON(http.StatusBadRequest, shortenURLResponse{
	// 			Message: ErrCodeNotFound.Error(),
	// 		})
	// 		return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, shortenURLResponse{
	// 		Message: InternalServerError.Error(),
	// 	})
	// 	return
	// }

	c.Redirect(http.StatusFound, url)
}
