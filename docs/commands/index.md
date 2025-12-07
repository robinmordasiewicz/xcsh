# Command Reference

This section provides detailed documentation for all vesctl commands.

## Command Groups

vesctl organizes commands into logical groups:

| Command Group | Description |
|---------------|-------------|
| [configuration](configuration.md) | Manage F5 XC resources (CRUD operations) |
| [request](request.md) | Low-level API requests and RPC calls |
| [site](site.md) | Manage cloud and edge sites |
| [api-endpoint](api-endpoint.md) | API endpoint discovery |
| [utilities](utilities.md) | Version, completion, and configuration commands |

## Command Structure

All vesctl commands follow a consistent pattern:

```
vesctl <command-group> <action> [resource-type] [name] [flags]
```

### Examples

```bash
# List all namespaces
vesctl configuration list namespace

# Get a specific resource
vesctl configuration get http_loadbalancer my-lb -n my-namespace

# Create from file
vesctl configuration create origin_pool -i pool.yaml

# Delete a resource
vesctl configuration delete healthcheck my-hc -n my-namespace
```

## Global Flags

These flags are available for all commands:

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | | Path to config file (default: ~/.vesconfig) |
| `--server-urls` | `-u` | API server URL(s) |
| `--cert` | `-c` | Path to client certificate file |
| `--key` | `-k` | Path to client key file |
| `--cacert` | | Path to CA certificate file |
| `--p12-bundle` | | Path to P12 certificate bundle |
| `--outfmt` | `-o` | Output format: json, yaml, table |
| `--verbose` | `-v` | Enable verbose output |
| `--debug` | | Enable debug logging |

## Getting Help

Every command supports the `--help` flag:

```bash
# Top-level help
vesctl --help

# Command group help
vesctl configuration --help

# Specific action help
vesctl configuration list --help
vesctl configuration get --help
```

## Output Formats

Control output format with `--outfmt`:

```bash
# Table format (default)
vesctl configuration list namespace

# JSON output
vesctl configuration list namespace --outfmt json

# YAML output
vesctl configuration list namespace --outfmt yaml
```

## Command Index

### Configuration Commands

- `vesctl configuration list <type>` - List resources
- `vesctl configuration get <type> <name>` - Get resource details
- `vesctl configuration create <type>` - Create a resource
- `vesctl configuration replace <type>` - Replace a resource
- `vesctl configuration delete <type> <name>` - Delete a resource

### Request Commands

- `vesctl request rpc` - Execute custom RPC calls
- `vesctl request secrets` - Manage secrets

### Site Commands

- `vesctl site` - Manage cloud and edge sites

### Utility Commands

- `vesctl version` - Show version information
- `vesctl completion` - Generate shell completions
- `vesctl configure` - Interactive configuration setup
