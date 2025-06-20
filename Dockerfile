# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder


WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./main.go \
    && [ -f server ] || (echo "Build failed: server binary not found" && exit 1)

FROM alpine:latest
WORKDIR /app

RUN apk --no-cache add ca-certificates wget

ARG MIGRATE_VERSION=v4.16.2

RUN wget -O migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz \
  && tar -xzf migrate.tar.gz -C /usr/local/bin \
  && rm migrate.tar.gz \
  && chmod +x /usr/local/bin/migrate

COPY --from=builder /app/server ./server
COPY scripts/migrations.sh ./migrations.sh
COPY migrations ./migrations

RUN if [ ! -f ./migrations.sh ]; then echo "Error: scripts/migrations.sh not found." >&2; exit 1; fi \
    && chmod +x ./migrations.sh

EXPOSE 8080

CMD sh -c './migrations.sh && ./server'
