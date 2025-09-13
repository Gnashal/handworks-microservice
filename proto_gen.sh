#!/bin/bash
set -e

# path to the booking proto file
PROTO=./common/.protobuf/booking.proto
OUT_DIR=./common/grpc/booking

mkdir -p $OUT_DIR

# generating go code only for booking.proto
protoc \
  --proto_path=./common/.protobuf \
  --go_out=$OUT_DIR \
  --go_opt=paths=source_relative \
  --go-grpc_out=$OUT_DIR \
  --go-grpc_opt=paths=source_relative \
  $PROTO

echo "Booking protobuf generation complete!"
