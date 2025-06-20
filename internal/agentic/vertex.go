package agentic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type VertexAIClient struct {
	ProjectID string
	Location  string
	Model     string
	AuthToken string
}

func NewVertexAIClient() *VertexAIClient {
	projectID := os.Getenv("VERTEX_PROJECT_ID")
	location := os.Getenv("VERTEX_LOCATION")
	model := os.Getenv("VERTEX_MODEL")
	token := os.Getenv("VERTEX_AUTH_TOKEN")
	return &VertexAIClient{
		ProjectID: projectID,
		Location:  location,
		Model:     model,
		AuthToken: token,
	}
}

func (v *VertexAIClient) SendCommand(ctx context.Context, command string) (*ParsedPrompt, error) {
	url := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1beta1/projects/%s/locations/%s/publishers/google/models/%s:generateContent", v.Location, v.ProjectID, v.Location, v.Model)
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

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+v.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
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

// stripCodeBlock removes code block markers and trims whitespace.
func stripCodeBlock(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		// Remove triple backticks and optional language
		lines := strings.SplitN(s, "\n", 2)
		if len(lines) == 2 {
			s = lines[1]
		}
		s = strings.TrimPrefix(s, "```")
		if idx := strings.Index(s, "```"); idx >= 0 {
			s = s[:idx]
		}
	}
	s = strings.Trim(s, "`") // Remove any stray single backticks
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
