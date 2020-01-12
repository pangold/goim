#!/bin/bash


output=goim
build_path=build
api_path=./protocol/api.proto
message_path=./protocol/message.proto

# check build folder
check_folder() {
  if [ ! -d $build_path ]; then
    mkdir -p $build_path
  fi
}

# compile protobuf files
compile_protobuf_files() {
  if [ -e $1 ]; then
    echo "compling $1"
    protoc -I=$(dirname $1) --go_out=$(dirname $1) $(basename $1)
  fi
  if [ -e $1 -a -e $2 ]; then
    echo "compling $2"
    protoc -I $(dirname $1) -I $(dirname $2) --go_out=plugins=grpc:$(dirname $2) $(basename $2)
  fi
}

# compile this go project
complile_project() {
  if [ "$2" = "linux" ]; then
    echo "compiling linux binary ..."
    env GOOS=linux GOARCH=amd64 go build -o $build_path/$1 main.go
  else
    echo "compiling current platform binary ..."
    go build -o $build_path/$1 main.go
  fi
}

#
check_folder
compile_protobuf_files $message_path $api_path
complile_project $output $1

# docker image