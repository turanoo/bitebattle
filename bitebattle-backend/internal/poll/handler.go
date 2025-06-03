package poll

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/turanoo/bitebattle/bitebattle-backend/internal/auth"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	polls := rg.Group("/polls")
	polls.Use(auth.AuthMiddleware())
	polls.POST("", h.CreatePollHandler)
	polls.GET("", h.GetUserPolls)
	polls.POST("/join/:inviteCode", h.JoinPoll)
	polls.GET("/:pollId", h.GetPollHandler)
	polls.POST("/:pollId/options", h.AddOptionHandler)
	polls.POST("/:pollId/vote", h.CastVoteHandler)
	polls.GET("/:pollId/results", h.GetResultsHandler)

}

func (h *Handler) CreatePollHandler(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userIDStr, ok := auth.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	poll, err := h.Service.CreatePoll(req.Name, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create poll"})
		return
	}

	c.JSON(http.StatusCreated, poll)
}

func (h *Handler) GetPollHandler(c *gin.Context) {
	pollID, err := uuid.Parse(c.Param("pollId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll ID"})
		return
	}

	poll, err := h.Service.GetPoll(pollID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get poll"})
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (h *Handler) JoinPoll(c *gin.Context) {
	inviteCode := c.Param("inviteCode")
	if inviteCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invite code is required"})
		return
	}

	userIDStr, ok := auth.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	poll, err := h.Service.JoinPoll(inviteCode, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join poll"})
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (h *Handler) AddOptionHandler(c *gin.Context) {
	pollID, err := uuid.Parse(c.Param("pollId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll ID"})
		return
	}

	var req struct {
		RestaurantID string `json:"restaurant_id"`
		Name         string `json:"name"`
		ImageURL     string `json:"image_url"`
		MenuURL      string `json:"menu_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	option, err := h.Service.AddOption(pollID, req.RestaurantID, req.Name, req.ImageURL, req.MenuURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add option"})
		return
	}

	c.JSON(http.StatusCreated, option)
}

func (h *Handler) CastVoteHandler(c *gin.Context) {
	pollID, err := uuid.Parse(c.Param("pollId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll ID"})
		return
	}

	var req struct {
		OptionID string `json:"option_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userIDStr, ok := auth.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	optionID, err := uuid.Parse(req.OptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid option ID"})
		return
	}

	vote, err := h.Service.CastVote(pollID, optionID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cast vote"})
		return
	}

	c.JSON(http.StatusOK, vote)
}

func (h *Handler) GetResultsHandler(c *gin.Context) {
	pollID, err := uuid.Parse(c.Param("pollId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll ID"})
		return
	}

	results, err := h.Service.GetResults(pollID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get results"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) GetUserPolls(c *gin.Context) {
	userIDStr := c.MustGet("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	polls, err := h.Service.GetPolls(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch polls"})
		return
	}

	c.JSON(http.StatusOK, polls)
}
