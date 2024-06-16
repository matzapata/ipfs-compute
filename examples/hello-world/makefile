
BINARY_NAME=main
GOOS=linux
GOARCH=amd64

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Build the binary
build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME) cmd/main.go

# Deploy the binary to ipfs
deploy: build
	ipfs-compute ${DEPLOYMENT_OWNER_PK}

# Clean build files
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -fr dist

# Install dependencies
deps:
	go get -v

# Build and run the binary
run: 
	go run cmd/main.go

# Help target to display available targets
help:
	@echo "Available targets:"
	@echo "  make          - Build the binary (default target)"
	@echo "  make build    - Build the binary"
	@echo "  make clean    - Remove the binary and other generated files"
	@echo "  make deps     - Install dependencies"
	@echo "  make run      - Run cmd/main.go"
	@echo "  make help     - Display this help message"

