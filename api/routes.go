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
	api.GET("/users/:id", userHandler.GetUser)
	api.GET("/users", userHandler.GetUserByQuery)
	api.POST("/users", userHandler.CreateUser)

	authHandler := auth.NewHandler(userService)
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	pollService := poll.NewService(db)
	pollHandler := poll.NewHandler(pollService)
	api.POST("/polls", pollHandler.CreatePoll)
	api.GET("/polls", pollHandler.GetPolls)
	api.POST("/polls/:pollId/join", pollHandler.JoinPoll)
	api.GET("/polls/:pollId", pollHandler.GetPoll)
	api.DELETE("/polls/:pollId", pollHandler.DeletePoll)
	api.PUT("/polls/:pollId", pollHandler.UpdatePoll)
	api.POST("/polls/:pollId/options", pollHandler.AddOption)
	api.POST("/polls/:pollId/vote", pollHandler.CastVote)
	api.POST("/polls/:pollId/unvote", pollHandler.UncastVote)
	api.GET("/polls/:pollId/results", pollHandler.GetResults)

	restaurantService := restaurant.NewService()
	restaurantHandler := restaurant.NewHandler(restaurantService)
	api.GET("/restaurants/search", restaurantHandler.SearchRestaurants)

	accountService := account.NewService(db)
	accountHandler := account.NewHandler(accountService)
	api.GET("/account", accountHandler.GetProfile)
	api.PUT("/account", accountHandler.UpdateProfile)

	h2hService := head2head.NewService(db)
	h2hHandler := head2head.NewHandler(h2hService)
	api.POST("/h2h/match", h2hHandler.CreateMatch)
	api.POST("/h2h/match/:id/accept", h2hHandler.AcceptMatch)
	api.POST("/h2h/match/:id/swipe", h2hHandler.SubmitSwipe)
	api.GET("/h2h/match/:id/results", h2hHandler.GetMatchResults)

	// Health check route
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
