package logger

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.DebugLevel)
}

func sanitize(args ...interface{}) []interface{} {
	safe := make([]interface{}, len(args))
	for i, arg := range args {
		if str, ok := arg.(string); ok && containsSensitive(str) {
			safe[i] = maskSensitive(str)
		} else {
			safe[i] = arg
		}
	}
	return safe
}

func sanitizef(format string, args ...interface{}) (string, []interface{}) {
	if containsSensitive(format) {
		format = maskSensitive(format)
	}
	return format, sanitize(args...)
}

func containsSensitive(s string) bool {
	low := strings.ToLower(s)
	return strings.Contains(low, "password") || strings.Contains(low, "passwd")
}

func maskSensitive(s string) string {
	return strings.ReplaceAll(s, "password", "******")
}

func Info(args ...interface{}) {
	Log.Info(sanitize(args...)...)
}

func Infof(format string, args ...interface{}) {
	f, a := sanitizef(format, args...)
	Log.Infof(f, a...)
}

func Warn(args ...interface{}) {
	Log.Warn(sanitize(args...)...)
}

func Warnf(format string, args ...interface{}) {
	f, a := sanitizef(format, args...)
	Log.Warnf(f, a...)
}

func Error(args ...interface{}) {
	Log.Error(sanitize(args...)...)
}

func Errorf(format string, args ...interface{}) {
	f, a := sanitizef(format, args...)
	Log.Errorf(f, a...)
}

func Debug(args ...interface{}) {
	Log.Debug(sanitize(args...)...)
}

func Debugf(format string, args ...interface{}) {
	f, a := sanitizef(format, args...)
	Log.Debugf(f, a...)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	safeFields := logrus.Fields{}
	for k, v := range fields {
		if containsSensitive(k) {
			safeFields[k] = "******"
		} else if str, ok := v.(string); ok && containsSensitive(str) {
			safeFields[k] = maskSensitive(str)
		} else {
			safeFields[k] = v
		}
	}
	return Log.WithFields(safeFields)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)
		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		requestID := ""
		if l, exists := c.Get(LoggerKey); exists {
			if entry, ok := l.(*logrus.Entry); ok {
				if rid, ok := entry.Data["requestId"].(string); ok {
					requestID = rid
				}
			}
		}
		if requestID == "" {
			requestID = c.GetHeader("X-Request-Id")
		}
		if requestID == "" {
			requestID = "unknown"
		}

		logFields := logrus.Fields{
			"method":    method,
			"path":      path,
			"status":    status,
			"duration":  duration.String(),
			"clientIP":  clientIP,
			"requestId": requestID,
		}
		if errMsg != "" {
			logFields["error"] = errMsg
			Log.WithFields(logFields).Error("Request completed with error")
		} else {
			Log.WithFields(logFields).Info("Request completed")
		}
	}
}

func ErrorRecovery() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(os.Stderr, func(c *gin.Context, recovered interface{}) {
		requestID := ""
		if l, exists := c.Get(LoggerKey); exists {
			if entry, ok := l.(*logrus.Entry); ok {
				if rid, ok := entry.Data["requestId"].(string); ok {
					requestID = rid
				}
			}
		}
		if requestID == "" {
			requestID = c.GetHeader("X-Request-Id")
		}
		if requestID == "" {
			requestID = "unknown"
		}

		logFields := logrus.Fields{
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"clientIP":  c.ClientIP(),
			"error":     recovered,
			"requestId": requestID,
		}
		Log.WithFields(logFields).Error("PANIC recovered")
		c.AbortWithStatusJSON(500, gin.H{
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"clientIP":  c.ClientIP(),
			"error":     "internal server error",
			"details":   recovered,
			"requestId": requestID,
		})
	})
}
