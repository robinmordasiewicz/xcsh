# Getting Started

This guide will help you get started with vesctl, the command-line interface for F5 Distributed Cloud.

## Prerequisites

Before installing vesctl, ensure you have:

- A terminal or command prompt
- Network access to download the binary
- (Optional) Go 1.22+ if building from source

## Installation Methods

Choose the installation method that works best for you:

| Method | Best For |
|--------|----------|
| [Quick Install](installation.md#quick-install) | Most users |
| [Manual Download](installation.md#manual-installation) | Custom locations |
| [Build from Source](installation.md#build-from-source) | Development |

## First Steps

1. **[Install vesctl](installation.md)** - Get the binary installed on your system
2. **[Configure Authentication](authentication.md)** - Set up your API credentials
3. **[Quick Start Tutorial](quickstart.md)** - Learn the basics with hands-on examples

## System Requirements

| Platform | Architecture | Status |
|----------|--------------|--------|
| Linux | amd64 | Supported |
| Linux | arm64 | Supported |
| macOS | amd64 | Supported |
| macOS | arm64 (Apple Silicon) | Supported |
| Windows | amd64 | Supported |

## Getting Help

If you encounter issues:

```bash
# View general help
vesctl --help

# View help for specific commands
vesctl configuration --help
vesctl configuration list --help
```

For additional support, visit the [GitHub repository](https://github.com/robinmordasiewicz/vesctl).
