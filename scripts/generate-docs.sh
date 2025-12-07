#!/bin/bash
# generate-docs.sh - Generate command reference documentation from vesctl CLI
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
DOCS_DIR="${PROJECT_ROOT}/docs"

# Build vesctl if needed
if [[ ! -x "${PROJECT_ROOT}/vesctl" ]]; then
    echo "Building vesctl..."
    cd "$PROJECT_ROOT"
    go build -o vesctl .
fi

VESCTL="${PROJECT_ROOT}/vesctl"

echo "Generating command documentation..."

# Generate configuration command docs
cat > "${DOCS_DIR}/commands/configuration.md" << 'HEADER'
# Configuration Commands

The `vesctl configuration` command group manages F5 Distributed Cloud resources.

## Overview

```text
HEADER

"$VESCTL" configuration --help >> "${DOCS_DIR}/commands/configuration.md" 2>&1 || true

cat >> "${DOCS_DIR}/commands/configuration.md" << 'MIDDLE'
```

## Subcommands

### List Resources

List all resources of a specific type:

```text
MIDDLE

"$VESCTL" configuration list --help >> "${DOCS_DIR}/commands/configuration.md" 2>&1 || true

cat >> "${DOCS_DIR}/commands/configuration.md" << 'FOOTER'
```

### Get Resource

Get details of a specific resource:

```text
FOOTER

"$VESCTL" configuration get --help >> "${DOCS_DIR}/commands/configuration.md" 2>&1 || true

cat >> "${DOCS_DIR}/commands/configuration.md" << 'END'
```

### Create Resource

Create a new resource from a YAML/JSON file:

```text
END

"$VESCTL" configuration create --help >> "${DOCS_DIR}/commands/configuration.md" 2>&1 || true

cat >> "${DOCS_DIR}/commands/configuration.md" << 'END2'
```

### Replace Resource

Replace an existing resource:

```text
END2

"$VESCTL" configuration replace --help >> "${DOCS_DIR}/commands/configuration.md" 2>&1 || true

cat >> "${DOCS_DIR}/commands/configuration.md" << 'END3'
```

### Delete Resource

Delete a resource:

```text
END3

"$VESCTL" configuration delete --help >> "${DOCS_DIR}/commands/configuration.md" 2>&1 || true

echo '```' >> "${DOCS_DIR}/commands/configuration.md"

# Generate request command docs
cat > "${DOCS_DIR}/commands/request.md" << 'HEADER'
# Request Commands

The `vesctl request` command group provides low-level API access.

## Overview

```text
HEADER

"$VESCTL" request --help >> "${DOCS_DIR}/commands/request.md" 2>&1 || true

cat >> "${DOCS_DIR}/commands/request.md" << 'MIDDLE'
```

## Subcommands

### RPC Request

Execute custom API RPC calls:

```text
MIDDLE

"$VESCTL" request rpc --help >> "${DOCS_DIR}/commands/request.md" 2>&1 || true

cat >> "${DOCS_DIR}/commands/request.md" << 'FOOTER'
```

### Secrets

Manage secrets:

```text
FOOTER

"$VESCTL" request secrets --help >> "${DOCS_DIR}/commands/request.md" 2>&1 || true

echo '```' >> "${DOCS_DIR}/commands/request.md"

# Generate site command docs
cat > "${DOCS_DIR}/commands/site.md" << 'HEADER'
# Site Commands

The `vesctl site` command group manages cloud and edge sites.

## Overview

```text
HEADER

"$VESCTL" site --help >> "${DOCS_DIR}/commands/site.md" 2>&1 || true

echo '```' >> "${DOCS_DIR}/commands/site.md"

# Generate api-endpoint command docs
cat > "${DOCS_DIR}/commands/api-endpoint.md" << 'HEADER'
# API Endpoint Commands

The `vesctl api-endpoint` command group manages API endpoint discovery.

## Overview

```text
HEADER

"$VESCTL" api-endpoint --help >> "${DOCS_DIR}/commands/api-endpoint.md" 2>&1 || true

echo '```' >> "${DOCS_DIR}/commands/api-endpoint.md"

# Generate utilities docs
cat > "${DOCS_DIR}/commands/utilities.md" << 'HEADER'
# Utility Commands

Additional utility commands for vesctl.

## Version

Display version information:

```text
HEADER

"$VESCTL" version --help >> "${DOCS_DIR}/commands/utilities.md" 2>&1 || true

cat >> "${DOCS_DIR}/commands/utilities.md" << 'MIDDLE'
```

### Example Output

```text
MIDDLE

"$VESCTL" version >> "${DOCS_DIR}/commands/utilities.md" 2>&1 || true

cat >> "${DOCS_DIR}/commands/utilities.md" << 'FOOTER'
```

## Completion

Generate shell completion scripts:

```text
FOOTER

"$VESCTL" completion --help >> "${DOCS_DIR}/commands/utilities.md" 2>&1 || true

cat >> "${DOCS_DIR}/commands/utilities.md" << 'END'
```

## Configure

Interactive configuration setup:

```text
END

"$VESCTL" configure --help >> "${DOCS_DIR}/commands/utilities.md" 2>&1 || true

echo '```' >> "${DOCS_DIR}/commands/utilities.md"

# Generate resources list
cat > "${DOCS_DIR}/reference/resources.md" << 'HEADER'
# Supported Resources

vesctl supports the following F5 Distributed Cloud resource types.

## Resource Types

The following resources can be managed with `vesctl configuration` commands:

| Resource Type | Description |
|---------------|-------------|
HEADER

# Extract resource types from configuration list help
"$VESCTL" configuration list --help 2>&1 | grep -E '^\s+[a-z]' | while read -r line; do
    resource=$(echo "$line" | awk '{print $1}')
    if [[ -n "$resource" ]]; then
        echo "| \`$resource\` | F5 XC resource type |" >> "${DOCS_DIR}/reference/resources.md"
    fi
done

cat >> "${DOCS_DIR}/reference/resources.md" << 'FOOTER'

## Usage Examples

### List all namespaces

```bash
vesctl configuration list namespace
```

### Get a specific namespace

```bash
vesctl configuration get namespace my-namespace
```

### Create a resource from file

```bash
vesctl configuration create namespace -i namespace.yaml
```
FOOTER

echo "Documentation generated successfully!"
