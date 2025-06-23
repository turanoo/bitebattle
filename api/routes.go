package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/bitebattle/internal/account"
	"github.com/turanoo/bitebattle/internal/agentic"
	"github.com/turanoo/bitebattle/internal/auth"
	"github.com/turanoo/bitebattle/internal/head2head"
	"github.com/turanoo/bitebattle/internal/poll"
	"github.com/turanoo/bitebattle/internal/restaurant"
	"github.com/turanoo/bitebattle/pkg/config"
)

func SetupRoutes(router *gin.Engine, db *sql.DB, cfg *config.Config) {
	api := router.Group("/v1")

	// Public routes
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Protected routes
	protected := api.Group("")
	protected.Use(auth.Auth0Middleware())

	accountService := account.NewService(db, cfg)
	accountHandler := account.NewHandler(accountService)
	protected.GET("/account", accountHandler.GetProfile)
	protected.PUT("/account", accountHandler.UpdateProfile)
	protected.POST("/account/profile-pic/upload-url", accountHandler.GetProfilePicUploadURL)
	protected.GET("/account/profile-pic/access-url", accountHandler.GetProfilePicAccessURL)

	pollService := poll.NewService(db, cfg)
	pollHandler := poll.NewHandler(pollService)
	protected.POST("/polls", pollHandler.CreatePoll)
	protected.GET("/polls", pollHandler.GetPolls)
	protected.POST("/polls/join", pollHandler.JoinPoll)
	protected.GET("/polls/:pollId", pollHandler.GetPoll)
	protected.DELETE("/polls/:pollId", pollHandler.DeletePoll)
	protected.PUT("/polls/:pollId", pollHandler.UpdatePoll)
	protected.POST("/polls/:pollId/options", pollHandler.AddOption)
	protected.POST("/polls/:pollId/vote", pollHandler.CastVote)
	protected.POST("/polls/:pollId/unvote", pollHandler.UncastVote)
	protected.GET("/polls/:pollId/results", pollHandler.GetResults)

	restaurantService := restaurant.NewService(cfg)
	restaurantHandler := restaurant.NewHandler(restaurantService)
	protected.GET("/restaurants/search", restaurantHandler.SearchRestaurants)

	h2hService := head2head.NewService(db, cfg)
	h2hHandler := head2head.NewHandler(h2hService)
	protected.POST("/h2h/match", h2hHandler.CreateMatch)
	protected.POST("/h2h/match/:id/accept", h2hHandler.AcceptMatch)
	protected.POST("/h2h/match/:id/swipe", h2hHandler.SubmitSwipe)
	protected.GET("/h2h/match/:id/results", h2hHandler.GetMatchResults)

	agenticVertex := agentic.NewVertexAIClient(cfg)
	agenticService := agentic.NewService(agenticVertex, *pollService, *restaurantService)
	agenticHandler := agentic.NewHandler(agenticService)
	protected.POST("/agentic/command", agenticHandler.Command)
}
