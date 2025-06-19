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

func fetchFromGooglePlaces(query string, location string, radius string) ([]Place, error) {
	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("google places API key not set")
	}

	endpoint := "https://maps.googleapis.com/maps/api/place/textsearch/json"
	params := url.Values{}
	params.Add("query", query)
	params.Add("location", location) // format: "lat,lng"
	params.Add("type", "restaurant")
	params.Add("radius", radius) // in meters
	params.Add("key", apiKey)

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call Google Places API: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	var result googlePlacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("google places API error: %s", result.Status)
	}

	return result.Results, nil
}
