# Set environment variables
ENV_FILE = .env
COMPOSE = docker-compose
DB_CONTAINER = biteboard-db

# Migration tool binary (assuming you're using golang-migrate)
MIGRATE = migrate
MIGRATIONS_DIR = ./migrations
DB_URL = postgres://$(shell grep DB_USER $(ENV_FILE) | cut -d '=' -f2):$(shell grep DB_PASSWORD $(ENV_FILE) | cut -d '=' -f2)@localhost:5432/$(shell grep DB_NAME $(ENV_FILE) | cut -d '=' -f2)?sslmode=disable

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make up              Start Postgres container"
	@echo "  make migrate         Run DB migrations"
	@echo "  make run             Run backend server"
	@echo "  make dev             Full local dev (up + migrate + run)"
	@echo "  make stop            Stop containers"
	@echo "  make build           Build Go binary"

.PHONY: up
up:
	$(COMPOSE) up -d db
	@echo "Waiting for Postgres to be ready..."
	@sleep 5

.PHONY: migrate
migrate:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: dev
dev: up migrate run

.PHONY: stop
stop:
	$(COMPOSE) down

.PHONY: build
build:
	go build -o bin/server cmd/server/main.go
