# xcsh CLI Makefile (TypeScript/Node.js)
#
# Usage:
#   make build        - Build the CLI for current platform
#   make build-all    - Build binaries for all platforms
#   make test         - Run all tests
#   make lint         - Run linter
#   make clean        - Clean build artifacts
#   make docs         - Generate documentation

# Include branding configuration (single source of truth)
include branding.mk

# Use branding variables
BINARY_NAME=$(CLI_NAME)
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build output directory
DIST_DIR=binaries

# Node.js/npm settings
NODE_VERSION_MIN=18
NPM?=npm
NPX?=npx

.PHONY: all build build-all test clean lint fmt install help \
        docs docs-all docs-clean docs-serve docs-build \
        download-specs download-specs-force check-upstream \
        generate generate-domains validate-generated \
        ts ts-build ts-test ts-lint ts-check ts-install ts-generate \
        ci pre-commit pre-push version

# Default target
all: build

# =============================================================================
# Build Targets
# =============================================================================

# Build the CLI for current platform
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@$(NPM) run build
	@echo "Build complete: dist/"

# Build standalone binaries for all platforms
build-all: clean-dist
	@echo "Building binaries for all platforms..."
	@$(NPM) run build:binaries
	@echo ""
	@echo "All builds complete:"
	@ls -lh $(DIST_DIR)/
	@echo ""
	@echo "Checksums:"
	@cd $(DIST_DIR) && shasum -a 256 * 2>/dev/null || sha256sum * 2>/dev/null || echo "Install shasum or sha256sum for checksums"

# =============================================================================
# Test Targets
# =============================================================================

# Run all tests
test:
	@echo "Running all tests..."
	@$(NPM) test -- --run

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@$(NPM) test -- --run --coverage
	@echo "Coverage report generated"

# Run tests in watch mode
test-watch:
	@echo "Running tests in watch mode..."
	@$(NPM) test

# =============================================================================
# Quality Targets
# =============================================================================

# Run linter
lint:
	@echo "Running linter..."
	@$(NPM) run lint
	@$(NPM) run format:check
	@echo "Lint passed"

# Run type checking
typecheck:
	@echo "Running type check..."
	@$(NPM) run typecheck
	@echo "Type check passed"

# Format code
fmt:
	@echo "Formatting code..."
	@$(NPM) run format
	@echo "Code formatted"

# Run all validation (typecheck + lint)
check: typecheck lint
	@echo "All checks passed!"

# Python linting (requires ruff)
lint-python:
	@echo "Linting Python scripts..."
	@if command -v ruff > /dev/null; then \
		ruff check scripts/*.py; \
	else \
		echo "ruff not installed. Install with: pip install ruff"; \
	fi

# Format Python scripts
fmt-python:
	@echo "Formatting Python scripts..."
	@if command -v ruff > /dev/null; then \
		ruff format scripts/*.py; \
		ruff check --fix scripts/*.py; \
	else \
		echo "ruff not installed. Install with: pip install ruff"; \
	fi

# Format shell scripts (requires shfmt)
fmt-shell:
	@echo "Formatting shell scripts..."
	@if command -v shfmt > /dev/null; then \
		shfmt -i 2 -ci -bn -w scripts/*.sh install.sh; \
	else \
		echo "shfmt not installed. Install with: brew install shfmt"; \
	fi

# Check shell script formatting
check-shell:
	@echo "Checking shell script formatting..."
	@if command -v shfmt > /dev/null; then \
		shfmt -i 2 -ci -bn -d scripts/*.sh install.sh; \
	else \
		echo "shfmt not installed. Install with: brew install shfmt"; \
	fi

# =============================================================================
# Clean Targets
# =============================================================================

# Clean build artifacts
clean: clean-dist
	@echo "Cleaning..."
	@rm -rf dist/
	@rm -rf node_modules/.cache/
	@echo "Clean complete"

clean-dist:
	@rm -rf $(DIST_DIR)/

# =============================================================================
# Install Targets
# =============================================================================

# Install dependencies
install:
	@echo "Installing dependencies..."
	@$(NPM) ci
	@echo "Dependencies installed"

# Link for local development
link:
	@echo "Linking for local development..."
	@$(NPM) link
	@echo "Linked: $(BINARY_NAME)"

# =============================================================================
# API Specifications
# =============================================================================

# Download enriched API specifications from GitHub releases
download-specs:
	@echo "Downloading enriched API specifications..."
	@./scripts/download-specs.sh

# Force re-download enriched API specifications (bypasses cache)
download-specs-force:
	@echo "Force downloading enriched API specifications..."
	@rm -rf .specs
	@./scripts/download-specs.sh

# Check for upstream API specification updates
check-upstream:
	@./scripts/check-upstream.sh

# Check for upstream updates (JSON output for scripting)
check-upstream-json:
	@./scripts/check-upstream.sh --json

# =============================================================================
# Code Generation
# =============================================================================

# Generate domain registry from upstream specs
generate: download-specs generate-domains validate-generated
	@echo ""
	@echo "✅ Code generation complete!"
	@echo "   Domains: src/types/domains_generated.ts"

# Generate TypeScript domain registry from specs
generate-domains: download-specs
	@echo "🏗️  Generating domains from upstream specs..."
	@$(NPX) tsx scripts/generate-domains.ts
	@echo "✓ Generated: src/types/domains_generated.ts"

# Validate generated files are present
validate-generated:
	@echo "🔍 Validating generated code..."
	@test -f src/types/domains_generated.ts || (echo "❌ domains_generated.ts missing" && exit 1)
	@echo "✓ All generated files present"

# =============================================================================
# Documentation Generation
# =============================================================================

PYTHON ?= python3
DOCS_OUTPUT = docs/commands
DOCS_TEMPLATES = scripts/templates

# Generate documentation from CLI
docs: build
	@echo "Generating documentation..."
	@VERSION=$$(node -p "require('./package.json').version"); \
	NODE_VER=$$(node --version); \
	$(PYTHON) scripts/generate-homebrew-docs.py \
		--version "$$VERSION" \
		--node-version "$$NODE_VER" \
		--output docs/install/homebrew.md; \
	$(PYTHON) scripts/generate-source-docs.py \
		--node-version "$$NODE_VER" \
		--output docs/install/source.md
	@echo "Documentation generated!"

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

# Generate ALL documentation
docs-all: build
	@echo "Generating all documentation..."
	@VERSION=$$(node -p "require('./package.json').version"); \
	NODE_VER=$$(node --version); \
	echo ""; \
	echo "Step 1/2: Homebrew documentation..."; \
	$(PYTHON) scripts/generate-homebrew-docs.py \
		--version "$$VERSION" \
		--node-version "$$NODE_VER" \
		--output docs/install/homebrew.md; \
	echo ""; \
	echo "Step 2/2: Source build documentation..."; \
	$(PYTHON) scripts/generate-source-docs.py \
		--node-version "$$NODE_VER" \
		--output docs/install/source.md
	@echo ""
	@echo "✅ All documentation generated!"

# Build the MkDocs documentation site
docs-build: docs-all
	@echo "Building MkDocs site..."
	@if command -v mkdocs > /dev/null; then \
		mkdocs build --strict; \
		echo "✅ Site built successfully in site/"; \
	else \
		echo "mkdocs not installed. Install with: pip install -r requirements-docs.txt"; \
		exit 1; \
	fi

# =============================================================================
# CI/CD Consistency Targets
# =============================================================================

# Run the full CI pipeline locally
ci: lint typecheck test build
	@echo ""
	@echo "✅ CI pipeline completed successfully!"

# Pre-commit hook: fast checks before committing
pre-commit: fmt lint typecheck
	@echo ""
	@echo "✅ Pre-commit checks passed!"

# Pre-push hook: comprehensive checks before pushing
pre-push: ci docs-all
	@echo ""
	@echo "✅ Pre-push checks passed!"
	@echo "   Your changes are ready to push."

# =============================================================================
# TypeScript CLI Aliases (for compatibility)
# =============================================================================

ts: check test build
	@echo "✅ TypeScript CLI all checks passed"

ts-build: build

ts-test: test

ts-lint: lint

ts-check: check

ts-install: install

ts-generate: generate-domains

ts-check-upstream: check-upstream

# =============================================================================
# Version Info
# =============================================================================

version:
	@echo "Version: $(VERSION)"
	@node -p "require('./package.json').version" 2>/dev/null || echo "Package version: unknown"
	@node --version 2>/dev/null || echo "Node.js: not installed"
	@$(NPM) --version 2>/dev/null | xargs -I{} echo "npm: {}" || echo "npm: not installed"

# =============================================================================
# Help
# =============================================================================

help:
	@echo "$(CLI_NAME) CLI Makefile"
	@echo ""
	@echo "=== CI/CD Consistency Targets (RECOMMENDED) ==="
	@echo "  make ci             - Run full CI pipeline locally"
	@echo "  make pre-commit     - Fast checks before committing (fmt, lint, typecheck)"
	@echo "  make pre-push       - Comprehensive checks before pushing (ci + docs)"
	@echo ""
	@echo "=== Build Commands ==="
	@echo "  make build          - Build for current platform"
	@echo "  make build-all      - Build binaries for all platforms"
	@echo "  make install        - Install npm dependencies"
	@echo "  make link           - Link for local development"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "=== Test Commands ==="
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make test-watch     - Run tests in watch mode"
	@echo ""
	@echo "=== Quality Commands ==="
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make typecheck      - Run TypeScript type checking"
	@echo "  make check          - Run all checks (typecheck + lint)"
	@echo ""
	@echo "=== Documentation Commands ==="
	@echo "  make docs           - Generate documentation"
	@echo "  make docs-all       - Generate ALL documentation"
	@echo "  make docs-build     - Generate docs and build MkDocs site"
	@echo "  make docs-clean     - Clean generated documentation"
	@echo "  make docs-serve     - Generate docs and serve locally"
	@echo ""
	@echo "=== API Specifications ==="
	@echo "  make download-specs     - Download latest enriched API specs"
	@echo "  make download-specs-force - Force re-download specs"
	@echo "  make check-upstream     - Check if upstream specs have updates"
	@echo ""
	@echo "=== Code Generation ==="
	@echo "  make generate           - Run full generation pipeline"
	@echo "  make generate-domains   - Generate domain registry from specs"
	@echo "  make validate-generated - Validate generated files are present"
	@echo ""
	@echo "=== Compatibility Aliases ==="
	@echo "  make ts                 - Run all checks (alias for check + test + build)"
	@echo "  make ts-build           - Build (alias)"
	@echo "  make ts-test            - Test (alias)"
	@echo "  make ts-lint            - Lint (alias)"
	@echo "  make ts-check           - Check (alias)"
	@echo "  make ts-install         - Install (alias)"
	@echo "  make ts-generate        - Generate domains (alias)"
	@echo ""
	@echo "=== Development Workflow ==="
	@echo "  Before committing:  make pre-commit"
	@echo "  Before pushing:     make pre-push"
	@echo "  Full CI locally:    make ci"
	@echo ""
	@echo "=== Version Info ==="
	@echo "  make version        - Show version information"
