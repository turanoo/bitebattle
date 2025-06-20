#!/bin/bash
# Generate .env file for docker-compose and migrations.sh from config/local.yaml

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
YAML_PATH="$PROJECT_ROOT/config/local.yaml"
ENV_PATH="$PROJECT_ROOT/.env"

awk '
    $1 == "db:" {in_db=1; next}
    in_db && /^[^ ]/ {in_db=0}
    in_db && $1 == "host:" {print "DB_HOST="$2}
    in_db && $1 == "port:" {print "DB_PORT="$2}
    in_db && $1 == "user:" {print "DB_USER="$2}
    in_db && $1 == "pass:" {print "DB_PASS="$2}
    in_db && $1 == "name:" {print "DB_NAME="$2}
' "$YAML_PATH" > "$ENV_PATH"

echo ".env file generated for docker-compose:"
cat "$ENV_PATH"