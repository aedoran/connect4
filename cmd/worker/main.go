package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	workers "github.com/jrallison/go-workers"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func embeddingJob(msg *workers.Msg) {
	logger.Info("embedding job", "args", msg.Args())
}

func linkJob(msg *workers.Msg) {
	logger.Info("link job", "args", msg.Args())
}

func main() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	workers.Configure(map[string]string{
		"server":  addr,
		"process": "1",
		"pool":    "1",
	})

	workers.Process("embeddings", embeddingJob, 1)
	workers.Process("links", linkJob, 1)

	go workers.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	workers.Quit()
}
