package vector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Config holds Qdrant connection settings.
type Config struct {
	// Port is the HTTP port Qdrant listens on.
	Port string
}

// LoadConfig reads settings from environment variables with fallbacks.
func LoadConfig() Config {
	port := os.Getenv("QDRANT_PORT")
	if port == "" {
		port = "6333"
	}
	return Config{Port: port}
}

// Client provides helpers for interacting with Qdrant.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Connect initializes a client using the given config.
func Connect(_ context.Context, cfg Config) (*Client, error) {
	return &Client{
		baseURL:    fmt.Sprintf("http://localhost:%s", cfg.Port),
		httpClient: &http.Client{},
	}, nil
}

// Point represents a single vector with optional payload.
type Point struct {
	ID      string                 `json:"id"`
	Vector  []float32              `json:"vector"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// Upsert writes points to the given collection.
func (c *Client) Upsert(ctx context.Context, collection string, pts []Point) error {
	body, err := json.Marshal(struct {
		Points []Point `json:"points"`
	}{pts})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut,
		fmt.Sprintf("%s/collections/%s/points?wait=true", c.baseURL, collection), bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("qdrant status %d", resp.StatusCode)
	}
	return nil
}

// QueryResult is a single vector search match.
type QueryResult struct {
	ID      string                 `json:"id"`
	Score   float32                `json:"score"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// Query searches for similar vectors in a collection.
func (c *Client) Query(ctx context.Context, collection string, vector []float32, limit int) ([]QueryResult, error) {
	body, err := json.Marshal(struct {
		Vector []float32 `json:"vector"`
		Limit  int       `json:"limit"`
	}{vector, limit})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/collections/%s/points/search", c.baseURL, collection), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("qdrant status %d", resp.StatusCode)
	}
	var out struct {
		Result []QueryResult `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Result, nil
}
