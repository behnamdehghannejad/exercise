# Makefile for Go project

.PHONY: help build run test clean lint fmt vet

# Default target
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make lint     - Run linter"
	@echo "  make fmt      - Format code"
	@echo "  make vet      - Run go vet"
	@echo "  make all      - Run fmt, vet, lint, test, build"

# Build binary
build:
	@echo "Building..."
	go build -o bin/myapp ./...

# Run application
run:
	@echo "Running..."
	go run ./...

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run all checks
all: fmt vet lint test build
	@echo "All checks passed!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy