package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/bitebattle/bitebattle-backend/pkg/logger"
	"github.com/turanoo/bitebattle/bitebattle-backend/pkg/utils"
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
	users.GET("", h.GetUserByQuery)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		logger.Warnf("Invalid input in CreateUser: %v", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	if _, err := h.Service.CreateUser(c.Request.Context(), &u); err != nil {
		logger.Errorf("Failed to create user: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
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
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	logger.Infof("User fetched: %s", id)
	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserByQuery(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		logger.Warn("email query param required in GetUserByQuery")
		utils.ErrorResponse(c, http.StatusBadRequest, "email query param required")
		return
	}
	user, err := h.Service.GetUserByEmail(c.Request.Context(), email)
	if err != nil || user == nil {
		logger.Warnf("User not found by email: %s", email)
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}
	logger.Infof("User fetched by email: %s", email)
	c.JSON(http.StatusOK, user)
}
