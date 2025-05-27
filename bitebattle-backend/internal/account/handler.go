package account

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
	account := rg.Group("/account")
	account.Use(auth.AuthMiddleware())

	account.GET("/", h.GetProfile)
	account.PUT("/", h.UpdateProfile)
	account.GET("/groups", h.GetUserGroups)
	account.GET("/polls", h.GetUserPolls)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userIDStr := c.MustGet("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	profile, err := h.Service.GetUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userIDStr := c.MustGet("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err = h.Service.UpdateUserProfile(userID, req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

func (h *Handler) GetUserGroups(c *gin.Context) {
	userIDStr := c.MustGet("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	groups, err := h.Service.GetUserGroups(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

func (h *Handler) GetUserPolls(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	polls, err := h.Service.GetUserPolls(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch polls"})
		return
	}

	c.JSON(http.StatusOK, polls)
}
