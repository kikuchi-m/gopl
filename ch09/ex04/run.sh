#!/bin/bash

pushd $(cd $(dirname $0) && pwd)
go build -o ex04
./ex04 | tee log
echo "- - - stopped - - -"
sort -n log | tail -n 100
