# xcsh

F5 Distributed Cloud Shell - A command-line interface for managing F5 Distributed Cloud resources.

## Documentation

Full documentation is available at **[robinmordasiewicz.github.io/f5xc-xcsh](https://robinmordasiewicz.github.io/f5xc-xcsh)**

## Installation

### Homebrew

```bash
brew tap robinmordasiewicz/tap
brew install --cask xcsh
```

### Install Script

```bash
curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/f5xc-xcsh/main/install.sh | sh
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

## Development & Domain System

### Automated Domain Synchronization

xcsh uses an **automated CI/CD-driven system** to keep domain definitions synchronized with upstream F5 Distributed Cloud API changes. This ensures the CLI always reflects the latest API structure without manual intervention.

#### How It Works

1. **Daily Checks**: GitHub Actions workflow (`sync-upstream-specs.yml`) checks for new upstream spec versions daily at 6 AM UTC
2. **Automatic Regeneration**: When updates are detected, the system:
   - Downloads latest enriched API specifications
   - Regenerates domain and resource registries
   - Validates code quality and tests pass
   - Creates a pull request with all changes
3. **Idempotent Generation**: Code generation is deterministic - running it twice with identical inputs produces byte-for-byte identical output
4. **CI/CD Validation**: Every commit validates that generated files match upstream specs and are reproducible

#### Domain Registry

The domain registry (`pkg/types/domains_generated.go`) is **automatically generated** from upstream API specifications (``.specs/index.json`). It currently contains **40 domains** organized by functional area:

- **Infrastructure**: cloud_infrastructure, site, site_management, container_services
- **Security**: waf, bot_and_threat_defense, network_security
- **Networking**: network, dns, network_connectivity, vpn
- **Observability**: observability_and_analytics, telemetry_and_insights, statistics
- **Identity**: identity, user_and_account_management, users
- **And 24 more...**

#### Manual Domain Configuration

Team-specific domain customization is managed in `.specs/domain_config.yaml`:

```yaml
# Domain aliases (short command shortcuts)
aliases:
  load_balancer: [lb]
  security: [sec]
  networking: [net]
  infrastructure: [infra]

# Deprecated domains with migration guidance
deprecated_domains:
  config:
    maps_to: system
    reason: "Configuration management merged into system domain"
    deprecated_since: "v1.0.25"

# Missing metadata requiring upstream attention
missing_metadata:
  - domain: api_security
    missing_field: "is_preview"
    reason: "Need to mark preview/beta domains"
```

This file is **version-controlled** and survives automated spec updates, allowing teams to maintain consistent domain aliases across releases.

#### For Developers

To regenerate domain definitions:

```bash
# Full generation pipeline
make generate

# Just regenerate domains
make generate-domains

# Verify idempotency (CI safety check)
make ci-generate
```

The generation pipeline:

1. Downloads latest specs via `scripts/download-specs.sh`
2. Runs `scripts/generate-domains.go` to create domain registry
3. Runs `scripts/generate-schemas.go` to create resource schemas
4. Validates against `scripts/validate-specs.go`

#### Upstream Spec Quality

When spec organization issues are detected, the system can automatically report them to the upstream repository. Use the GitHub issue template to report specification problems:

- Missing metadata fields
- Resource classification issues
- Domain organization concerns

See `.github/ISSUE_TEMPLATE/upstream-spec-quality.md` for the standardized reporting format.

## License

This project is open source. See the LICENSE file for details.
