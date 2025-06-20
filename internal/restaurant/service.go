package restaurant

import "os"

type Service struct {
	Endpoint string
	APIKey   string
}

func NewService() *Service {
	return &Service{
		Endpoint: os.Getenv("GOOGLE_PLACES_API_ENDPOINT"),
		APIKey:   os.Getenv("GOOGLE_PLACES_API_KEY"),
	}
}

func (s *Service) SearchRestaurants(query, location, radius string) ([]Place, error) {
	return s.fetchFromGooglePlaces(query, location, radius)
}
