# vesctl Documentation

**vesctl** is an open-source command-line interface for managing F5 Distributed Cloud (formerly Volterra) resources.

## Features

- **Full API Coverage** - Manage all F5 XC resources from the command line
- **Multiple Output Formats** - JSON, YAML, table, and TSV output formats
- **Shell Completion** - Auto-completion for Bash, Zsh, Fish, and PowerShell
- **Secure Authentication** - P12 bundle or certificate/key pair authentication
- **Cross-Platform** - Available for Linux, macOS, and Windows

## Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | sh
```

## Basic Usage

```bash
# Configure authentication
vesctl configure

# List namespaces
vesctl configuration list namespace

# Get resource details
vesctl configuration get namespace my-namespace --outfmt json

# Create a resource from file
vesctl configuration create http_loadbalancer -i lb.yaml
```

## Command Structure

vesctl commands follow a consistent pattern:

```
vesctl <command-group> <action> [resource-type] [name] [flags]
```

### Command Groups

| Group | Description |
|-------|-------------|
| `configuration` | Manage F5 XC resources (CRUD operations) |
| `request` | Low-level API requests |
| `site` | Manage cloud and edge sites |
| `api-endpoint` | API endpoint discovery |
| `configure` | Interactive configuration setup |
| `version` | Display version information |
| `completion` | Generate shell completions |

## Getting Help

```bash
# General help
vesctl --help

# Command-specific help
vesctl configuration --help
vesctl configuration list --help
```

## Next Steps

- [Installation Guide](getting-started/installation.md) - Detailed installation instructions
- [Quick Start](getting-started/quickstart.md) - Get up and running quickly
- [Authentication](getting-started/authentication.md) - Configure API credentials
- [Command Reference](commands/index.md) - Complete command documentation
