# Set environment variables
COMPOSE = docker-compose


# Default target
.PHONY: help up migrate run dev stop destroy build fresh docker-build docker-push lint test

help:
	@echo "Available commands:"
	@echo "  make up              Start Postgres container"
	@echo "  make run             Run backend server"
	@echo "  make dev             Full local dev (up + run)"
	@echo "  make stop            Stop containers"
	@echo "  make destroy         Stop and destroy containers (data will be deleted from db)"
	@echo "  make build           Build Go binary"
	@echo "  make fresh           Stop, destroy, up, and run"
	@echo "  make lint            Run golangci-lint to check code quality"
	@echo "  make test            Run tests"
	@echo "  make docker-build    Build Docker image for the server"
	@echo "  make docker-push     Push Docker image to Google Container Registry"

up:
	bash scripts/env.sh
	$(COMPOSE) up -d db
	@echo "Waiting for Postgres to be ready..."
	@sleep 5

run:
	go run main.go

dev: up run

stop:
	$(COMPOSE) down

destroy:
	$(COMPOSE) down -v && rm -f .env

build:
	go build -o bin/server main.go

lint:
	golangci-lint run ./...

test:
	cd tests && go test ./... && cd ..

fresh: destroy up run

docker-build:
	docker build -t gcr.io/bitebattle/server .

docker-push:
	docker push gcr.io/bitebattle/server
