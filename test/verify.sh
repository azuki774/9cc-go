#!/bin/bash
BINARY_PATH=`dirname ${0}`/../build/9cc-go
INPUT=${1}
ANSWER=${2}

${BINARY_PATH} -o ./test/tmp.s ${INPUT} 

cc -c -o ./test/tmp.o ./test/tmp.s
rm -rf ./test/tmp.s

cc -c -o ./test/src/show.o ./test/src/show.c
cc -o ./test/tmp ./test/src/show.o ./test/tmp.o
rm -rf ./test/src/show.o

ACTUAL=`./test/tmp`
rm -rf ./test/tmp.o
rm -rf ./test/tmp

if [ ${ACTUAL} = ${ANSWER} ]; then
echo "[Pass] ${INPUT}"
else
echo "[Failed] ${INPUT} : output = ${ACTUAL} , answer = ${ANSWER}"
exit 1
fi
