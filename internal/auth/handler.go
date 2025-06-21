package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/bitebattle/internal/user"
	"github.com/turanoo/bitebattle/pkg/logger"
	"github.com/turanoo/bitebattle/pkg/utils"
)

type Handler struct {
	userService *user.Service
}

func NewHandler(userService *user.Service) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) Register(c *gin.Context) {
	log := logger.FromContext(c)
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("Invalid register request")
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	if !utils.IsPasswordStrong(req.Password) {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrWeakPassword.Error())
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.WithError(err).Error("Failed to hash password")
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user_ := &user.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
	}

	ctx := c.Request.Context()
	createdUser, err := h.userService.CreateUser(ctx, user_)
	if err != nil {
		if errors.Is(err, user.ErrUserExists) {
			utils.ErrorResponse(c, http.StatusConflict, "User with this email already exists.")
		} else {
			log.WithError(err).Warn("Failed to create user")
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user.")
		}
		return
	}

	token, err := GenerateToken(createdUser.ID)
	if err != nil {
		log.WithError(err).Errorf("Failed to generate token for user %s", createdUser.ID)
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	log.Infof("User registered: %s", createdUser.ID)
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *Handler) Login(c *gin.Context) {
	log := logger.FromContext(c)
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("Invalid login request")
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	ctx := c.Request.Context()
	u, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.WithError(err).Warnf("Login failed for email %s", req.Email)
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if err := utils.CheckPasswordHash(u.PasswordHash, req.Password); err != nil {
		log.Warnf("Invalid password for user %s", u.ID)
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := GenerateToken(u.ID)
	if err != nil {
		log.WithError(err).Errorf("Failed to generate token for user %s", u.ID)
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	log.Infof("User logged in: %s", u.ID)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
