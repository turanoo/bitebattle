package head2head

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
	h2h := rg.Group("/h2h")
	h2h.Use(auth.AuthMiddleware())

	h2h.POST("/match", h.CreateMatchHandler)
	h2h.POST("/match/:id/accept", h.AcceptMatchHandler)
	h2h.POST("/match/:id/swipe", h.SubmitSwipeHandler)
	h2h.GET("/match/:id/results", h.GetMatchResultsHandler)
}

func (h *Handler) CreateMatchHandler(c *gin.Context) {
	var req struct {
		InviteeID  string   `json:"invitee_id"`
		Categories []string `json:"categories"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	inviterID := c.MustGet("userID").(uuid.UUID)
	inviteeID, err := uuid.Parse(req.InviteeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invitee ID"})
		return
	}

	match, err := h.Service.CreateMatch(inviterID, inviteeID, req.Categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create match"})
		return
	}

	c.JSON(http.StatusCreated, match)
}

func (h *Handler) AcceptMatchHandler(c *gin.Context) {
	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match ID"})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	if err := h.Service.AcceptMatch(matchID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not accept match"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "match accepted"})
}

func (h *Handler) SubmitSwipeHandler(c *gin.Context) {
	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match ID"})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	var req struct {
		RestaurantID   string `json:"restaurant_id"`
		RestaurantName string `json:"restaurant_name"`
		Liked          bool   `json:"liked"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	swipe, err := h.Service.SubmitSwipe(matchID, userID, req.RestaurantID, req.RestaurantName, req.Liked)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record swipe"})
		return
	}

	c.JSON(http.StatusOK, swipe)
}

func (h *Handler) GetMatchResultsHandler(c *gin.Context) {
	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match ID"})
		return
	}

	matches, err := h.Service.GetMutualLikes(matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch match results"})
		return
	}

	c.JSON(http.StatusOK, matches)
}
