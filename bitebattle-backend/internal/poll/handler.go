package poll

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	polls := rg.Group("/polls")
	polls.POST("/", h.CreatePollHandler)
	polls.POST("/:pollId/options", h.AddOptionHandler)
	polls.POST("/:pollId/vote", h.CastVoteHandler)
	polls.GET("/:pollId/results", h.GetResultsHandler)
}

func (h *Handler) CreatePollHandler(c *gin.Context) {
	var req struct {
		GroupID   string `json:"group_id"`
		CreatedBy string `json:"created_by"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	groupID, err := uuid.Parse(req.GroupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}
	userID, err := uuid.Parse(req.CreatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	poll, err := h.Service.CreatePoll(groupID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create poll"})
		return
	}

	c.JSON(http.StatusCreated, poll)
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
		UserID   string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	optionID, err := uuid.Parse(req.OptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid option ID"})
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
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
