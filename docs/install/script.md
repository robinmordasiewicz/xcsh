# Install Script

Install f5xcctl on Linux or macOS using the universal install script.

## Quick Install

```bash
curl -fsSL https://robinmordasiewicz.github.io/f5xcctl/install.sh | sh
```

## What the Script Does

1. Detects your operating system and architecture
2. Downloads the appropriate binary from GitHub releases
3. Verifies the checksum for security
4. Installs the binary to an appropriate location
5. Sets up shell completions for your shell
6. Configures your PATH if needed

## Installation Locations

### With Administrator Access

- Binary: `/usr/local/bin/f5xcctl`
- Completions: System-wide locations

### Without Administrator Access

- Binary: `~/.local/bin/f5xcctl`
- Completions: User directories (`~/.zsh/completions/`, `~/.local/share/bash-completion/`)

## Environment Variables

Customize the installation with these environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `F5XC_NO_SUDO` | Force user-local installation | Detect |
| `F5XC_INSTALL_DIR` | Custom binary location | Auto |
| `F5XC_VERSION` | Specific version to install | Latest |

### Example: User-Local Installation

```bash
F5XC_NO_SUDO=1 curl -fsSL https://robinmordasiewicz.github.io/f5xcctl/install.sh | sh
```

### Example: Specific Version

```bash
F5XC_VERSION=4.39.0 curl -fsSL https://robinmordasiewicz.github.io/f5xcctl/install.sh | sh
```

## Post-Installation

After installation, you may need to restart your shell or source your profile:

```bash
source ~/.zshrc  # or ~/.bashrc
```

Verify the installation:

```bash
f5xcctl version
```

## Uninstall

To uninstall, remove the binary and completions:

```bash
# Remove binary
rm -f ~/.local/bin/f5xcctl  # or /usr/local/bin/f5xcctl

# Remove completions
rm -f ~/.zsh/completions/_f5xcctl
rm -f ~/.local/share/bash-completion/completions/f5xcctl
rm -f ~/.config/fish/completions/f5xcctl.fish
```
