# Load Balancer Examples

Examples for configuring HTTP and TCP load balancers with xcsh.

## Minimum Configuration Requirements

### HTTP Load Balancer - Required Fields

| Field | Required | Description |
|-------|----------|-------------|
| `metadata.name` | ✅ | Unique name for the load balancer |
| `metadata.namespace` | ✅ | Target namespace |
| `spec.domains` | ✅ | List of domain names to serve |
| `spec.http` OR `spec.https` OR `spec.https_auto_cert` | ✅ | Listener type (choose one) |
| `spec.advertise_on_public_default_vip` | ✅ | Required for public internet access |

### Minimal HTTP Load Balancer Example

```yaml
metadata:
  name: my-http-lb
  namespace: shared
spec:
  domains:
    - myapp.example.com
  http:
    port: 80
  advertise_on_public_default_vip: {}
```

**Create command:**

```bash
xcsh cdn create http_loadbalancer -n shared -i http-lb.yaml
```

## CRUD Quick Reference

### Create

```bash
xcsh cdn create http_loadbalancer -n <namespace> -i <file.yaml>
```

### Read (List)

```bash
xcsh cdn list http_loadbalancer -n <namespace> --output-format json
```

### Read (Get)

```bash
xcsh cdn get http_loadbalancer <name> -n <namespace> --output-format yaml
```

### Update (Replace)

```bash
xcsh cdn replace http_loadbalancer <name> -n <namespace> -i <file.yaml>
```

### Delete

```bash
xcsh cdn delete http_loadbalancer <name> -n <namespace>
```

## HTTP Load Balancer

### Basic HTTP Load Balancer

Create a simple HTTP load balancer with a single origin pool.

**origin-pool.yaml:**

```yaml
metadata:
  name: example-origin-pool
  namespace: example-namespace
spec:
  origin_servers:
    - public_ip:
        ip: 192.168.1.100
  port: 8080
  loadbalancer_algorithm: ROUND_ROBIN
```

**http-lb.yaml:**

```yaml
metadata:
  name: example-http-lb
  namespace: example-namespace
spec:
  domains:
    - example.com
  http:
    port: 80
  advertise_on_public_default_vip: {}
  default_route_pools:
    - pool:
        name: example-origin-pool
        namespace: example-namespace
```

**Deploy:**

```bash
# Create origin pool
xcsh cdn create origin_pool -n example-namespace -i origin-pool.yaml

# Create load balancer
xcsh cdn create http_loadbalancer -n example-namespace -i http-lb.yaml

# Verify
xcsh cdn get http_loadbalancer example-http-lb -n example-namespace
```

### HTTPS Load Balancer

Add TLS termination to your load balancer with automatic certificate management.

**https-lb.yaml:**

```yaml
metadata:
  name: example-https-lb
  namespace: example-namespace
spec:
  domains:
    - example.com
  https_auto_cert:
    http_redirect: true
    add_hsts: true
  advertise_on_public_default_vip: {}
  default_route_pools:
    - pool:
        name: example-origin-pool
        namespace: example-namespace
```

### Load Balancer with WAF

Add Web Application Firewall protection.

**waf-lb.yaml:**

```yaml
metadata:
  name: example-waf-lb
  namespace: example-namespace
spec:
  domains:
    - example.com
  https_auto_cert:
    http_redirect: true
  advertise_on_public_default_vip: {}
  default_route_pools:
    - pool:
        name: example-origin-pool
        namespace: example-namespace
  app_firewall:
    name: example-waf-policy
    namespace: example-namespace
```

## TCP Load Balancer

### Basic TCP Load Balancer

**tcp-lb.yaml:**

```yaml
metadata:
  name: example-tcp-lb
  namespace: example-namespace
spec:
  listen_port: 3306
  origin_pools:
    - pool:
        name: example-db-pool
        namespace: example-namespace
```

**Deploy:**

```bash
xcsh virtual create tcp_loadbalancer -n example-namespace -i tcp-lb.yaml
```

## Health Checks

### HTTP Health Check

**healthcheck.yaml:**

```yaml
metadata:
  name: example-healthcheck
  namespace: example-namespace
spec:
  http_health_check:
    path: /health
    expected_status_codes:
      - "200"
  interval: 30
  timeout: 10
  unhealthy_threshold: 3
  healthy_threshold: 2
```

### TCP Health Check

**tcp-healthcheck.yaml:**

```yaml
metadata:
  name: example-tcp-healthcheck
  namespace: example-namespace
spec:
  tcp_health_check: {}
  interval: 15
  timeout: 5
```

## Origin Pool Options

### Multiple Origin Servers

```yaml
metadata:
  name: multi-origin-pool
  namespace: example-namespace
spec:
  origin_servers:
    - public_ip:
        ip: 192.168.1.100
    - public_ip:
        ip: 192.168.1.101
    - public_ip:
        ip: 192.168.1.102
  port: 8080
  loadbalancer_algorithm: ROUND_ROBIN
  healthcheck:
    - name: example-healthcheck
      namespace: example-namespace
```

### Origin Pool with DNS

```yaml
metadata:
  name: dns-origin-pool
  namespace: example-namespace
spec:
  origin_servers:
    - public_name:
        dns_name: backend.example.com
  port: 443
  use_tls:
    use_host_header_as_sni: {}
```

## Management Commands

### List Load Balancers

```bash
# List all HTTP load balancers
xcsh cdn list http_loadbalancer -n example-namespace

# List all TCP load balancers
xcsh virtual list tcp_loadbalancer -n example-namespace

# Output as JSON
xcsh cdn list http_loadbalancer -n example-namespace --output-format json
```

### Get Details

```bash
# Get as table (default)
xcsh cdn get http_loadbalancer example-lb -n example-namespace

# Get as YAML
xcsh cdn get http_loadbalancer example-lb -n example-namespace --output-format yaml

# Get as JSON
xcsh cdn get http_loadbalancer example-lb -n example-namespace --output-format json
```

### Update Load Balancer

```bash
# Export current config
xcsh cdn get http_loadbalancer example-lb -n example-namespace --output-format yaml > lb.yaml

# Edit lb.yaml...

# Apply changes (requires confirmation)
xcsh cdn replace http_loadbalancer example-lb -n example-namespace -i lb.yaml
```

### Delete Load Balancer

```bash
# With confirmation prompt
xcsh cdn delete http_loadbalancer example-lb -n example-namespace
```
