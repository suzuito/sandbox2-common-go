#!/bin/sh

if [ z$1 = "z" ]; then
    echo "wait-until-http-health <URL>"
    echo "URL is empty"
    exit 1
fi

while [ true ]; do
    curl $1
    if [ $? -eq 0 ]; then
        echo 'helth check is OK'
        break
    fi
    echo 'waiting until helth check is OK'
    sleep 1
done
