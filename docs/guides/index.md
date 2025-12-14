# Examples

Real-world examples for common f5xcctl use cases.

## Contents

| Example | Description |
|---------|-------------|
| [Load Balancers](load-balancer.md) | HTTP and TCP load balancer configuration |
| [Cloud Sites](cloud-sites.md) | AWS and Azure site deployment |

## Quick Examples

### List All Namespaces

```bash
f5xcctl configuration list namespace
```

**Output:**

```text
NAME           DESCRIPTION
default        Default namespace
example-namespace   My application namespace
system         System namespace
```

### Create Resource from File

```bash
f5xcctl configuration create http_loadbalancer -i lb.yaml
```

### Get Resource as YAML

```bash
f5xcctl configuration get namespace example-namespace --outfmt yaml
```

**Output:**

```yaml
metadata:
  name: example-namespace
  uid: abc123
spec:
  description: My application namespace
```

### Delete Resource

```bash
f5xcctl configuration delete http_loadbalancer example-lb -n example-namespace --yes
```

## Common Workflows

### Deploy HTTP Load Balancer

1. Create an origin pool:

```bash
f5xcctl configuration create origin_pool -i origin-pool.yaml
```

2. Create a health check:

```bash
f5xcctl configuration create healthcheck -i healthcheck.yaml
```

3. Create the load balancer:

```bash
f5xcctl configuration create http_loadbalancer -i lb.yaml
```

4. Verify deployment:

```bash
f5xcctl configuration get http_loadbalancer example-lb -n example-namespace
```

### Update Configuration

1. Export current configuration:

```bash
f5xcctl configuration get http_loadbalancer example-lb -n example-namespace --outfmt yaml > lb.yaml
```

2. Edit the file as needed

3. Apply changes:

```bash
f5xcctl configuration replace http_loadbalancer -i lb.yaml
```

### Cleanup Resources

```bash
# Delete load balancer
f5xcctl configuration delete http_loadbalancer example-lb -n example-namespace --yes

# Delete origin pool
f5xcctl configuration delete origin_pool example-pool -n example-namespace --yes

# Delete health check
f5xcctl configuration delete healthcheck example-hc -n example-namespace --yes
```
