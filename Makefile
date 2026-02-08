.PHONY: build test lint clean dev

# Build variables
BINARY_NAME=lango
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build targets
build:
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/lango

build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/lango

build-darwin:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/lango

build-all: build-linux build-darwin

# Development
dev:
	$(GOBUILD) -o bin/$(BINARY_NAME) ./cmd/lango && ./bin/$(BINARY_NAME) serve

# Testing
test:
	$(GOTEST) -v -race -cover ./...

test-short:
	$(GOTEST) -v -short ./...

bench:
	$(GOTEST) -bench=. -benchmem ./...

# Linting
lint:
	golangci-lint run ./...

# Dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Cleaning
clean:
	$(GOCLEAN)
	rm -rf bin/

# Generate
generate:
	$(GOCMD) generate ./...

# Docker
docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .

# Help
help:
	@echo "Available targets:"
	@echo "  build       - Build binary for current platform"
	@echo "  build-linux - Build for Linux amd64"
	@echo "  build-darwin- Build for macOS arm64"
	@echo "  build-all   - Build for all platforms"
	@echo "  dev         - Build and run locally"
	@echo "  test        - Run tests with race detector"
	@echo "  lint        - Run linter"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  clean       - Remove build artifacts"
