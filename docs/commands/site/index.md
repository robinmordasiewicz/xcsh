---
title: "vesctl site"
description: "Deploy and manage F5 XC sites on public cloud providers."
keywords:
  - F5 Distributed Cloud
  - vesctl
  - F5 XC
  - site
command: "vesctl site"
command_group: "site"
aliases:
  - "s"
---

# vesctl site

> Deploy and manage F5 XC sites on public cloud providers.

## Synopsis

```bash
vesctl site <command> [flags]
```

## Aliases

This command can also be invoked as:

- `vesctl s`

## Available Commands

| Command | Description |
|---------|-------------|
| [aws_vpc](aws_vpc/index.md) | Manage AWS VPC site creation through view apis |
| [azure_vnet](azure_vnet/index.md) | Manage Azure VNet site creation through view apis |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--log-color` |  | bool | true | Enable colored log output. |
| `--log-fabulous` |  | bool | true | Enable enhanced log formatting. |
| `--log-level` |  | int | 3 | Set the logging verbosity level for site operations. |

## Examples

```bash
vesctl site aws_vpc create
```

## See Also

- [Command Reference](../index.md)
- [vesctl site aws_vpc](aws_vpc/index.md)
- [vesctl site azure_vnet](azure_vnet/index.md)
