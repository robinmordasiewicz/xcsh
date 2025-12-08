# Global Flags

These flags are available for all vesctl commands.

## Connection Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--config` | | Path to config file | `~/.vesconfig` |
| `--server-urls` | `-u` | API server URL(s) | From config |
| `--cert` | `-c` | Path to client certificate | From config |
| `--key` | `-k` | Path to client key | From config |
| `--cacert` | | Path to CA certificate | System CA |
| `--p12-bundle` | | Path to P12 certificate bundle | From config |

## Output Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--outfmt` | `-o` | Output format (json, yaml, table) | `table` |
| `--query` | | JMESPath query for filtering | |

## Logging Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--verbose` | `-v` | Increase logging verbosity | `false` |
| `--debug` | | Enable debug logging | `false` |
| `--only-show-errors` | | Only show errors | `false` |

## Operation Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--no-wait` | Don't wait for long-running operations | `false` |
| `--yes` | Skip confirmation prompts | `false` |

## Examples

### Specify Config File

```bash
vesctl --config /path/to/custom-config.yaml configuration list namespace
```

### Override Server URL

```bash
vesctl --server-urls https://custom.api.url/api configuration list namespace
```

### Use Custom Certificates

```bash
vesctl --cert /path/to/cert.pem --key /path/to/key.pem configuration list namespace
```

### Enable Debug Output

```bash
vesctl --debug configuration list namespace
```

### JSON Output

```bash
vesctl --outfmt json configuration list namespace
```

### Combine Flags

```bash
vesctl --verbose --outfmt yaml configuration get namespace example-ns
```

## Precedence

Flags override configuration in this order (highest to lowest):

1. Command-line flags
2. Environment variables
3. Configuration file (`~/.vesconfig`)
4. Default values
