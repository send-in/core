#!/usr/bin/env bash

set -e

echo "📖 Generating swagger..."
swag init \
  -g cmd/server/main.go \
  --parseDependency \
  --parseInternal

echo "✅ Swagger updated"