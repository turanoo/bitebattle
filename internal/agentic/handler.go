package agentic

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turanoo/bitebattle/internal/auth"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) AgenticCommandHandler(c *gin.Context) {
	var req AgenticCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AgenticCommandResponse{Success: false, Message: "invalid request: " + err.Error()})
		return
	}

	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, AgenticCommandResponse{Success: false, Message: "unauthorized: " + err.Error()})
		return
	}

	result, err := h.Service.OrchestrateCommand(context.Background(), userID, req.Command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AgenticCommandResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, AgenticCommandResponse{Success: true, Message: "ok", Data: result})
}
