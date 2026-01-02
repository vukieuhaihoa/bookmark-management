package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/bookmark-management/internal/service"
)

type ShortenURL interface {
}

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

func (h *shortenURLHandler) ShortenURL(c *gin.Context) {
	// Implementation goes here
	input := &shortenURLRequest{}
	err := c.ShouldBindJSON(input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	code, err := h.shortenURLSvc.ShortenURL(c, input.URL)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, shortenURLResponse{
		Code:    code,
		Message: "Shorten URL generated successfully!",
	})
}
