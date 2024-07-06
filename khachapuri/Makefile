PROVIDER_BINARY_NAME=provider
CLI_BINARY_NAME=khachapuri
GATEWAY_BINARY_NAME=gateway
# GOOS=linux
# GOARCH=amd64

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# provider ========================================================================

provider-run: 
	go run cmd/provider/main.go

provider-build:
	go build -o ./build/$(PROVIDER_BINARY_NAME) cmd/provider/main.go

# gateway ========================================================================

gateway-run:
	go run cmd/gateway/main.go

gateway-build:
	go build -o ./build/$(GATEWAY_BINARY_NAME) cmd/gateway/main.go

# cli ========================================================================

cli-run:
	go run cmd/cli/main.go

cli-build:
	go build -o ./build/$(CLI_BINARY_NAME) cmd/cli/main.go

cli-install: cli-build
	sudo mv ./build/$(CLI_BINARY_NAME) /usr/local/bin/$(CLI_BINARY_NAME)
	sudo chmod +x /usr/local/bin/$(CLI_BINARY_NAME)

# all ========================================================================

clean:
	go clean
	rm -fr build

deps:
	go get -v

.PHONY: build clean deps run help

