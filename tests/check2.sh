#!/bin/bash

set -e

metricstest -test.v -test.run=^TestIteration2[AB]*$ \
        -source-path=. \
        -agent-binary-path=cmd/agent/agent
