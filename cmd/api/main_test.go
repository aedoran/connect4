package main

import (
    "net/http"
    "net/http/httptest"
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
