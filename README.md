# xcsh

F5 Distributed Cloud Shell - A command-line interface for managing F5 Distributed Cloud resources.

## Documentation

Full documentation is available at **[robinmordasiewicz.github.io/xcsh](https://robinmordasiewicz.github.io/xcsh)**

## Installation

### Homebrew

```bash
brew tap robinmordasiewicz/tap
brew install --cask xcsh
```

### Install Script

```bash
curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/xcsh/main/install.sh | sh
```

## Usage

### Basic Command Structure

The CLI is organized around **domains** matching F5 Distributed Cloud API structure:

```bash
xcsh <domain> <operation> <resource-type> [resource-name] [flags]
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
xcsh load_balancer list http_loadbalancer

# Using alias for shorter command
xcsh lb list http_loadbalancer

# List in specific namespace
xcsh lb list http_loadbalancer -n production
```

#### Get a Specific Resource

```bash
# Get a load balancer configuration
xcsh load_balancer get http_loadbalancer example-lb

# Get from specific namespace
xcsh lb get http_loadbalancer example-lb -n production
```

#### Create a Resource

```bash
# Create from YAML file
xcsh load_balancer create http_loadbalancer -i lb-config.yaml

# Create from inline JSON
xcsh lb create origin_pool --json-data '{"metadata":{"name":"example-pool"},...}'
```

#### Delete a Resource

```bash
# Delete with confirmation
xcsh load_balancer delete http_loadbalancer example-lb

# Delete without confirmation (for scripts)
xcsh lb delete http_loadbalancer example-lb --yes
```

#### Apply (Create or Update)

```bash
# Apply from YAML (creates if not exists, updates if does)
xcsh load_balancer apply http_loadbalancer -i lb-config.yaml
```

#### Get Help

```bash
# Show available domains
xcsh --help

# Show domain-specific operations
xcsh load_balancer --help

# Show operation-specific help
xcsh load_balancer list --help

# Show resource-type help
xcsh load_balancer list http_loadbalancer --help
```

## License

This project is open source. See the LICENSE file for details.
