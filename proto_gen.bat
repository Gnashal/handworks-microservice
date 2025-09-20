@echo off



@REM start protoc --proto_path=.protobuf ".protobuf/account.proto" ^
@REM     --go_out=common/grpc/account ^
@REM     --go_opt=paths=source_relative ^
@REM     --go-grpc_out=common/grpc/account ^
@REM     --go-grpc_opt=paths=source_relative ^

@REM This works in mine since windows
protoc --proto_path=common/.protobuf common/.protobuf/booking.proto ^
    --go_out=common/grpc/booking ^
    --go_opt=paths=source_relative ^
    --go-grpc_out=common/grpc/booking ^
    --go-grpc_opt=paths=source_relative

pause