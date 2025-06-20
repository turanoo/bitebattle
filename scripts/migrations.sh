#!/bin/sh
set -e

# Load .env file if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Determine environment
ENV=${APP_ENV:-local}

if [ "$ENV" = "prod" ]; then
  echo "[MIGRATE DEBUG] Environment is prod. Fetching secrets from GCP Secret Manager..."

  # Fetch secrets from GCP Secret Manager (assumes gcloud CLI is authenticated)
  DB_USER=$(gcloud secrets versions access latest --secret="DB_USER")
  DB_PASS=$(gcloud secrets versions access latest --secret="DB_PASS")
  DB_NAME=$(gcloud secrets versions access latest --secret="DB_NAME")
  INSTANCE_CONNECTION_NAME=$(gcloud secrets versions access latest --secret="INSTANCE_CONNECTION_NAME")

  echo "[MIGRATE DEBUG] Using Cloud SQL connection for migrations"
  DB_URL="postgres://$DB_USER:$DB_PASS@/$DB_NAME?host=/cloudsql/$INSTANCE_CONNECTION_NAME&sslmode=disable"
  MIGRATIONS_PATH="/app/migrations"

else
  echo "[MIGRATE DEBUG] Environment is local."

  # Use environment variables set by env.sh
  DB_USER=${DB_USER}
  DB_PASS=${DB_PASS}
  DB_NAME=${DB_NAME}
  DB_HOST=${DB_HOST:-localhost}
  DB_PORT=${DB_PORT:-5432}

  echo "[MIGRATE DEBUG] DB_USER=$DB_USER"
  echo "[MIGRATE DEBUG] DB_PASS=$DB_PASS"
  echo "[MIGRATE DEBUG] DB_NAME=$DB_NAME"
  echo "[MIGRATE DEBUG] DB_HOST=$DB_HOST"
  echo "[MIGRATE DEBUG] DB_PORT=$DB_PORT"

  echo "[MIGRATE DEBUG] Using TCP connection for migrations"
  DB_URL="postgres://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
  MIGRATIONS_PATH="./migrations"
fi

echo "[MIGRATE DEBUG] DB_URL=$DB_URL"
migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" up
