---
title: "vesctl completion"
description: "Generate shell completion scripts for bash or zsh."
keywords:
  - F5 Distributed Cloud
  - vesctl
  - completion
  - F5 XC
command: "vesctl completion"
command_group: "completion"
---

# vesctl completion

> Generate shell completion scripts for bash or zsh.

## Synopsis

```bash
vesctl completion <command> [flags]
```

## Description

To load completions:
Bash:

$ source <(vesctl completion bash)

# To load completions for each session, execute once:
Linux:
  $ vesctl completion bash > /etc/bash_completion.d/yourprogram
MacOS:
  $ vesctl completion bash > /usr/local/etc/bash_completion.d/yourprogram

Zsh:

 $ source <(vesctl completion zsh)

 # To load completions for each session, execute once:
 $ vesctl completion zsh > "${fpath[1]}/_vesctl"


## Available Commands

| Command | Description |
|---------|-------------|

## See Also

- [Command Reference](../index.md)
