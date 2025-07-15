#!/bin/bash
set -e

echo "Waiting for PostgreSQL to be ready..."
until goose postgres "host=$DB_HOST port=$DB_PORT dbname=$DB_NAME user=$DB_USER password=$DB_PASSWORD sslmode=disable" status; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

echo "Running migrations..."
goose -dir migrations postgres "host=$DB_HOST port=$DB_PORT dbname=$DB_NAME user=$DB_USER password=$DB_PASSWORD sslmode=disable" up