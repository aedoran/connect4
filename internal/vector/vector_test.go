package vector

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpsertAndQuery(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/collections/test/points", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/collections/test/points/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		resp := struct {
			Result []QueryResult `json:"result"`
		}{Result: []QueryResult{{ID: "1", Score: 0.9}}}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("encode: %v", err)
		}
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := &Client{baseURL: srv.URL, httpClient: srv.Client()}
	if err := c.Upsert(context.Background(), "test", []Point{{ID: "1", Vector: []float32{1, 2}}}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	res, err := c.Query(context.Background(), "test", []float32{1, 2}, 1)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(res) != 1 || res[0].ID != "1" {
		t.Fatalf("unexpected result: %#v", res)
	}
}
