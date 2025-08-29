#!/bin/bash

export MIGRATION_DSN="host=${PG_HOST} port=5432 dbname=${PG_DATABASE_NAME} user=${PG_USER} password=${PG_PASSWORD} sslmode=disable"

echo "Applying migrations with DSN: host=${PG_HOST} dbname=${PG_DATABASE_NAME} user=${PG_USER}"

# Ждём PostgreSQL
echo "Waiting for PostgreSQL to be ready..."
until pg_isready -h "${PG_HOST}" -p 5432 -U "${PG_USER}" -d "${PG_DATABASE_NAME}"; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 2
done

echo "PostgreSQL is ready. Running migrations..."
goose -dir "./migrations" postgres "${MIGRATION_DSN}" up -v

echo "✅ Migrations applied successfully."