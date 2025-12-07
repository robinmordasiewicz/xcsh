# Cloud Site Examples

Examples for deploying F5 Distributed Cloud sites in AWS and Azure.

## AWS VPC Site

### Basic AWS VPC Site

Deploy an F5 XC site in an AWS VPC.

**aws-vpc-site.yaml:**

```yaml
metadata:
  name: my-aws-site
  namespace: system
spec:
  aws:
    region: us-east-1
  vpc:
    vpc_id: vpc-12345678
  ingress_egress_gw:
    aws_certified_hw: aws-byol-voltmesh
    instance_type: t3.xlarge
  nodes:
    - aws_az_name: us-east-1a
      reserved_inside_subnet:
        existing_subnet_id: subnet-inside-123
```

**Deploy:**

```bash
# Create cloud credentials first
vesctl configuration create cloud_credentials -i aws-creds.yaml -n system

# Create site
vesctl configuration create aws_vpc_site -i aws-vpc-site.yaml -n system

# Check status
vesctl configuration get aws_vpc_site my-aws-site -n system
```

### AWS Site with Multiple Nodes

**multi-az-site.yaml:**

```yaml
metadata:
  name: my-multi-az-site
  namespace: system
spec:
  aws:
    region: us-east-1
  vpc:
    vpc_id: vpc-12345678
  ingress_egress_gw:
    aws_certified_hw: aws-byol-voltmesh
    instance_type: t3.xlarge
  nodes:
    - aws_az_name: us-east-1a
      reserved_inside_subnet:
        existing_subnet_id: subnet-inside-1a
    - aws_az_name: us-east-1b
      reserved_inside_subnet:
        existing_subnet_id: subnet-inside-1b
    - aws_az_name: us-east-1c
      reserved_inside_subnet:
        existing_subnet_id: subnet-inside-1c
```

## Azure VNet Site

### Basic Azure VNet Site

Deploy an F5 XC site in an Azure VNet.

**azure-vnet-site.yaml:**

```yaml
metadata:
  name: my-azure-site
  namespace: system
spec:
  azure:
    region: eastus
  resource_group: my-resource-group
  vnet:
    vnet_name: my-vnet
  ingress_egress_gw:
    azure_certified_hw: azure-byol-voltmesh
    instance_type: Standard_D3_v2
  nodes:
    - azure_az: "1"
      inside_subnet:
        subnet_name: inside-subnet
```

**Deploy:**

```bash
# Create cloud credentials
vesctl configuration create cloud_credentials -i azure-creds.yaml -n system

# Create site
vesctl configuration create azure_vnet_site -i azure-vnet-site.yaml -n system
```

## Cloud Credentials

### AWS Credentials

**aws-creds.yaml:**

```yaml
metadata:
  name: my-aws-creds
  namespace: system
spec:
  aws_secret_key:
    access_key: AKIAIOSFODNN7EXAMPLE
    secret_key:
      blindfold_secret_info:
        location: string:///base64-encoded-secret
```

### Azure Credentials

**azure-creds.yaml:**

```yaml
metadata:
  name: my-azure-creds
  namespace: system
spec:
  azure_client_secret:
    subscription_id: your-subscription-id
    tenant_id: your-tenant-id
    client_id: your-client-id
    client_secret:
      blindfold_secret_info:
        location: string:///base64-encoded-secret
```

## Site Management Commands

### List Sites

```bash
# List AWS sites
vesctl configuration list aws_vpc_site -n system

# List Azure sites
vesctl configuration list azure_vnet_site -n system
```

### Get Site Details

```bash
# Get AWS site
vesctl configuration get aws_vpc_site my-aws-site -n system --outfmt yaml

# Get Azure site
vesctl configuration get azure_vnet_site my-azure-site -n system --outfmt yaml
```

### Site Status

```bash
# Check site registration status
vesctl site status my-aws-site -n system
```

### Delete Site

```bash
# Delete AWS site
vesctl configuration delete aws_vpc_site my-aws-site -n system --yes

# Delete Azure site
vesctl configuration delete azure_vnet_site my-azure-site -n system --yes
```

## Troubleshooting

### Check Site Status

```bash
vesctl configuration get aws_vpc_site my-site -n system --outfmt json | jq '.status'
```

### List All Cloud Resources

```bash
# List all cloud credentials
vesctl configuration list cloud_credentials -n system

# List all sites
vesctl configuration list aws_vpc_site -n system
vesctl configuration list azure_vnet_site -n system
```

### Debug Mode

```bash
vesctl --debug configuration get aws_vpc_site my-site -n system
```
