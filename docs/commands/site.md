# Site Commands

The `vesctl site` command group manages cloud and edge sites.

## Overview

```text
vesctl site --help
```

## Site Types

vesctl supports managing various site types:

| Site Type | Description |
|-----------|-------------|
| AWS VPC Site | Sites deployed in AWS VPCs |
| Azure VNet Site | Sites deployed in Azure VNets |
| GCP VPC Site | Sites deployed in Google Cloud VPCs |
| Edge Sites | Physical or virtual edge deployments |

## Common Operations

### List Sites

```bash
# List all sites
vesctl site list

# List sites in a namespace
vesctl site list -n system
```

### Get Site Details

```bash
vesctl site get <site-name> -n <namespace>
```

### Site Status

Check the status and health of deployed sites:

```bash
vesctl site status <site-name> -n <namespace>
```

## Cloud Site Management

### AWS VPC Sites

```bash
# Get AWS VPC site help
vesctl site aws-vpc --help

# List AWS VPC sites
vesctl configuration list aws_vpc_site -n system
```

### Azure VNet Sites

```bash
# Get Azure VNet site help
vesctl site azure-vnet --help

# List Azure VNet sites
vesctl configuration list azure_vnet_site -n system
```

## Related Commands

For detailed site resource management, use the `configuration` commands:

```bash
# Create AWS VPC site
vesctl configuration create aws_vpc_site -i aws-site.yaml -n system

# Get site details
vesctl configuration get aws_vpc_site my-aws-site -n system --outfmt yaml
```
