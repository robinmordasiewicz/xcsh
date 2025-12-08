# Environment Variables

vesctl can be configured using environment variables.

## Authentication Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `VES_P12_PASSWORD` | Password for P12 bundle | `export VES_P12_PASSWORD="secret"` |
| `VES_P12_BUNDLE` | Path to P12 bundle | `export VES_P12_BUNDLE="/path/to/creds.p12"` |
| `VES_CERT` | Path to client certificate | `export VES_CERT="/path/to/cert.pem"` |
| `VES_KEY` | Path to client key | `export VES_KEY="/path/to/key.pem"` |

## Connection Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `VES_SERVER_URLS` | API server URL(s) | `export VES_SERVER_URLS="https://tenant.console.ves.volterra.io/api"` |
| `VES_CACERT` | Path to CA certificate | `export VES_CACERT="/path/to/ca.pem"` |

## Output Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `VES_OUTPUT` | Default output format | `export VES_OUTPUT="json"` |

## Configuration Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `VES_CONFIG` | Path to config file | `export VES_CONFIG="/path/to/config.yaml"` |

## Usage Examples

### Basic Setup

```bash
# Set server URL
export VES_SERVER_URLS="https://your-tenant.console.ves.volterra.io/api"

# Set P12 credentials
export VES_P12_BUNDLE="/path/to/api-creds.p12"
export VES_P12_PASSWORD="your-password"

# Run command
vesctl configuration list namespace
```

### Certificate Authentication

```bash
export VES_SERVER_URLS="https://your-tenant.console.ves.volterra.io/api"
export VES_CERT="/path/to/cert.pem"
export VES_KEY="/path/to/key.pem"

vesctl configuration list namespace
```

### JSON Output Default

```bash
export VES_OUTPUT="json"

# All commands now output JSON by default
vesctl configuration list namespace
```

## Shell Configuration

Add to your shell profile for persistent configuration:

### Bash (~/.bashrc)

```bash
export VES_SERVER_URLS="https://your-tenant.console.ves.volterra.io/api"
export VES_P12_BUNDLE="$HOME/.vesconfig/api-creds.p12"
export VES_P12_PASSWORD="your-password"
```

### Zsh (~/.zshrc)

```bash
export VES_SERVER_URLS="https://your-tenant.console.ves.volterra.io/api"
export VES_P12_BUNDLE="$HOME/.vesconfig/api-creds.p12"
export VES_P12_PASSWORD="your-password"
```

### Fish (~/.config/fish/config.fish)

```fish
set -x VES_SERVER_URLS "https://your-tenant.console.ves.volterra.io/api"
set -x VES_P12_BUNDLE "$HOME/.vesconfig/api-creds.p12"
set -x VES_P12_PASSWORD "your-password"
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
