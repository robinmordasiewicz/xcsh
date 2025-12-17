# Shell Completion

Enable tab completion for your shell.

## Bash

```bash
# Current session
source <(f5xcctl completion bash)

# Permanent (Linux)
f5xcctl completion bash > /etc/bash_completion.d/f5xcctl

# Permanent (macOS with Homebrew)
f5xcctl completion bash > $(brew --prefix)/etc/bash_completion.d/f5xcctl
```

## Zsh

```bash
# Enable completion system (if not already enabled)
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Install completion
f5xcctl completion zsh > "${fpath[1]}/_f5xcctl"
```

## Fish

```bash
f5xcctl completion fish > ~/.config/fish/completions/f5xcctl.fish
```

## PowerShell

```powershell
f5xcctl completion powershell | Out-String | Invoke-Expression
```

## Troubleshooting

### Completions not working after upgrade

Zsh caches completion functions in `~/.zcompdump*` files. After upgrading f5xcctl, the stale cache may prevent new completions from loading.

**Fix:** Clear the cache and restart your shell:

```bash
rm -f ~/.zcompdump* && exec zsh
```

### Completions not loading at all

Ensure the completion system is enabled in your `~/.zshrc`:

```bash
autoload -Uz compinit && compinit
```

Then restart your shell or run:

```bash
source ~/.zshrc
```
