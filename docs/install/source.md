# Build from Source

Build vesctl directly from source code.

## Prerequisites

Building from source requires:

| Requirement | Version | Check Command |
|-------------|---------|---------------|
| Go | go1.23.4 | `go version` |
| Git | any | `git --version` |

**Verify prerequisites:**

```text
$ go version
go version go1.23.4 darwin/arm64

$ git --version
git version 2.44.0
```

## Clone Repository

```bash
git clone https://github.com/robinmordasiewicz/vesctl.git
cd vesctl
```

**Example output:**

```text
$ git clone https://github.com/robinmordasiewicz/vesctl.git
Cloning into 'vesctl'...
remote: Enumerating objects: 1234, done.
remote: Counting objects: 100% (234/234), done.
remote: Compressing objects: 100% (156/156), done.
remote: Total 1234 (delta 89), reused 180 (delta 72), pack-reused 1000
Receiving objects: 100% (1234/1234), 256.00 KiB | 2.00 MiB/s, done.
Resolving deltas: 100% (678/678), done.
$ cd vesctl
```

## Build

```bash
go build -o vesctl .
```

Git commit and build date are automatically embedded when building from a git repository.

**Example output:**

```text
$ go build -o vesctl .
(no output - build successful)
```

## Verify Build

```bash
./vesctl version
```

**Example output:**

```text
$ ./vesctl version
vesctl version dev
  commit:   935d038
  built:    2025-12-10T03:35:08Z
  go:       go1.23.4
  platform: darwin/arm64
```

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

## Build with Custom Version

For release-quality builds with a custom version string and stripped binaries:

```bash
go build -ldflags="-s -w -X github.com/robinmordasiewicz/vesctl/cmd.Version=1.0.0" \
  -o vesctl .
```

**Example output:**

```text
$ ./vesctl version
vesctl version 1.0.0
  commit:   935d038
  built:    2025-12-10T03:35:08Z
  go:       go1.23.4
  platform: darwin/arm64
```

The `-s -w` flags strip debug symbols for smaller binaries.
