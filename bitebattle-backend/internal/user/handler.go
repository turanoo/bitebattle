package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/bitebattle/bitebattle-backend/pkg/logger"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.GET(":id", h.GetUser)
	users.GET("", h.GetUserByQuery) // Use query param for email
}

func (h *Handler) CreateUser(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		logger.Warnf("Invalid input in CreateUser: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if _, err := h.Service.CreateUser(c.Request.Context(), &u); err != nil {
		logger.Errorf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	logger.Infof("User created: %s", u.Email)
	c.JSON(http.StatusCreated, u)
}

func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.Service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		logger.Warnf("User not found: %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	logger.Infof("User fetched: %s", id)
	c.JSON(http.StatusOK, user)
}

// New handler for query param based email lookup
func (h *Handler) GetUserByQuery(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		logger.Warn("email query param required in GetUserByQuery")
		c.JSON(http.StatusBadRequest, gin.H{"error": "email query param required"})
		return
	}
	user, err := h.Service.GetUserByEmail(c.Request.Context(), email)
	if err != nil || user == nil {
		logger.Warnf("User not found by email: %s", email)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	logger.Infof("User fetched by email: %s", email)
	c.JSON(http.StatusOK, user)
}
