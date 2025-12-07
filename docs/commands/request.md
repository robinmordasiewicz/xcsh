# Request Commands

The `vesctl request` command group provides low-level API access for advanced operations.

## Overview

```text
vesctl request --help
```

The request commands enable direct interaction with the F5 Distributed Cloud API,
allowing you to execute custom RPC calls and manage secrets programmatically.

## Subcommands

### RPC Request

Execute custom API RPC calls for operations not covered by standard commands:

```bash
vesctl request rpc <method> <path> [flags]
```

**Supported Methods:**

| Method | Description |
|--------|-------------|
| `GET` | Retrieve resources or data |
| `POST` | Create resources or execute actions |
| `PUT` | Update existing resources |
| `DELETE` | Remove resources |

## Common Use Cases

### Retrieve Tenant Information

```bash
vesctl request rpc GET /api/web/tenant
```

**Example Output:**

```json
{
  "name": "my-tenant",
  "uid": "abc123-def456",
  "tenant_type": "CUSTOMER"
}
```

### List Available API Endpoints

```bash
vesctl request rpc GET /api
```

### Query Namespace Details

```bash
vesctl request rpc GET /api/web/namespaces/my-namespace
```

### Execute Custom Queries

For complex queries, use a request body file:

**query.json:**

```json
{
  "namespace": "my-namespace",
  "label_selector": "app=web"
}
```

```bash
vesctl request rpc POST /api/custom/query -i query.json
```

## Advanced Examples

### GET with Output Formatting

```bash
# JSON output
vesctl request rpc GET /api/web/namespaces --outfmt json

# YAML output
vesctl request rpc GET /api/web/namespaces --outfmt yaml
```

### POST with Request Body

Create a custom resource using a JSON payload:

**payload.json:**

```json
{
  "metadata": {
    "name": "my-resource",
    "namespace": "my-namespace"
  },
  "spec": {
    "description": "Created via RPC"
  }
}
```

```bash
vesctl request rpc POST /api/config/namespaces/my-namespace/resources -i payload.json
```

### Debug API Calls

Enable verbose output to troubleshoot API interactions:

```bash
vesctl --debug request rpc GET /api/web/tenant
```

This displays:

- Request URL and headers
- Response status code
- Response body

## Secrets Management

Manage secrets through the request interface:

```bash
vesctl request secrets [subcommand] [flags]
```

### List Secrets

```bash
vesctl request secrets list -n my-namespace
```

### Get Secret Details

```bash
vesctl request secrets get my-secret -n my-namespace
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `-i, --input-file` | `-i` | Request body from file (JSON or YAML) |
| `--outfmt` | `-o` | Output format (json, yaml, table) |
| `-n, --namespace` | `-n` | Target namespace |

## Best Practices

1. **Use standard commands first** - The `configuration` commands handle most operations
2. **Test with GET** - Verify endpoints before making changes
3. **Enable debug mode** - Use `--debug` when troubleshooting
4. **Save responses** - Redirect output to files for complex operations

```bash
vesctl request rpc GET /api/web/namespaces --outfmt json > namespaces.json
```

## Related Commands

- [configuration](configuration.md) - Standard resource management
- [api-endpoint](api-endpoint.md) - API endpoint discovery
