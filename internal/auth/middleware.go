package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/bitebattle/pkg/utils"
)

const userIDContextKey = "userID"

// Auth0Middleware validates Auth0 JWT and extracts user info from claims
func Auth0Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "authorization header missing")
			c.Abort()
			return
		}

		tokenStr, err := utils.ExtractBearerToken(authHeader)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenStr)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		// Auth0 user ID is the 'sub' claim (string)
		userID := claims.UserID
		if userID == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid user id in token")
			c.Abort()
			return
		}
		c.Set(userIDContextKey, userID)

		// Optionally set email and name if present in claims
		// (Add these fields to Claims struct if needed)

		c.Next()
	}
}

// UserIDFromContext returns the Auth0 user ID (string) from context
func UserIDFromContext(c *gin.Context) (string, error) {
	userID, ok := c.Get(userIDContextKey)
	if !ok {
		return "", http.ErrNoCookie
	}
	idStr, ok := userID.(string)
	if !ok {
		return "", http.ErrNoCookie
	}
	return idStr, nil
}
