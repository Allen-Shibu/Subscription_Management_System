#!/usr/bin/env bash

set -euo pipefail

# Resolve project root (parent directory of scripts/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Check if golang-migrate CLI is installed
if ! command -v migrate >/dev/null 2>&1; then
    echo "Error: golang-migrate CLI is not installed."
    exit 1
fi

# Check if .env exists in project root
if [ ! -f "$PROJECT_ROOT/.env" ]; then
    echo "Error: .env file not found at $PROJECT_ROOT/.env"
    exit 1
fi

# Load environment variables
set -a
source "$PROJECT_ROOT/.env"
set +a

# Ensure DATABASE_URL is available
if [ -z "${DATABASE_URL:-}" ]; then
    echo "Error: DATABASE_URL is not set in .env."
    exit 1
fi

# Migrations directory
MIGRATIONS_DIR="$PROJECT_ROOT/migrations"

command="${1:-}"
arg="${2:-}"

case "$command" in
    up)
        migrate \
            -path "$MIGRATIONS_DIR" \
            -database "$DATABASE_URL" \
            up
        ;;

    down)
        count="${arg:-1}"

        read -rp "Roll back $count migration(s)? [y/N] " confirm

        if [[ "$confirm" =~ ^[Yy]$ ]]; then
            migrate \
                -path "$MIGRATIONS_DIR" \
                -database "$DATABASE_URL" \
                down "$count"
        else
            echo "Rollback cancelled."
        fi
        ;;

    create)
        if [ -z "$arg" ]; then
            echo "Error: Migration name is required."
            echo "Usage: ./scripts/migrate.sh create <migration_name>"
            exit 1
        fi

        migrate create \
            -ext sql \
            -dir "$MIGRATIONS_DIR" \
            -seq \
            "$arg"
        ;;

    force)
        if [ -z "$arg" ]; then
            echo "Error: Migration version is required."
            echo "Usage: ./scripts/migrate.sh force <version>"
            exit 1
        fi

        migrate \
            -path "$MIGRATIONS_DIR" \
            -database "$DATABASE_URL" \
            force "$arg"
        ;;

    *)
        cat <<EOF
Usage:

  ./scripts/migrate.sh up

      Apply all pending migrations.

  ./scripts/migrate.sh down [count]

      Roll back the last migration (default: 1).

  ./scripts/migrate.sh create <migration_name>

      Create a new sequential migration.

  ./scripts/migrate.sh force <version>

      Force the migration version (use only to recover from a dirty state).
EOF
        exit 1
        ;;
esac