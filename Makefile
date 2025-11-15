# Task Manager API - Makefile
# A convenient way to manage the project during hackathon development

.PHONY: help run build test clean dev install format lint

# Default target
help:
	@echo "🚀 Task Manager API - Development Commands"
	@echo "=========================================="
	@echo ""
	@echo "Available commands:"
	@echo "  make run       - Run the development server"
	@echo "  make build     - Build the production binary"
	@echo "  make test      - Run tests (when implemented)"
	@echo "  make dev       - Run with auto-reload (requires air)"
	@echo "  make install   - Install dependencies"
	@echo "  make format    - Format Go code"
	@echo "  make lint      - Run linter (requires golangci-lint)"
	@echo "  make clean     - Clean build artifacts"
	@echo ""

# Run the development server
run:
	@echo "🚀 Starting Task Manager API..."
	go run cmd/server/main.go

# Build the production binary
build:
	@echo "🔨 Building production binary..."
	go build -ldflags="-s -w" -o bin/task-manager cmd/server/main.go
	@echo "✅ Binary created at bin/task-manager"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Development with auto-reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "🔄 Starting development server with auto-reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Falling back to regular run..."; \
		make run; \
	fi

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

# Format Go code
format:
	@echo "🎨 Formatting Go code..."
	go fmt ./...
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	fi

# Run linter
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "❌ golangci-lint not found. Install from https://golangci-lint.run/"; \
		echo "Falling back to go vet..."; \
		go vet ./...; \
	fi

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Create necessary directories
dirs:
	@mkdir -p bin logs

# Quick setup for new developers
setup: install format
	@echo "✅ Project setup complete!"
	@echo "Run 'make run' to start the server"
