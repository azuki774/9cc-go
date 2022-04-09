#!/bin/bash
NOWDIR=`dirname ${0}`
BINARY_PATH=`dirname ${0}`/../build/9cc-go

echo "Use ${BINARY_PATH}"

# ${NOWDIR}/verify.sh <input c file> <answer> || exit <error code>
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/0-0.c 123 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/0-1.c 10 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/0-2.c 10 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/0-3.c 0 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/0-4.c 5 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/0-5.c 26 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/0-6.c 56 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/0-7.c 90 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/0-8.c 6 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/0-9.c 2 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-0-0.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-0-1.c 0 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-0-2.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-1-0.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-1-1.c 0 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-1-2.c 0 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-2-0.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-2-1.c 0 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-2-2.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-3-0.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-3-1.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-3-2.c 0 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-4-0.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-4-1.c 0 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-4-2.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-5-0.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-5-1.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/1-5-2.c 0 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-0-1.c 123 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-1-0.c 90 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-1-1.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-1-2.c 75 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-2-0.c 100 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-2-1.c 100 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-2-2.c 100 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-3-0.c 10 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-3-1.c 14 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-3-2.c 3 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-4-0.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-4-1.c 3 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-4-2.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-4-3.c 2 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-4-4.c 3 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-5-0.c 10 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-6-0.c 32 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-6-1.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-7-0.c 1 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-7-1.c 2 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-7-2.c 97 || exit 1
${NOWDIR}/verify-nomain.sh ${NOWDIR}/src/2-7-3.c 2 || exit 1
${NOWDIR}/verify.sh ${NOWDIR}/src/3-0-0.c 15 || exit 1
${NOWDIR}/verify.sh ${NOWDIR}/src/3-0-1.c 6 || exit 1
${NOWDIR}/verify.sh ${NOWDIR}/src/3-0-2.c 55 || exit 1
${NOWDIR}/verify.sh ${NOWDIR}/src/3-1-0.c 5 || exit 1
${NOWDIR}/verify.sh ${NOWDIR}/src/3-1-1.c 6 || exit 1
${NOWDIR}/verify.sh ${NOWDIR}/src/3-1-2.c 55 || exit 1
${NOWDIR}/verify.sh ${NOWDIR}/src/3-1-3.c 100 || exit 1
${NOWDIR}/verify.sh ${NOWDIR}/src/3-1-4.c 13 || exit 1
${NOWDIR}/verify.sh ${NOWDIR}/src/3-1-5.c 55 || exit 1
echo "All test passed"
