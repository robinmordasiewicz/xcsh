# Quick Start

Get up and running with vesctl in minutes.

## Step 1: Install vesctl

```bash
curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | sh
```

## Step 2: Configure Authentication

### Option A: Interactive Setup

Run the configure command to set up your credentials interactively:

```bash
vesctl configure
```

### Option B: Configuration File

Create `~/.vesconfig`:

```yaml
server-urls:
  - https://your-tenant.console.ves.volterra.io/api
p12-bundle: ~/.vesconfig/api-creds.p12
```

Set the P12 password:

```bash
export VES_P12_PASSWORD="your-p12-password"
```

## Step 3: Verify Connection

Test your configuration by listing namespaces:

```bash
vesctl configuration list namespace
```

You should see a list of available namespaces in your tenant.

## Step 4: Explore Resources

### List All Namespaces

```bash
vesctl configuration list namespace
```

### Get Namespace Details

```bash
vesctl configuration get namespace default --outfmt yaml
```

### List HTTP Load Balancers

```bash
vesctl configuration list http_loadbalancer -n my-namespace
```

## Common Commands

| Task | Command |
|------|---------|
| List namespaces | `vesctl configuration list namespace` |
| Get resource details | `vesctl configuration get <type> <name> -n <namespace>` |
| Create resource | `vesctl configuration create <type> -i file.yaml` |
| Delete resource | `vesctl configuration delete <type> <name> -n <namespace>` |
| Show version | `vesctl version` |

## Output Formats

Control output format with the `--outfmt` flag:

```bash
# Default table format
vesctl configuration list namespace

# JSON output
vesctl configuration list namespace --outfmt json

# YAML output
vesctl configuration list namespace --outfmt yaml
```

## Next Steps

- [Authentication Details](authentication.md) - Learn about all authentication methods
- [Command Reference](../commands/index.md) - Full command documentation
- [Examples](../examples/index.md) - Real-world usage examples
