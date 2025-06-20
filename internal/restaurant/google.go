package restaurant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (s *Service) fetchFromGooglePlaces(query string, location string, radius string) ([]Place, error) {

	params := url.Values{}
	params.Add("query", query)
	params.Add("location", location) // format: "lat,lng"
	params.Add("type", "restaurant")
	params.Add("radius", radius) // in meters
	params.Add("key", s.APIKey)

	fullURL := fmt.Sprintf("%s?%s", s.Endpoint, params.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call Google Places API: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	var result GooglePlacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("google places API error: %s", result.Status)
	}

	return result.Results, nil
}
