# CLI Makefile
#
# Usage:
#   make build        - Build the CLI binary for current platform
#   make build-all    - Build binaries for all platforms (linux/darwin/windows)
#   make test         - Run all tests
#   make test-unit    - Run unit tests only
#   make test-int     - Run integration tests only
#   make clean        - Clean build artifacts
#   make lint         - Run linter
#   make fmt          - Format code
#   make install      - Install binary to GOPATH/bin
#   make release-dry  - Test GoReleaser without publishing

# Include branding configuration (single source of truth)
include branding.mk

# Use branding variables
BINARY_NAME=$(CLI_NAME)
MODULE=$(MODULE_PATH)
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

# LLM description generation settings
LLM_WORKERS?=8

.PHONY: all build build-all test test-unit test-int clean lint fmt install help \
        release-dry release-snapshot verify check watch \
        build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 \
        docs docs-all docs-nav docs-clean docs-serve docs-check docs-build generate-examples \
        generate-schemas validate-schemas report-schemas generate-schemas-strict \
        generate-llm-descriptions generate-schemas-with-llm maybe-llm-descriptions \
        ci pre-commit pre-push verify-schemas-ci verify-lint-config \
        download-specs download-specs-force

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
	@if [ -z "$(F5XC_API_URL)" ]; then \
		echo "Error: F5XC_API_URL not set"; \
		echo ""; \
		echo "Set these environment variables:"; \
		echo "  export F5XC_API_URL=\"https://tenant.staging.volterra.us\""; \
		echo "  export F5XC_API_P12_FILE=\"/path/to/cert.p12\""; \
		echo "  export F5XC_P12_PASSWORD=\"password\""; \
		exit 1; \
	fi
	F5XC_P12_PASSWORD="$(F5XC_P12_PASSWORD)" go test -v ./tests/integration/...

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

# Canonical golangci-lint version - single source of truth
# This MUST match .github/workflows/ci.yml env.GOLANGCI_LINT_VERSION
GOLANGCI_LINT_VERSION := v2.7.2

# Verify lint configuration consistency (mirrors CI verify-lint-config job)
verify-lint-config:
	@echo "Verifying golangci-lint version consistency..."
	@PRECOMMIT_VERSION=$$(grep -A2 'repo: https://github.com/golangci/golangci-lint' .pre-commit-config.yaml | grep 'rev:' | sed 's/.*rev: //'); \
	CI_VERSION=$$(grep 'GOLANGCI_LINT_VERSION:' .github/workflows/ci.yml | head -1 | sed 's/.*: //'); \
	echo "  Makefile version:    $(GOLANGCI_LINT_VERSION)"; \
	echo "  Pre-commit version:  $$PRECOMMIT_VERSION"; \
	echo "  CI workflow version: $$CI_VERSION"; \
	if [ "$(GOLANGCI_LINT_VERSION)" != "$$PRECOMMIT_VERSION" ]; then \
		echo ""; \
		echo "❌ ERROR: Pre-commit version mismatch!"; \
		echo "   Expected: $(GOLANGCI_LINT_VERSION)"; \
		echo "   Found:    $$PRECOMMIT_VERSION"; \
		echo "   Update .pre-commit-config.yaml to use $(GOLANGCI_LINT_VERSION)"; \
		exit 1; \
	fi; \
	if [ "$(GOLANGCI_LINT_VERSION)" != "$$CI_VERSION" ]; then \
		echo ""; \
		echo "❌ ERROR: CI workflow version mismatch!"; \
		echo "   Expected: $(GOLANGCI_LINT_VERSION)"; \
		echo "   Found:    $$CI_VERSION"; \
		echo "   Update .github/workflows/ci.yml GOLANGCI_LINT_VERSION"; \
		exit 1; \
	fi
	@echo "✅ All golangci-lint versions are consistent"

# Run linter (requires golangci-lint)
lint: verify-lint-config
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

# Generate documentation from CLI spec
docs: build
	@echo "Generating documentation from CLI spec..."
	@$(PYTHON) scripts/generate-docs.py \
		--cli-binary ./$(BINARY_NAME) \
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
		--cli-binary ./$(BINARY_NAME) \
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

# Download enriched API specifications from GitHub releases
# This fetches the latest enriched specs and caches them locally
download-specs:
	@echo "Downloading enriched API specifications..."
	@./scripts/download-specs.sh

# Force re-download enriched API specifications (bypasses cache)
download-specs-force:
	@echo "Force downloading enriched API specifications..."
	@rm -rf .specs
	@./scripts/download-specs.sh

# Generate examples from OpenAPI specifications
# This creates pkg/types/examples_generated.go with JSON examples for CLI help
generate-examples: download-specs
	@echo "Generating CLI examples from OpenAPI specifications..."
	@go run scripts/generate-examples.go -output pkg/types/examples_generated.go
	@echo "Examples generated successfully!"

# Generate resource schemas from OpenAPI specifications
# This creates pkg/types/schemas_generated.go with AI-friendly schema intelligence
generate-schemas: download-specs
	@echo "Generating resource schemas from OpenAPI specifications..."
	@go run scripts/generate-schemas.go -v -update-resources
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

# =============================================================================
# CI/CD Consistency Targets
# These targets mirror the GitHub Actions CI pipeline for local development
# =============================================================================

# Run the full CI pipeline locally (mirrors ci.yml workflow)
# This is what GitHub Actions runs on every push/PR
ci: lint test-unit verify-schemas-ci build release-dry
	@echo ""
	@echo "✅ CI pipeline completed successfully!"
	@echo "   This mirrors the GitHub Actions CI workflow."

# Verify schemas match CI expectations (mirrors verify-schemas job in ci.yml)
verify-schemas-ci:
	@echo "Verifying schemas (CI mode)..."
	@cp pkg/types/schemas_generated.go /tmp/schemas_before.go
	@$(MAKE) generate-schemas 2>/dev/null
	@if ! diff -q pkg/types/schemas_generated.go /tmp/schemas_before.go > /dev/null 2>&1; then \
		echo "::error::Schema file is out of sync with OpenAPI specs!"; \
		echo "Run 'make generate-schemas' locally and commit the changes."; \
		diff pkg/types/schemas_generated.go /tmp/schemas_before.go || true; \
		mv /tmp/schemas_before.go pkg/types/schemas_generated.go; \
		exit 1; \
	fi
	@rm /tmp/schemas_before.go
	@echo "✅ Schemas are up to date and idempotent"

# Pre-commit hook: fast checks before committing
# Run this before every commit to catch issues early
pre-commit: fmt lint test-unit
	@echo ""
	@echo "✅ Pre-commit checks passed!"

# Pre-push hook: comprehensive checks before pushing
# Run this before pushing to ensure CI will pass
pre-push: ci docs-all
	@echo ""
	@echo "✅ Pre-push checks passed!"
	@echo "   Your changes are ready to push."

# =============================================================================
# Documentation Generation (mirrors docs.yml workflow)
# =============================================================================

# Generate ALL documentation (mirrors docs.yml generate job)
docs-all: build
	@echo "Generating all documentation..."
	@echo ""
	@echo "Step 1/4: Command documentation..."
	@$(PYTHON) scripts/generate-docs.py \
		--cli-binary ./$(BINARY_NAME) \
		--output $(DOCS_OUTPUT) \
		--templates $(DOCS_TEMPLATES) \
		--clean \
		--update-mkdocs
	@echo ""
	@echo "Step 2/4: CloudStatus documentation..."
	@$(PYTHON) scripts/generate-cloudstatus-docs.py \
		--cli-binary ./$(BINARY_NAME) \
		--output docs/commands/cloudstatus \
		--clean
	@echo ""
	@echo "Step 3/4: Homebrew documentation..."
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
		--output docs/install/homebrew.md
	@echo ""
	@echo "Step 4/4: Source build documentation..."
	@$(PYTHON) scripts/generate-source-docs.py \
		--go-version "$$(go version | grep -oE 'go[0-9]+\.[0-9]+\.[0-9]+')" \
		--output docs/install/source.md
	@echo ""
	@echo "✅ All documentation generated!"
	@echo "   Command docs: $(DOCS_OUTPUT)"
	@echo "   Files: $$(find $(DOCS_OUTPUT) -name '*.md' | wc -l) markdown files"

# Build the MkDocs documentation site (mirrors docs.yml build job)
docs-build: docs-all
	@echo "Building MkDocs site..."
	@if command -v mkdocs > /dev/null; then \
		mkdocs build --strict; \
		echo "✅ Site built successfully in site/"; \
	else \
		echo "mkdocs not installed. Install with: pip install -r requirements-docs.txt"; \
		exit 1; \
	fi

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
	@echo "$(CLI_NAME) CLI Makefile"
	@echo ""
	@echo "=== CI/CD Consistency Targets (RECOMMENDED) ==="
	@echo "  make ci             - Run full CI pipeline locally (mirrors GitHub Actions)"
	@echo "  make pre-commit     - Fast checks before committing (fmt, lint, test)"
	@echo "  make pre-push       - Comprehensive checks before pushing (ci + docs)"
	@echo ""
	@echo "=== Build Commands ==="
	@echo "  make build          - Build binary for current platform"
	@echo "  make build-all      - Build binaries for all platforms (linux/darwin/windows)"
	@echo "  make install        - Install binary to GOPATH/bin"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "=== Test Commands ==="
	@echo "  make test           - Run all tests"
	@echo "  make test-unit      - Run unit tests only"
	@echo "  make test-int       - Run integration tests (requires env vars)"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo ""
	@echo "=== Quality Commands ==="
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter (includes version check)"
	@echo "  make verify         - Verify code compiles"
	@echo "  make check          - Run all checks (fmt, vet, test)"
	@echo "  make verify-lint-config - Verify golangci-lint versions match"
	@echo "  make verify-schemas-ci - Verify schemas match CI expectations"
	@echo ""
	@echo "=== Release Commands ==="
	@echo "  make release-dry    - Test GoReleaser without publishing"
	@echo "  make release-snapshot - Build snapshot release"
	@echo "  make version        - Show version info"
	@echo ""
	@echo "=== Documentation Commands ==="
	@echo "  make docs           - Generate command documentation only"
	@echo "  make docs-all       - Generate ALL documentation (mirrors CI)"
	@echo "  make docs-build     - Generate docs and build MkDocs site"
	@echo "  make docs-nav       - Update mkdocs.yml navigation only"
	@echo "  make docs-clean     - Clean generated documentation"
	@echo "  make docs-serve     - Generate docs and serve locally"
	@echo "  make docs-check     - Show current spec hash"
	@echo ""
	@echo "=== API Specifications ==="
	@echo "  make download-specs     - Download latest enriched API specs (auto-cached)"
	@echo "  make download-specs-force - Force re-download specs (bypass cache)"
	@echo ""
	@echo "=== Code Generation Commands ==="
	@echo "  make generate-examples - Generate CLI examples from OpenAPI specs"
	@echo "  make generate-schemas  - Generate resource schemas from OpenAPI specs"
	@echo "  make validate-schemas  - Validate schemas without regenerating"
	@echo "  make report-schemas    - Report missing specs and coverage"
	@echo "  make generate-schemas-strict - Generate schemas, fail on missing critical"
	@echo "  make generate-llm-descriptions - Generate LLM descriptions (requires Ollama)"
	@echo "  make generate-schemas-with-llm - Regenerate schemas with LLM descriptions"
	@echo "  make maybe-llm-descriptions    - Auto-detect Ollama and regenerate if available"
	@echo ""
	@echo "=== Development Workflow ==="
	@echo "  make watch          - Rebuild on file changes"
	@echo ""
	@echo "=== Recommended Workflow ==="
	@echo "  Before committing:  make pre-commit"
	@echo "  Before pushing:     make pre-push"
	@echo "  Full CI locally:    make ci"
	@echo ""
	@echo "=== Environment Variables ==="
	@echo "  F5XC_API_URL        - API URL (for integration tests)"
	@echo "  F5XC_API_P12_FILE   - Path to P12 certificate bundle"
	@echo "  F5XC_P12_PASSWORD   - Password for P12 bundle"
	@echo "  LLM_WORKERS         - Number of parallel LLM workers (default: 8)"
	@echo ""
	@echo "=== Creating a Release ==="
	@echo "  Releases are fully automated via GitHub Actions:"
	@echo "  1. Commit changes to main branch"
	@echo "  2. Push to origin (git push)"
	@echo "  3. GitHub Actions automatically:"
	@echo "     - Runs CI tests"
	@echo "     - Creates semantic version tag"
	@echo "     - Generates LLM descriptions (if needed)"
	@echo "     - Builds and signs binaries"
	@echo "     - Updates Homebrew tap"
	@echo "     - Generates and deploys documentation"
