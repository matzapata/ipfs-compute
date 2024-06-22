BINARY_NAME=khachapuri-provider
GOOS=linux
GOARCH=amd64

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# build the binary
build:
	go build -o ./build/$(BINARY_NAME) cmd/api/main.go

# clean build files
clean:
	go clean
	rm -fr build

# install dependencies
deps:
	go get -v

# build and run the binary
run: 
	go run cmd/api/main.go

# help target to display available targets
help:
	@echo "Available targets:"
	@echo "  make          - Build the binary (default target)"
	@echo "  make build    - Build the binary"
	@echo "  make clean    - Remove the binary and other generated files"
	@echo "  make deps     - Install dependencies"
	@echo "  make run      - Run cmd/main.go"
	@echo "  make help     - Display this help message"

.PHONY: build clean deps run help

