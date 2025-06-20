# Set environment variables
COMPOSE = docker-compose


# Default target
.PHONY: help up migrate run dev stop destroy build fresh docker-build docker-push lint test

help:
	@echo "Available commands:"
	@echo "  make up              Start Postgres container"
	@echo "  make migrate         Run DB migrations"
	@echo "  make run             Run backend server"
	@echo "  make dev             Full local dev (up + migrate + run)"
	@echo "  make stop            Stop containers"
	@echo "  make destroy         Stop and destroy containers (data will be deleted from db)"
	@echo "  make build           Build Go binary"
	@echo "  make fresh           Stop, destroy, start, migrate, and run"
	@echo "  make lint            Run golangci-lint to check code quality"
	@echo "  make test            Run tests"
	@echo "  make docker-build    Build Docker image for the server"
	@echo "  make docker-push     Push Docker image to Google Container Registry"

up:
	$(COMPOSE) up -d db
	@echo "Waiting for Postgres to be ready..."
	@sleep 5

migrate:
	bash scripts/run_migrations.sh

run:
	go run cmd/server/main.go

dev: up migrate run

stop:
	$(COMPOSE) down

destroy:
	$(COMPOSE) down -v

build:
	go build -o bin/server cmd/server/main.go

lint:
	golangci-lint run ./...

test:
	cd tests && go test ./... && cd ..

fresh: destroy up migrate run

docker-build:
	docker build -t gcr.io/bitebattle/server .

docker-push:
	docker push gcr.io/bitebattle/server
