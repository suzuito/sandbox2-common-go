#!/bin/sh

export FILE_PATH_BIN=$(realpath $1)
export GOCOVERDIR=$(realpath $2)

mkdir -p ${GOCOVERDIR} && rm ${GOCOVERDIR}/*
go test -count=1 -v ./e2e/release/increment-release-version/...
EXIT_CODE=$?
sh report-e2e.sh ${GOCOVERDIR}

exit $EXIT_CODE
