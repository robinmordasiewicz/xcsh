#!/usr/bin/env bash
# Build standalone binaries using Bun
# ink v5+ requires ESM with top-level await, which pkg doesn't support
# Bun's compile feature handles ESM natively
#
# IMPORTANT: Bun cross-compilation does NOT work reliably!
# macOS binaries must be built on macOS, Linux on Linux, etc.
# This script builds only for the current platform by default.
# Use --platform to specify which platforms to build.
set -eo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DIST_DIR="$PROJECT_ROOT/binaries"
ENTRY="$PROJECT_ROOT/dist/index.js"

# Parse arguments
PLATFORMS=""
while [[ $# -gt 0 ]]; do
  case $1 in
    --platform)
      PLATFORMS="$2"
      shift 2
      ;;
    --all)
      PLATFORMS="all"
      shift
      ;;
    *)
      echo "Unknown argument: $1"
      echo "Usage: $0 [--platform linux|darwin|windows|all] [--all]"
      exit 1
      ;;
  esac
done

# Detect current platform if not specified
if [[ -z "$PLATFORMS" ]]; then
  case "$(uname -s)" in
    Linux*) PLATFORMS="linux" ;;
    Darwin*) PLATFORMS="darwin" ;;
    MINGW* | MSYS* | CYGWIN*) PLATFORMS="windows" ;;
    *)
      echo "Unknown platform: $(uname -s)"
      exit 1
      ;;
  esac
  echo "Auto-detected platform: $PLATFORMS"
fi

# Check for Bun
if ! command -v bun &>/dev/null; then
  if [[ -x "$HOME/.bun/bin/bun" ]]; then
    BUN="$HOME/.bun/bin/bun"
  else
    echo "Error: Bun is not installed. Install with: curl -fsSL https://bun.sh/install | bash"
    exit 1
  fi
else
  BUN="bun"
fi

echo "Using Bun: $BUN ($($BUN --version))"

# Ensure dist/index.js exists
if [[ ! -f "$ENTRY" ]]; then
  echo "Error: dist/index.js not found. Run 'npm run build' first."
  exit 1
fi

# Create output directory
mkdir -p "$DIST_DIR"

echo ""
echo "Building binaries for: $PLATFORMS"
echo ""

# Build for a specific target
build_target() {
  local target="$1"
  local output="$2"
  echo "Building: $output (target: bun-$target)"
  if $BUN build "$ENTRY" --compile --target "bun-$target" --outfile "$DIST_DIR/$output" 2>&1; then
    echo "  ✓ Built successfully"
    return 0
  else
    echo "  ✗ Failed to build"
    return 1
  fi
}

# Build based on platform selection
build_linux() {
  echo "=== Building Linux binaries ==="
  build_target "linux-x64" "xcsh-linux-x64"
  build_target "linux-arm64" "xcsh-linux-arm64"
}

build_darwin() {
  echo "=== Building macOS binaries ==="
  build_target "darwin-x64" "xcsh-macos-x64"
  build_target "darwin-arm64" "xcsh-macos-arm64"
}

build_windows() {
  echo "=== Building Windows binaries ==="
  build_target "windows-x64" "xcsh-win-x64.exe"
}

# Execute builds
case "$PLATFORMS" in
  all)
    build_linux
    build_darwin
    build_windows
    ;;
  linux)
    build_linux
    ;;
  darwin)
    build_darwin
    ;;
  windows)
    build_windows
    ;;
  *)
    echo "Unknown platform: $PLATFORMS"
    exit 1
    ;;
esac

echo ""
echo "Build complete. Binaries:"
ls -lh "$DIST_DIR"/ 2>/dev/null || echo "  (no binaries built)"

# Generate checksums
echo ""
echo "Checksums:"
cd "$DIST_DIR"
if command -v sha256sum &>/dev/null; then
  sha256sum * 2>/dev/null || true
elif command -v shasum &>/dev/null; then
  shasum -a 256 * 2>/dev/null || true
else
  echo "  (install shasum or sha256sum for checksums)"
fi
