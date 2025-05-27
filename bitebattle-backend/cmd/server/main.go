package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/turanoo/bitebattle/bitebattle-backend/api"
	"github.com/turanoo/bitebattle/bitebattle-backend/pkg/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables")
	}

	if err := db.Init(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	database := db.GetDB()

	router := gin.Default()

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
