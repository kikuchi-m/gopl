#!/bin/bash
CWD=$(cd $(dirname $0) && pwd)

$CWD/clockwall \
  NewYork=localhost:8001 \
  Tokyo=localhost:8002 \
  London=localhost:8003 \
