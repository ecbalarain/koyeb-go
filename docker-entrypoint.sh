#!/bin/sh
# Docker entrypoint script for OXLOOK API
# Runs database migrations before starting the application

set -e

echo "🚀 Starting OXLOOK API deployment..."

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
  echo "❌ ERROR: DATABASE_URL environment variable is not set"
  exit 1
fi

echo "📊 Running database migrations..."
if ./migrate; then
  echo "✅ Migrations completed successfully"
else
  echo "❌ Migration failed"
  exit 1
fi

echo "🌐 Starting API server..."
exec ./main
