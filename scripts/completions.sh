#!/bin/sh
# Generate shell completion files for GoReleaser builds
set -e

# Clean and create completions directory
rm -rf completions
mkdir -p completions

# Build binary first to ensure it exists
# go run . can have issues in CI environments
echo "Building vesctl for completion generation..."
go build -o /tmp/vesctl-completions .

# Generate completions for all supported shells using built binary
for sh in bash zsh fish; do
    echo "Generating ${sh} completions..."
    /tmp/vesctl-completions completion "${sh}" > "completions/vesctl.${sh}"
done

# Cleanup temp binary
rm -f /tmp/vesctl-completions

echo "Shell completions generated in completions/"
ls -la completions/
