#!/bin/sh
set -e

# Load environment variables from .env file if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

if [ -n "$INSTANCE_CONNECTION_NAME" ]; then
  echo "[MIGRATE DEBUG] Using Cloud SQL connection for migrations"
  DB_URL="postgres://$DB_USER:$DB_PASS@/$DB_NAME?host=/cloudsql/$INSTANCE_CONNECTION_NAME&sslmode=disable"
  MIGRATIONS_PATH="/app/migrations"
else
  echo "[MIGRATE DEBUG] Using TCP connection for migrations"
  DB_HOST=${DB_HOST:-localhost}
  DB_PORT=${DB_PORT:-5432}
  DB_URL="postgres://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
  MIGRATIONS_PATH="./migrations"
fi

echo "[MIGRATE DEBUG] DB_URL=$DB_URL"
migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" up
