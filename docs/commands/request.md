# Request Commands

The `vesctl request` command group provides low-level API access.

## Overview

```text
vesctl request --help
```

## Subcommands

### RPC Request

Execute custom API RPC calls:

```bash
vesctl request rpc <method> <path> [flags]
```

**Examples:**

```bash
# GET request
vesctl request rpc GET /api/web/namespaces

# POST request with body
vesctl request rpc POST /api/some/endpoint -i request.json
```

### Secrets

Manage secrets:

```bash
vesctl request secrets [subcommand] [flags]
```

## Common Use Cases

### Custom API Calls

For API endpoints not covered by standard commands:

```bash
# List available APIs
vesctl request rpc GET /api

# Get tenant information
vesctl request rpc GET /api/web/tenant
```

### Working with Secrets

```bash
# Get secret help
vesctl request secrets --help
```

## Flags

| Flag | Description |
|------|-------------|
| `-i, --input-file` | Request body from file |
| `--outfmt` | Output format |
