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
echo "Building xcsh for completion generation..."
go build -o ./xcsh-completions .

if [ ! -f "./xcsh-completions" ]; then
    echo "ERROR: Failed to build xcsh-completions binary"
    exit 1
fi

echo "Binary built successfully"

# Generate completions for all supported shells using built binary
for sh in bash zsh fish; do
    echo "Generating ${sh} completions..."
    ./xcsh-completions completion "${sh}" > "completions/xcsh.${sh}"

    if [ ! -s "completions/xcsh.${sh}" ]; then
        echo "ERROR: Failed to generate ${sh} completions (file empty or missing)"
        exit 1
    fi
done

# Cleanup temp binary
rm -f ./xcsh-completions

echo "=== Completion generation finished ==="
echo "Generated files:"
ls -la completions/

# Verify files exist and have content
for sh in bash zsh fish; do
    if [ ! -s "completions/xcsh.${sh}" ]; then
        echo "ERROR: completions/xcsh.${sh} is missing or empty"
        exit 1
    fi
    echo "  completions/xcsh.${sh}: $(wc -c < completions/xcsh.${sh}) bytes"
done

echo "All completions generated successfully"
