package account

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
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
	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		logger.Warnf("Invalid user id in GetProfile token: %v", err)
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
	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		logger.Warnf("Invalid user id in UpdateProfile token: %v", err)
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

func (h *Handler) GetProfilePicUploadURL(c *gin.Context) {
	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	bucket := os.Getenv("GCS_PROFILE_BUCKET")
	object := "profile_pics/" + userID.String() + "_" + time.Now().Format("20060102150405") + ".jpg"
	contentType := "image/jpeg"

	url, err := generateSignedUploadURL(c.Request.Context(), bucket, object, contentType)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to generate upload url")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_url": url,
		"object_url": "https://storage.googleapis.com/" + bucket + "/" + object,
	})
}

func (h *Handler) GetProfilePicAccessURL(c *gin.Context) {
	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	profile, err := h.Service.GetUserProfile(userID)
	if err != nil || profile.ProfilePicURL == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "profile picture not found")
		return
	}

	bucket, object, err := parseGCSURL(*profile.ProfilePicURL)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "invalid profile picture url")
		return
	}

	url, err := generateSignedAccessURL(c.Request.Context(), bucket, object)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to generate access url")
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_url": url})
}

func generateSignedUploadURL(ctx context.Context, bucket, object, contentType string) (string, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		cerr := client.Close()
		if cerr != nil {
			fmt.Printf("Failed to close storage client: %v\n", cerr)
		}
	}()

	url, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		Method:      "PUT",
		Expires:     time.Now().Add(15 * time.Minute),
		ContentType: contentType,
	})
	return url, err
}

func generateSignedAccessURL(ctx context.Context, bucket, object string) (string, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		cerr := client.Close()
		if cerr != nil {
			fmt.Printf("Failed to close storage client: %v\n", cerr)
		}
	}()

	url, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	})
	return url, err
}

func parseGCSURL(gcsURL string) (bucket, object string, err error) {
	// Expects format: https://storage.googleapis.com/bucket/object
	const prefix = "https://storage.googleapis.com/"
	if !strings.HasPrefix(gcsURL, prefix) {
		return "", "", fmt.Errorf("invalid GCS URL")
	}
	trimmed := strings.TrimPrefix(gcsURL, prefix)
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid GCS URL")
	}
	return parts[0], parts[1], nil
}
