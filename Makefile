.PHONY: help build test test-verbose test-coverage lint fmt vet clean install run-examples docs

# Variables
BINARY_NAME=better-auth-sdk-go
GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
GOFMT=$(GO) fmt
GOLINT=golangci-lint

# Default target
help:
	@echo "Better Auth SDK for Go - Makefile commands:"
	@echo ""
	@echo "  make build          - Build the project"
	@echo "  make test           - Run tests"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install        - Install dependencies"
	@echo "  make run-examples   - Run example programs"
	@echo "  make docs           - Generate documentation"
	@echo "  make all            - Run fmt, vet, lint, and test"
	@echo ""

# Build the project
build:
	@echo "Building..."
	$(GO) build -v ./...

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	$(GOTEST) -v -race -cover ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@which $(GOLINT) > /dev/null || (echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/" && exit 1)
	$(GOLINT) run ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -f coverage.out coverage.html
	rm -rf bin/
	rm -rf dist/

# Install dependencies
install:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Update dependencies
update:
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy

# Run example programs
run-examples:
	@echo "Running basic auth example..."
	$(GO) run examples/basic_auth/main.go
	@echo ""
	@echo "Running session management example..."
	$(GO) run examples/session_management/main.go

# Generate documentation
docs:
	@echo "Generating documentation..."
	$(GO) doc -all > docs.txt
	@echo "Documentation generated: docs.txt"

# Run all checks
all: fmt vet lint test
	@echo "All checks passed!"

# Check dependencies
check-deps:
	@echo "Checking dependencies..."
	$(GO) mod verify

# Security audit
security:
	@echo "Running security audit..."
	@which gosec > /dev/null || (echo "gosec not installed. Install it with: go install github.com/securego/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec ./...

# Benchmark tests
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Initialize project
init:
	@echo "Initializing project..."
	$(GO) mod init github.com/medapsis/better-auth-sdk-go || true
	$(GO) mod tidy

# Create release build
release:
	@echo "Creating release build..."
	mkdir -p bin
	GOOS=linux GOARCH=amd64 $(GO) build -o bin/$(BINARY_NAME)-linux-amd64
	GOOS=darwin GOARCH=amd64 $(GO) build -o bin/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(GO) build -o bin/$(BINARY_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 $(GO) build -o bin/$(BINARY_NAME)-windows-amd64.exe
	@echo "Release builds created in bin/"
