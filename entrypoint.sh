#! /bin/sh

apk add -u --no-cache tzdata make

go mod vendor

make all

make build-darwin-amd64
make build-darwin-arm64
make build-windows-386
make build-windows-amd64
make build-linux-386
make build-linux-amd64
make build-linux-arm64
