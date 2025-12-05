# vesctl Compatibility Testing

This directory contains the compatibility testing framework for comparing our vesctl implementation against the original F5 proprietary binary.

## Original Binary

The original proprietary vesctl binary is stored locally for comparison testing:
- **Location**: `./vesctl-0.2.47-original` (in project root)
- **Version**: 0.2.47
- **Source**: https://vesio.azureedge.net/releases/vesctl/0.2.47/vesctl.darwin-arm64.gz

> **Note**: This file is gitignored (`vesctl*` pattern) and must be downloaded for testing.

### Download Original Binary

```bash
# Using the provided script
./claudedocs/compatibility/scripts/download-original.sh

# Or manually
curl -L https://vesio.azureedge.net/releases/vesctl/0.2.47/vesctl.darwin-arm64.gz -o /tmp/vesctl.gz
gunzip /tmp/vesctl.gz
chmod +x /tmp/vesctl
mv /tmp/vesctl ./vesctl-0.2.47-original

# Verify
./vesctl-0.2.47-original version
```

## Running Tests

### Prerequisites

1. Original binary downloaded: `./vesctl-0.2.47-original`
2. Our binary built: `go build -o vesctl.darwin-arm64 .`

### Quick Start

```bash
# Set environment variables (use local paths)
export ORIGINAL_VESCTL=./vesctl-0.2.47-original
export OUR_VESCTL=./vesctl.darwin-arm64

# Run all behavioral tests
./claudedocs/compatibility/tests/phase1-configure/test-behavior.sh

# Run version/completion tests
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
