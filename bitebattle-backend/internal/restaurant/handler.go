package restaurant

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
	routes := rg.Group("/restaurants")
	routes.GET("/search", h.SearchRestaurants)
}

func (h *Handler) SearchRestaurants(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		logger.Warn("query parameter 'q' is required in SearchRestaurants")
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	location := c.DefaultQuery("location", "37.7749,-122.4194") // Default: SF

	places, err := h.Service.SearchRestaurants(query, location, "10000")
	if err != nil {
		logger.Errorf("Failed to fetch restaurants: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch restaurants"})
		return
	}

	logger.Infof("Restaurants search: query=%s, location=%s, found=%d", query, location, len(places))
	c.JSON(http.StatusOK, places)
}
