#!/bin/bash

set -e

echo "Starting test iteration 2"

metricstest -test.v -test.run=^TestIteration2[AB]*$ \
        -source-path=. \
        -agent-binary-path=cmd/agent/agent
