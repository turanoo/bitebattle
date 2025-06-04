#!/bin/bash

# Exit if anything fails
set -e

# Ensure this is set — change if needed
DB_URL=${DB_URL:-"postgres://postgres:postgres@db:5432/bitebattle?sslmode=disable"}

docker run --rm \
  -v $(pwd)/migrations:/migrations \
  migrate/migrate \
  -path=/migrations \
  -database "$DB_URL" \
  up
