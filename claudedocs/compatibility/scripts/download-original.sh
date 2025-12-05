#!/bin/bash
# download-original.sh - Download the original vesctl binary for comparison testing

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"

VESCTL_VERSION="${VESCTL_VERSION:-0.2.47}"
VESCTL_ARCH="${VESCTL_ARCH:-darwin-arm64}"
VESCTL_URL="https://vesio.azureedge.net/releases/vesctl/${VESCTL_VERSION}/vesctl.${VESCTL_ARCH}.gz"

# Default install path is project root with version in filename
INSTALL_PATH="${INSTALL_PATH:-${PROJECT_ROOT}/vesctl-${VESCTL_VERSION}-original}"

echo "Downloading original vesctl v${VESCTL_VERSION} for ${VESCTL_ARCH}..."
echo "URL: ${VESCTL_URL}"
echo "Target: ${INSTALL_PATH}"
echo ""

# Check if already exists
if [[ -x "$INSTALL_PATH" ]]; then
    echo "Original binary already exists at ${INSTALL_PATH}"
    echo "To re-download, remove it first: rm ${INSTALL_PATH}"
    exit 0
fi

# Create temp directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Download
echo "Downloading..."
curl -L "${VESCTL_URL}" -o "${TEMP_DIR}/vesctl.gz"

# Decompress
echo "Decompressing..."
gunzip "${TEMP_DIR}/vesctl.gz"

# Make executable
chmod +x "${TEMP_DIR}/vesctl"

# Verify it works
echo "Verifying..."
"${TEMP_DIR}/vesctl" version || echo "(Version check completed)"

# Install to project root (no sudo needed)
echo ""
echo "Installing to ${INSTALL_PATH}..."
mv "${TEMP_DIR}/vesctl" "${INSTALL_PATH}"

echo ""
echo "Done! Original vesctl installed at: ${INSTALL_PATH}"
echo ""
echo "This file is gitignored (vesctl* pattern)."
echo ""
echo "Run tests with:"
echo "  export ORIGINAL_VESCTL=${INSTALL_PATH}"
echo "  export OUR_VESCTL=${PROJECT_ROOT}/vesctl.darwin-arm64"
echo "  ./claudedocs/compatibility/run-all-tests.sh"
