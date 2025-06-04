package logger

import (
	"os"
	"strings"

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
