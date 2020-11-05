#!/bin/sh
# The entrypoint for container
set -e
if [ -z "$DATA_FILE_PATH" ]; then
    mkdir -p /app/data
    export DATA_FILE_PATH="/app/data/data.json"
fi
exec /app/ghstsbot
