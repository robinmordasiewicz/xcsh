# Authentication

vesctl supports multiple authentication methods to connect to the F5 Distributed Cloud API.

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

**Using configuration file (~/.vesconfig):**

```yaml
server-urls:
  - https://your-tenant.console.ves.volterra.io/api
p12-bundle: /path/to/api-creds.p12
```

**Set the password as environment variable:**

```bash
export VES_P12_PASSWORD="your-p12-password"
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

**Using configuration file (~/.vesconfig):**

```yaml
server-urls:
  - https://your-tenant.console.ves.volterra.io/api
cert: /path/to/cert.pem
key: /path/to/key.pem
```

**Using command-line flags:**

```bash
vesctl --cert /path/to/cert.pem --key /path/to/key.pem configuration list namespace
```

## Configuration File

The default configuration file location is `~/.vesconfig`.

### Full Example

```yaml
# API server URL
server-urls:
  - https://your-tenant.console.ves.volterra.io/api

# P12 bundle authentication
p12-bundle: /path/to/api-creds.p12

# OR certificate/key authentication
# cert: /path/to/cert.pem
# key: /path/to/key.pem

# Optional CA certificate for custom trust
# cacert: /path/to/ca.pem
```

### Custom Config Location

Specify a custom configuration file:

```bash
vesctl --config /path/to/custom-config.yaml configuration list namespace
```

## Environment Variables

Override configuration with environment variables:

| Variable | Description |
|----------|-------------|
| `VES_P12_PASSWORD` | Password for P12 bundle |
| `VES_SERVER_URLS` | API server URL(s) |
| `VES_P12_BUNDLE` | Path to P12 bundle |
| `VES_CERT` | Path to client certificate |
| `VES_KEY` | Path to client key |

### Example

```bash
export VES_SERVER_URLS="https://your-tenant.console.ves.volterra.io/api"
export VES_P12_BUNDLE="/path/to/api-creds.p12"
export VES_P12_PASSWORD="your-password"

vesctl configuration list namespace
```

## Interactive Configuration

Use `vesctl configure` for interactive setup:

```bash
vesctl configure
```

This will prompt for:

1. API server URL
2. Authentication method (P12 or cert/key)
3. File paths
4. (Optional) P12 password

## Verifying Authentication

Test your configuration:

```bash
# List namespaces to verify connection
vesctl configuration list namespace

# Check with verbose output for troubleshooting
vesctl --verbose configuration list namespace
```

## Security Best Practices

1. **Protect your credentials** - Keep P12 files and keys secure with appropriate file permissions
2. **Use environment variables** - Never commit passwords to version control
3. **Rotate credentials** - Periodically rotate your API certificates
4. **Limit scope** - Use credentials with minimum required permissions

### File Permissions

```bash
# Restrict access to credential files
chmod 600 ~/.vesconfig
chmod 600 /path/to/api-creds.p12
chmod 600 /path/to/key.pem
```

## Troubleshooting

### Certificate Errors

```
Error: x509: certificate signed by unknown authority
```

Add the CA certificate:

```yaml
cacert: /path/to/ca-bundle.pem
```

### Invalid Password

```
Error: pkcs12: decryption password incorrect
```

Verify `VES_P12_PASSWORD` is set correctly.

### Connection Refused

```
Error: dial tcp: connection refused
```

Check the server URL and network connectivity.
