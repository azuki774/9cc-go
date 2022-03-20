#!/bin/bash
NOWDIR=`dirname ${0}`
BINARY_PATH=`dirname ${0}`/../build/9cc-go

echo "Use ${BINARY_PATH}"

# ${NOWDIR}/verify.sh <input c file> <answer> || exit <error code>

echo "All test passed"
