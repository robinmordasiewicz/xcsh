# Environment Variables

f5xcctl can be configured using environment variables. Run `f5xcctl --help` to see all available environment variables.

## Authentication Variables

| Variable | Description | Related Flag |
|----------|-------------|--------------|
| `F5XC_API_TOKEN` | API token for authenticating with F5 Distributed Cloud services. | `--api-token` |
| `F5XC_CERT` | Path to the client certificate file for mTLS authentication. | `--cert` |
| `F5XC_KEY` | Path to the client private key file for mTLS authentication. | `--key` |
| `F5XC_P12_FILE` | Path to the PKCS#12 bundle file containing client certificate and key. | `--p12-bundle` |
| `F5XC_P12_PASSWORD` | Password for decrypting the PKCS#12 bundle file. | - |

## Connection Variables

| Variable | Description | Related Flag |
|----------|-------------|--------------|
| `F5XC_API_URL` | F5 Distributed Cloud API endpoint URL override. | `--server-url` |
| `F5XC_CACERT` | Path to the CA certificate file for TLS server verification. | `--cacert` |

## Output Variables

| Variable | Description | Related Flag |
|----------|-------------|--------------|
| `F5XC_OUTPUT` | Default output format for command results (text, json, yaml, or table). | `--output-format` |

## Configuration Variables

| Variable | Description | Related Flag |
|----------|-------------|--------------|
| `F5XC_CONFIG` | Path to the f5xcctl configuration file. | `--config` |

## Usage Examples

### Basic Setup

=== "Mac/Linux"

    ```bash
    # Set server URL
    export F5XC_API_URL="https://your-tenant.console.ves.volterra.io/api"

    # Set P12 credentials
    export F5XC_P12_FILE="/path/to/api-creds.p12"
    export F5XC_P12_PASSWORD="your-password"

    # Run command
    f5xcctl configuration list namespace
    ```

=== "Windows"

    ```powershell
    # Set server URL
    $env:F5XC_API_URL = "https://your-tenant.console.ves.volterra.io/api"

    # Set P12 credentials
    $env:F5XC_P12_FILE = "C:\path\to\api-creds.p12"
    $env:F5XC_P12_PASSWORD = "your-password"

    # Run command
    f5xcctl configuration list namespace
    ```

### Certificate Authentication

=== "Mac/Linux"

    ```bash
    export F5XC_API_URL="https://your-tenant.console.ves.volterra.io/api"
    export F5XC_CERT="/path/to/cert.pem"
    export F5XC_KEY="/path/to/key.pem"

    f5xcctl configuration list namespace
    ```

=== "Windows"

    ```powershell
    $env:F5XC_API_URL = "https://your-tenant.console.ves.volterra.io/api"
    $env:F5XC_CERT = "C:\path\to\cert.pem"
    $env:F5XC_KEY = "C:\path\to\key.pem"

    f5xcctl configuration list namespace
    ```

### JSON Output Default

=== "Mac/Linux"

    ```bash
    export F5XC_OUTPUT="json"

    # All commands now output JSON by default
    f5xcctl configuration list namespace
    ```

=== "Windows"

    ```powershell
    $env:F5XC_OUTPUT = "json"

    # All commands now output JSON by default
    f5xcctl configuration list namespace
    ```

## Shell Configuration

Add to your shell profile for persistent configuration:

### Bash (~/.bashrc)

```bash
export F5XC_API_URL="https://your-tenant.console.ves.volterra.io/api"
export F5XC_P12_FILE="$HOME/api-creds.p12"
export F5XC_P12_PASSWORD="your-password"
```

### Zsh (~/.zshrc)

```bash
export F5XC_API_URL="https://your-tenant.console.ves.volterra.io/api"
export F5XC_P12_FILE="$HOME/api-creds.p12"
export F5XC_P12_PASSWORD="your-password"
```

### Fish (~/.config/fish/config.fish)

```fish
set -x F5XC_API_URL "https://your-tenant.console.ves.volterra.io/api"
set -x F5XC_P12_FILE "$HOME/api-creds.p12"
set -x F5XC_P12_PASSWORD "your-password"
```

### PowerShell ($PROFILE)

```powershell
$env:F5XC_API_URL = "https://your-tenant.console.ves.volterra.io/api"
$env:F5XC_P12_FILE = "$env:USERPROFILE\api-creds.p12"
$env:F5XC_P12_PASSWORD = "your-password"
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
