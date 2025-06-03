package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/turanoo/bitebattle/bitebattle-backend/api"
	"github.com/turanoo/bitebattle/bitebattle-backend/pkg/db"
	"github.com/turanoo/bitebattle/bitebattle-backend/pkg/logger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables")
	}

	if err := db.Init(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	logger.Init()

	database := db.GetDB()

	router := gin.Default()
	router.Use(RequestLogger())

	api.SetupRoutes(router, database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)
		logger.Infof("%s %s %d %s", method, path, status, duration)
	}
}
