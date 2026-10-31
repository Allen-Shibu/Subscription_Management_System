#!/usr/bin/env bash
# Applies or rolls back migration/*_up.sql / *_down.sql against DATABASE_URL.
# Applied versions are tracked in schema_migrations, so re-running is safe.
set -euo pipefail

cd "$(dirname "$0")/.."

if [[ -f .env ]]; then
	set -a
	# shellcheck disable=SC1091
	source .env
	set +a
fi

if [[ -z "${DATABASE_URL:-}" ]]; then
	echo "DATABASE_URL is not set" >&2
	exit 1
fi

direction="${1:-up}"
psql "${DATABASE_URL}" -v ON_ERROR_STOP=1 -q \
	-c "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW())"

applied() {
	psql "${DATABASE_URL}" -tAq -c "SELECT 1 FROM schema_migrations WHERE version = '$1'"
}

case "$direction" in
up)
	for file in migration/*_up.sql; do
		version="$(basename "$file" _up.sql)"
		[[ -n "$(applied "$version")" ]] && continue
		echo "applying $version"
		{ cat "$file"; printf "\nINSERT INTO schema_migrations (version) VALUES ('%s');\n" "$version"; } |
			psql "${DATABASE_URL}" -v ON_ERROR_STOP=1 -q -1
	done
	;;
down)
	for file in $(ls -r migration/*_down.sql); do
		version="$(basename "$file" _down.sql)"
		[[ -z "$(applied "$version")" ]] && continue
		echo "reverting $version"
		{ cat "$file"; printf "\nDELETE FROM schema_migrations WHERE version = '%s';\n" "$version"; } |
			psql "${DATABASE_URL}" -v ON_ERROR_STOP=1 -q -1
	done
	;;
*)
	echo "usage: $0 [up|down]" >&2
	exit 1
	;;
esac
