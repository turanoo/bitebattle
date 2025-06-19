package account

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/bitebattle/pkg/logger"
	"github.com/turanoo/bitebattle/pkg/utils"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
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
	userID, _ := c.Get("userID")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	ctx := c.Request.Context()
	err := h.Service.UpdateProfile(ctx, userID.(string), req.Name, req.Email)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			utils.ErrorResponse(c, http.StatusConflict, "User with this email already exists.")
		} else {
			logger.Errorf("Failed to update profile for user %s: %v", userID, err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile.")
		}
		return
	}

	c.Status(http.StatusOK)
}
