@echo off
setlocal enabledelayedexpansion
    :: github版本的protoc 生成命令
    ::protoc --proto_path=%%~dpa/../ --proto_path=%%~dpa --go_out=plugins=grpc:%%~dpa %%~nxa

    :: google接手后的protoc 生成命令 新版需要继承UnimplementedUserSvcServer 方法
    :: protoc --proto_path=%%~dpa/../ --proto_path=%%~dpa --go_out=:%%~dpa %%~nxa  %%~nxa --go-grpc_out=:%%~dpa %%~nxa
for /r "." %%a in (*.proto) do (
     protoc --proto_path=%%~dpa/../ --proto_path=%%~dpa --go_out=:%%~dpa %%~nxa  %%~nxa --go-grpc_out=:%%~dpa %%~nxa
)

@REM proto 引入外部proto时使用以下命令 第三方proto需要在GOPATH内 https://github.com/protocolbuffers/protobuf
set "protobuf_files="
 for %%f in (*.proto) do (
    set "protobuf_files=!protobuf_files!%%f "
)

set "proto_path=%GOPATH%\src"
set "output_path=."

if not exist "%proto_path%" (
    echo "Invalid proto path"
    exit /b 1
)

if not exist "%output_path%" (
    echo "Invalid output path"
    exit /b 1
)

protoc --proto_path="%proto_path%" --go_out="%output_path%" --go_opt=paths=source_relative --proto_path=. %protobuf_files%

endlocal