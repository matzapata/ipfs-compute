API_BINARY_NAME=api
CLI_BINARY_NAME=api
GOOS=linux
GOARCH=amd64

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run-api: 
	go run cmd/api/main.go

run-cli:
	go run cmd/cli/main.go

build-api:
	go build -o ./build/$(API_BINARY_NAME) cmd/api/main.go

build-cli:
	go build -o ./build/$(CLI_BINARY_NAME) cmd/cli/main.go

clean:
	go clean
	rm -fr build

deps:
	go get -v

.PHONY: build clean deps run help

