package restaurant

import (
	"github.com/turanoo/bitebattle/pkg/config"
)

type Service struct {
	Endpoint string
	APIKey   string
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		Endpoint: cfg.GooglePlaces.APIEndpoint,
		APIKey:   cfg.GooglePlaces.APIKey,
	}
}

func (s *Service) SearchRestaurants(query, location, radius string) ([]Place, error) {
	return s.fetchFromGooglePlaces(query, location, radius)
}
