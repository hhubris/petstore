#!/usr/bin/env bash
# Migration runner for the Pet Store database.
#
# Wraps golang-migrate to inject the petstore_password
# session variable required by migration 000001.
#
# Required env vars (injected by mise):
#   POSTGRES_PASSWORD  — postgres superuser password
#   PETSTORE_PASSWORD  — petstore application user password
#
# Optional env vars:
#   DATABASE_URL — full connection string override

set -euo pipefail

: "${POSTGRES_PASSWORD:?POSTGRES_PASSWORD must be set}"
: "${PETSTORE_PASSWORD:?PETSTORE_PASSWORD must be set}"

if [ -z "${DATABASE_URL:-}" ]; then
    DATABASE_URL="postgres://postgres:${POSTGRES_PASSWORD}@localhost:5432/postgres?sslmode=disable&options=-c%20migration.petstore_password%3D${PETSTORE_PASSWORD}"
fi

# Wait for PostgreSQL to accept connections (up to 30s)
for i in $(seq 1 30); do
    if docker exec petstore-db pg_isready -U postgres \
        >/dev/null 2>&1; then
        break
    fi
    if [ "$i" -eq 30 ]; then
        echo "error: timed out waiting for PostgreSQL" >&2
        exit 1
    fi
    sleep 1
done

exec migrate -path migrations -database "$DATABASE_URL" up
