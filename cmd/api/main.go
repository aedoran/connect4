package main

import (
	"fmt"
	"log"
	"net/http"

	"mem0-go/internal/config"
)

func setupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, `{"status":"ok"}`); err != nil {
			log.Printf("write healthz response: %v", err)
		}
	})
	return mux
}

func main() {
	cfg := config.Load()
	addr := ":" + cfg.HTTPPort

	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, setupRoutes()); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
