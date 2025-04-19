.PHONY: all lint test

all: lint test

lint:
	golangci-lint run ./... 

test:
	go test -v ./...