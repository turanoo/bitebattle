package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/bitebattle/bitebattle-backend/internal/user"
	"github.com/turanoo/bitebattle/bitebattle-backend/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	userService *user.Service
}

func NewHandler(userService *user.Service) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.RegisterHandler)
	rg.POST("/login", h.LoginHandler)
}

func (h *Handler) RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid register request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &user.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
	}

	ctx := c.Request.Context()
	createdUser, err := h.userService.CreateUser(ctx, user)
	if err != nil {
		logger.Warnf("Failed to create user: %v", err)
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	token, err := GenerateToken(createdUser.ID)
	if err != nil {
		logger.Errorf("Failed to generate token for user %s: %v", createdUser.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	logger.Infof("User registered: %s", createdUser.ID)
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	u, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Warnf("Login failed for email %s: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		logger.Warnf("Invalid password for user %s", u.ID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := GenerateToken(u.ID)
	if err != nil {
		logger.Errorf("Failed to generate token for user %s: %v", u.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	logger.Infof("User logged in: %s", u.ID)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
