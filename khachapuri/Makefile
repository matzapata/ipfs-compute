PROVIDER_BINARY_NAME=provider
CLI_BINARY_NAME=khachapuri
GATEWAY_BINARY_NAME=gateway
GOOS=linux
GOARCH=amd64

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run-provider: 
	go run cmd/provider/main.go

run-cli:
	go run cmd/cli/main.go

run-gateway:
	go run cmd/gateway/main.go

build-provider:
	go build -o ./build/$(PROVIDER_BINARY_NAME) cmd/provider/main.go

build-cli:
	go build -o ./build/$(CLI_BINARY_NAME) cmd/cli/main.go

build-gateway:
	go build -o ./build/$(GATEWAY_BINARY_NAME) cmd/gateway/main.go

install-cli: build-cli
	sudo mv ./build/$(CLI_BINARY_NAME) /usr/local/bin/$(CLI_BINARY_NAME)

clean:
	go clean
	rm -fr build

deps:
	go get -v

.PHONY: build clean deps run help

