#!/bin/bash

 check folder
if test ! -d "build"; then
  mkdir -p build
fi

# check output name
if test "$1" = ""; then
  output="goim"
else
  output="$1"
fi

# compile protobuf files
compile_protobuf_files() {
  if [ -e $1 ]; then
    echo "compling $1"
    protoc --proto_path=$(dirname $1) --go_out=$(dirname $1) $(basename $1)
  fi
  if [ -e $1 -a -e $2 ]; then
    echo "compling $2"
    protoc --proto_path=$(dirname $1) --proto_path=$(dirname $2) --go_out=plugins=grpc:$(dirname $2) $(basename $2)
  fi
}

complile_project() {
  if [ "$1" = "linux" ]; then
    echo "compiling linux binary ..."
    env GOOS=linux GOARCH=amd64 go build -o build/$2 main.go
  else
    echo "compiling current platform binary ..."
    go build -o build/$2 main.go
  fi
}

compile_protobuf_files ./codec/protobuf/message.proto ./api/grpc/im.proto
complile_project $2

# compile


# docker image