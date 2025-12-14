# Shell Completion

Enable tab completion for your shell:

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
# Enable completion system
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Install completion
f5xcctl completion zsh > "${fpath[1]}/_vesctl"
```

## Fish

```bash
f5xcctl completion fish > ~/.config/fish/completions/f5xcctl.fish
```

## PowerShell

```powershell
f5xcctl completion powershell | Out-String | Invoke-Expression
```
