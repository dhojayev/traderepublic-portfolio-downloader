.PHONY: all init generate lint test build-darwin-amd64 build-darwin-arm64 build-windows-386 build-windows-amd64 build-linux-386 build-linux-amd64 build-linux-arm64

all: lint test

init:
	go run -mod=mod github.com/google/wire/cmd/wire ./...

generate:
	go generate ./...
	
lint:
	go run -mod=mod github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./... 

test:
	go test -v ./...

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -v -o /tmp/portfoliodownloader ./cmd/portfoliodownloader/public

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -v -o /tmp/portfoliodownloader ./cmd/portfoliodownloader/public

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -v -o /tmp/portfoliodownloader.exe ./cmd/portfoliodownloader/public

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -v -o /tmp/portfoliodownloader ./cmd/portfoliodownloader/public

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -v -o /tmp/portfoliodownloader ./cmd/portfoliodownloader/public