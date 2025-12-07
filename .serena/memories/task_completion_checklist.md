# Task Completion Checklist for vesctl

## Before Committing Code

### 1. Format and Lint
```bash
make fmt                # Format Go code
make lint               # Run golangci-lint
```

### 2. Verify Compilation
```bash
make verify             # Run go build and go vet
```

### 3. Run Tests
```bash
make test-unit          # Run unit tests with race detection
```

### 4. Quick Full Check
```bash
make check              # Runs fmt, verify, and test-unit
```

## For Significant Changes

### Build and Test Binary
```bash
make build              # Build the binary
./vesctl --help         # Verify CLI works
./vesctl version        # Check version info
```

### Run Pre-commit Hooks
```bash
pre-commit run --all-files
```

### Update Documentation (if needed)
- Update docs/ markdown files
- Run `mkdocs serve` to preview

## Integration Testing (Optional)
Requires environment variables:
```bash
export VES_API_URL="https://tenant.staging.volterra.us"
export VES_API_P12_FILE="/path/to/cert.p12"
export VES_P12_PASSWORD="password"
make test-int
```

## Release Preparation
```bash
make release-dry        # Test GoReleaser
git tag v1.x.x          # Create version tag
git push origin v1.x.x  # Push tag (triggers GitHub Actions release)
```

## Commit Message Format
- Use conventional commits style when appropriate
- Be descriptive about what changed and why
- Reference issues if applicable
