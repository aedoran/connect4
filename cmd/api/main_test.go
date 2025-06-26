package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"log/slog"
	"os"
)

func TestHealthz(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := setupApp(logger)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to get healthz: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestGraphQLUpsertAndSearch(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := setupApp(logger)

	// upsert a memory
	body := `{"query":"mutation { upsertMemory(userID:1, content:\"hi\", vector:[1,2]) { id }}"}`
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}

	// query search
	body = `{"query":"query { search(vector:[1,2], limit:1) { id score }}"}`
	req = httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func TestRESTCreateAndSearch(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := setupApp(logger)

	body := `{"userID":1,"content":"hi","vector":[1,2]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/memories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}

	body = `{"vector":[1,2],"limit":1}`
	req = httptest.NewRequest(http.MethodPost, "/api/v1/memories/search", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func TestRESTGetMemory(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := setupApp(logger)

	body := `{"userID":1,"content":"hello","vector":[1,2]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/memories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	var create struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&create); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/memories/"+strconv.FormatInt(create.ID, 10), nil)
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	var m struct {
		ID      int64  `json:"id"`
		UserID  int64  `json:"userID"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		t.Fatalf("decode memory: %v", err)
	}
	if m.ID != create.ID || m.Content != "hello" {
		t.Fatalf("unexpected memory %+v", m)
	}
}

func TestGraphQLBadRequest(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := setupApp(logger)

	body := `{"query":"{ unknown }"}`
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestOpenAPIDocs(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := setupApp(logger)

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("docs: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
}
