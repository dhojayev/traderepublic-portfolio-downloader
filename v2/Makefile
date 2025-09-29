.PHONY: all generate lint test

all: lint test

generate:
	set -e
	~/go/bin/oapi-codegen -config pkg/traderepublic/openapi-config.yml pkg/traderepublic/openapi.yml > pkg/traderepublic/openapi_gen.go
	go fmt pkg/traderepublic/openapi_gen.go
	go generate ./...

lint:
	go tool golangci-lint run ./... 

test:
	go test -v ./...