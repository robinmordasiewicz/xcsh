# Load Balancer Examples

Examples for configuring HTTP and TCP load balancers.

## HTTP Load Balancer

### Basic HTTP Load Balancer

Create a simple HTTP load balancer with a single origin pool.

**origin-pool.yaml:**

```yaml
metadata:
  name: my-origin-pool
  namespace: my-namespace
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

**Deploy:**

```bash
# Create origin pool
vesctl configuration create origin_pool -i origin-pool.yaml

# Create load balancer
vesctl configuration create http_loadbalancer -i http-lb.yaml

# Verify
vesctl configuration get http_loadbalancer my-http-lb -n my-namespace
```

### HTTPS Load Balancer

Add TLS termination to your load balancer.

**https-lb.yaml:**

```yaml
metadata:
  name: my-https-lb
  namespace: my-namespace
spec:
  domains:
    - example.com
  https_auto_cert:
    http_redirect: true
    add_hsts: true
  default_route_pools:
    - pool:
        name: my-origin-pool
        namespace: my-namespace
```

### Load Balancer with WAF

Add Web Application Firewall protection.

**waf-lb.yaml:**

```yaml
metadata:
  name: my-waf-lb
  namespace: my-namespace
spec:
  domains:
    - example.com
  https_auto_cert:
    http_redirect: true
  default_route_pools:
    - pool:
        name: my-origin-pool
        namespace: my-namespace
  app_firewall:
    name: my-waf-policy
    namespace: my-namespace
```

## TCP Load Balancer

### Basic TCP Load Balancer

**tcp-lb.yaml:**

```yaml
metadata:
  name: my-tcp-lb
  namespace: my-namespace
spec:
  listen_port: 3306
  origin_pools:
    - pool:
        name: my-db-pool
        namespace: my-namespace
```

**Deploy:**

```bash
vesctl configuration create tcp_loadbalancer -i tcp-lb.yaml
```

## Health Checks

### HTTP Health Check

**healthcheck.yaml:**

```yaml
metadata:
  name: my-healthcheck
  namespace: my-namespace
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
  name: my-tcp-healthcheck
  namespace: my-namespace
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
  namespace: my-namespace
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
    - name: my-healthcheck
      namespace: my-namespace
```

### Origin Pool with DNS

```yaml
metadata:
  name: dns-origin-pool
  namespace: my-namespace
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
vesctl configuration list http_loadbalancer -n my-namespace

# List all TCP load balancers
vesctl configuration list tcp_loadbalancer -n my-namespace
```

### Get Details

```bash
# Get as table
vesctl configuration get http_loadbalancer my-lb -n my-namespace

# Get as YAML
vesctl configuration get http_loadbalancer my-lb -n my-namespace --outfmt yaml

# Get as JSON
vesctl configuration get http_loadbalancer my-lb -n my-namespace --outfmt json
```

### Update Load Balancer

```bash
# Export current config
vesctl configuration get http_loadbalancer my-lb -n my-namespace --outfmt yaml > lb.yaml

# Edit lb.yaml...

# Apply changes
vesctl configuration replace http_loadbalancer -i lb.yaml
```

### Delete Load Balancer

```bash
# With confirmation
vesctl configuration delete http_loadbalancer my-lb -n my-namespace

# Skip confirmation
vesctl configuration delete http_loadbalancer my-lb -n my-namespace --yes
```
