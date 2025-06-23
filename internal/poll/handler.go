package poll

import (
	"database/sql"
	"errors"
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

func (h *Handler) CreatePoll(c *gin.Context) {
	log := logger.FromContext(c)
	var req CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		log.WithError(err).Warn("Invalid user id in CreatePoll token")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	poll, err := h.Service.CreatePoll(req.Name, userID)
	if err != nil {
		log.WithError(err).Errorf("Failed to create poll for user %s", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create poll"})
		return
	}

	log.Infof("Poll created: %s by user %s", poll.ID, userID)
	c.JSON(http.StatusCreated, poll)
}

func (h *Handler) GetPolls(c *gin.Context) {
	log := logger.FromContext(c)
	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		log.WithError(err).Warn("Invalid user id in GetPolls")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	polls, err := h.Service.GetPolls(userID)
	if err != nil {
		log.WithError(err).Errorf("Failed to fetch polls for user %s", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch polls"})
		return
	}

	c.JSON(http.StatusOK, polls)
}

func (h *Handler) JoinPoll(c *gin.Context) {
	log := logger.FromContext(c)
	var req JoinPollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		log.WithError(err).Warn("Invalid user id in JoinPoll token")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	poll, err := h.Service.JoinPoll(req.InviteCode, userID)
	if err != nil {
		if errors.Is(err, ErrInvalidInviteCode) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, ErrAlreadyMember) {
			utils.ErrorResponse(c, http.StatusConflict, "User is already a member or owner of this poll.")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to join poll.")
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (h *Handler) GetPoll(c *gin.Context) {
	log := logger.FromContext(c)
	pollID := c.Param("pollId")

	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		log.WithError(err).Warn("Invalid user id in GetPoll token")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	poll, err := h.Service.GetPoll(pollID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(c, http.StatusNotFound, "Poll not found.")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve poll.")
		}
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (h *Handler) DeletePoll(c *gin.Context) {
	log := logger.FromContext(c)
	pollID := c.Param("pollId")

	if err := h.Service.DeletePoll(pollID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.WithError(err).Warnf("Poll not found: %s", pollID)
			utils.ErrorResponse(c, http.StatusNotFound, "Poll not found.")
		} else {
			log.WithError(err).Errorf("Failed to delete poll: %s", pollID)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete poll.")
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdatePoll(c *gin.Context) {
	log := logger.FromContext(c)
	var req UpdatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("Invalid update poll request")
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	pollID := c.Param("pollId")

	poll, err := h.Service.UpdatePoll(pollID, req.Name)
	if err != nil {
		log.WithError(err).Errorf("Failed to update poll %s", pollID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update poll"})
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (h *Handler) AddOption(c *gin.Context) {
	log := logger.FromContext(c)
	var req AddOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	pollID := c.Param("pollId")

	var addedOptions []PollOption
	for _, opt := range req {
		option, err := h.Service.AddOption(pollID, opt.RestaurantID, opt.Name, opt.ImageURL, opt.MenuURL)
		if err != nil {
			log.WithError(err).Errorf("Failed to add option %s to poll %s", opt.Name, pollID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add option"})
			return
		}
		addedOptions = append(addedOptions, *option)
	}

	c.JSON(http.StatusCreated, addedOptions)
}

func (h *Handler) CastVote(c *gin.Context) {
	log := logger.FromContext(c)
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	pollID := c.Param("pollId")

	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		log.WithError(err).Warn("Invalid user id in CastVote token")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	vote, err := h.Service.CastVote(pollID, req.OptionID, userID)
	if err != nil {
		if errors.Is(err, ErrOptionNotInPoll) {
			utils.ErrorResponse(c, http.StatusBadRequest, "Option does not exist for this poll.")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cast vote"})
		return
	}

	c.JSON(http.StatusOK, vote)
}

func (h *Handler) UncastVote(c *gin.Context) {
	log := logger.FromContext(c)
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	pollID := c.Param("pollId")

	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		log.WithError(err).Warn("Invalid user id in UncastVote token")
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.Service.RemoveVote(pollID, req.OptionID, userID); err != nil {
		if errors.Is(err, ErrOptionNotInPoll) {
			utils.ErrorResponse(c, http.StatusBadRequest, "Option does not exist for this poll.")
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(c, http.StatusNotFound, "Vote not found.")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove vote.")
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetResults(c *gin.Context) {
	log := logger.FromContext(c)
	pollID := c.Param("pollId")

	results, err := h.Service.GetResults(pollID)
	if err != nil {
		log.WithError(err).Errorf("Failed to get results for poll %s", pollID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get results"})
		return
	}

	c.JSON(http.StatusOK, results)
}
