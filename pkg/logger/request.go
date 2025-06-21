package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const LoggerKey = "request_logger"

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		entry := Log.WithFields(logrus.Fields{
			"requestId": requestID,
		})
		c.Set(LoggerKey, entry)
		c.Next()
	}
}

func FromContext(c *gin.Context) *logrus.Entry {
	if l, exists := c.Get(LoggerKey); exists {
		if entry, ok := l.(*logrus.Entry); ok {
			return entry
		}
	}
	return Log.WithFields(logrus.Fields{})
}
