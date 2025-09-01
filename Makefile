.PHONY: build test lint clean run docker-build docker-push help

# Binary name
BINARY_NAME=keydra
MAIN_PACKAGE=./cmd/keydra

# Docker image settings
DOCKER_REGISTRY ?= docker.io
DOCKER_IMAGE ?= $(DOCKER_REGISTRY)/m4rk0g/keydra
DOCKER_TAG ?= latest

# Go settings
GOOS ?= linux
GOARCH ?= amd64

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-s -w" -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

run: ## Run the application locally
	@echo "Running $(BINARY_NAME) locally..."
	go run $(MAIN_PACKAGE) -config=config.yaml.example

docker-build: ## Build Docker image
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-push: docker-build ## Build and push Docker image
	@echo "Pushing Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

install-deps: ## Install development dependencies
	@echo "Installing dependencies..."
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

deps: ## Update dependencies
	@echo "Updating dependencies..."
	go mod tidy
	go mod verify
