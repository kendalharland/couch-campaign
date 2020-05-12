#!/bin/bash

SRC_DIR="$(pwd)/proto"
DART_DST_DIR="$(pwd)/app/lib/src/api/"
GO_DST_DIR="$(pwd)"

protoc -I="$SRC_DIR" --dart_out="$DART_DST_DIR" "$SRC_DIR/api.proto"
protoc -I="$SRC_DIR" --go_out="$GO_DST_DIR" "$SRC_DIR/api.proto"
