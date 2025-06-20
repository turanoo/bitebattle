# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder


WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./main.go

FROM alpine:latest
WORKDIR /app

RUN apk --no-cache add ca-certificates wget

# Install migrate CLI
RUN wget -O migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz \
  && tar -xzf migrate.tar.gz -C /usr/local/bin \
  && rm migrate.tar.gz \
  && chmod +x /usr/local/bin/migrate

COPY --from=builder /app/server ./server
COPY migrations ./migrations
COPY scripts/migrations.sh ./migrations.sh

RUN chmod +x ./migrations.sh

EXPOSE 8080

CMD sh -c './migrations.sh && ./server'
