package account

import (
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

func (h *Handler) GetProfile(c *gin.Context) {
	log := logger.FromContext(c)
	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		log.WithError(err).Warn("Invalid user id in GetProfile token")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	profile, err := h.Service.GetUserProfile(userID)
	if err != nil {
		log.WithError(err).Errorf("Failed to get profile for user %s", userID)
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to get profile")
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	log := logger.FromContext(c)
	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		log.WithError(err).Warn("Invalid user id in UpdateProfile token")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	ctx := c.Request.Context()
	err = h.Service.UpdateProfile(ctx, userID, req.Name, req.Email)
	if err != nil {
		log.WithError(err).Errorf("Failed to update profile for user %s", userID)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile.")
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetProfilePicUploadURL(c *gin.Context) {
	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	uploadURL, objectURL, err := h.Service.GenerateProfilePicUploadURL(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to generate upload url")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_url": uploadURL,
		"object_url": objectURL,
	})
}

func (h *Handler) GetProfilePicAccessURL(c *gin.Context) {
	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	accessURL, err := h.Service.GenerateProfilePicAccessURL(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "profile picture not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_url": accessURL})
}
