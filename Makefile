# Set environment variables
ENV_FILE = .env
COMPOSE = docker-compose
DB_CONTAINER = bitebattle-db


# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make up              Start Postgres container"
	@echo "  make migrate         Run DB migrations"
	@echo "  make run             Run backend server"
	@echo "  make dev             Full local dev (up + migrate + run)"
	@echo "  make stop            Stop containers"
	@echo "  make destroy         Stop and destroy containers (data will be deleted from db)"
	@echo "  make build           Build Go binary"

.PHONY: up
up:
	$(COMPOSE) up -d db
	@echo "Waiting for Postgres to be ready..."
	@sleep 5

.PHONY: migrate
migrate:
	bash scripts/run_migrations.sh

.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: dev
dev: up migrate run

.PHONY: stop
stop:
	$(COMPOSE) down

.PHONY: destroy
destroy:
	$(COMPOSE) down -v

.PHONY: build
build:
	go build -o bin/server cmd/server/main.go


.PHONY: debug
debug: destroy up migrate

.PHONY: docker-build
docker-build:
	docker build -t gcr.io/bitebattle/server .

.PHONY: docker-push
docker-push:
	docker push gcr.io/bitebattle/server
