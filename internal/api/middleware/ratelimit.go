package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/repository/ratelimit"
)

// RateLimit defines the contract for the rate limiting middleware.
type RateLimit interface {
	RateLimit() gin.HandlerFunc
}

// rateLimit is the concrete implementation of the RateLimit interface.
type rateLimit struct {
	repo ratelimit.Repository
}

// NewRateLimit creates a new instance of the rate limiting middleware.
//
// Parameters:
//   - repo: The rate limit repository used to manage rate limit counters.
//
// Returns:
//   - RateLimit: A new rate limiting middleware instance.
func NewRateLimit(repo ratelimit.Repository) RateLimit {
	return &rateLimit{repo: repo}
}

const (
	RateLimitInterval  = 10 * time.Second
	RateLimitMaxCount  = 20
	RateLimitKeyFormat = "rate_limit:%s"
)

// RateLimit returns a Gin middleware handler function for rate limiting.
//
// The middleware checks the current request count for the client's IP address.
// If the count exceeds the defined maximum, it responds with a 429 Too Many Requests status.
// Otherwise, it increments the request count and allows the request to proceed.
//
// Returns:
//   - gin.HandlerFunc: The Gin middleware handler function for rate limiting.
func (r *rateLimit) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf(RateLimitKeyFormat, c.ClientIP())

		currRate, err := r.repo.GetCurrentRateLimit(c, key)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get current rate limit")
		}

		if currRate != -1 && currRate >= RateLimitMaxCount {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please try again later."})
			c.Abort()
			return
		}

		if err := r.repo.IncreaseRateLimit(c, key, RateLimitInterval); err != nil {
			log.Error().Err(err).Msg("Failed to increase rate limit")
		}

		c.Next()
	}
}
