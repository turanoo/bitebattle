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
	polls.GET("", h.GetPollsHandler)
	polls.POST("/join/:inviteCode", h.JoinPollHandler)
	polls.GET("/:pollId", h.GetPollHandler)
	polls.DELETE("/:pollId", h.DeletePollHandler)
	polls.POST("/:pollId/options", h.AddOptionHandler)
	polls.POST("/:pollId/vote", h.CastVoteHandler)
	polls.POST("/:pollId/unvote", h.UncastVoteHandler)
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

func (h *Handler) GetPollsHandler(c *gin.Context) {
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

func (h *Handler) GetPollHandler(c *gin.Context) {
	pollID, err := uuid.Parse(c.Param("pollId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll ID"})
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

	poll, err := h.Service.GetPoll(pollID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get poll"})
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (h *Handler) DeletePollHandler(c *gin.Context) {
	pollID, err := uuid.Parse(c.Param("pollId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll ID"})
		return
	}

	err = h.Service.DeletePoll(pollID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete poll"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) JoinPollHandler(c *gin.Context) {
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

	var req []struct {
		RestaurantID string `json:"restaurant_id"`
		Name         string `json:"name"`
		ImageURL     string `json:"image_url"`
		MenuURL      string `json:"menu_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if len(req) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no options provided"})
		return
	}

	var addedOptions []PollOption
	for _, opt := range req {
		option, err := h.Service.AddOption(pollID, opt.RestaurantID, opt.Name, opt.ImageURL, opt.MenuURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add option"})
			return
		}
		addedOptions = append(addedOptions, *option)
	}

	c.JSON(http.StatusCreated, addedOptions)
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

func (h *Handler) UncastVoteHandler(c *gin.Context) {
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

	err = h.Service.RemoveVote(pollID, optionID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to uncast vote"})
		return
	}

	c.Status(http.StatusNoContent)
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
