.PHONY: all init generate lint test run-dev run-prod

all: generate lint test

init:
	go run -mod=mod github.com/google/wire/cmd/wire ./...

generate:
	go generate ./...
	
lint:
	go run -mod=mod github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2 run ./...

test:
	go test -v ./...

run-dev:
	go run ./cmd/portfoliodownloader/dev -l --trace +490123456789

run-public:
	go run ./cmd/portfoliodownloader/public -w --trace +490123456789