#!/bin/bash

# Development Environment Setup Script
# This script sets up the development environment for the Salesforce OAuth CLI

set -e

echo "🚀 Setting up development environment for Salesforce OAuth CLI..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.19+ first."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.19"

if ! printf '%s\n%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V -C; then
    echo "❌ Go version $GO_VERSION is too old. Please upgrade to Go $REQUIRED_VERSION or later."
    exit 1
fi

echo "✅ Go version $GO_VERSION is compatible"

# Install golangci-lint
echo "📦 Installing golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    echo "✅ golangci-lint is already installed"
else
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    echo "✅ golangci-lint installed successfully"
fi

# Download dependencies
echo "📦 Downloading Go dependencies..."
go mod download
go mod tidy

# Run initial quality checks
echo "🔍 Running initial quality checks..."

# Format code
echo "  - Formatting code..."
go fmt ./...

# Run tests
echo "  - Running tests..."
go test -v ./...

# Run linter
echo "  - Running linter..."
golangci-lint run --timeout=5m

echo ""
echo "🎉 Development environment setup complete!"
echo ""
echo "Available commands:"
echo "  make help        - Show all available make targets"
echo "  make pre-commit  - Run all quality checks before committing"
echo "  make test        - Run tests"
echo "  make lint        - Run linter (MANDATORY - zero errors required)"
echo "  make build       - Build the application"
echo ""
echo "⚠️  IMPORTANT: All code must pass linting with zero errors before any work is considered complete!"
