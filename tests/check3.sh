#!/bin/bash

set -e

echo "Starting test iteration 3"

metricstest -test.v -test.run=^TestIteration3[AB]*$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server
