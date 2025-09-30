.PHONY: all generate lint test

all: lint test

generate:
	set -e
	go tool oapi-codegen -config pkg/traderepublic/openapi-config.yml pkg/traderepublic/openapi.yml > pkg/traderepublic/openapi_gen.go
	
	go tool go-jsonschema -p traderepublic \
	--schema-output=https://raw.githubusercontent.com/dhojayev/traderepublic-portfolio-downloader/refs/heads/main/v2/pkg/common.json=pkg/traderepublic/common_gen.go \
	--schema-output=https://raw.githubusercontent.com/dhojayev/traderepublic-portfolio-downloader/refs/heads/main/v2/pkg/ws_connect_request.json=pkg/traderepublic/ws_connect_request_gen.go \
	--schema-output=https://raw.githubusercontent.com/dhojayev/traderepublic-portfolio-downloader/refs/heads/main/v2/pkg/traderepublic/ws_sub_request.json=pkg/traderepublic/ws_sub_request_gen.go \
	--schema-output=https://raw.githubusercontent.com/dhojayev/traderepublic-portfolio-downloader/refs/heads/main/v2/pkg/traderepublic/ws_response.json=pkg/traderepublic/ws_response_gen.go \
	--schema-output=https://raw.githubusercontent.com/dhojayev/traderepublic-portfolio-downloader/refs/heads/main/v2/pkg/traderepublic/timeline_transactions.json=pkg/traderepublic/timeline_transactions_gen.go \
	--schema-output=https://raw.githubusercontent.com/dhojayev/traderepublic-portfolio-downloader/refs/heads/main/v2/pkg/traderepublic/timeline_detail.json=pkg/traderepublic/timeline_detail_gen.go \
	--schema-output=https://raw.githubusercontent.com/dhojayev/traderepublic-portfolio-downloader/refs/heads/main/v2/pkg/traderepublic/instrument.json=pkg/traderepublic/instrument_gen.go \
	pkg/traderepublic/common.json \
	pkg/traderepublic/ws_connect_request.json \
	pkg/traderepublic/ws_sub_request.json \
	pkg/traderepublic/ws_response.json \
	pkg/traderepublic/timeline_transactions.json \
	pkg/traderepublic/timeline_detail.json \
	pkg/traderepublic/instrument.json

	go fmt ./pkg/traderepublic/...
	go generate ./...

lint:
	go tool golangci-lint run ./... 

test:
	go test -v ./...