---
title: "vesctl api-endpoint"
description: "Discover and manage API endpoints within F5 XC service mesh."
keywords:
  - api-endpoint
  - vesctl
  - F5 Distributed Cloud
  - F5 XC
command: "vesctl api-endpoint"
command_group: "api-endpoint"
---

# vesctl api-endpoint

> Discover and manage API endpoints within F5 XC service mesh.

## Synopsis

```bash
vesctl api-endpoint <command> [flags]
```

## Available Commands

| Command | Description |
|---------|-------------|
| [control](control/index.md) | Create layer 7 policies based on discovered API endpoints. |
| [discover](discover/index.md) | Discover API endpoints between services in a service mesh. |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--app-type` |  | string |  | App type name labeled on vK8s services or HTTP load balancer objects. |
| `--log-color` |  | bool | true | Enable colored log output. |
| `--log-fabulous` |  | bool | true | Enable enhanced log formatting. |
| `--log-level` |  | int | 3 | Set the logging verbosity level (1-5). |
| `--range` | `-r` | string | 1h | Time range for querying service mesh data (e.g., '1h', '24h'). |

## Examples

```bash
vesctl api-endpoint discover --namespace default
```

## See Also

- [Command Reference](../index.md)
- [vesctl api-endpoint control](control/index.md)
- [vesctl api-endpoint discover](discover/index.md)
