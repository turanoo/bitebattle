package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/biteboard/biteboard-backend/internal/user"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	api := router.Group("/api")

	// User routes
	userService := user.NewService(db)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(api)

	// Future: Add other domain routes like auth, group, poll, etc.
}
