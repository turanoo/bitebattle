package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func ExtractBearerToken(header string) (string, error) {
	parts := strings.Split(header, "Bearer ")
	if len(parts) != 2 {
		return "", errors.New("invalid authorization header format")
	}
	return parts[1], nil
}
