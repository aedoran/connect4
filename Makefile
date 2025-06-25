.PHONY: run dev test lint docker-build ci

run:
	go run ./cmd/api

dev: run

test:
	go test ./...

lint:
	golangci-lint run

docker-build:
	docker build -t connect4:latest .

ci: lint test

