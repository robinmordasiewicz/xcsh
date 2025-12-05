#!/bin/bash
# test.sh - Idempotent compatibility test runner for vesctl
#
# This script provides a single entry point for running all compatibility tests.
# It auto-detects the platform and sets appropriate binary paths.
#
# Usage:
#   ./test.sh                     # Run phases 1-3 (no API required)
#   ./test.sh --with-api          # Run all phases including namespace CRUD
#   ./test.sh --phase 1           # Run only phase 1
#   ./test.sh --phase 4 --with-api # Run only namespace CRUD tests
#   ./test.sh --help              # Show all options
#
# Environment Variables:
#   VES_API_URL        - F5 XC API URL (e.g., https://tenant.console.ves.volterra.io/api)
#   VES_API_P12_FILE   - Path to P12 certificate file
#   VES_P12_PASSWORD   - Password for P12 certificate

set -e

# Change to script directory (project root)
cd "$(dirname "$0")"

# Setup Go path if not already in PATH
if ! command -v go &> /dev/null; then
    if [[ -d "$HOME/go_install/go/bin" ]]; then
        export PATH="$HOME/go_install/go/bin:$PATH"
        export GOPATH="$HOME/go"
    elif [[ -d "/usr/local/go/bin" ]]; then
        export PATH="/usr/local/go/bin:$PATH"
    fi
fi

# Detect platform
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
[[ "$ARCH" == "x86_64" ]] && ARCH="amd64"
[[ "$ARCH" == "aarch64" ]] && ARCH="arm64"

echo "=========================================="
echo "vesctl Compatibility Test Runner"
echo "=========================================="
echo ""
echo "Platform: ${OS}-${ARCH}"
echo ""

# Get absolute project root path
PROJECT_ROOT="$(pwd)"

# Set binary paths based on platform (use absolute paths for consistency)
if [[ "$OS" == "linux" && "$ARCH" == "amd64" ]]; then
    # Linux amd64: use vesctl-0.2.35 (last working version)
    export ORIGINAL_VESCTL="${ORIGINAL_VESCTL:-${PROJECT_ROOT}/vesctl-0.2.35}"
    export OUR_VESCTL="${OUR_VESCTL:-${PROJECT_ROOT}/vesctl}"
else
    # Other platforms: use 0.2.47 (only for offline tests - has TLS bug)
    export ORIGINAL_VESCTL="${ORIGINAL_VESCTL:-${PROJECT_ROOT}/vesctl-0.2.47-original}"
    export OUR_VESCTL="${OUR_VESCTL:-${PROJECT_ROOT}/vesctl.${OS}-${ARCH}}"
    echo "WARNING: Only Linux amd64 has a working original binary for API tests"
    echo "         API tests (--with-api) will fail on this platform"
    echo ""
fi

# Build our binary if missing or source files are newer
build_if_needed() {
    if [[ ! -f "$OUR_VESCTL" ]]; then
        echo "Binary not found: $OUR_VESCTL"
        echo "Building vesctl..."
        go build -o "$OUR_VESCTL" .
        echo "Built: $OUR_VESCTL"
        echo ""
    elif [[ $(find . -name "*.go" -newer "$OUR_VESCTL" 2>/dev/null | head -1) ]]; then
        echo "Source files newer than binary, rebuilding..."
        go build -o "$OUR_VESCTL" .
        echo "Rebuilt: $OUR_VESCTL"
        echo ""
    fi
}

# Check for original binary
check_original() {
    if [[ ! -x "$ORIGINAL_VESCTL" ]]; then
        echo "ERROR: Original vesctl not found at: $ORIGINAL_VESCTL"
        echo ""
        echo "Download with:"
        if [[ "$OS" == "linux" && "$ARCH" == "amd64" ]]; then
            echo "  curl -sL https://vesio.azureedge.net/releases/vesctl/0.2.35/vesctl.linux-amd64.gz | gunzip > vesctl-0.2.35"
            echo "  chmod +x vesctl-0.2.35"
        else
            echo "  ./claudedocs/compatibility/scripts/download-original.sh"
        fi
        exit 1
    fi
}

# Build our binary if needed
build_if_needed

# Check original binary exists
check_original

echo "Original: $ORIGINAL_VESCTL"
echo "Ours:     $OUR_VESCTL"
echo ""

# Run the compatibility test suite
exec ./claudedocs/compatibility/run-all-tests.sh "$@"
