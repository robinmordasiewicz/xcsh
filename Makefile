# f5xcctl CLI Makefile
#
# Usage:
#   make build        - Build the f5xcctl binary for current platform
#   make build-all    - Build binaries for all platforms (linux/darwin/windows)
#   make test         - Run all tests
#   make test-unit    - Run unit tests only
#   make test-int     - Run integration tests only
#   make clean        - Clean build artifacts
#   make lint         - Run linter
#   make fmt          - Format code
#   make install      - Install binary to GOPATH/bin
#   make release-dry  - Test GoReleaser without publishing

BINARY_NAME=f5xcctl
MODULE=github.com/robinmordasiewicz/f5xcctl
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT?=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BRANCH?=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_AUTHOR?=$(shell git log -1 --format='%an' 2>/dev/null || echo "unknown")
BUILD_NUMBER?=0
LDFLAGS=-ldflags "-s -w -X ${MODULE}/cmd.Version=${VERSION} -X ${MODULE}/cmd.GitCommit=${GIT_COMMIT} -X ${MODULE}/cmd.BuildDate=${BUILD_TIME} -X ${MODULE}/cmd.Branch=${BRANCH} -X '${MODULE}/cmd.BuildAuthor=${BUILD_AUTHOR}' -X ${MODULE}/cmd.BuildNumber=${BUILD_NUMBER}"

# Build output directory
DIST_DIR=dist

# Platforms for cross-compilation
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: all build build-all test test-unit test-int clean lint fmt install help \
        release-dry release-snapshot verify check watch \
        build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 \
        docs docs-nav docs-clean docs-serve docs-check generate-examples \
        generate-schemas validate-schemas report-schemas generate-schemas-strict

# Default target
all: build

# Build the binary for current platform
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

# Build binaries for all platforms
build-all: clean-dist build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64
	@echo ""
	@echo "All builds complete:"
	@ls -lh $(DIST_DIR)/
	@echo ""
	@echo "Checksums:"
	@cd $(DIST_DIR) && shasum -a 256 * 2>/dev/null || sha256sum * 2>/dev/null || echo "Install shasum or sha256sum for checksums"

# Individual platform builds
build-linux-amd64:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .

build-linux-arm64:
	@echo "Building for Linux (arm64)..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .

build-darwin-amd64:
	@echo "Building for macOS (amd64)..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .

build-darwin-arm64:
	@echo "Building for macOS (arm64)..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .

build-windows-amd64:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Run all tests
test: build
	@echo "Running all tests..."
	@if [ -f ./scripts/test.sh ]; then \
		./scripts/test.sh; \
	else \
		go test -v ./pkg/...; \
	fi

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test -v -race ./pkg/...

# Run integration tests only (requires environment variables)
test-int: build
	@echo "Running integration tests..."
	@if [ -z "$(VES_API_URL)" ]; then \
		echo "Error: VES_API_URL not set"; \
		echo ""; \
		echo "Set these environment variables:"; \
		echo "  export VES_API_URL=\"https://tenant.staging.volterra.us\""; \
		echo "  export VES_API_P12_FILE=\"/path/to/cert.p12\""; \
		echo "  export VES_P12_PASSWORD=\"password\""; \
		exit 1; \
	fi
	VES_P12_PASSWORD="$(VES_P12_PASSWORD)" go test -v ./tests/integration/...

# Run tests with coverage
test-coverage: build
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out -race ./pkg/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Clean build artifacts
clean: clean-dist
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

clean-dist:
	@rm -rf $(DIST_DIR)/

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "  brew install golangci-lint"; \
		echo "  # or"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Install binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@if [ -z "$(GOPATH)" ]; then \
		echo "GOPATH not set, installing to ~/go/bin"; \
		mkdir -p ~/go/bin; \
		cp $(BINARY_NAME) ~/go/bin/; \
		echo "Installed to ~/go/bin/$(BINARY_NAME)"; \
	else \
		cp $(BINARY_NAME) $(GOPATH)/bin/; \
		echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"; \
	fi

# Verify code compiles
verify:
	@echo "Verifying code..."
	go build ./...
	go vet ./...

# Run quick validation (format, vet, test)
check: fmt verify test-unit
	@echo "All checks passed!"

# GoReleaser: dry run (test release without publishing)
release-dry:
	@echo "Running GoReleaser dry run..."
	@if command -v goreleaser > /dev/null; then \
		goreleaser check && \
		goreleaser release --snapshot --skip=publish --clean; \
	else \
		echo "goreleaser not installed. Install with:"; \
		echo "  brew install goreleaser"; \
		echo "  # or"; \
		echo "  go install github.com/goreleaser/goreleaser@latest"; \
	fi

# GoReleaser: snapshot build
release-snapshot:
	@echo "Building snapshot release..."
	@if command -v goreleaser > /dev/null; then \
		goreleaser release --snapshot --clean; \
	else \
		echo "goreleaser not installed. Install with:"; \
		echo "  brew install goreleaser"; \
	fi

# Development helper - rebuild on changes (requires fswatch)
watch:
	@if command -v fswatch > /dev/null; then \
		fswatch -o . -e ".*" -i "\\.go$$" | xargs -n1 -I{} make build; \
	else \
		echo "fswatch not installed. Install with: brew install fswatch"; \
	fi

# Documentation targets
PYTHON ?= python3
DOCS_OUTPUT = docs/commands
DOCS_TEMPLATES = scripts/templates

# Generate documentation from f5xcctl --spec
docs: build
	@echo "Generating documentation from CLI spec..."
	@$(PYTHON) scripts/generate-docs.py \
		--f5xcctl ./$(BINARY_NAME) \
		--output $(DOCS_OUTPUT) \
		--templates $(DOCS_TEMPLATES) \
		--clean \
		--update-mkdocs
	@echo "Generating install documentation..."
	@VERSION_OUTPUT=$$(./$(BINARY_NAME) version); \
	VERSION=$$(echo "$$VERSION_OUTPUT" | head -1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' || echo "dev"); \
	COMMIT=$$(echo "$$VERSION_OUTPUT" | grep 'commit:' | awk '{print $$2}' || echo "local"); \
	BUILT=$$(echo "$$VERSION_OUTPUT" | grep 'built:' | awk '{print $$2}' || echo "now"); \
	GO_VER=$$(echo "$$VERSION_OUTPUT" | grep 'go:' | awk '{print $$2}' || echo ""); \
	PLATFORM=$$(echo "$$VERSION_OUTPUT" | grep 'platform:' | awk '{print $$2}' || echo ""); \
	$(PYTHON) scripts/generate-homebrew-docs.py \
		--version "$$VERSION" \
		--commit "$$COMMIT" \
		--built "$$BUILT" \
		--go-version "$$GO_VER" \
		--platform "$$PLATFORM" \
		--output docs/install/homebrew.md; \
	$(PYTHON) scripts/generate-source-docs.py \
		--go-version "$$(go version | grep -oE 'go[0-9]+\.[0-9]+\.[0-9]+')" \
		--output docs/install/source.md
	@echo "Documentation generated successfully!"
	@echo "  Output: $(DOCS_OUTPUT)"
	@echo "  Files: $$(find $(DOCS_OUTPUT) -name '*.md' | wc -l) markdown files"

# Generate navigation only (faster, for mkdocs.yml updates)
docs-nav: build
	@echo "Generating navigation structure..."
	@$(PYTHON) scripts/generate-docs.py \
		--f5xcctl ./$(BINARY_NAME) \
		--nav-only \
		--update-mkdocs
	@echo "Navigation updated in mkdocs.yml"

# Clean generated documentation
docs-clean:
	@echo "Cleaning generated documentation..."
	@rm -rf $(DOCS_OUTPUT)/*
	@echo "Generated docs cleaned"

# Serve documentation locally with hot-reload
docs-serve: docs
	@echo "Starting documentation server..."
	@if command -v mkdocs > /dev/null; then \
		mkdocs serve; \
	else \
		echo "mkdocs not installed. Install with: pip install mkdocs mkdocs-material"; \
	fi

# Check current spec hash (useful for debugging idempotency)
docs-check: build
	@echo "Computing spec hash..."
	@./$(BINARY_NAME) --spec | sha256sum | cut -d' ' -f1
	@echo ""
	@echo "Spec statistics:"
	@echo "  Commands: $$(./$(BINARY_NAME) --spec | jq '.commands | length')"
	@echo "  Size: $$(./$(BINARY_NAME) --spec | wc -c) bytes"

# Generate examples from OpenAPI specifications
# This creates pkg/types/examples_generated.go with JSON examples for CLI help
generate-examples:
	@echo "Generating CLI examples from OpenAPI specifications..."
	@go run scripts/generate-examples.go -output pkg/types/examples_generated.go
	@echo "Examples generated successfully!"

# Generate resource schemas from OpenAPI specifications
# This creates pkg/types/schemas_generated.go with AI-friendly schema intelligence
generate-schemas:
	@echo "Generating resource schemas from OpenAPI specifications..."
	@go run scripts/generate-schemas.go -v
	@echo "Schema generation complete!"

# Validate schemas without regenerating (useful for CI)
validate-schemas:
	@echo "Validating schema generation..."
	@go run scripts/generate-schemas.go -validate -v

# Report missing specs (useful for debugging)
report-schemas:
	@echo "Generating schema coverage report..."
	@go run scripts/generate-schemas.go -report -v

# Strict schema generation (fails on missing critical resources)
generate-schemas-strict:
	@echo "Generating schemas (strict mode)..."
	@go run scripts/generate-schemas.go -v -strict

# Show version info
version:
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Branch: $(BRANCH)"
	@echo "Build Author: $(BUILD_AUTHOR)"
	@echo "Build Number: $(BUILD_NUMBER)"

# Show help
help:
	@echo "f5xcctl CLI Makefile"
	@echo ""
	@echo "Build Commands:"
	@echo "  make build          - Build binary for current platform"
	@echo "  make build-all      - Build binaries for all platforms (linux/darwin/windows)"
	@echo "  make install        - Install binary to GOPATH/bin"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "Test Commands:"
	@echo "  make test           - Run all tests"
	@echo "  make test-unit      - Run unit tests only"
	@echo "  make test-int       - Run integration tests (requires env vars)"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo ""
	@echo "Quality Commands:"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make verify         - Verify code compiles"
	@echo "  make check          - Run all checks (fmt, vet, test)"
	@echo ""
	@echo "Release Commands:"
	@echo "  make release-dry    - Test GoReleaser without publishing"
	@echo "  make release-snapshot - Build snapshot release"
	@echo "  make version        - Show version info"
	@echo ""
	@echo "Documentation Commands:"
	@echo "  make docs           - Generate documentation from CLI spec"
	@echo "  make docs-nav       - Update mkdocs.yml navigation only"
	@echo "  make docs-clean     - Clean generated documentation"
	@echo "  make docs-serve     - Generate docs and serve locally"
	@echo "  make docs-check     - Show current spec hash"
	@echo ""
	@echo "Code Generation Commands:"
	@echo "  make generate-examples - Generate CLI examples from OpenAPI specs"
	@echo "  make generate-schemas  - Generate resource schemas from OpenAPI specs"
	@echo "  make validate-schemas  - Validate schemas without regenerating"
	@echo "  make report-schemas    - Report missing specs and coverage"
	@echo "  make generate-schemas-strict - Generate schemas, fail on missing critical"
	@echo ""
	@echo "Development Commands:"
	@echo "  make watch          - Rebuild on file changes"
	@echo ""
	@echo "Environment Variables (for integration tests):"
	@echo "  VES_API_URL        - API URL"
	@echo "  VES_API_P12_FILE   - Path to P12 certificate bundle"
	@echo "  VES_P12_PASSWORD   - Password for P12 bundle"
	@echo ""
	@echo "Creating a Release:"
	@echo "  1. Update version: git tag v1.0.0"
	@echo "  2. Push tag: git push origin v1.0.0"
	@echo "  3. GitHub Actions will automatically build and publish"
