# Configuration Commands

The `vesctl configuration` command group manages F5 Distributed Cloud resources.

## Overview

```text
vesctl configuration --help
```

!!! note "Auto-generated"
    This page is automatically updated when the documentation is built.
    The content below shows the command help output.

## Subcommands

### List Resources

List all resources of a specific type:

```bash
vesctl configuration list <resource-type> [flags]
```

**Examples:**

```bash
# List all namespaces
vesctl configuration list namespace

# List HTTP load balancers in a namespace
vesctl configuration list http_loadbalancer -n my-namespace

# List with JSON output
vesctl configuration list namespace --outfmt json
```

### Get Resource

Get details of a specific resource:

```bash
vesctl configuration get <resource-type> <name> [flags]
```

**Examples:**

```bash
# Get namespace details
vesctl configuration get namespace my-namespace

# Get with YAML output
vesctl configuration get http_loadbalancer my-lb -n my-namespace --outfmt yaml
```

### Create Resource

Create a new resource from a YAML/JSON file:

```bash
vesctl configuration create <resource-type> -i <file> [flags]
```

**Examples:**

```bash
# Create from YAML file
vesctl configuration create http_loadbalancer -i lb.yaml

# Create with namespace
vesctl configuration create origin_pool -i pool.yaml -n my-namespace
```

### Replace Resource

Replace an existing resource:

```bash
vesctl configuration replace <resource-type> -i <file> [flags]
```

**Examples:**

```bash
# Replace from file
vesctl configuration replace http_loadbalancer -i lb-updated.yaml
```

### Delete Resource

Delete a resource:

```bash
vesctl configuration delete <resource-type> <name> [flags]
```

**Examples:**

```bash
# Delete a resource
vesctl configuration delete http_loadbalancer my-lb -n my-namespace

# Delete with confirmation bypass
vesctl configuration delete origin_pool my-pool -n my-namespace --yes
```

## Common Flags

| Flag | Short | Description |
|------|-------|-------------|
| `-n, --namespace` | `-n` | Namespace for the resource |
| `-i, --input-file` | `-i` | Input file (YAML or JSON) |
| `--outfmt` | `-o` | Output format (json, yaml, table) |
