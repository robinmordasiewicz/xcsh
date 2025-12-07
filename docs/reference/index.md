# Reference

This section contains detailed reference documentation for vesctl.

## Contents

| Reference | Description |
|-----------|-------------|
| [Global Flags](global-flags.md) | Command-line flags available for all commands |
| [Environment Variables](environment-variables.md) | Environment variables for configuration |
| [Resources](resources.md) | Supported F5 XC resource types |

## Quick Reference

### Authentication

```bash
# P12 bundle
export VES_P12_PASSWORD="your-password"

# Or certificate/key in config
# cert: /path/to/cert.pem
# key: /path/to/key.pem
```

### Common Commands

```bash
vesctl configuration list <type>           # List resources
vesctl configuration get <type> <name>     # Get resource
vesctl configuration create <type> -i f    # Create from file
vesctl configuration delete <type> <name>  # Delete resource
```

### Output Formats

```bash
--outfmt json   # JSON output
--outfmt yaml   # YAML output
--outfmt table  # Table output (default)
```
