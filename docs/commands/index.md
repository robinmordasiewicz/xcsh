---
title: "Command Reference"
description: "Complete reference for all vesctl CLI commands"
keywords:
  - vesctl
  - CLI
  - F5 XC
  - F5 Distributed Cloud
  - command reference
---

# Command Reference

All vesctl commands follow a consistent pattern:

```bash
vesctl <command-group> <action> [resource-type] [name] [flags]
```

## Command Groups

vesctl organizes commands into logical groups:

| Command | Description |
|---------|-------------|
| [api-endpoint](api-endpoint/index.md) | Discover and manage API endpoints within F5 XC service mesh. |
| [completion](completion/index.md) | Generate shell completion scripts for bash or zsh. |
| [configuration](configuration/index.md) | Manage F5 XC configuration objects using CRUD operations. |
| [help](help/index.md) | Help about any command |
| [request](request/index.md) | Execute custom API requests to F5 Distributed Cloud. |
| [site](site/index.md) | Deploy and manage F5 XC sites on public cloud providers. |
| [version](version/index.md) | Display vesctl version and build information |

## Global Flags

These flags are available for all commands:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--api-token` |  | bool | false | Authenticate using the API token from VES_API_TOKEN environment variable. |
| `--cacert` | `-a` | string |  | Path to the server CA certificate file for TLS verification. |
| `--cert` | `-c` | string |  | Path to the client certificate file for mTLS authentication. |
| `--config` |  | string |  | Path to the configuration file containing API URL and credentials (default "/Users/r.mordasiewicz/.vesconfig"). |
| `--hardwareKey` |  | bool | false | Use a YubiKey hardware security module for TLS authentication. |
| `--key` | `-k` | string |  | Path to the client private key file for mTLS authentication. |
| `--non-interactive` |  | bool | false | Disable interactive prompts and fail if required arguments are missing. |
| `--output` | `-o` | string | ./ | Directory path for command output files. |
| `--output-format` |  | string |  | Set the output format to text, json, yaml, or table. |
| `--p12-bundle` |  | string |  | Path to PKCS#12 bundle file containing client certificate and key. Set password in VES_P12_PASSWORD. |
| `--server-url` | `-u` | string |  | F5 Distributed Cloud API endpoint URL. |
| `--show-curl` |  | bool | false | Output equivalent curl commands for each API request. |
| `--spec` |  | bool | false | Output machine-readable CLI specification (JSON/YAML) |
| `--timeout` |  | int | 5 | Maximum time in seconds to wait for command completion. |

## Quick Start

### List resources

```bash
vesctl configuration list namespace
```

### Get a specific resource

```bash
vesctl configuration get namespace example-namespace
```

### Create from file

```bash
vesctl configuration create http_loadbalancer -i lb.yaml
```

### Delete a resource

```bash
vesctl configuration delete http_loadbalancer example-lb -n example-namespace
```

## Output Formats

vesctl supports multiple output formats:

| Format | Flag | Description |
|--------|------|-------------|
| Text | `--output-format text` | Human-readable text output (default) |
| JSON | `--output-format json` | JSON format for scripting |
| YAML | `--output-format yaml` | YAML format for configuration |
| Table | `--output-format table` | Tabular format for lists |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `VES_API_TOKEN` | API token for authentication (use with `--api-token` flag) |
| `VES_P12_PASSWORD` | Password for P12 bundle file |

## Version

Built from version: `v4.15.2-3-g3a4e3ba`