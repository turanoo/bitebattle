package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/bitebattle/internal/account"
	"github.com/turanoo/bitebattle/internal/auth"
	"github.com/turanoo/bitebattle/internal/head2head"
	"github.com/turanoo/bitebattle/internal/poll"
	"github.com/turanoo/bitebattle/internal/restaurant"
	"github.com/turanoo/bitebattle/internal/user"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	api := router.Group("/v1")

	userService := user.NewService(db)
	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(userService)

	// Public routes
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	// Protected routes
	protected := api.Group("")
	protected.Use(auth.AuthMiddleware())

	protected.GET("/users/:id", userHandler.GetUser)
	protected.GET("/users", userHandler.GetUserByQuery)
	protected.POST("/users", userHandler.CreateUser)

	accountService := account.NewService(db)
	accountHandler := account.NewHandler(accountService)
	protected.GET("/account", accountHandler.GetProfile)
	protected.PUT("/account", accountHandler.UpdateProfile)

	pollService := poll.NewService(db)
	pollHandler := poll.NewHandler(pollService)
	protected.POST("/polls", pollHandler.CreatePoll)
	protected.GET("/polls", pollHandler.GetPolls)
	protected.POST("/polls/:pollId/join", pollHandler.JoinPoll)
	protected.GET("/polls/:pollId", pollHandler.GetPoll)
	protected.DELETE("/polls/:pollId", pollHandler.DeletePoll)
	protected.PUT("/polls/:pollId", pollHandler.UpdatePoll)
	protected.POST("/polls/:pollId/options", pollHandler.AddOption)
	protected.POST("/polls/:pollId/vote", pollHandler.CastVote)
	protected.POST("/polls/:pollId/unvote", pollHandler.UncastVote)
	protected.GET("/polls/:pollId/results", pollHandler.GetResults)

	restaurantService := restaurant.NewService()
	restaurantHandler := restaurant.NewHandler(restaurantService)
	protected.GET("/restaurants/search", restaurantHandler.SearchRestaurants)

	h2hService := head2head.NewService(db)
	h2hHandler := head2head.NewHandler(h2hService)
	protected.POST("/h2h/match", h2hHandler.CreateMatch)
	protected.POST("/h2h/match/:id/accept", h2hHandler.AcceptMatch)
	protected.POST("/h2h/match/:id/swipe", h2hHandler.SubmitSwipe)
	protected.GET("/h2h/match/:id/results", h2hHandler.GetMatchResults)
}
