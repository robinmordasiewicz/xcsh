#!/bin/bash
# Check for upstream API specification updates without downloading
# Source: https://github.com/robinmordasiewicz/f5xc-api-enriched/releases
#
# This script checks if a newer version of the enriched API specs is available
# and reports the status without performing a download.
#
# Exit codes:
#   0 - Up to date (or update available with --quiet)
#   1 - Error fetching release information
#   2 - Update available (only without --quiet)

set -e

SPECS_DIR=".specs"
ENRICHED_REPO="robinmordasiewicz/f5xc-api-enriched"
API_URL="https://api.github.com/repos/${ENRICHED_REPO}/releases/latest"

# Parse arguments
QUIET=false
JSON_OUTPUT=false
for arg in "$@"; do
  case $arg in
    --quiet | -q)
      QUIET=true
      ;;
    --json)
      JSON_OUTPUT=true
      ;;
  esac
done

# GitHub authentication (optional, but recommended for CI to avoid rate limits)
GITHUB_AUTH_HEADER=""
if [ -n "$GITHUB_TOKEN" ]; then
  GITHUB_AUTH_HEADER="Authorization: Bearer $GITHUB_TOKEN"
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
  if [ "$QUIET" = false ] && [ "$JSON_OUTPUT" = false ]; then
    echo -e "${GREEN}info${NC} $1"
  fi
}

log_warn() {
  if [ "$QUIET" = false ] && [ "$JSON_OUTPUT" = false ]; then
    echo -e "${YELLOW}warn${NC} $1"
  fi
}

log_error() {
  if [ "$JSON_OUTPUT" = false ]; then
    echo -e "${RED}error${NC} $1" >&2
  fi
}

output_json() {
  local current="$1"
  local latest="$2"
  local update_available="$3"
  local status="$4"

  cat <<EOF
{
  "current_version": "$current",
  "latest_version": "$latest",
  "update_available": $update_available,
  "status": "$status",
  "repository": "$ENRICHED_REPO"
}
EOF
}

# Build curl command with optional auth header
CURL_OPTS=(-s -L --max-time 30 --connect-timeout 10)
if [ -n "$GITHUB_AUTH_HEADER" ]; then
  CURL_OPTS+=(-H "$GITHUB_AUTH_HEADER")
fi

# Fetch latest release information
RELEASE_JSON=$(curl "${CURL_OPTS[@]}" "$API_URL" 2>/dev/null)

if [ -z "$RELEASE_JSON" ]; then
  log_error "Failed to fetch release information from GitHub API"
  if [ "$JSON_OUTPUT" = true ]; then
    output_json "" "" "false" "error"
  fi
  exit 1
fi

# Check for API errors
if echo "$RELEASE_JSON" | grep -q '"message"'; then
  ERROR_MSG=$(echo "$RELEASE_JSON" | grep -o '"message"[[:space:]]*:[[:space:]]*"[^"]*"' | head -1 | cut -d'"' -f4)
  log_error "GitHub API error: $ERROR_MSG"
  if [ "$JSON_OUTPUT" = true ]; then
    output_json "" "" "false" "error"
  fi
  exit 1
fi

LATEST_VERSION=$(echo "$RELEASE_JSON" | grep -o '"tag_name"[[:space:]]*:[[:space:]]*"[^"]*"' | head -1 | cut -d'"' -f4)
if [ -z "$LATEST_VERSION" ]; then
  log_error "Could not determine latest version from API response"
  if [ "$JSON_OUTPUT" = true ]; then
    output_json "" "" "false" "error"
  fi
  exit 1
fi

# Get current version
CURRENT_VERSION=""
if [ -f "$SPECS_DIR/.version" ]; then
  CURRENT_VERSION=$(cat "$SPECS_DIR/.version")
fi

# Determine status
if [ -z "$CURRENT_VERSION" ]; then
  if [ "$JSON_OUTPUT" = true ]; then
    output_json "" "$LATEST_VERSION" "true" "not_installed"
  else
    log_warn "No specs installed"
    log_info "Latest available: $LATEST_VERSION"
    log_info "Run 'make download-specs' to download"
  fi
  if [ "$QUIET" = false ]; then
    exit 2
  fi
  exit 0
elif [ "$CURRENT_VERSION" = "$LATEST_VERSION" ]; then
  if [ "$JSON_OUTPUT" = true ]; then
    output_json "$CURRENT_VERSION" "$LATEST_VERSION" "false" "up_to_date"
  else
    log_info "Specs are up to date ($CURRENT_VERSION)"
  fi
  exit 0
else
  if [ "$JSON_OUTPUT" = true ]; then
    output_json "$CURRENT_VERSION" "$LATEST_VERSION" "true" "update_available"
  else
    log_warn "Update available!"
    log_info "Current: $CURRENT_VERSION"
    log_info "Latest:  $LATEST_VERSION"
    log_info "Run 'make download-specs' to update"
  fi
  if [ "$QUIET" = false ]; then
    exit 2
  fi
  exit 0
fi
