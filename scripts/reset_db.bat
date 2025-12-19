@echo off
echo === Reset LinkBio Database ===

echo.
echo Step 1: Drop and recreate database...
dropdb --if-exists linkbio
createdb linkbio

echo.
echo Step 2: Run migrations...
psql -d linkbio -f apps/api/migrations/0001_init.sql

echo.
echo Step 3: Seed data...
psql -d linkbio -f scripts/seed.sql

echo.
echo === Done! Database ready ===
pause
