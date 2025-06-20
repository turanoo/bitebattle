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

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server ./server
COPY migrations ./migrations

EXPOSE 8080

CMD ["./server"]
