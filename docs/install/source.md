# Build from Source

Build f5xcctl from source for development or to get the latest unreleased features.

## Prerequisites

- Go 1.21 or later
- Git

### Install Go

=== "macOS"

    ```bash
    brew install go
    ```

=== "Linux"

    ```bash
    # Download from https://go.dev/dl/
    wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    ```

=== "Windows"

    Download and run the installer from [go.dev/dl](https://go.dev/dl/)

## Clone and Build

```bash
# Clone the repository
git clone https://github.com/robinmordasiewicz/f5xcctl.git
cd f5xcctl

# Build the binary
go build -o f5xcctl .
```

## Install Locally

Move the binary to a directory in your PATH:

```bash
# macOS/Linux
sudo mv f5xcctl /usr/local/bin/

# Or for user-local installation
mv f5xcctl ~/.local/bin/
```

## Development Setup

For active development:

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Run the CLI directly
go run . --help

# Build with version information
go build -ldflags "-X main.version=dev" -o f5xcctl .
```

## Generate Completions

After building, generate shell completions:

```bash
# Zsh
./f5xcctl completion zsh > ~/.zsh/completions/_f5xcctl

# Bash
./f5xcctl completion bash > ~/.local/share/bash-completion/completions/f5xcctl

# Fish
./f5xcctl completion fish > ~/.config/fish/completions/f5xcctl.fish
```

## Verify Build

```bash
./f5xcctl version
```
