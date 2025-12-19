#!/bin/bash
# Run database migrations

set -e

DB_URL="${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/linkbio?sslmode=disable}"

echo "Running migrations..."
psql "$DB_URL" -f apps/api/migrations/0001_init.sql

echo "Migrations completed!"
