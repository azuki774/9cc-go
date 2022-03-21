#!/bin/bash
NOWDIR=`dirname ${0}`
BINARY_PATH=`dirname ${0}`/../build/9cc-go

echo "Use ${BINARY_PATH}"

# ${NOWDIR}/verify.sh <input c file> <answer> || exit <error code>
# ${NOWDIR}/verify.sh ${NOWDIR}/src/0-1.c 42 || exit 1
# ${NOWDIR}/verify.sh ${NOWDIR}/src/0-2.c 10 || exit 1
# ${NOWDIR}/verify.sh ${NOWDIR}/src/0-3.c -10 || exit 1

echo "All test passed"
