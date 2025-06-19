package user

import (
	"database/sql"
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

func (h *Handler) CreateUser(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	ctx := c.Request.Context()
	createdUser, err := h.Service.CreateUser(ctx, &u)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			utils.ErrorResponse(c, http.StatusConflict, "User with this email already exists.")
		} else {
			logger.Errorf("Failed to create user: %v", err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user.")
		}
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	user, err := h.Service.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found.")
		} else {
			logger.Errorf("Failed to get user %s: %v", id, err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve user.")
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserByQuery(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Email query parameter is required.")
		return
	}

	ctx := c.Request.Context()
	user, err := h.Service.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found.")
		} else {
			logger.Errorf("Failed to get user by email %s: %v", email, err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve user.")
		}
		return
	}

	c.JSON(http.StatusOK, user)
}
