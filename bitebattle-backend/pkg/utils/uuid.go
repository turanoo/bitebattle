package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDStr := c.MustGet("userID").(string)
	return uuid.Parse(userIDStr)
}
