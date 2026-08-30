#!/bin/bash

set -e

TEMP_FILE=/tmp/tmpfile

echo > $TEMP_FILE

find_unused_port() {
    local port
    while true; do
        port=$((RANDOM % 49152 + 1024))  # Генерация случайного порта в диапазоне от 1024 до 49151
        if ! netstat -nl | grep -qE "[\s:]:${port}\s"; then
            echo "$port"
            return 0
        fi
    done
}

echo "Starting test iteration 9"

SERVER_PORT=$(find_unused_port)
          ADDRESS="localhost:${SERVER_PORT}"
          metricstest -test.v -test.run=^TestIteration9$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -file-storage-path="$TEMP_FILE" \
            -server-port="$SERVER_PORT" \
            -source-path=.

