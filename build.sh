#!/usr/bin/env bash

PLATFORM=$1

set -e

ROOT=$(cd "$(dirname "$0")";cd .;pwd)

BIN=${ROOT}/git
rm -rf ${BIN}
CGO_ENABLED=0 GOOS=${PLATFORM} GOARCH=amd64 go build -o ${BIN} ${ROOT}/main.go
chmod +x ${BIN}
