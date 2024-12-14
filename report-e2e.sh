#!/bin/sh

set -e

DIR_PATH=$1

go tool covdata func -i=${DIR_PATH} > ${DIR_PATH}/covdata.func.txt
go tool covdata textfmt -i=${DIR_PATH} -o=${DIR_PATH}/textfmt.txt
go tool cover -html=${DIR_PATH}/textfmt.txt -o=${DIR_PATH}/index.html
