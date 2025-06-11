#!/bin/sh
set -e

if [ -n "$INSTANCE_CONNECTION_NAME" ]; then
  # Cloud SQL Unix socket for migrate CLI
  DB_URL="postgres://$DB_USER:$DB_PASS@/$DB_NAME?host=/cloudsql/$INSTANCE_CONNECTION_NAME&sslmode=disable"
else
  # TCP (local/dev)
  DB_HOST=${DB_HOST:-localhost}
  DB_PORT=${DB_PORT:-5432}
  DB_URL="postgres://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
fi

echo "[MIGRATE DEBUG] DB_URL=$DB_URL"
migrate -path=/app/migrations -database "$DB_URL" up