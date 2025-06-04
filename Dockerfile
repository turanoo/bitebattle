# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder


WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server ./server
COPY migrations ./migrations
COPY config ./config

EXPOSE 8080

CMD ["./server"]
