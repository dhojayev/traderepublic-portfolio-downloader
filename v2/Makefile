.PHONY: all generate lint test

all: lint test

generate:
	set -e
	go tool oapi-codegen -config pkg/traderepublic/openapi-config.yml pkg/traderepublic/openapi.yml > pkg/traderepublic/openapi_gen.go
	go tool go-jsonschema -p traderepublic pkg/traderepublic/timeline_transactions.schema.json > pkg/traderepublic/timeline_transactions_gen.go
	go tool go-jsonschema -p traderepublic pkg/traderepublic/websocket_response.schema.json > pkg/traderepublic/websocket_response_gen.go
	go tool go-jsonschema -p traderepublic pkg/traderepublic/websocket_sub_request.schema.json > pkg/traderepublic/websocket_sub_response_gen.go
	go fmt ./pkg/traderepublic/...
	go generate ./...

lint:
	go tool golangci-lint run ./... 

test:
	go test -v ./...