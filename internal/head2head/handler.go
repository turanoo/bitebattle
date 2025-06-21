package head2head

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *Handler) CreateMatch(c *gin.Context) {
	log := logger.FromContext(c)
	var req CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		log.WithError(err).Warn("Invalid user id in CreateMatch token")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	inviteeID, err := uuid.Parse(req.InviteeID)
	if err != nil {
		log.WithError(err).Warn("Invalid invitee ID in CreateMatch")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid invitee ID")
		return
	}

	match, err := h.Service.CreateMatch(userID, inviteeID, req.Categories)
	if err != nil {
		log.WithError(err).Error("Could not create match")
		utils.ErrorResponse(c, http.StatusInternalServerError, "could not create match")
		return
	}

	log.Infof("Head2Head match created: %s by %s", match.ID, userID)
	c.JSON(http.StatusCreated, match)
}

func (h *Handler) AcceptMatch(c *gin.Context) {
	log := logger.FromContext(c)
	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid match ID")
		return
	}

	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		log.WithError(err).Warn("Invalid user id in AcceptMatch")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.Service.AcceptMatch(matchID, userID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "could not accept match")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "match accepted"})
}

func (h *Handler) SubmitSwipe(c *gin.Context) {
	log := logger.FromContext(c)
	var req SubmitSwipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	matchIDStr := c.Param("id")
	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid match ID")
		return
	}

	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		log.WithError(err).Warn("Invalid user id in SubmitSwipe")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	swipe, err := h.Service.SubmitSwipe(matchID, userID, req.RestaurantID, req.RestaurantName, req.Liked)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to record swipe")
		return
	}

	c.JSON(http.StatusOK, swipe)
}

func (h *Handler) GetMatchResults(c *gin.Context) {
	log := logger.FromContext(c)
	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.WithError(err).Warn("Invalid match ID in GetMatchResults")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid match ID")
		return
	}

	matches, err := h.Service.GetMutualLikes(matchID)
	if err != nil {
		log.WithError(err).Error("Failed to fetch match results")
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch match results")
		return
	}

	c.JSON(http.StatusOK, matches)
}
