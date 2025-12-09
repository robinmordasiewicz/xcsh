---
title: "vesctl request"
description: "Execute custom API requests to F5 Distributed Cloud."
keywords:
  - F5 XC
  - request
  - vesctl
  - F5 Distributed Cloud
command: "vesctl request"
command_group: "request"
aliases:
  - "req"
  - "r"
---

# vesctl request

> Execute custom API requests to F5 Distributed Cloud.

## Synopsis

```bash
vesctl request <command> [flags]
```

## Aliases

This command can also be invoked as:

- `vesctl req`
- `vesctl r`

## Available Commands

| Command | Description |
|---------|-------------|
| [command-sequence](command-sequence/index.md) | Execute a command sequence |
| [rpc](rpc/index.md) | RPC Invocation |
| [secrets](secrets/index.md) | Execute commands for secret_management |

## Examples

```bash
vesctl request secrets encrypt --policy-doc temp_policy --public-key pub_key secret
```

## See Also

- [Command Reference](../index.md)
- [vesctl request command-sequence](command-sequence/index.md)
- [vesctl request rpc](rpc/index.md)
- [vesctl request secrets](secrets/index.md)
