#!/bin/bash

set -e

if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

if [ -z "$DATABASE_URL" ]; then
    echo "DATABASE_URL is not set"
    exit 1
fi

ACTION=${1:-up}

case "$ACTION" in
    up)
        migrate \
            -path internal/migrations \
            -database "$DATABASE_URL" \
            up
        ;;
    down)
        migrate \
            -path internal/migrations \
            -database "$DATABASE_URL" \
            down 1
        ;;
    version)
        migrate \
            -path internal/migrations \
            -database "$DATABASE_URL" \
            version
        ;;
    force)
        if [ -z "$2" ]; then
            echo "Usage: ./migrate.sh force <version>"
            exit 1
        fi

        migrate \
            -path internal/migrations \
            -database "$DATABASE_URL" \
            force "$2"
        ;;
    *)
        echo "Usage: ./migrate.sh [up|down|version|force]"
        exit 1
        ;;
esac