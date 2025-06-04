package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/turanoo/bitebattle/bitebattle-backend/internal/auth"
	"github.com/turanoo/bitebattle/bitebattle-backend/pkg/logger"
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

	account.GET("", h.GetProfile)
	account.PUT("", h.UpdateProfile)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userIDStr := c.MustGet("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Warnf("Invalid user id in GetProfile: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	profile, err := h.Service.GetUserProfile(userID)
	if err != nil {
		logger.Errorf("Failed to get profile for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
		return
	}

	logger.Infof("Profile fetched for user %s", userID)
	c.JSON(http.StatusOK, profile)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userIDStr := c.MustGet("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Warnf("Invalid user id in UpdateProfile: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid request in UpdateProfile: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err = h.Service.UpdateUserProfile(userID, req.Name, req.Email, req.Password)
	if err != nil {
		logger.Errorf("Failed to update profile for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	logger.Infof("Profile updated for user %s", userID)
	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}
