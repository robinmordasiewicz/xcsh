# vesctl Release Checklist

Pre-release validation checklist for vesctl compatibility and quality.

## Build Verification

- [ ] `go build -o vesctl .` completes without errors
- [ ] `go test ./...` passes all unit tests
- [ ] `./vesctl version` shows correct version info
- [ ] Binary runs on target platforms (linux-amd64, darwin-arm64, darwin-amd64)

## Compatibility Testing

### Offline Tests (No API Required)

- [ ] **Phase 1**: Configure/Authentication commands
  - `./run-all-tests.sh --phase 1`
  - Expected: All tests PASS

- [ ] **Phase 5**: Multi-Resource Validation (help text)
  - `./run-all-tests.sh --phase 5`
  - Expected: Most tests PASS (some WARN for new resources)

- [ ] **Phase 6**: Request Commands
  - `./run-all-tests.sh --phase 6`
  - Expected: PASS for common commands, documented differences for enhancements

- [ ] **Phase 7**: Site Management
  - `./run-all-tests.sh --phase 7`
  - Expected: PASS for base help, documented differences for create flags

### API Tests (Requires Credentials)

- [ ] **Phase 4**: Namespace CRUD
  - `./run-all-tests.sh --phase 4 --with-api`
  - Expected: All CRUD operations functional

- [ ] **Phase 8**: Output Formats
  - `./run-all-tests.sh --phase 8 --with-api`
  - Expected: Table/YAML identical, JSON semantically equivalent

- [ ] **Phase 9**: Error Handling
  - `./run-all-tests.sh --phase 9 --with-api`
  - Expected: Consistent error behavior (WARN for improved exit codes)

### Full Test Suite

- [ ] Run all phases with API:
  ```bash
  export F5XC_API_URL="https://your-tenant.console.ves.volterra.io"
  export F5XC_P12_FILE="/path/to/credentials.p12"
  export F5XC_P12_PASSWORD="your-password"
  export ORIGINAL_VESCTL="./vesctl-0.2.35"
  export OUR_F5XCCTL="./vesctl"
  ./run-all-tests.sh --with-api
  ```
- [ ] Review `summary.md` report
- [ ] Verify pass rate meets threshold (target: >95%)

## Documentation

- [ ] `ACCEPTABLE_DIFFERENCES.md` is up to date
- [ ] Any new differences are documented and justified
- [ ] README reflects current functionality
- [ ] CHANGELOG updated with new features/fixes

## Code Quality

- [ ] No `TODO` or `FIXME` comments for critical issues
- [ ] No hardcoded credentials or secrets
- [ ] Error messages are clear and actionable
- [ ] Help text is accurate and complete

## Security

- [ ] Credentials handled securely (not logged, not stored in plaintext)
- [ ] Certificate validation working correctly
- [ ] No command injection vulnerabilities
- [ ] P12/PEM file permissions respected

## CI/CD

- [ ] All GitHub Actions workflows pass
- [ ] Compatibility tests workflow runs successfully
- [ ] Release workflow configured correctly
- [ ] GoReleaser configuration valid

## Final Verification

- [ ] Test key workflows manually:
  ```bash
  # Authentication
  ./vesctl configure
  ./vesctl whoami

  # Basic CRUD
  ./vesctl configuration list namespace
  ./vesctl configuration get namespace default --outfmt json

  # Request commands
  ./vesctl request rpc --help
  ./vesctl request secrets --help
  ```

- [ ] Compare output with original vesctl for key commands
- [ ] Verify no regressions from previous release

## Release Steps

1. [ ] Update version in `cmd/version.go`
2. [ ] Create release tag: `git tag -a v0.x.x -m "Release v0.x.x"`
3. [ ] Push tag: `git push origin v0.x.x`
4. [ ] Verify GoReleaser builds succeed
5. [ ] Test downloaded binary from release artifacts
6. [ ] Update release notes with compatibility summary

## Post-Release

- [ ] Announce release (if applicable)
- [ ] Monitor for issues in first 24-48 hours
- [ ] Update any external documentation

---

## Quick Reference

### Environment Setup

```bash
# Required for API tests
export F5XC_API_URL="https://your-tenant.console.ves.volterra.io"
export F5XC_P12_FILE="/path/to/credentials.p12"
export F5XC_P12_PASSWORD="your-password"

# Required for comparison tests
export ORIGINAL_VESCTL="./vesctl-0.2.35"
export OUR_F5XCCTL="./vesctl"
```

### Common Commands

```bash
# Run all offline tests
./run-all-tests.sh

# Run all tests with API
./run-all-tests.sh --with-api

# Run specific phase
./run-all-tests.sh --phase 5

# Run with verbose output
./run-all-tests.sh --with-api -v
```

### Acceptable Test Results

| Result | Meaning | Action |
|--------|---------|--------|
| PASS | Output identical | None |
| WARN | Minor difference | Review, document if acceptable |
| FAIL | Significant difference | Investigate and fix |
| SKIP | Test not applicable | Verify skip reason |

---

*Last updated: 2024-12-06*
