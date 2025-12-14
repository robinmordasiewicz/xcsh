# Suggested Commands for f5xcctl Development

## Build Commands
```bash
make build              # Build binary for current platform (darwin-arm64)
make build-all          # Build for all platforms (linux/darwin/windows)
make install            # Install binary to GOPATH/bin
make clean              # Clean build artifacts
go build -o f5xcctl .    # Direct Go build
```

## Test Commands
```bash
make test               # Run all tests
make test-unit          # Run unit tests only (go test -v -race ./pkg/...)
make test-int           # Run integration tests (requires F5XC_API_URL, F5XC_API_P12_FILE, F5XC_P12_PASSWORD)
make test-coverage      # Run tests with coverage report
```

## Quality Commands
```bash
make fmt                # Format Go code (go fmt ./...)
make lint               # Run golangci-lint (timeout=5m)
make verify             # Verify code compiles (go build + go vet)
make check              # Run all checks (fmt, vet, test)
```

## Pre-commit Hooks
```bash
pre-commit install      # Install pre-commit hooks
pre-commit run --all-files  # Run all hooks manually
```

## Release Commands
```bash
make release-dry        # Test GoReleaser without publishing
make release-snapshot   # Build snapshot release
make version            # Show version info
```

## Development Commands
```bash
make watch              # Rebuild on file changes (requires fswatch)
./f5xcctl --help         # Test the built CLI
./f5xcctl version        # Check version info
```

## Documentation
```bash
pip install -r requirements-docs.txt  # Install MkDocs dependencies
mkdocs serve            # Serve docs locally at http://127.0.0.1:8000
mkdocs build            # Build static documentation
```

## System Utilities (Darwin)
```bash
git status              # Check git status
git branch              # List branches
ls -la                  # List files
grep -r "pattern" .     # Search in files
find . -name "*.go"     # Find Go files
```
