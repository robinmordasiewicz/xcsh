# Homebrew

Install f5xcctl on macOS using Homebrew.

## Installation

```bash
brew tap robinmordasiewicz/f5xcctl
brew install --cask f5xcctl
```

## Verify Installation

```bash
f5xcctl version
```

## Update

```bash
brew upgrade --cask f5xcctl
```

## Uninstall

```bash
brew uninstall --cask f5xcctl
```

## What Gets Installed

- Binary at `/opt/homebrew/Caskroom/f5xcctl/{version}/f5xcctl`
- Symlink at `/opt/homebrew/bin/f5xcctl`
- Shell completions for bash, zsh, and fish

## Troubleshooting

### Command Not Found

If `f5xcctl` is not found after installation, ensure Homebrew's bin directory is in your PATH:

```bash
export PATH="/opt/homebrew/bin:$PATH"
```

Add this to your shell profile (`~/.zshrc` or `~/.bashrc`) to make it permanent.

### Tap Not Found

If the tap command fails, verify network connectivity and try:

```bash
brew tap robinmordasiewicz/f5xcctl --force-auto-update
```
