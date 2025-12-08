# Supported Resources

vesctl supports the following F5 Distributed Cloud resource types.

## Resource Types

### Load Balancing

| Resource Type | Description |
|---------------|-------------|
| `http_loadbalancer` | HTTP Load Balancer for L7 traffic |
| `tcp_loadbalancer` | TCP Load Balancer for L4 traffic |
| `origin_pool` | Origin Pool for backend servers |
| `healthcheck` | Health Check for monitoring backends |

### Security

| Resource Type | Description |
|---------------|-------------|
| `app_firewall` | Application Firewall (WAF) policy |
| `service_policy` | Service Policy for traffic control |
| `rate_limiter` | Rate Limiter policy |
| `ip_prefix_set` | IP Prefix Set for allowlist/blocklist |

### API Security

| Resource Type | Description |
|---------------|-------------|
| `api_definition` | API Definition for OpenAPI specs |
| `api_endpoint` | Discovered API endpoints |

### Infrastructure

| Resource Type | Description |
|---------------|-------------|
| `namespace` | Namespace for resource organization |
| `certificate` | TLS Certificate management |
| `virtual_host` | Virtual Host configuration |

### Cloud Sites

| Resource Type | Description |
|---------------|-------------|
| `aws_vpc_site` | AWS VPC Site for cloud connectivity |
| `azure_vnet_site` | Azure VNet Site for cloud connectivity |
| `gcp_vpc_site` | GCP VPC Site for cloud connectivity |

### Credentials

| Resource Type | Description |
|---------------|-------------|
| `cloud_credentials` | Cloud provider credentials |

### Identity

| Resource Type | Description |
|---------------|-------------|
| `user` | User account (read-only) |

## Usage

### List Resources

```bash
vesctl configuration list <resource-type> [flags]
```

Examples:

```bash
# List all namespaces
vesctl configuration list namespace

# List HTTP load balancers in a namespace
vesctl configuration list http_loadbalancer -n example-namespace
```

### Get Resource

```bash
vesctl configuration get <resource-type> <name> [flags]
```

Examples:

```bash
# Get namespace details
vesctl configuration get namespace example-namespace

# Get with JSON output
vesctl configuration get http_loadbalancer example-lb -n example-ns --outfmt json
```

### Create Resource

```bash
vesctl configuration create <resource-type> -i <file> [flags]
```

Examples:

```bash
# Create HTTP load balancer from file
vesctl configuration create http_loadbalancer -i lb.yaml

# Create in specific namespace
vesctl configuration create origin_pool -i pool.yaml -n example-namespace
```

### Delete Resource

```bash
vesctl configuration delete <resource-type> <name> [flags]
```

Examples:

```bash
# Delete with confirmation
vesctl configuration delete http_loadbalancer example-lb -n example-namespace

# Delete without confirmation
vesctl configuration delete origin_pool example-pool -n example-namespace --yes
```

## Resource YAML Format

Resources are defined in YAML format:

```yaml
metadata:
  name: resource-name
  namespace: example-namespace
  labels:
    key: value
spec:
  # Resource-specific configuration
```

### Example: HTTP Load Balancer

```yaml
metadata:
  name: example-http-lb
  namespace: example-namespace
spec:
  domains:
    - example.com
  http:
    port: 80
  default_route_pools:
    - pool:
        name: example-origin-pool
        namespace: example-namespace
```

### Example: Origin Pool

```yaml
metadata:
  name: example-origin-pool
  namespace: example-namespace
spec:
  origin_servers:
    - public_ip:
        ip: 192.168.1.100
  port: 8080
  healthcheck:
    - name: example-healthcheck
      namespace: example-namespace
```
