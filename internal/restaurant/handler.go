package restaurant

import (
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

func (h *Handler) SearchRestaurants(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		logger.Warn("query parameter 'q' is required in SearchRestaurants")
		utils.ErrorResponse(c, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	location := c.DefaultQuery("location", "37.7749,-122.4194") // Default: SF

	places, err := h.Service.SearchRestaurants(query, location, "10000")
	if err != nil {
		logger.Errorf("Failed to fetch restaurants: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch restaurants")
		return
	}

	logger.Infof("Restaurants search: query=%s, location=%s, found=%d", query, location, len(places))
	c.JSON(http.StatusOK, places)
}
