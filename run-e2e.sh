#!/bin/sh

export FILE_PATH_BIN=$(realpath $1)
export GOCOVERDIR=$(realpath $2)

mkdir -p ${GOCOVERDIR} && rm ${GOCOVERDIR}/*
go test -count=1 -v $3
EXIT_CODE=$?
sh report-gocovdir.sh ${GOCOVERDIR}

exit $EXIT_CODE
