package restaurant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type googlePlacesResponse struct {
	Results []Place `json:"results"`
	Status  string  `json:"status"`
}

func fetchFromGooglePlaces(query string, location string) ([]Place, error) {
	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Google Places API key not set")
	}

	endpoint := "https://maps.googleapis.com/maps/api/place/textsearch/json"
	params := url.Values{}
	params.Add("query", query)
	params.Add("location", location) // format: "lat,lng"
	params.Add("type", "restaurant")
	params.Add("radius", "10000") // in meters
	params.Add("key", apiKey)

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call Google Places API: %w", err)
	}
	defer resp.Body.Close()

	var result googlePlacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("Google Places API error: %s", result.Status)
	}

	return result.Results, nil
}
