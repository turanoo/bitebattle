package agentic

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/bitebattle/internal/auth"
	"github.com/turanoo/bitebattle/pkg/logger"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Command(c *gin.Context) {
	log := logger.FromContext(c)
	var req AgenticRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("invalid request")
		c.JSON(http.StatusBadRequest, AgenticResponse{Success: false, Message: "invalid request: " + err.Error()})
		return
	}

	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		log.WithError(err).Warn("unauthorized")
		c.JSON(http.StatusUnauthorized, AgenticResponse{Success: false, Message: "unauthorized: " + err.Error()})
		return
	}

	result, err := h.Service.OrchestrateCommand(context.Background(), userID, req.Command)
	if err != nil {
		log.WithError(err).Error("command failed")
		c.JSON(http.StatusInternalServerError, AgenticResponse{Success: false, Message: err.Error()})
		return
	}
	log.Info("command succeeded")
	c.JSON(http.StatusOK, AgenticResponse{Success: true, Message: "ok", Data: result})
}
