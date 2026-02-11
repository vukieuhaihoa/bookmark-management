package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
	"github.com/vukieuhaihoa/bookmark-management/pkg/jwtutils"
)

// JWTAuth defines the interface for JWT authentication middleware.
type JWTAuth interface {
	JWTAuth() gin.HandlerFunc
}

// jwtAuth is the concrete implementation of the JWTAuth interface.
type jwtAuth struct {
	jwtValidator jwtutils.JWTValidator
}

// NewJWTAuth creates a new instance of the JWT authentication middleware.
//
// Parameters:
//   - jwtValidator: The JWT validator used to validate tokens.
//
// Returns:
//   - JWTAuth: A new JWT authentication middleware instance.
func NewJWTAuth(jwtValidator jwtutils.JWTValidator) JWTAuth {
	return &jwtAuth{
		jwtValidator: jwtValidator,
	}
}

// JWTAuth returns a Gin middleware handler function for JWT authentication.
//
// The middleware extracts the JWT token from the Authorization header,
// validates it, and sets the user ID in the Gin context if valid.
// If the token is missing or invalid, it aborts the request with a 401 Unauthorized response.
//
// Returns:
//   - gin.HandlerFunc: The Gin middleware handler function for JWT authentication.
func (j *jwtAuth) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := parts[1]

		// validate token
		claims, err := j.jwtValidator.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.InvalidTokenResponse)
			return
		}

		// set claims to context
		c.Set("claims", claims)

		// proceed to the next handler
		c.Next()
	}
}
