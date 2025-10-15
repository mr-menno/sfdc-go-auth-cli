# Variables
BINARY_NAME=sfdc-auth
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

# Default target
.PHONY: all
all: clean test lint build

# Clean build artifacts
.PHONY: clean
clean:
	rm -f ${BINARY_NAME}
	rm -rf dist/
	go clean

# Run tests
.PHONY: test
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
.PHONY: test-coverage
test-coverage: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build for current platform
.PHONY: build
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} main.go

# Build for all platforms
.PHONY: build-all
build-all: clean
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe main.go

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Run linter (MANDATORY - zero errors required)
.PHONY: lint
lint:
	@echo "🔍 Running golangci-lint (zero errors required)..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "❌ golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi
	@echo "✅ Linting passed!"

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Run the application
.PHONY: run
run: build
	./${BINARY_NAME}

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t ${BINARY_NAME}:${VERSION} .
	docker build -t ${BINARY_NAME}:latest .

# Run Docker container
.PHONY: docker-run
docker-run: docker-build
	docker run -it --rm -p 8080:8080 ${BINARY_NAME}:latest

# Create release archive
.PHONY: release
release: build-all
	cd dist && \
	for file in *; do \
		if [[ $$file == *.exe ]]; then \
			zip $${file%.exe}.zip $$file; \
		else \
			tar -czf $$file.tar.gz $$file; \
		fi; \
	done

# Pre-commit validation (runs all quality checks)
.PHONY: pre-commit
pre-commit: fmt test lint
	@echo "✅ All quality checks passed! Ready to commit."

# Development setup
.PHONY: dev-setup
dev-setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go mod download

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, test, and build"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  build        - Build for current platform"
	@echo "  build-all    - Build for all platforms"
	@echo "  deps         - Install dependencies"
	@echo "  lint         - Run linter (MANDATORY - zero errors required)"
	@echo "  fmt          - Format code"
	@echo "  run          - Build and run the application"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Build and run Docker container"
	@echo "  release      - Create release archives"
	@echo "  pre-commit   - Run all quality checks (fmt, test, lint)"
	@echo "  dev-setup    - Set up development environment"
	@echo "  help         - Show this help message"
