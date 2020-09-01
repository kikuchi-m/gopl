#!/bin/bash
CWD=$(cd $(dirname $0) && pwd)

TZ=US/Eastern     $PWD/clock2 -port 8001 &
TZ=Asia/Tokyo     $PWD/clock2 -port 8002 &
TZ=Europe/London  $PWD/clock2 -port 8003 &
ps -ef | grep clock2
