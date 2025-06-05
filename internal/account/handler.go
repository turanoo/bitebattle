package account

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	account := rg.Group("/account")
	account.Use(auth.AuthMiddleware())

	account.GET("", h.GetProfile)
	account.PUT("", h.UpdateProfile)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, err := utils.UserIDFromContext(c)
	if err != nil {
		logger.Warnf("Invalid user id in GetProfile: %v", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	profile, err := h.Service.GetUserProfile(userID)
	if err != nil {
		logger.Errorf("Failed to get profile for user %s: %v", userID, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to get profile")
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, err := utils.UserIDFromContext(c)
	if err != nil {
		logger.Warnf("Invalid user id in UpdateProfile: %v", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req struct {
		Name            *string `json:"name"`
		Email           *string `json:"email"`
		CurrentPassword *string `json:"current_password"`
		NewPassword     *string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid request in UpdateProfile: %v", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	err = h.Service.UpdateUserProfile(userID, req.Name, req.Email, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			logger.Warnf("Incorrect current password for user %s", userID)
			utils.ErrorResponse(c, http.StatusUnauthorized, "current password is incorrect")
			return
		}
		logger.Errorf("Failed to update profile for user %s: %v", userID, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "update failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}
