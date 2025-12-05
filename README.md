# F5 Distributed Cloud CLI (f5xc)

An open-source command-line interface for managing F5 Distributed Cloud (formerly Volterra) resources.

## Installation

## Quick Start

### 1. Configure Authentication

Create a configuration file at `~/.vesconfig`:

```yaml
server-urls:
  - https://your-tenant.console.ves.volterra.io/api
p12-bundle: ~/.f5xc/api-creds.p12
```

Or use certificate and key files:

```yaml
server-urls:
  - https://your-tenant.console.ves.volterra.io/api
cert: ~/.f5xc/cert.pem
key: ~/.f5xc/key.pem
```

### 2. Set P12 Password (if using P12 bundle)

```bash
export VES_P12_PASSWORD="your-p12-password"
```

### 3. Verify Configuration

```bash
f5xc version
f5xc namespace list
```

## Usage

### Command Structure

f5xc follows Azure CLI-style command patterns:

```bash
f5xc <resource-type> <operation> [options]
```

### Common Operations

```bash
# List resources
f5xc http-loadbalancer list --namespace my-namespace
f5xc origin-pool list -n my-namespace -o table

# Show resource details
f5xc http-loadbalancer show my-lb --namespace my-namespace
f5xc origin-pool show my-pool -n my-namespace -o json

# Create resource from file
f5xc http-loadbalancer create --file lb.yaml
f5xc origin-pool create -f pool.yaml -n my-namespace

# Update resource
f5xc http-loadbalancer update --file lb.yaml
f5xc origin-pool update -f pool.yaml

# Delete resource
f5xc http-loadbalancer delete my-lb --namespace my-namespace
f5xc origin-pool delete my-pool -n my-namespace --yes
```

### Available Resources

| Resource | Description |
|----------|-------------|
| `http-loadbalancer` | HTTP Load Balancer for L7 traffic management |
| `tcp-loadbalancer` | TCP Load Balancer for L4 traffic management |
| `origin-pool` | Origin Pool for backend server groups |
| `healthcheck` | Health Check for monitoring backend health |
| `namespace` | Namespace for resource organization |
| `app-firewall` | Application Firewall (WAF) policy |
| `service-policy` | Service Policy for traffic control |
| `certificate` | TLS Certificate management |
| `aws-vpc-site` | AWS VPC Site for cloud connectivity |
| `azure-vnet-site` | Azure VNet Site for cloud connectivity |
| `virtual-host` | Virtual Host configuration |
| `rate-limiter` | Rate Limiter policy |
| `cloud-credentials` | Cloud provider credentials |
| `api-definition` | API Definition for OpenAPI specs |
| `user` | User account management (read-only) |

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | | Path to config file (default: ~/.vesconfig) |
| `--server-urls` | `-u` | API server URL(s) |
| `--cert` | `-c` | Path to client certificate file |
| `--key` | `-k` | Path to client key file |
| `--cacert` | | Path to CA certificate file |
| `--p12-bundle` | | Path to P12 certificate bundle |
| `--output` | `-o` | Output format: json, yaml, table, tsv, none |
| `--query` | | JMESPath query string for filtering output |
| `--debug` | | Show all debug logs |
| `--verbose` | | Increase logging verbosity |
| `--only-show-errors` | | Only show errors, suppressing warnings |
| `--no-wait` | | Do not wait for long-running operations to finish |

### Output Formats

```bash
# YAML output (default)
f5xc namespace list

# JSON output
f5xc namespace list -o json

# Table output
f5xc namespace list -o table

# TSV output (for scripting)
f5xc namespace list -o tsv
```

## Shell Completion

### Bash

```bash
# For current session
source <(f5xc completion bash)

# Permanent (Linux)
f5xc completion bash > /etc/bash_completion.d/f5xc

# Permanent (macOS with Homebrew)
f5xc completion bash > $(brew --prefix)/etc/bash_completion.d/f5xc
```

### Zsh

```bash
# Enable completion
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Install completion
f5xc completion zsh > "${fpath[1]}/_f5xc"
```

### Fish

```bash
f5xc completion fish > ~/.config/fish/completions/f5xc.fish
```

### PowerShell

```powershell
f5xc completion powershell | Out-String | Invoke-Expression
```

## Authentication

### P12 Bundle

The recommended authentication method uses a P12 bundle downloaded from the F5 XC Console:

1. Log in to F5 XC Console
2. Navigate to **Administration** > **Personal Management** > **Credentials**
3. Create credentials with type "API Certificate"
4. Download the P12 file
5. Set the password: `export VES_P12_PASSWORD="your-password"`

### Certificate and Key

Alternatively, extract certificate and key from P12:

```bash
# Extract certificate
openssl pkcs12 -in api-creds.p12 -nodes -nokeys -out cert.pem

# Extract private key
openssl pkcs12 -in api-creds.p12 -nodes -nocerts -out key.pem
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `VES_P12_PASSWORD` | Password for P12 bundle |
| `F5XC_SERVER_URLS` | API server URL(s) |
| `F5XC_CERT` | Path to client certificate |
| `F5XC_KEY` | Path to client key |
| `F5XC_P12_BUNDLE` | Path to P12 bundle |
| `F5XC_OUTPUT` | Default output format |

## Examples

### Create an HTTP Load Balancer

```yaml
# lb.yaml
metadata:
  name: my-http-lb
  namespace: my-namespace
spec:
  domains:
    - example.com
  http:
    port: 80
  default_route_pools:
    - pool:
        name: my-origin-pool
        namespace: my-namespace
```

```bash
f5xc http-loadbalancer create -f lb.yaml
```

### Create an Origin Pool

```yaml
# pool.yaml
metadata:
  name: my-origin-pool
  namespace: my-namespace
spec:
  origin_servers:
    - public_ip:
        ip: 192.168.1.100
  port: 8080
  healthcheck:
    - name: my-healthcheck
      namespace: my-namespace
```

```bash
f5xc origin-pool create -f pool.yaml
```

### List and Filter Resources

```bash
# List all load balancers in namespace
f5xc http-loadbalancer list -n my-namespace

# Show specific load balancer as JSON
f5xc http-loadbalancer show my-lb -n my-namespace -o json

# List with table format
f5xc namespace list -o table
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source. See the LICENSE file for details.

## Related Documentation

- [F5 Distributed Cloud Documentation](https://docs.cloud.f5.com/)
- [F5 XC API Reference](https://docs.cloud.f5.com/docs/api)
