# Install Script

Install xcsh on Linux or macOS using the universal install script.

## Quick Install

```bash
curl -fsSL https://robinmordasiewicz.github.io/xcsh/install.sh | sh
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

- Binary: `/usr/local/bin/xcsh`
- Completions: System-wide locations

### Without Administrator Access

- Binary: `~/.local/bin/xcsh`
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
F5XC_NO_SUDO=1 curl -fsSL https://robinmordasiewicz.github.io/xcsh/install.sh | sh
```

### Example: Specific Version

```bash
F5XC_VERSION=4.39.0 curl -fsSL https://robinmordasiewicz.github.io/xcsh/install.sh | sh
```

## Post-Installation

After installation, you may need to restart your shell or source your profile:

```bash
source ~/.zshrc  # or ~/.bashrc
```

Verify the installation:

```bash
xcsh version
```

## Uninstall

To uninstall, remove the binary and completions:

```bash
# Remove binary
rm -f ~/.local/bin/xcsh  # or /usr/local/bin/xcsh

# Remove completions
rm -f ~/.zsh/completions/_xcsh
rm -f ~/.local/share/bash-completion/completions/xcsh
rm -f ~/.config/fish/completions/xcsh.fish
```
