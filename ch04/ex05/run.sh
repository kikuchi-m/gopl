#!/bin/bash
set -e

if [ $# -eq 0 ]; then
  echo "usage: $0 ARG [ARG...]"
  exit 1
fi

go run "$(cd $(dirname $0) && pwd)/removeadjacent.go" $*
