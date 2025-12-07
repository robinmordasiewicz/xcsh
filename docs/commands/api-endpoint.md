# API Endpoint Commands

The `vesctl api-endpoint` command group manages API endpoint discovery.

## Overview

```text
vesctl api-endpoint --help
```

## Description

API endpoint commands help discover and manage API endpoints in your F5 XC tenant. This is useful for API security and management features.

## Common Operations

### List Discovered Endpoints

```bash
vesctl api-endpoint list -n <namespace>
```

### Get Endpoint Details

```bash
vesctl api-endpoint get <endpoint-name> -n <namespace>
```

## Related Resources

API endpoint discovery works with:

- HTTP Load Balancers
- API Definitions
- Service Policies

## Configuration via Resources

For full API definition management, use configuration commands:

```bash
# List API definitions
vesctl configuration list api_definition -n my-namespace

# Create API definition
vesctl configuration create api_definition -i api-def.yaml
```
