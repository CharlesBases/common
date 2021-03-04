#!/usr/bin/env zsh

set -e

function main() {
  protoc --proto_path=$GOPATH/src:. --stack_out=. --go_out=. websocket.proto
}

main