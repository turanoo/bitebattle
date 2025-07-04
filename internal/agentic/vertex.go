package agentic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/turanoo/bitebattle/pkg/config"
)

var (
	cachedToken     string
	tokenExpiryTime int64
)

type VertexAIClient struct {
	Url string
}

func NewVertexAIClient(cfg *config.Config) *VertexAIClient {
	url := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1beta1/projects/%s/locations/%s/publishers/google/models/%s:generateContent",
		cfg.Vertex.Location, cfg.Vertex.ProjectID, cfg.Vertex.Location, cfg.Vertex.Model)
	return &VertexAIClient{
		Url: url,
	}
}

func (v *VertexAIClient) SendCommand(ctx context.Context, command string) (*ParsedPrompt, error) {
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": buildPrompt(command)},
				},
			},
		},
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequestWithContext(ctx, "POST", v.Url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	token, err := getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := resp.Body.Close()
		_ = cerr
	}()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vertex AI error: %s", string(respBody))
	}

	var vertexResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&vertexResp); err != nil {
		return nil, err
	}
	if len(vertexResp.Candidates) == 0 || len(vertexResp.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("no candidates from Vertex AI")
	}

	jsonText := stripCodeBlock(vertexResp.Candidates[0].Content.Parts[0].Text)

	var intent ParsedPrompt
	if err := json.Unmarshal([]byte(jsonText), &intent); err != nil {
		return nil, fmt.Errorf("failed to parse intent: %w", err)
	}
	return &intent, nil
}

func stripCodeBlock(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		lines := strings.SplitN(s, "\n", 2)
		if len(lines) == 2 {
			s = lines[1]
		}
		s = strings.TrimPrefix(s, "```")
		if idx := strings.Index(s, "```"); idx >= 0 {
			s = s[:idx]
		}
	}
	s = strings.Trim(s, "`")
	return strings.TrimSpace(s)
}

func buildPrompt(command string) string {
	return `You are a strict JSON-generating agent for a food poll app.

Your task is to extract and return only a valid JSON object from a user's command using this exact structure:

{
  "food": <string, required — type of food or restaurant>,
  "location": <string, required — must be in "lat,lng" format>,
  "radius": <string, required — search radius in meters>
}

Rules:
- **"food"**: Extract what kind of food or restaurant is mentioned (e.g., pizza, sushi, Indian food).
- **"location"**: If the command includes a specific coordinate (e.g., 37.7749,-122.4194), extract it as-is.  
  If it includes a city, town, neighborhood, or zip code (e.g., "in Austin", "near 94107"), convert that to a best estimate of coordinates in "lat,lng" format.
  If no location is mentioned, default to "40.7128,-74.0060" (New York City).
- **"radius"**: Extract if mentioned like "within 3000 meters" or "in a 5km radius". Otherwise, default to "10000".

You must return only the JSON output. No explanations, no comments.

Examples:
Command: create a poll for tacos in San Diego within 5000 meters  
Output: {"food": "tacos", "location": "32.7157,-117.1611", "radius": "5000"}

Command: create a poll for ramen near Austin  
Output: {"food": "ramen", "location": "30.2672,-97.7431", "radius": "10000"}

Command: create a poll for Chinese food around 94107  
Output: {"food": "Chinese food", "location": "37.7691,-122.3933", "radius": "10000"}

Command: create a poll for sushi restaurants at 37.7749,-122.4194 within 5000 meters  
Output: {"food": "sushi restaurants", "location": "37.7749,-122.4194", "radius": "5000"}

Command: create a poll for vegan food  
Output: {"food": "vegan food", "location": "40.7128,-74.0060", "radius": "10000"}

Command: ` + command + `
Output:`
}

func getAccessToken() (string, error) {
	// If we have a cached token and it's not expired, return it
	if cachedToken != "" && tokenExpiryTime > time.Now().Unix()+60 {
		return cachedToken, nil
	}

	req, err := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Metadata-Flavor", "Google")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		cerr := resp.Body.Close()
		_ = cerr
	}()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	cachedToken = tokenResp.AccessToken
	tokenExpiryTime = time.Now().Unix() + int64(tokenResp.ExpiresIn)

	return cachedToken, nil
}
