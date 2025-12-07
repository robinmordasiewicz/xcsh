# vesctl Project Overview

## Purpose
vesctl (also known as f5xc CLI) is an open-source command-line interface for managing F5 Distributed Cloud (formerly Volterra) resources. It follows Azure CLI-style command patterns for intuitive resource management.

## Tech Stack
- **Language**: Go 1.22
- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **Configuration**: Viper (github.com/spf13/viper)
- **YAML Parsing**: gopkg.in/yaml.v3
- **Certificate Handling**: go-pkcs12 for P12 bundle authentication
- **Build System**: Make + GoReleaser
- **Documentation**: MkDocs with Material theme

## Key Features
- CRUD operations for F5 XC resources (HTTP/TCP load balancers, origin pools, etc.)
- Multiple output formats: JSON, YAML, table, TSV
- P12 bundle or certificate/key authentication
- Shell completion for bash, zsh, fish, PowerShell
- Cross-platform builds (Linux, macOS, Windows)

## Authentication Methods
1. **P12 Bundle** (recommended): Downloaded from F5 XC Console
2. **Certificate + Key**: Extracted from P12 or separate files
3. **Environment Variables**: VES_P12_PASSWORD, F5XC_SERVER_URLS, etc.

## Configuration
- Default config file: `~/.vesconfig`
- Supports YAML format with server-urls, cert, key, p12-bundle options
