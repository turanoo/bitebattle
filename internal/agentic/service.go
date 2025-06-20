package agentic

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/turanoo/bitebattle/internal/poll"
	"github.com/turanoo/bitebattle/internal/restaurant"
)

type Service struct {
	Vertex *VertexAIClient
	Poll   *poll.Service
	Rest   *restaurant.Service
}

func NewService(vertex *VertexAIClient, pollSvc poll.Service, restSvc restaurant.Service) *Service {
	return &Service{
		Vertex: vertex,
		Poll:   &pollSvc,
		Rest:   &restSvc,
	}
}

func (s *Service) OrchestrateCommand(ctx context.Context, userID uuid.UUID, command string) (interface{}, error) {
	// Step 1: Parse command for query/location
	parsedPrompt, err := s.Vertex.SendCommand(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("vertex AI error: %w", err)
	}

	food := parsedPrompt.Food
	location := parsedPrompt.Location
	if food == "" {
		return nil, fmt.Errorf("could not extract restaurant query from command")
	}
	if location == "" {
		return nil, fmt.Errorf("could not extract location from command")
	}

	places, err := s.Rest.SearchRestaurants(food, location, "5000")
	if err != nil {
		return nil, fmt.Errorf("failed to search restaurants: %w", err)
	}
	maxOptions := 7
	if len(places) < maxOptions {
		maxOptions = len(places)
	}

	title := food + " poll"
	poll, err := s.Poll.CreatePoll(title, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create poll: %w", err)
	}

	var addedOptions []string
	for i := 0; i < maxOptions; i++ {
		place := places[i]
		imageURL := ""
		if len(place.Photos) > 0 {
			imageURL = place.Photos[0].PhotoReference
		}
		_, err := s.Poll.AddOption(poll.ID, place.PlaceID, place.Name, imageURL, "")
		if err == nil {
			addedOptions = append(addedOptions, place.Name)
		}
	}

	return map[string]interface{}{"poll_id": poll.ID, "title": poll.Name, "options": addedOptions}, nil
}
