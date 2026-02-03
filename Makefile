# Makefile for Ngatex API Gateway

BINARY_NAME=ngatex
GO_FILES=$(shell find . -name "*.go")

.PHONY: all build run test clean docker-build help

all: build

build: ## Build the gateway binary
	@echo "Building binary..."
	go build -o $(BINARY_NAME) ./cmd/gateway

run: build ## Build and run the gateway with default config
	./$(BINARY_NAME) --config=config.yaml

mock: ## Run the mock services for testing
	go run ./cmd/run-mocks/main.go

test: ## Run all tests
	go test -v ./...

clean: ## Remove binary and build artifacts
	rm -f $(BINARY_NAME)

docker-build: ## Build the production Docker image
	docker build -t ngatex:latest .

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'