# Authentication

xcsh supports multiple authentication methods to connect to the F5 Distributed Cloud API.

## P12 Bundle (Recommended)

The recommended authentication method uses a P12 certificate bundle.

### Obtaining Credentials

1. Log in to the F5 XC Console
2. Navigate to **Administration** > **Personal Management** > **Credentials**
3. Click **Create Credentials**
4. Select type **API Certificate**
5. Download the P12 file
6. Note the password provided

### Configuration

**Using configuration file (~/.f5xcconfig):**

```yaml
server-url: https://your-tenant.console.ves.volterra.io/api
p12-bundle: /path/to/api-creds.p12
```

**Set the password as environment variable:**

```bash
export F5XC_P12_PASSWORD="your-p12-password"
```

## Certificate and Key

Use separate certificate and key files instead of a P12 bundle.

### Extract from P12

If you have a P12 file, extract the certificate and key:

```bash
# Extract certificate
openssl pkcs12 -in api-creds.p12 -nodes -nokeys -out cert.pem

# Extract private key
openssl pkcs12 -in api-creds.p12 -nodes -nocerts -out key.pem
```

### Configuration

**Using configuration file (~/.f5xcconfig):**

```yaml
server-url: https://your-tenant.console.ves.volterra.io/api
cert: /path/to/cert.pem
key: /path/to/key.pem
```

**Using command-line flags:**

```bash
xcsh --cert /path/to/cert.pem --key /path/to/key.pem configuration list namespace
```

## API Token

Use an API token for authentication without managing certificate files. Ideal for CI/CD pipelines and automation.

### Obtaining an API Token

1. Log in to the F5 XC Console
2. Navigate to **Administration** > **Personal Management** > **Credentials**
3. Click **Create Credentials**
4. Select type **API Token**
5. Copy the generated token

### Configuration

**Using environment variables (recommended):**

```bash
export F5XC_API_TOKEN="your-api-token"
export F5XC_API_URL="https://your-tenant.console.ves.volterra.io"  # Optional, overrides config

xcsh identity list namespace
```

**Using configuration file (~/.f5xcconfig):**

```yaml
server-url: https://your-tenant.console.ves.volterra.io/api
api-token: true  # Token value from F5XC_API_TOKEN environment variable
```

**Using interactive configuration:**

```bash
xcsh configure
# Select option 3: API Token
```

**Using command-line flag:**

```bash
xcsh --api-token configuration list namespace
```

**Using login command:**

```bash
export F5XC_API_TOKEN='your-api-token'
xcsh login --tenant example-tenant --api-token
```

## Configuration File

The default configuration file location is `~/.f5xcconfig`.

### Full Example

```yaml
# API server URL
server-url: https://your-tenant.console.ves.volterra.io/api

# P12 bundle authentication
p12-bundle: /path/to/api-creds.p12

# OR certificate/key authentication
# cert: /path/to/cert.pem
# key: /path/to/key.pem

# OR API token authentication (token from F5XC_API_TOKEN env var)
# api-token: true

# Optional CA certificate for custom trust
# cacert: /path/to/ca.pem
```

### Custom Config Location

Specify a custom configuration file:

```bash
xcsh --config /path/to/custom-config.yaml configuration list namespace
```

## Environment Variables

Override configuration with environment variables:

| Variable | Description |
|----------|-------------|
| `F5XC_API_TOKEN` | API token for authentication |
| `F5XC_API_URL` | API server URL (overrides server-url in config) |
| `F5XC_P12_PASSWORD` | Password for P12 bundle |
| `F5XC_API_URL` | API server URL(s) |
| `F5XC_P12_FILE` | Path to P12 bundle |
| `F5XC_CERT` | Path to client certificate |
| `F5XC_KEY` | Path to client key |

### Example

```bash
export F5XC_API_URL="https://your-tenant.console.ves.volterra.io/api"
export F5XC_P12_FILE="/path/to/api-creds.p12"
export F5XC_P12_PASSWORD="your-password"

xcsh identity list namespace
```

## Interactive Configuration

Use `xcsh configure` for interactive setup:

```bash
xcsh configure
```

This will prompt for:

1. API server URL
2. Authentication method (P12, cert/key, or API token)
3. File paths (for P12 or cert/key) or environment variable instructions (for API token)
4. Default output format

## Verifying Authentication

Test your configuration:

```bash
# List namespaces to verify connection
xcsh identity list namespace

# Check with verbose output for troubleshooting
xcsh --verbose configuration list namespace
```
