#!/bin/bash
# One-time cleanup of old v6.x.x releases
# This script deletes all GitHub releases and tags matching v6.* pattern

set -e

echo "================================================================"
echo "  GitHub Release Cleanup Script"
echo "  This will delete all old releases (v0.x through v6.x)"
echo "================================================================"
echo ""

# Check if gh CLI is available
if ! command -v gh &>/dev/null; then
  echo "Error: gh CLI is not installed. Install it from https://cli.github.com/"
  exit 1
fi

# Check if authenticated
if ! gh auth status &>/dev/null; then
  echo "Error: Not authenticated with gh CLI. Run 'gh auth login' first."
  exit 1
fi

# List all releases matching v6.*
echo "Fetching releases..."
# Match old version patterns (v0.x through v6.x)
RELEASES=$(gh release list --limit 200 --json tagName -q '.[].tagName' | grep -E '^v[0-6]\.' || true)

if [ -z "$RELEASES" ]; then
  echo "No old releases (v0.x through v6.x) found. Nothing to clean up."
  exit 0
fi

# Count releases
RELEASE_COUNT=$(echo "$RELEASES" | wc -l | tr -d ' ')

echo ""
echo "Found $RELEASE_COUNT releases to delete:"
echo "----------------------------------------"
echo "$RELEASES"
echo "----------------------------------------"
echo ""
echo "Press Ctrl+C to cancel, or Enter to continue..."
read

echo ""
echo "Starting cleanup..."
echo ""

# Delete each release and its tag
for TAG in $RELEASES; do
  echo "Deleting release and tag: $TAG"
  if gh release delete "$TAG" --yes --cleanup-tag 2>/dev/null; then
    echo "  ✓ Deleted $TAG"
  else
    echo "  ✗ Failed to delete $TAG (may already be deleted)"
  fi
done

echo ""
echo "================================================================"
echo "  Cleanup complete!"
echo "  Deleted $RELEASE_COUNT releases"
echo "================================================================"
