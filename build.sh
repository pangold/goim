#!/bin/bash

# check folder
if test ! -d "build"; then
  mkdir -p build
fi

# check output name
if test "$1" = ""; then
  output="goim"
else
  output="$1"
fi

# compile
if test "$2" = "linux"; then
  echo "compiling linux binary ..."
  env GOOS=linux GOARCH=amd64 go build -o build/"${output}" main.go
else
  echo "compiling current platform binary ..."
  go build -o build/"${output}" main.go
fi

# docker image