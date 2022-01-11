#!/bin/bash
# shellcheck disable=SC2038
find ./ -name "*.proto$" | xargs -I {} protoc --go_out=plugns=grpc:. {}