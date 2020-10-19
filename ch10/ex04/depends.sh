#!/bin/sh
go list -f '{{.ImportPath}} {{join .Deps " "}} ' '...' | awk '/ '"`echo "$@" | sed -e 's/\//\\\\\//g' -e 's/ / \/ || \/ /g'`"' /{ print $1 }'
