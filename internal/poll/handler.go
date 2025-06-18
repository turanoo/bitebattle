package poll

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/turanoo/bitebattle/internal/auth"
	"github.com/turanoo/bitebattle/pkg/logger"
	"github.com/turanoo/bitebattle/pkg/utils"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) CreatePoll(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid request in CreatePoll: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userIDStr, ok := auth.GetUserIDFromContext(c)
	if !ok {
		logger.Warn("Unauthorized access in CreatePoll")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Warnf("Invalid user ID in CreatePoll: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	poll, err := h.Service.CreatePoll(req.Name, userID)
	if err != nil {
		logger.Errorf("Failed to create poll for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create poll"})
		return
	}

	logger.Infof("Poll created: %s by user %s", poll.ID, userID)
	c.JSON(http.StatusCreated, poll)
}

func (h *Handler) GetPolls(c *gin.Context) {
	userID, err := utils.UserIDFromContext(c)
	if err != nil {
		logger.Warnf("Invalid user id in GetPolls: %v", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	polls, err := h.Service.GetPolls(userID)
	if err != nil {
		logger.Errorf("Failed to fetch polls for user %s: %v", userID, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch polls")
		return
	}

	c.JSON(http.StatusOK, polls)
}

func (h *Handler) GetPoll(c *gin.Context) {
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

func (h *Handler) DeletePoll(c *gin.Context) {
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

func (h *Handler) UpdatePoll(c *gin.Context) {
	pollID, err := uuid.Parse(c.Param("pollId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll ID"})
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required in body"})
		return
	}

	poll, err := h.Service.UpdatePoll(pollID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update poll"})
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (h *Handler) JoinPoll(c *gin.Context) {
	pollId := c.Param("pollId")
	if pollId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pollId is required"})
		return
	}
	var req struct {
		InviteCode string `json:"invite_code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.InviteCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invite_code is required in body"})
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
	poll, err := h.Service.JoinPoll(req.InviteCode, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join poll"})
		return
	}
	c.JSON(http.StatusOK, poll)
}

func (h *Handler) AddOption(c *gin.Context) {
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

func (h *Handler) CastVote(c *gin.Context) {
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

func (h *Handler) UncastVote(c *gin.Context) {
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

func (h *Handler) GetResults(c *gin.Context) {
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
