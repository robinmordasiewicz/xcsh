#!/bin/sh
# Generate shell completion files for GoReleaser builds
set -e

# Clean and create completions directory
rm -rf completions
mkdir -p completions

# Generate completions for all supported shells
for sh in bash zsh fish; do
    echo "Generating ${sh} completions..."
    go run . completion "${sh}" > "completions/vesctl.${sh}"
done

echo "Shell completions generated in completions/"
ls -la completions/
