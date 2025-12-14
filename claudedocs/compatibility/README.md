# vesctl Compatibility Testing

This directory contains the compatibility testing framework for comparing our vesctl implementation against the original F5 proprietary binary.

## Original Binary

The original proprietary vesctl binary is stored locally for comparison testing.

### IMPORTANT: Known Bug in Original vesctl

**Versions 0.2.36 and later have a critical TLS bug** that causes crashes with the following error:

```
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation]
...
crypto/tls.(*CertificateRequestInfo).SupportsCertificate(0x0, 0x1400010e040)
```

This bug affects ALL platforms (macOS arm64, macOS amd64, Linux amd64) and occurs during any authenticated API operation.

### Working Version for Testing

**Use version 0.2.35 on Linux amd64** for compatibility testing:
- **Version**: 0.2.35 (last working version)
- **Platform**: Linux amd64 only
- **Source**: https://vesio.azureedge.net/releases/vesctl/0.2.35/vesctl.linux-amd64.gz

> **Note**: There is NO working macOS arm64 version of the original vesctl. All compatibility testing with authenticated API operations must be performed on a Linux amd64 host.

### Download Original Binary (Linux amd64)

```bash
# Download working version 0.2.35 for Linux amd64
curl -sL https://vesio.azureedge.net/releases/vesctl/0.2.35/vesctl.linux-amd64.gz | gunzip > vesctl-0.2.35
chmod +x vesctl-0.2.35

# Verify it works
./vesctl-0.2.35 version

# Create symlink for test scripts
ln -sf vesctl-0.2.35 vesctl-original
```

### Version Compatibility Matrix

| Version | Status | Issue |
|---------|--------|-------|
| 0.2.35 | **WORKS** | Last working version |
| 0.2.36+ | BROKEN | TLS CertManager nil pointer crash |
| 0.2.47 | BROKEN | Same TLS bug as 0.2.36 |

### Legacy Download (Non-functional for API tests)

For non-authenticated tests only (help, completion, etc.):

```bash
# macOS arm64 (0.2.47 - has TLS bug, only for offline tests)
curl -L https://vesio.azureedge.net/releases/vesctl/0.2.47/vesctl.darwin-arm64.gz -o /tmp/vesctl.gz
gunzip /tmp/vesctl.gz
chmod +x /tmp/vesctl
mv /tmp/vesctl ./vesctl-0.2.47-original
```

## Running Tests

### Prerequisites

**For authenticated API tests (Phase 3+):**
1. Linux amd64 host (e.g., Ubuntu VM)
2. Original binary: `vesctl-0.2.35` (download from link above)
3. Our binary: `GOOS=linux GOARCH=amd64 go build -o vesctl-linux-amd64 .`

**For offline tests (Phase 1-2):**
1. Original binary: Any version (e.g., `./vesctl-0.2.47-original`)
2. Our binary: `go build -o vesctl .`

### Quick Start (Linux amd64)

```bash
# On your Linux test host
cd ~/GIT/robinmordasiewicz/vesctl

# Set environment variables
export ORIGINAL_VESCTL=./vesctl-0.2.35
export OUR_F5XCCTL=./f5xcctl-linux-amd64

# Create convenience symlinks
ln -sf vesctl-0.2.35 vesctl-original
ln -sf vesctl-linux-amd64 vesctl-ours

# Run namespace list comparison
./vesctl-original configuration list namespace
./vesctl-ours configuration list namespace
```

### Quick Start (Local macOS - offline tests only)

```bash
# Set environment variables (use local paths)
export ORIGINAL_VESCTL=./vesctl-0.2.47-original
export OUR_F5XCCTL=./f5xcctl.darwin-arm64

# Run all behavioral tests (no API calls)
./claudedocs/compatibility/tests/phase1-configure/test-behavior.sh

# Run version/completion tests (no API calls)
./claudedocs/compatibility/tests/phase2-simple/test-simple.sh
```

### Test Phases

| Phase | Description | Requires API |
|-------|-------------|--------------|
| Phase 1 | Behavioral consistency (help, flags, structure) | No |
| Phase 2 | Version, completion, configure commands | No |
| Phase 3 | Configuration CRUD operations | Yes |
| Phase 4 | Other commands (request, site, api-endpoint) | Yes |

## Reports

- `BEHAVIORAL_REPORT.md` - Phase 1 behavioral test results
- `PHASE2_REPORT.md` - Phase 2 simple command test results
- `COMPATIBILITY_REPORT.md` - Overall compatibility summary

## Test Results Summary

| Phase | Pass Rate | Status |
|-------|-----------|--------|
| Phase 1: Behavioral | 94% (34/36) | ✅ |
| Phase 2: Simple | 100% (7/7) | ✅ |

## Directory Structure

```
claudedocs/compatibility/
├── README.md                 # This file
├── BEHAVIORAL_REPORT.md      # Phase 1 results
├── PHASE2_REPORT.md          # Phase 2 results
├── COMPATIBILITY_REPORT.md   # Overall summary
├── lib/
│   ├── common.sh             # Shared test utilities
│   ├── compare.sh            # Comparison functions
│   └── report.sh             # Report generation
├── tests/
│   ├── phase1-configure/
│   │   ├── test-configure.sh # Content comparison tests
│   │   └── test-behavior.sh  # Behavioral consistency tests
│   ├── phase2-simple/
│   │   └── test-simple.sh    # Version/completion tests
│   └── results/              # Test output directory
└── run-all-tests.sh          # Master test runner
```

## Key Findings

### Identical Behavior
- Global flags (--config, --p12-bundle, --cert, --key, --cacert, --outfmt)
- Help text structure and formatting
- Command aliases (configuration, cfg, c)
- Completion scripts (bash and zsh)
- Error message patterns

### Expected Differences
- Resource list (API spec evolution: 22 deprecated, 151 new)
- Version output format (HTML vs YAML)

### Enhancements in Our Implementation
- `configure` command - Interactive CLI setup
- `login`/`logout`/`whoami` commands - Authentication workflow
