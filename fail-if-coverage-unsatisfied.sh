#!/bin/sh

T=$1

COV=$(tail -1 gocovfunc.txt | awk '{print $3}' | sed 's/.$//g')

R=$(echo "${COV} < ${T}" | bc)

if [ "$R" = "1" ];then
    echo "coverage must larger than ${T}"
    exit 1
fi
