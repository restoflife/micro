#!/bin/bash
find ./ -name "*.proto$" | xargs -I {} protoc --go_out=plugns=grpc:. {}