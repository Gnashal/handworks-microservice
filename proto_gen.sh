#!/bin/bash
set -e

# path to the payment proto file
PROTO=./common/.protobuf/payment.proto
OUT_DIR=./common/grpc/payment

mkdir -p $OUT_DIR

# generating go code only for payment.proto
protoc \
  --proto_path=./common/.protobuf \
  --go_out=$OUT_DIR \
  --go_opt=paths=source_relative \
  --go-grpc_out=$OUT_DIR \
  --go-grpc_opt=paths=source_relative \
  $PROTO

echo "payment protobuf generation complete!"
