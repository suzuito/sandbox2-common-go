#!/bin/sh

if [ z"$1" = "z" ]; then
    echo "arg1 file path of binary is required"
fi

if [ z"$2" = "z" ]; then
    echo "arg2 dir path of gocov is required"
fi

if [ z"$3" = "z" ]; then
    echo "arg3 test target is required"
fi

make $1
if [ $? -ne 0 ]; then
    echo "failed to build: make $1"
    exit 1
fi
export FILE_PATH_BIN=$(realpath $1)

mkdir -p $2
rm -f $2/*
export GOCOVERDIR=$(realpath $2)

go test -count=1 -v $3
EXIT_CODE=$?
sh report-gocovdir.sh ${GOCOVERDIR}

exit $EXIT_CODE
