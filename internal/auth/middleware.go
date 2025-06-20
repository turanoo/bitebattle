package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/turanoo/bitebattle/pkg/utils"
)

const userIDContextKey = "userID"

func AuthMiddleware() gin.HandlerFunc {
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

		userUUID, err := uuid.Parse(claims.UserID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid user id in token")
			c.Abort()
			return
		}
		c.Set(userIDContextKey, userUUID)

		c.Next()
	}
}

func UserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID := c.MustGet(userIDContextKey).(uuid.UUID)
	return userID, nil
}
