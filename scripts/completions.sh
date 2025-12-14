#!/bin/sh
# Generate shell completion files for GoReleaser builds
set -e

# Output to stderr so GoReleaser captures it
exec 1>&2

echo "=== Starting completion generation ==="
echo "Working directory: $(pwd)"

# Clean and create completions directory
rm -rf completions
mkdir -p completions

# Build binary first to ensure it exists
# go run . can have issues in CI environments
echo "Building f5xcctl for completion generation..."
go build -o ./f5xcctl-completions .

if [ ! -f "./f5xcctl-completions" ]; then
    echo "ERROR: Failed to build f5xcctl-completions binary"
    exit 1
fi

echo "Binary built successfully"

# Generate completions for all supported shells using built binary
for sh in bash zsh fish; do
    echo "Generating ${sh} completions..."
    ./f5xcctl-completions completion "${sh}" > "completions/f5xcctl.${sh}"

    if [ ! -s "completions/f5xcctl.${sh}" ]; then
        echo "ERROR: Failed to generate ${sh} completions (file empty or missing)"
        exit 1
    fi
done

# Cleanup temp binary
rm -f ./f5xcctl-completions

echo "=== Completion generation finished ==="
echo "Generated files:"
ls -la completions/

# Verify files exist and have content
for sh in bash zsh fish; do
    if [ ! -s "completions/f5xcctl.${sh}" ]; then
        echo "ERROR: completions/f5xcctl.${sh} is missing or empty"
        exit 1
    fi
    echo "  completions/f5xcctl.${sh}: $(wc -c < completions/f5xcctl.${sh}) bytes"
done

echo "All completions generated successfully"
