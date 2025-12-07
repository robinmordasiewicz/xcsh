# Install Script Reference

The `install.sh` script provides automated installation of vesctl with platform detection, checksum verification, and shell completion setup.

## Quick Usage

Install the latest version with a single command:

=== "curl"

    ```bash
    curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | sh
    ```

=== "wget"

    ```bash
    wget -qO- https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | sh
    ```

## Features

| Feature | Description |
|---------|-------------|
| **Automatic Platform Detection** | Detects OS and architecture automatically |
| **HTTP Client Fallback** | Works with curl or wget |
| **Retry Logic** | 3 retries with 2-second delays for resilient downloads |
| **Checksum Verification** | SHA256 verification using sha256sum, shasum, or openssl |
| **Shell Completion** | Automatic setup for bash, zsh, and fish |
| **Upgrade Detection** | Skips installation if the same version is already installed |

## Supported Platforms

| Operating System | Architectures |
|------------------|---------------|
| **Linux** | amd64 (x86_64), arm64 (aarch64) |
| **macOS** | amd64 (Intel), arm64 (Apple Silicon) |
| **Windows** | amd64, arm64 (manual download required) |

!!! note "Windows Users"
    The install script does not support Windows directly. Download the appropriate zip file from the [GitHub Releases](https://github.com/robinmordasiewicz/vesctl/releases/latest) page.

## Environment Variables

Configure the installer behavior with these environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `VESCTL_VERSION` | Specific version to install | latest |
| `VESCTL_INSTALL_DIR` | Installation directory | `/usr/local/bin` |
| `VESCTL_NO_SUDO` | Skip sudo usage if set | unset |
| `VESCTL_NO_VERIFY` | Skip checksum verification | unset |

### Examples

**Install a specific version:**

```bash
curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | VESCTL_VERSION=1.1.0 sh
```

**Install to a custom directory (no sudo required):**

```bash
curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | VESCTL_INSTALL_DIR=$HOME/.local/bin sh
```

**Skip checksum verification (not recommended):**

```bash
curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | VESCTL_NO_VERIFY=1 sh
```

## Command Line Options

| Option | Description |
|--------|-------------|
| `--help`, `-h` | Display usage information |
| `--uninstall` | Remove vesctl and shell completions |

### Uninstall

Remove vesctl and all shell completions:

```bash
curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | sh -s -- --uninstall
```

## How It Works

### Platform Detection

The script detects your platform using standard Unix commands:

1. **Operating System**: Uses `uname -s` to identify Linux, macOS, or other systems
2. **Architecture**: Uses `uname -m` to detect x86_64/amd64 or arm64/aarch64

### HTTP Client Selection

The script automatically selects an available HTTP client:

1. **curl** (preferred) - with retry and timeout options
2. **wget** (fallback) - if curl is not available

Both clients are configured with:

- 3 retry attempts
- 2-second delay between retries
- Silent/quiet mode for clean output

### Checksum Verification

Downloaded archives are verified using SHA256 checksums. The script tries multiple tools in order:

1. `sha256sum` (Linux standard)
2. `shasum -a 256` (macOS standard)
3. `openssl dgst -sha256` (fallback)

If no SHA256 tool is available, verification is skipped with a warning.

### Installation Process

1. **Detect platform** - Identify OS and architecture
2. **Fetch version** - Get latest version from GitHub API (or use specified version)
3. **Check existing** - Skip if same version already installed
4. **Download archive** - Fetch the appropriate tar.gz file
5. **Verify checksum** - Validate SHA256 hash
6. **Extract binary** - Unpack the archive
7. **Install** - Move binary to installation directory
8. **Setup completion** - Configure shell completion for detected shell

## Shell Completion

The installer automatically sets up shell completion based on your current shell:

### Bash

Completion files are installed to one of:

- `/etc/bash_completion.d/vesctl` (system-wide, if writable)
- `~/.local/share/bash-completion/completions/vesctl` (user)

### Zsh

Completion files are installed to:

- `~/.zsh/completions/_vesctl`

Add to your `~/.zshrc`:

```bash
fpath=(~/.zsh/completions $fpath)
autoload -Uz compinit && compinit
```

### Fish

Completion files are installed to:

- `~/.config/fish/completions/vesctl.fish`

## Troubleshooting

### Permission Denied

If you encounter permission errors:

```bash
# Option 1: Use a custom directory
VESCTL_INSTALL_DIR=$HOME/.local/bin sh install.sh

# Option 2: Run as root (not recommended)
sudo sh install.sh
```

### HTTP Client Not Found

Install curl or wget:

```bash
# Debian/Ubuntu
sudo apt-get install curl

# RHEL/CentOS
sudo yum install curl

# Alpine
apk add curl

# macOS (curl is pre-installed)
```

### Checksum Verification Failed

If checksum verification fails:

1. Try downloading again - the file may have been corrupted
2. Check your internet connection
3. Verify the release exists on [GitHub](https://github.com/robinmordasiewicz/vesctl/releases)

### Version Not Found

If the specified version doesn't exist:

```bash
# List available versions
curl -s https://api.github.com/repos/robinmordasiewicz/vesctl/releases | grep tag_name

# Install latest instead
curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | sh
```

### Binary Not in PATH

If `vesctl` is not found after installation:

```bash
# Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
export PATH="$PATH:/usr/local/bin"

# Or if using custom directory
export PATH="$PATH:$HOME/.local/bin"
```
