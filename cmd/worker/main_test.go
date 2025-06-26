package main

import (
	"testing"

	workers "github.com/jrallison/go-workers"
)

func TestEmbeddingJob(t *testing.T) {
	msg := workers.NewMsg([]interface{}{"hi"})
	embeddingJob(msg)
}

func TestLinkJob(t *testing.T) {
	msg := workers.NewMsg([]interface{}{"a", "b"})
	linkJob(msg)
}
