package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/turanoo/bitebattle/bitebattle-backend/api"
	"github.com/turanoo/bitebattle/bitebattle-backend/pkg/db"
	"github.com/turanoo/bitebattle/bitebattle-backend/pkg/logger"
)

func main() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, using system env variables")
	}

	if err := db.Init(); err != nil {
		logger.Errorf("Failed to connect to DB: %v", err)
		os.Exit(1)
	}

	database := db.GetDB()

	router := gin.New()
	router.Use(RequestLogger())
	router.Use(ErrorRecovery())

	api.SetupRoutes(router, database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Infof("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Errorf("Server failed: %v", err)
		os.Exit(1)
	}
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
		if errMsg != "" {
			logger.Warnf("%s %s %d %s %s | ERR: %s", method, path, status, duration, clientIP, errMsg)
		} else {
			logger.Infof("%s %s %d %s %s", method, path, status, duration, clientIP)
		}
	}
}

func ErrorRecovery() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(os.Stderr, func(c *gin.Context, recovered interface{}) {
		logger.Errorf("PANIC: %v | Path: %s | Method: %s | IP: %s", recovered, c.Request.URL.Path, c.Request.Method, c.ClientIP())
		c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
	})
}
