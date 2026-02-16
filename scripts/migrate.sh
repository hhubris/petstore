#!/usr/bin/env bash
# Migration runner for the Pet Store database.
#
# Creates the petstore role and database if they don't
# exist, then runs golang-migrate against the petstore
# database.
#
# Role and database creation happen via psql (inside the
# container) before migrate runs, because CREATE DATABASE
# cannot execute inside a transaction (which
# golang-migrate uses).
#
# Required env vars (injected by mise):
#   POSTGRES_PASSWORD  — postgres superuser password
#   PETSTORE_PASSWORD  — petstore application user password
#
# Optional env vars:
#   DB_HOST — database host (default: localhost)
#   DB_PORT — database port (default: 5432)

set -euo pipefail

: "${POSTGRES_PASSWORD:?POSTGRES_PASSWORD must be set}"
: "${PETSTORE_PASSWORD:?PETSTORE_PASSWORD must be set}"

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

CONTAINER="petstore-db"

# Helper: run psql inside the container as postgres.
# The -i flag allows heredocs to be piped via stdin.
run_psql() {
    docker exec -i \
        -e PGPASSWORD="$POSTGRES_PASSWORD" \
        -e PETSTORE_PASSWORD="$PETSTORE_PASSWORD" \
        "$CONTAINER" psql -U postgres "$@"
}

# Wait for PostgreSQL to accept connections (up to 30s)
for i in $(seq 1 30); do
    if docker exec "$CONTAINER" pg_isready -U postgres \
        >/dev/null 2>&1; then
        break
    fi
    if [ "$i" -eq 30 ]; then
        echo "error: timed out waiting for PostgreSQL" >&2
        exit 1
    fi
    sleep 1
done

# Create the petstore role if it doesn't exist.
run_psql -v ON_ERROR_STOP=1 <<SQL
DO \$\$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_roles
        WHERE rolname = 'petstore'
    ) THEN
        EXECUTE format(
            'CREATE ROLE petstore LOGIN PASSWORD %L',
            '${PETSTORE_PASSWORD}'
        );
    END IF;
END
\$\$;
SQL

# Create the petstore database if it doesn't exist.
if ! run_psql -tAc \
    "SELECT 1 FROM pg_database WHERE datname = 'petstore'" \
    | grep -q 1; then
    run_psql -v ON_ERROR_STOP=1 \
        -c "CREATE DATABASE petstore OWNER petstore"
fi

# Run migrations against the petstore database.
PETSTORE_URL="postgres://postgres:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/petstore?sslmode=disable"

exec migrate -path migrations -database "$PETSTORE_URL" up
