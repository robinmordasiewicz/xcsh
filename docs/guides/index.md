# Examples

Real-world examples for common xcsh use cases.

## Contents

| Example | Description |
|---------|-------------|
| [Load Balancers](load-balancer.md) | HTTP and TCP load balancer configuration |
| [Cloud Sites](cloud-sites.md) | AWS and Azure site deployment |

## Quick Examples

### List All Namespaces

```bash
xcsh identity list namespace
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
xcsh load_balancer create http_loadbalancer -i lb.yaml
```

### Get Resource as YAML

```bash
xcsh <domain> get namespace example-namespace --outfmt yaml
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
xcsh load_balancer delete http_loadbalancer example-lb -n example-namespace --yes
```

## Common Workflows

### Deploy HTTP Load Balancer

1. Create an origin pool:

```bash
xcsh load_balancer create origin_pool -i origin-pool.yaml
```

2. Create a health check:

```bash
xcsh <domain> create healthcheck -i healthcheck.yaml
```

3. Create the load balancer:

```bash
xcsh load_balancer create http_loadbalancer -i lb.yaml
```

4. Verify deployment:

```bash
xcsh load_balancer get http_loadbalancer example-lb -n example-namespace
```

### Update Configuration

1. Export current configuration:

```bash
xcsh load_balancer get http_loadbalancer example-lb -n example-namespace --outfmt yaml > lb.yaml
```

2. Edit the file as needed

3. Apply changes:

```bash
xcsh load_balancer replace http_loadbalancer -i lb.yaml
```

### Cleanup Resources

```bash
# Delete load balancer
xcsh load_balancer delete http_loadbalancer example-lb -n example-namespace --yes

# Delete origin pool
xcsh <domain> delete origin_pool example-pool -n example-namespace --yes

# Delete health check
xcsh <domain> delete healthcheck example-hc -n example-namespace --yes
```
