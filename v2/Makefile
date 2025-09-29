.PHONY: all generate lint test

all: lint test

generate:
	set -e
	~/go/bin/oapi-codegen -config internal/traderepublic/api/restclient/openapi-config.yaml openapi-rest.yaml > internal/traderepublic/api/restclient/openapi_gen.go
	go fmt internal/traderepublic/api/restclient/openapi_gen.go
	go generate ./...

lint:
	go tool golangci-lint run ./... 

test:
	go test -v ./...