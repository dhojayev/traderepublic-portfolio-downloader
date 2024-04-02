#! /bin/sh

apk add -u --no-cache tzdata make

go mod vendor

make wire
make lint
make test

make build-darwin-amd64
make build-darwin-arm64
make build-windows-386
make build-windows-amd64
make build-linux-386
make build-linux-amd64
make build-linux-arm64
