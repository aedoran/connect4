package main

import (
	"net/http"
	"net/http/httptest"
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
