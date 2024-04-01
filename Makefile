.PHONY: all wire lint test build-darwin-amd64 build-darwin-arm64 build-windows-386 build-windows-amd64 build-linux-386 build-linux-amd64 build-linux-arm64

all: wire lint test

wire:
	go run -mod=mod github.com/google/wire/cmd/wire ./...

lint:
	go run -mod=mod github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2 run ./... 

test:
	go test -v ./...

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -v -o /tmp/portfoliodownloader ./cmd/portfoliodownloader/public

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -v -o /tmp/portfoliodownloader ./cmd/portfoliodownloader/public

build-windows-386:
	GOOS=windows GOARCH=386 go build -v -o /tmp/portfoliodownloader ./cmd/portfoliodownloader/public

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -v -o /tmp/portfoliodownloader ./cmd/portfoliodownloader/public

build-linux-386:
	GOOS=linux GOARCH=386 go build -v -o /tmp/portfoliodownloader ./cmd/portfoliodownloader/public

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -v -o /tmp/portfoliodownloader ./cmd/portfoliodownloader/public

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -v -o /tmp/portfoliodownloader ./cmd/portfoliodownloader/public