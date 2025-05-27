package restaurant

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) SearchRestaurants(query, location, radius string) ([]Place, error) {
	return fetchFromGooglePlaces(query, location, radius)
}
