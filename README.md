# f5xcctl

A command-line interface for managing F5 Distributed Cloud resources.

## Documentation

Full documentation is available at **[robinmordasiewicz.github.io/f5xcctl](https://robinmordasiewicz.github.io/f5xcctl)**

## Installation

### Homebrew

```bash
brew tap robinmordasiewicz/tap
brew install --cask f5xcctl
```

### Install Script

```bash
curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/f5xcctl/main/install.sh | sh
```

## Usage

### Basic Command Structure

The CLI is organized around **domains** matching F5 Distributed Cloud API structure:

```bash
f5xcctl <domain> <operation> <resource-type> [resource-name] [flags]
```

### Domain-Based Commands

Common domains include:

| Domain | Alias | Purpose |
|--------|-------|---------|
| `load_balancer` | `lb` | Load balancing and origin pools |
| `infrastructure` | `infra` | Core infrastructure resources |
| `security` | `sec` | WAF, DDoS, bot defense |
| `networking` | `net` | Network routing and configuration |
| `observability` | `obs`, `o11y` | Monitoring and observability |
| `api_security` | `apisec` | API protection and security |
| `identity` | `iam` | Identity and access management |

### Examples

#### List Resources

```bash
# List HTTP load balancers in default namespace
f5xcctl load_balancer list http_loadbalancer

# Using alias for shorter command
f5xcctl lb list http_loadbalancer

# List in specific namespace
f5xcctl lb list http_loadbalancer -n production
```

#### Get a Specific Resource

```bash
# Get a load balancer configuration
f5xcctl load_balancer get http_loadbalancer example-lb

# Get from specific namespace
f5xcctl lb get http_loadbalancer example-lb -n production
```

#### Create a Resource

```bash
# Create from YAML file
f5xcctl load_balancer create http_loadbalancer -i lb-config.yaml

# Create from inline JSON
f5xcctl lb create origin_pool --json-data '{"metadata":{"name":"example-pool"},...}'
```

#### Delete a Resource

```bash
# Delete with confirmation
f5xcctl load_balancer delete http_loadbalancer example-lb

# Delete without confirmation (for scripts)
f5xcctl lb delete http_loadbalancer example-lb --yes
```

#### Apply (Create or Update)

```bash
# Apply from YAML (creates if not exists, updates if does)
f5xcctl load_balancer apply http_loadbalancer -i lb-config.yaml
```

#### Get Help

```bash
# Show available domains
f5xcctl --help

# Show domain-specific operations
f5xcctl load_balancer --help

# Show operation-specific help
f5xcctl load_balancer list --help

# Show resource-type help
f5xcctl load_balancer list http_loadbalancer --help
```

## License

This project is open source. See the LICENSE file for details.
