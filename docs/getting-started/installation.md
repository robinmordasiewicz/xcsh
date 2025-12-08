# Installation

This guide covers all methods for installing vesctl.

## Homebrew (Recommended for macOS/Linux)

The easiest way to install vesctl on macOS or Linux:

```bash
brew tap robinmordasiewicz/tap
brew install vesctl
```

**Upgrade to latest version:**

```bash
brew upgrade vesctl
```

**Uninstall:**

```bash
brew uninstall vesctl
```

## Quick Install

The easiest way to install vesctl is with the install script:

```bash
curl -fsSL https://robinmordasiewicz.github.io/vesctl/install.sh | sh
```

This will:

- Detect your platform (Linux/macOS, amd64/arm64)
- Download the latest release
- Verify the checksum
- Install to `/usr/local/bin`
- Set up shell completion

### Install Options

**Install a specific version:**

```bash
curl -fsSL https://robinmordasiewicz.github.io/vesctl/install.sh | VESCTL_VERSION=1.1.0 sh
```

**Install to a custom directory:**

```bash
curl -fsSL https://robinmordasiewicz.github.io/vesctl/install.sh | VESCTL_INSTALL_DIR=$HOME/.local/bin sh
```

**Uninstall:**

```bash
curl -fsSL https://robinmordasiewicz.github.io/vesctl/install.sh | sh -s -- --uninstall
```

## Manual Installation

Download the appropriate binary from [GitHub Releases](https://github.com/robinmordasiewicz/vesctl/releases).

## Build from Source

Build vesctl from source code:

```bash
# Clone the repository
git clone https://github.com/robinmordasiewicz/vesctl.git
cd vesctl

# Build the binary
go build -o vesctl .

# Install (optional)
sudo mv vesctl /usr/local/bin/
```

## Verify Installation

After installation, verify vesctl is working:

```bash
vesctl version
```

Expected output:

```text
vesctl version 0.1.0
  commit: abc1234
  built:  2024-01-01T00:00:00Z
```

## Shell Completion

Enable tab completion for your shell:

### Bash

```bash
# Current session
source <(vesctl completion bash)

# Permanent (Linux)
vesctl completion bash > /etc/bash_completion.d/vesctl

# Permanent (macOS with Homebrew)
vesctl completion bash > $(brew --prefix)/etc/bash_completion.d/vesctl
```

### Zsh

```bash
# Enable completion system
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Install completion
vesctl completion zsh > "${fpath[1]}/_vesctl"
```

### Fish

```bash
vesctl completion fish > ~/.config/fish/completions/vesctl.fish
```

### PowerShell

```powershell
vesctl completion powershell | Out-String | Invoke-Expression
```

## Troubleshooting

### Permission Denied

If you get a permission denied error:

```bash
# Make the binary executable
chmod +x vesctl

# Or use sudo for installation
sudo mv vesctl /usr/local/bin/
```

### Command Not Found

Ensure the installation directory is in your PATH:

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH="$PATH:/usr/local/bin"
```
