#!/bin/bash
set -e

# path to the inventory proto file
PROTO=./common/.protobuf/inventory.proto
OUT_DIR=./common/grpc/inventory

mkdir -p $OUT_DIR

# generating go code only for inventory.proto
protoc \
  --proto_path=./common/.protobuf \
  --go_out=$OUT_DIR \
  --go_opt=paths=source_relative \
  --go-grpc_out=$OUT_DIR \
  --go-grpc_opt=paths=source_relative \
  $PROTO

echo "inventory protobuf generation complete!"
