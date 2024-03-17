#!/bin/bash

set -e

metricstest -test.v -test.run=^TestIteration1$  \
                -binary-path=cmd/server/server
