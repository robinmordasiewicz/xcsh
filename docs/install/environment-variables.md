# Environment Variables

vesctl can be configured using environment variables.

## Authentication Variables

| Variable | Description |
|----------|-------------|
| `VES_P12_PASSWORD` | Password for P12 bundle |
| `VES_P12_FILE` | Path to P12 bundle |
| `VES_CERT` | Path to client certificate |
| `VES_KEY` | Path to client key |

## Connection Variables

| Variable | Description |
|----------|-------------|
| `VES_API_URL` | API server URL(s) |
| `VES_CACERT` | Path to CA certificate |

## Output Variables

| Variable | Description |
|----------|-------------|
| `VES_OUTPUT` | Default output format |

## Configuration Variables

| Variable | Description |
|----------|-------------|
| `VES_CONFIG` | Path to config file |

## Usage Examples

### Basic Setup

=== "Mac/Linux"

    ```bash
    # Set server URL
    export VES_API_URL="https://your-tenant.console.ves.volterra.io/api"

    # Set P12 credentials
    export VES_P12_FILE="/path/to/api-creds.p12"
    export VES_P12_PASSWORD="your-password"

    # Run command
    vesctl configuration list namespace
    ```

=== "Windows"

    ```powershell
    # Set server URL
    $env:VES_API_URL = "https://your-tenant.console.ves.volterra.io/api"

    # Set P12 credentials
    $env:VES_P12_FILE = "C:\path\to\api-creds.p12"
    $env:VES_P12_PASSWORD = "your-password"

    # Run command
    vesctl configuration list namespace
    ```

### Certificate Authentication

=== "Mac/Linux"

    ```bash
    export VES_API_URL="https://your-tenant.console.ves.volterra.io/api"
    export VES_CERT="/path/to/cert.pem"
    export VES_KEY="/path/to/key.pem"

    vesctl configuration list namespace
    ```

=== "Windows"

    ```powershell
    $env:VES_API_URL = "https://your-tenant.console.ves.volterra.io/api"
    $env:VES_CERT = "C:\path\to\cert.pem"
    $env:VES_KEY = "C:\path\to\key.pem"

    vesctl configuration list namespace
    ```

### JSON Output Default

=== "Mac/Linux"

    ```bash
    export VES_OUTPUT="json"

    # All commands now output JSON by default
    vesctl configuration list namespace
    ```

=== "Windows"

    ```powershell
    $env:VES_OUTPUT = "json"

    # All commands now output JSON by default
    vesctl configuration list namespace
    ```

## Shell Configuration

Add to your shell profile for persistent configuration:

### Bash (~/.bashrc)

```bash
export VES_API_URL="https://your-tenant.console.ves.volterra.io/api"
export VES_P12_FILE="$HOME/api-creds.p12"
export VES_P12_PASSWORD="your-password"
```

### Zsh (~/.zshrc)

```bash
export VES_API_URL="https://your-tenant.console.ves.volterra.io/api"
export VES_P12_FILE="$HOME/api-creds.p12"
export VES_P12_PASSWORD="your-password"
```

### Fish (~/.config/fish/config.fish)

```fish
set -x VES_API_URL "https://your-tenant.console.ves.volterra.io/api"
set -x VES_P12_FILE "$HOME/api-creds.p12"
set -x VES_P12_PASSWORD "your-password"
```

### PowerShell ($PROFILE)

```powershell
$env:VES_API_URL = "https://your-tenant.console.ves.volterra.io/api"
$env:VES_P12_FILE = "$env:USERPROFILE\api-creds.p12"
$env:VES_P12_PASSWORD = "your-password"
```

!!! note "PowerShell Profile"
    Add to your `$PROFILE` for persistent configuration. Create it first if needed:
    ```powershell
    New-Item -ItemType File -Force -Path $PROFILE
    ```

## Precedence

Environment variables override config file settings but are overridden by command-line flags:

1. Command-line flags (highest priority)
2. **Environment variables**
3. Configuration file
4. Default values (lowest priority)

## Security Considerations

- Never commit environment variables with secrets to version control
- Use secret management tools for production environments
- Consider using `.envrc` with direnv for project-specific settings
- Ensure proper file permissions on shell configuration files
