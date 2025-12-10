# Build from Source

Build vesctl directly from source code.

## Prerequisites

Building from source requires:

| Requirement | Version | Check Command |
|-------------|---------|---------------|
| Go | go1.25.5 | `go version` |
| Git | any | `git --version` |


## Clone Repository

```bash
git clone https://github.com/robinmordasiewicz/vesctl.git
cd vesctl
```


## Build

```bash
go build -o vesctl .
```


## Verify Build

```bash
./vesctl version
```

Expected output shows version, commit hash, build timestamp, Go version, and platform.

## Install (Optional)

Move the binary to your PATH:

=== "User Install"

    ```bash
    mkdir -p ~/.local/bin
    mv vesctl ~/.local/bin/
    ```

=== "System Install"

    ```bash
    sudo mv vesctl /usr/local/bin/
    ```

## Build with Version Info

For release-quality builds with embedded version information:

```bash
go build -ldflags="-X github.com/robinmordasiewicz/vesctl/cmd.Version=dev \
  -X github.com/robinmordasiewicz/vesctl/cmd.GitCommit=$(git rev-parse --short HEAD) \
  -X github.com/robinmordasiewicz/vesctl/cmd.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o vesctl .
```

