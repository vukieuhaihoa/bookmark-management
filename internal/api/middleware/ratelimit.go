package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/repository/ratelimit"
	"github.com/vukieuhaihoa/bookmark-management/pkg/utils"
)

const (
	RateLimitKeyFormat = "rate_limit:%s"

	RateLimitIPKey       = "ip"
	IP_RateLimitInterval = 1 * time.Minute
	IP_RateLimitMaxCount = 100

	RateLimitUserIDKey       = "user_id"
	UserID_RateLimitInterval = 1 * time.Second
	UserID_RateLimitMaxCount = 20
)

// keyFilter maps rate limit keys to functions that extract the relevant value and settings from the Gin context.
var keyFilter = map[string]func(c *gin.Context) (string, time.Duration, int){
	RateLimitIPKey: func(c *gin.Context) (string, time.Duration, int) {
		return c.ClientIP(), IP_RateLimitInterval, IP_RateLimitMaxCount
	},

	RateLimitUserIDKey: func(c *gin.Context) (string, time.Duration, int) {
		userID, err := utils.GetUserIDFromJWTClaims(c)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get user ID from JWT claims")
			return c.ClientIP(), IP_RateLimitInterval, IP_RateLimitMaxCount
		}
		return userID, UserID_RateLimitInterval, UserID_RateLimitMaxCount
	},
}

// RateLimit defines the contract for the rate limiting middleware.
type RateLimit interface {
	// RateLimit returns a Gin middleware handler function for rate limiting.
	//
	// The middleware checks the current request count for the client's IP address.
	// If the count exceeds the defined maximum, it responds with a 429 Too Many Requests status.
	// Otherwise, it increments the request count and allows the request to proceed.
	//
	// Parameters:
	//   - key: A string key used to identify the rate limit counter (e.g., based on client IP).
	//
	// Returns:
	//   - gin.HandlerFunc: The Gin middleware handler function for rate limiting.
	RateLimit(key string) gin.HandlerFunc
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

// RateLimit returns a Gin middleware handler function for rate limiting.
//
// The middleware checks the current request count for the client's IP address.
// If the count exceeds the defined maximum, it responds with a 429 Too Many Requests status.
// Otherwise, it increments the request count and allows the request to proceed.
//
// Returns:
//   - gin.HandlerFunc: The Gin middleware handler function for rate limiting.
func (r *rateLimit) RateLimit(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn, ok := keyFilter[key]
		if !ok {
			log.Warn().Str("key", key).Msg("Invalid rate limit key, defaulting to IP")
			fn = keyFilter[RateLimitIPKey]
		}

		value, exp, maxCount := fn(c)
		rateLimitKey := fmt.Sprintf(RateLimitKeyFormat, value)

		currRate, err := r.repo.GetCurrentRateLimit(c, rateLimitKey)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get current rate limit")
		}

		if currRate != -1 && currRate >= maxCount {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please try again later."})
			c.Abort()
			return
		}

		if err := r.repo.IncreaseRateLimit(c, rateLimitKey, exp); err != nil {
			log.Error().Err(err).Msg("Failed to increase rate limit")
		}

		c.Next()
	}
}
