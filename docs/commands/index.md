# Command Reference

All vesctl commands follow a consistent pattern:

```
vesctl <command-group> <action> [resource-type] [name] [flags]
```

## Command Groups

vesctl organizes commands into logical groups:

| Available Commands | |
|---------------|-------------|
| [api-endpoint](api-endpoint/index.md) | Execute commands for API Endpoint Discovery and Control |
| [completion](completion/index.md) | Generate completion script |
| [configuration](configuration/index.md) | Configure object |
| [help](help/index.md) | Help about any command |
| [request](request/index.md) | Execute Commands |
| [site](site/index.md) | Manage site creation through view public cloud site apis |
| [version](version/index.md) | Print build version |

## Flags

These flags are available for all commands:

| Short | Full | Type | Description |
|------|---|----|-------------|
| | `--api-token` | | Use API token from VES_API_TOKEN environment variable |
| `-a` | `--cacert` | string | Server CA cert file path |
| `-c` | `--cert` | string | Client cert file path |
| | `--config` | string | A configuration file to use for API gateway URL and credentials (default: ~/.vesconfig) |
| | `--hardwareKey` | Use yubikey for TLS connection |
| `-h` | `--help` | | help for vesctl |
| `-k` |`--key` | | Client key file path |
| | `--outfmt` | string | Output format for command |
| `-o` | `--output` | string | Output dir for command (default "./") |
| | `--p12-bundle` | string | Client P12 bundle (key+cert) file path. Any password for this file should be in environment variable VES_P12_PASSWORD |
| `-u` | `--server-urls` | strings | API endpoint URL (default [http://localhost:8001]) |
| | `--show-curl` | | Emit requests from program in CURL format |
| | `--timeout` | int | Timeout (in seconds) for command to finish (default 5) |

