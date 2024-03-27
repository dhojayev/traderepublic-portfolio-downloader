.PHONY: all install wire lint

all: install wire lint

install:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.1

wire:
	go run -mod=mod github.com/google/wire/cmd/wire ./...

lint:
	golangci-lint run ./...