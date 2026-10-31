#!/usr/bin/env bash

set -euo pipefail

# Resolve project root (parent directory of migration/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
# Load .env if present
if [ -f "$PROJECT_ROOT/.env" ]; then
    set -a
    # shellcheck disable=SC1090
    source "$PROJECT_ROOT/.env"
    set +a
elif [ -f ".env" ]; then
    set -a
    # shellcheck disable=SC1091
    source .env
    set +a
fi

# Ensure DATABASE_URL is available
if [ -z "${DATABASE_URL:-}" ]; then
    echo "Error: DATABASE_URL is not set." >&2
    exit 1
fi

# Migrations directory
if [ -d "$PROJECT_ROOT/migration" ]; then
    MIGRATIONS_DIR="$PROJECT_ROOT/migration"
elif [ -d "$PROJECT_ROOT/migrations" ]; then
    MIGRATIONS_DIR="$PROJECT_ROOT/migrations"
else
    MIGRATIONS_DIR="migration"
fi

command="${1:-up}"
arg="${2:-}"

# If golang-migrate CLI is installed, use it for migrate commands
if command -v migrate >/dev/null 2>&1; then
    case "$command" in
        up)
            migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
            ;;
        down)
            count="${arg:-1}"
            read -rp "Roll back $count migration(s)? [y/N] " confirm
            if [[ "$confirm" =~ ^[Yy]$ ]]; then
                migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down "$count"
            else
                echo "Rollback cancelled."
            fi
            ;;
        create)
            if [ -z "$arg" ]; then
                echo "Error: Migration name is required." >&2
                echo "Usage: ./migration/migrate.sh create <migration_name>" >&2
                exit 1
            fi
            migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$arg"
            ;;
        force)
            if [ -z "$arg" ]; then
                echo "Error: Migration version is required." >&2
                echo "Usage: ./migration/migrate.sh force <version>" >&2
                exit 1
            fi
            migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force "$arg"
            ;;
        *)
            echo "Unknown command: $command" >&2
            exit 1
            ;;
    esac
else
    # Fallback to psql when golang-migrate is not installed
    if ! command -v psql >/dev/null 2>&1; then
        echo "Error: Neither golang-migrate nor psql CLI is installed." >&2
        exit 1
    fi

    psql "${DATABASE_URL}" -v ON_ERROR_STOP=1 -q \
        -c "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW())"

    applied() {
        psql "${DATABASE_URL}" -tAq -c "SELECT 1 FROM schema_migrations WHERE version = '$1'"
    }

    case "$command" in
        up)
            for file in "$MIGRATIONS_DIR"/*_up.sql; do
                [ -e "$file" ] || continue
                version="$(basename "$file" _up.sql)"
                [[ -n "$(applied "$version")" ]] && continue
                echo "applying $version"
                { cat "$file"; printf "\nINSERT INTO schema_migrations (version) VALUES ('%s');\n" "$version"; } |
                    psql "${DATABASE_URL}" -v ON_ERROR_STOP=1 -q -1
            done
            ;;
        down)
            count="${arg:-1}"
            reverted=0
            for file in $(ls -r "$MIGRATIONS_DIR"/*_down.sql 2>/dev/null); do
                [ -e "$file" ] || continue
                version="$(basename "$file" _down.sql)"
                [[ -z "$(applied "$version")" ]] && continue
                echo "reverting $version"
                { cat "$file"; printf "\nDELETE FROM schema_migrations WHERE version = '%s';\n" "$version"; } |
                    psql "${DATABASE_URL}" -v ON_ERROR_STOP=1 -q -1
                reverted=$((reverted + 1))
                if [ "$count" -ne 0 ] && [ "$reverted" -ge "$count" ]; then
                    break
                fi
            done
            ;;
        *)
            echo "usage: $0 [up|down]" >&2
            exit 1
            ;;
    esac
fi