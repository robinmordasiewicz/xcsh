# Utility Commands

Additional utility commands for vesctl.

## Version

Display version information:

```bash
vesctl version
```

**Example Output:**

```text
vesctl version 0.1.0
  commit: abc1234
  built:  2024-01-01T00:00:00Z
```

### Version Flags

| Flag | Description |
|------|-------------|
| `--short` | Print version number only |
| `--json` | Output as JSON |

## Completion

Generate shell completion scripts for command auto-completion.

```bash
vesctl completion [bash|zsh|fish|powershell]
```

### Bash Completion

```bash
# Current session
source <(vesctl completion bash)

# Permanent (Linux)
vesctl completion bash > /etc/bash_completion.d/vesctl

# Permanent (macOS with Homebrew)
vesctl completion bash > $(brew --prefix)/etc/bash_completion.d/vesctl
```

### Zsh Completion

```bash
# Enable completion
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Install completion
vesctl completion zsh > "${fpath[1]}/_vesctl"
```

### Fish Completion

```bash
vesctl completion fish > ~/.config/fish/completions/vesctl.fish
```

### PowerShell Completion

```powershell
vesctl completion powershell | Out-String | Invoke-Expression
```

## Configure

Interactive configuration setup:

```bash
vesctl configure
```

This command guides you through:

1. Setting the API server URL
2. Choosing authentication method
3. Configuring credential file paths
4. Testing the connection

### Configuration File

The `configure` command creates or updates `~/.vesconfig`:

```yaml
server-urls:
  - https://your-tenant.console.ves.volterra.io/api
p12-bundle: /path/to/credentials.p12
```

## Help

Get help for any command:

```bash
# General help
vesctl --help

# Command-specific help
vesctl configuration --help
vesctl configuration list --help
```
