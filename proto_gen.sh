#!/bin/bash
set -e

# path to proto files
PROTO_DIR=./common/.protobuf
OUT_DIR=./common/grpc/account

mkdir -p $OUT_DIR

# generating go code
protoc \
  --proto_path=$PROTO_DIR \
  --go_out=$OUT_DIR \
  --go_opt=paths=source_relative \
  --go-grpc_out=$OUT_DIR \
  --go-grpc_opt=paths=source_relative \
  $PROTO_DIR/*.proto

echo "Protobuf generation complete!"
