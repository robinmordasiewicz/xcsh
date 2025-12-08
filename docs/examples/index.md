# Examples

Real-world examples for common vesctl use cases.

## Contents

| Example | Description |
|---------|-------------|
| [Load Balancers](load-balancer.md) | HTTP and TCP load balancer configuration |
| [Cloud Sites](cloud-sites.md) | AWS and Azure site deployment |

## Quick Examples

### List All Namespaces

```bash
vesctl configuration list namespace
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
vesctl configuration create http_loadbalancer -i lb.yaml
```

### Get Resource as YAML

```bash
vesctl configuration get namespace example-namespace --outfmt yaml
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
vesctl configuration delete http_loadbalancer example-lb -n example-namespace --yes
```

## Common Workflows

### Deploy HTTP Load Balancer

1. Create an origin pool:

```bash
vesctl configuration create origin_pool -i origin-pool.yaml
```

2. Create a health check:

```bash
vesctl configuration create healthcheck -i healthcheck.yaml
```

3. Create the load balancer:

```bash
vesctl configuration create http_loadbalancer -i lb.yaml
```

4. Verify deployment:

```bash
vesctl configuration get http_loadbalancer example-lb -n example-namespace
```

### Update Configuration

1. Export current configuration:

```bash
vesctl configuration get http_loadbalancer example-lb -n example-namespace --outfmt yaml > lb.yaml
```

2. Edit the file as needed

3. Apply changes:

```bash
vesctl configuration replace http_loadbalancer -i lb.yaml
```

### Cleanup Resources

```bash
# Delete load balancer
vesctl configuration delete http_loadbalancer example-lb -n example-namespace --yes

# Delete origin pool
vesctl configuration delete origin_pool example-pool -n example-namespace --yes

# Delete health check
vesctl configuration delete healthcheck example-hc -n example-namespace --yes
```
