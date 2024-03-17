#!/bin/bash

set -e

echo "Starting test iteration 1"

metricstest -test.v -test.run=^TestIteration1$  \
                -binary-path=cmd/server/server
