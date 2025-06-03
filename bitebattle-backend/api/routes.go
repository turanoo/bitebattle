package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/bitebattle/bitebattle-backend/internal/account"
	"github.com/turanoo/bitebattle/bitebattle-backend/internal/auth"
	"github.com/turanoo/bitebattle/bitebattle-backend/internal/head2head"
	"github.com/turanoo/bitebattle/bitebattle-backend/internal/poll"
	"github.com/turanoo/bitebattle/bitebattle-backend/internal/restaurant"
	"github.com/turanoo/bitebattle/bitebattle-backend/internal/user"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	api := router.Group("/api")

	// User routes
	userService := user.NewService(db)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(api)

	// Auth routes
	authHandler := auth.NewHandler(userService)
	authHandler.RegisterRoutes(api)

	// Poll routes
	pollService := poll.NewService(db)
	pollHandler := poll.NewHandler(pollService)
	pollHandler.RegisterRoutes(api)

	// Restaurant routes
	restaurantService := restaurant.NewService()
	restaurantHandler := restaurant.NewHandler(restaurantService)
	restaurantHandler.RegisterRoutes(api)

	// Account routes
	accountService := account.NewService(db)
	accountHandler := account.NewHandler(accountService)
	accountHandler.RegisterRoutes(api)

	// inside SetupRoutes
	h2hService := head2head.NewService(db)
	h2hHandler := head2head.NewHandler(h2hService)
	h2hHandler.RegisterRoutes(api)

}
