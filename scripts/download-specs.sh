#!/bin/bash
# Download enriched API specifications from GitHub releases
# Source: https://github.com/robinmordasiewicz/f5xc-api-enriched/releases
#
# This script downloads the latest enriched API specs and extracts them
# to .specs/domains/ for use during the build process.
#
# The downloaded specs include all domain resources with:
# - Complete OpenAPI 3.0 specifications
# - Example values using 'example-' prefix (standard naming convention)
# - Consolidated openapi.json/yaml for reference
#
# Features:
# - Exponential backoff retry logic for transient network errors
# - Automatic fallback on permanent errors
# - Detailed logging for debugging
# - Version caching in .specs/.version to track upstream releases

set -e

SPECS_DIR=".specs"
ENRICHED_REPO="robinmordasiewicz/f5xc-api-enriched"
API_URL="https://api.github.com/repos/${ENRICHED_REPO}/releases/latest"

# Retry configuration
MAX_RETRIES=5
INITIAL_BACKOFF=2
MAX_BACKOFF=60

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}ℹ${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

log_error() {
    echo -e "${RED}✗${NC} $1"
}

log_success() {
    echo -e "${GREEN}✓${NC} $1"
}

log_debug() {
    echo -e "${BLUE}⋯${NC} $1"
}

# Determine if an error is transient (worth retrying)
is_transient_error() {
    local error="$1"
    # Network timeouts, connection errors, temporary DNS issues
    if echo "$error" | grep -qiE "(timeout|connection|refused|temporarily|unavailable|timed out|resolving|name or service)"; then
        return 0
    fi
    # curl exit codes for transient errors
    # 7: connection failed, 28: operation timeout, 35: SSL connect error
    return 1
}

# Retry function with exponential backoff
retry_with_backoff() {
    local name="$1"
    shift
    local attempt=1
    local backoff=$INITIAL_BACKOFF

    while [ $attempt -le $MAX_RETRIES ]; do
        log_debug "[$name] Attempt $attempt/$MAX_RETRIES"

        # Execute the command and capture output and exit code
        if output=$("$@" 2>&1); then
            echo "$output"
            return 0
        fi

        local exit_code=$?
        local error="$output"

        if [ $attempt -eq $MAX_RETRIES ]; then
            # Last attempt - show full error
            log_error "[$name] Failed after $MAX_RETRIES attempts"
            log_error "Exit code: $exit_code"
            log_error "Error: $error"
            return $exit_code
        fi

        # Check if error is transient
        if is_transient_error "$error"; then
            log_warn "[$name] Transient error (attempt $attempt/$MAX_RETRIES): $error"
            log_warn "[$name] Retrying in ${backoff}s..."
            sleep "$backoff"

            # Exponential backoff: double the wait time, cap at MAX_BACKOFF
            backoff=$((backoff * 2))
            if [ $backoff -gt $MAX_BACKOFF ]; then
                backoff=$MAX_BACKOFF
            fi
        else
            # Permanent error - don't retry
            log_error "[$name] Permanent error (not retrying): $error"
            return $exit_code
        fi

        attempt=$((attempt + 1))
    done

    return 1
}

# Helper function to fetch with retry
fetch_with_retry() {
    local name="$1"
    local url="$2"
    local output_file="$3"

    if [ -n "$output_file" ]; then
        retry_with_backoff "$name" curl -f -L --max-time 30 --connect-timeout 10 "$url" -o "$output_file"
    else
        retry_with_backoff "$name" curl -f -sL --max-time 30 --connect-timeout 10 "$url"
    fi
}

# Fetch latest release information
log_info "Fetching latest enriched spec release..."
RELEASE_JSON=$(fetch_with_retry "GitHub API" "$API_URL")

if [ -z "$RELEASE_JSON" ]; then
    log_error "Failed to fetch release information from GitHub API"
    exit 1
fi

VERSION=$(echo "$RELEASE_JSON" | grep -o '"tag_name"[[:space:]]*:[[:space:]]*"[^"]*"' | head -1 | cut -d'"' -f4)
if [ -z "$VERSION" ]; then
    log_error "Could not determine latest version from API response"
    log_debug "API Response: $RELEASE_JSON"
    exit 1
fi

log_info "Latest version: $VERSION"

# Check if already downloaded AND required files exist
# Note: .specs/index.json and .specs/domains/ are gitignored, so they won't
# exist after a fresh git clone even if .specs/.version exists
if [ -f "$SPECS_DIR/.version" ]; then
    CURRENT_VERSION=$(cat "$SPECS_DIR/.version")
    if [ "$CURRENT_VERSION" = "$VERSION" ] && [ -f "$SPECS_DIR/index.json" ] && [ -d "$SPECS_DIR/domains" ]; then
        log_success "Specs already up to date ($VERSION)"
        exit 0
    fi
    if [ "$CURRENT_VERSION" != "$VERSION" ]; then
        log_info "Updating from $CURRENT_VERSION to $VERSION"
    else
        log_info "Re-downloading specs (required files missing)"
    fi
fi

# Find download URLs
ZIP_URL=$(echo "$RELEASE_JSON" | grep -o '"browser_download_url"[[:space:]]*:[[:space:]]*"[^"]*\.zip"' | head -1 | cut -d'"' -f4)
INDEX_URL=$(echo "$RELEASE_JSON" | grep -o '"browser_download_url"[[:space:]]*:[[:space:]]*"[^"]*index\.json"' | head -1 | cut -d'"' -f4)

if [ -z "$ZIP_URL" ]; then
    log_error "Could not find specs ZIP file in release assets"
    exit 1
fi

# Create specs directory
mkdir -p "$SPECS_DIR"

# Download ZIP file with retry logic
log_info "Downloading specs from: $ZIP_URL"
fetch_with_retry "ZIP Download" "$ZIP_URL" "$SPECS_DIR/specs.zip"

# Verify ZIP file is valid
log_debug "Verifying ZIP file integrity..."
if ! unzip -t "$SPECS_DIR/specs.zip" > /dev/null 2>&1; then
    log_error "Downloaded ZIP file is corrupted"
    rm -f "$SPECS_DIR/specs.zip"
    exit 1
fi
log_success "ZIP file verified"

# Clean up old domain specs to prevent stale files from accumulating
# when upstream reorganizes or removes domains
if [ -d "$SPECS_DIR/domains" ]; then
    log_info "Cleaning up old domain specs..."
    rm -rf "$SPECS_DIR/domains"
fi

# Extract domain specs
log_info "Extracting domain specifications..."
unzip -q -o "$SPECS_DIR/specs.zip" -d "$SPECS_DIR"
rm -f "$SPECS_DIR/specs.zip"

# Download index if available
if [ -n "$INDEX_URL" ]; then
    log_info "Downloading index metadata..."
    fetch_with_retry "Index Download" "$INDEX_URL" "$SPECS_DIR/index.json"
fi

# Record version
echo "$VERSION" > "$SPECS_DIR/.version"

# Verify domains directory
if [ ! -d "$SPECS_DIR/domains" ]; then
    log_error "Domains directory not found in downloaded specs"
    exit 1
fi

# Count domain files
DOMAIN_COUNT=$(find "$SPECS_DIR/domains" -name "*.json" -type f | wc -l)
if [ "$DOMAIN_COUNT" -eq 0 ]; then
    log_error "No domain spec files found"
    exit 1
fi

log_success "Downloaded enriched specs $VERSION"
log_info "Location: $SPECS_DIR/domains/"
log_info "Domain files: $DOMAIN_COUNT"
