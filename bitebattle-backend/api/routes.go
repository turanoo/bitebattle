package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/bitebattle/bitebattle-backend/internal/auth"
	"github.com/turanoo/bitebattle/bitebattle-backend/internal/group"
	"github.com/turanoo/bitebattle/bitebattle-backend/internal/user"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	api := router.Group("/api")

	// User routes
	userService := user.NewService(db)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(api)

	// Group routes
	groupService := group.NewService(db)
	groupHandler := group.NewHandler(groupService)
	groupHandler.RegisterRoutes(api)

	// Auth routes
	authHandler := auth.NewHandler(userService)
	authHandler.RegisterRoutes(api)

	// Future: Add other domain routes like auth, group, poll, etc.
}
