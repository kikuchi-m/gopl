#!/bin/bash
set -e

pushd "$(cd $(dirname $0) && pwd)" >/dev/null

go build -o listarchives
./listarchives x.tar x.zip
