.PHONY: all init generate lint test build-darwin-amd64 build-darwin-arm64 build-windows-386 build-windows-amd64 build-linux-386 build-linux-amd64 build-linux-arm64

all: lint test

init:
	go run -mod=mod github.com/google/wire/cmd/wire ./...

generate:
	rm -f **/*_mock.go
	go generate ./...
	go run ./cmd/example-generator
	
lint:
	golangci-lint run ./... 

test:
	go test -v ./...

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -v -o /tmp/portfoliodownloader/public/portfoliodownloader-darwin-amd64 ./cmd/portfoliodownloader/public

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -v -o /tmp/portfoliodownloader/public/portfoliodownloader-darwin-arm64 ./cmd/portfoliodownloader/public

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -v -o /tmp/portfoliodownloader/public/portfoliodownloader-windows-amd64.exe ./cmd/portfoliodownloader/public

build-windows-arm64:
	GOOS=windows GOARCH=arm64 go build -v -o /tmp/portfoliodownloader/public/portfoliodownloader-windows-arm64.exe ./cmd/portfoliodownloader/public

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -v -o /tmp/portfoliodownloader/public/portfoliodownloader-linux-amd64 ./cmd/portfoliodownloader/public

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -v -o /tmp/portfoliodownloader/public/portfoliodownloader-linux-arm64 ./cmd/portfoliodownloader/public