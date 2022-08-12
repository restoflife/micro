#!/bin/bash
# shellcheck disable=SC2038
find ./ -name "*.proto" | xargs -I {} protoc -I . --go_out=plugins=grpc:.  {}
