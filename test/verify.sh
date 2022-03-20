#!/bin/bash
BINARY_PATH=`dirname ${0}`/../build/9cc-go
${BINARY_PATH} -o tmp.s ${INPUT}
INPUT=${1}
ANSWER=${2}

cc -o ./test/tmp ./test/tmp.s
rm -rf ./test/tmp.s

./test/tmp
ACTUAL=${?}
rm -rf ./test/tmp

if [ ${ACTUAL} = ${ANSWER} ]; then
echo "[Pass] ${INPUT}"
else
echo "[Failed] ${INPUT} : output = ${ACTUAL} , answer = ${ANSWER}"
exit 1
fi
