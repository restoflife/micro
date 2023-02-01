@echo off
    :: github版本的protoc 生成命令
    ::protoc --proto_path=%%~dpa/../ --proto_path=%%~dpa --go_out=plugins=grpc:%%~dpa %%~nxa

    :: google接手后的protoc 生成命令 新版需要继承UnimplementedUserSvcServer 方法
    :: protoc --proto_path=%%~dpa/../ --proto_path=%%~dpa --go_out=:%%~dpa %%~nxa  %%~nxa --go-grpc_out=:%%~dpa %%~nxa
for /r "." %%a in (*.proto) do (
     protoc --proto_path=%%~dpa/../ --proto_path=%%~dpa --go_out=:%%~dpa %%~nxa  %%~nxa --go-grpc_out=:%%~dpa %%~nxa
)