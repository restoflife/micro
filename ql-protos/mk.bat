@echo off

for /r "." %%a in (*.proto) do (
    protoc --proto_path=%%~dpa/../ --proto_path=%%~dpa --go_out=plugins=grpc:%%~dpa %%~nxa
)