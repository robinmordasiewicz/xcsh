# API Endpoint Commands

The `vesctl api-endpoint` command group manages API endpoint discovery and security.

## Overview

```text
vesctl api-endpoint --help
```

API endpoint commands help you discover, monitor, and manage API endpoints in your
F5 Distributed Cloud tenant. This is essential for API security, governance,
and documentation.

## Why API Discovery Matters

- **Security visibility** - Identify all APIs exposed through your load balancers
- **Shadow API detection** - Discover undocumented or forgotten endpoints
- **Compliance** - Ensure all APIs are properly secured and monitored
- **Documentation** - Generate accurate API inventories

## Common Operations

### List Discovered Endpoints

View all discovered API endpoints in a namespace:

```bash
vesctl api-endpoint list -n example-namespace
```

**Example Output:**

```text
NAME                    METHOD    PATH              LOAD BALANCER
api-users-get           GET       /api/users        example-http-lb
api-users-post          POST      /api/users        example-http-lb
api-orders-get          GET       /api/orders       example-http-lb
api-health              GET       /health           example-http-lb
```

### Get Endpoint Details

Retrieve detailed information about a specific endpoint:

```bash
vesctl api-endpoint get api-users-get -n example-namespace
```

**With YAML output:**

```bash
vesctl api-endpoint get api-users-get -n example-namespace --outfmt yaml
```

### Filter by Load Balancer

List endpoints for a specific HTTP load balancer:

```bash
vesctl api-endpoint list -n example-namespace --lb example-http-lb
```

## Discovery Workflow

### Step 1: Enable API Discovery

API discovery is configured on your HTTP load balancer. Ensure your load balancer has API discovery enabled:

```yaml
metadata:
  name: example-http-lb
  namespace: example-namespace
spec:
  domains:
    - api.example.com
  api_discovery:
    enable: true
    discovered_api_settings:
      purge_duration_for_inactive_endpoints: 168h
```

### Step 2: Deploy the Load Balancer

```bash
vesctl configuration create http_loadbalancer -i lb-with-discovery.yaml
```

### Step 3: Generate Traffic

Send traffic through the load balancer to trigger endpoint discovery. The system
analyzes request patterns to identify API endpoints.

### Step 4: Review Discovered Endpoints

```bash
# List all discovered endpoints
vesctl api-endpoint list -n example-namespace

# Export to JSON for analysis
vesctl api-endpoint list -n example-namespace --outfmt json > endpoints.json
```

### Step 5: Create API Definition

Use discovered endpoints to create a formal API definition:

```bash
# Get endpoint details
vesctl api-endpoint get api-users-get -n example-namespace --outfmt yaml > endpoint.yaml

# Create API definition
vesctl configuration create api_definition -i api-def.yaml
```

## API Definition Management

For comprehensive API management, use configuration commands with API definitions:

### List API Definitions

```bash
vesctl configuration list api_definition -n example-namespace
```

### Create API Definition

**api-definition.yaml:**

```yaml
metadata:
  name: example-api-definition
  namespace: example-namespace
spec:
  swagger:
    openapi: "3.0.0"
    info:
      title: "My API"
      version: "1.0.0"
    paths:
      /api/users:
        get:
          summary: "List users"
        post:
          summary: "Create user"
```

```bash
vesctl configuration create api_definition -i api-definition.yaml
```

### Get API Definition

```bash
vesctl configuration get api_definition example-api-definition -n example-namespace --outfmt yaml
```

## Integration with Security Policies

Discovered endpoints can be protected with service policies:

### Create Endpoint-Specific Policy

```yaml
metadata:
  name: protect-user-api
  namespace: example-namespace
spec:
  algo: FIRST_MATCH
  rules:
    - metadata:
        name: block-unauthorized
      spec:
        action: DENY
        api_endpoint:
          methods:
            - POST
            - DELETE
          path:
            prefix: /api/users
```

```bash
vesctl configuration create service_policy -i policy.yaml
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `-n, --namespace` | `-n` | Target namespace |
| `--lb` | | Filter by load balancer name |
| `--outfmt` | `-o` | Output format (json, yaml, table) |

## Best Practices

1. **Enable discovery early** - Configure API discovery when creating load balancers
2. **Regular reviews** - Periodically review discovered endpoints for shadow APIs
3. **Document findings** - Export discovered endpoints and create formal API definitions
4. **Apply security** - Use service policies to protect discovered endpoints
5. **Monitor changes** - Track new endpoints appearing in your environment

## Troubleshooting

### No Endpoints Discovered

If no endpoints appear:

1. Verify API discovery is enabled on the load balancer
2. Confirm traffic is flowing through the load balancer
3. Check the namespace is correct
4. Allow time for discovery (endpoints appear after traffic analysis)

```bash
# Check load balancer configuration
vesctl configuration get http_loadbalancer example-lb -n example-namespace --outfmt yaml | grep -A5 api_discovery
```

### Debug Mode

Enable verbose output for troubleshooting:

```bash
vesctl --debug api-endpoint list -n example-namespace
```

## Related Commands

- [configuration](configuration.md) - Manage API definitions and policies
- [request](request.md) - Low-level API access
- [Load Balancer Examples](../examples/load-balancer.md) - HTTP load balancer configuration
