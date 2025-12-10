#!/bin/sh
# vesctl install script
# Usage: curl -fsSL https://robinmordasiewicz.github.io/vesctl/install.sh | sh
#
# Environment variables:
#   VES_VERSION      - Specific version to install (default: latest)
#   VES_INSTALL_DIR  - Installation directory (default: /usr/local/bin)
#   VES_NO_SUDO      - Skip sudo if set to any value
#   VES_NO_VERIFY    - Skip checksum verification if set

set -eu

# Configuration
GITHUB_REPO="robinmordasiewicz/vesctl"
GITHUB_API="https://api.github.com/repos/${GITHUB_REPO}/releases/latest"
GITHUB_RELEASES="https://github.com/${GITHUB_REPO}/releases/download"
DEFAULT_INSTALL_DIR="/usr/local/bin"
USER_INSTALL_DIR="$HOME/.local/bin"
BINARY_NAME="vesctl"

# Colors for output (check if terminal supports colors)
if [ -t 1 ] && [ -n "$(tput colors 2>/dev/null)" ] && [ "$(tput colors)" -ge 8 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    BLUE='\033[0;34m'
    CYAN='\033[0;36m'
    NC='\033[0m' # No Color
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    CYAN=''
    NC=''
fi

# Temporary directory for downloads
TEMP_DIR=""

# ============================================
# Utility Functions
# ============================================

status() {
    printf "${BLUE}==>${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}==>${NC} %s\n" "$1"
}

warning() {
    printf "${YELLOW}Warning:${NC} %s\n" "$1" >&2
}

error() {
    printf "${RED}Error:${NC} %s\n" "$1" >&2
    exit 1
}

info() {
    printf "${CYAN}   ${NC} %s\n" "$1"
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

cleanup() {
    if [ -n "$TEMP_DIR" ] && [ -d "$TEMP_DIR" ]; then
        rm -rf "$TEMP_DIR"
    fi
}

# Set up cleanup trap
trap cleanup EXIT

# ============================================
# HTTP Client Abstraction
# ============================================

# Determine available HTTP client
detect_http_client() {
    if command_exists curl; then
        echo "curl"
    elif command_exists wget; then
        echo "wget"
    else
        echo "none"
    fi
}

# Download a file using available HTTP client
# Usage: http_download <url> <output_file>
http_download() {
    URL="$1"
    OUTPUT="$2"
    HTTP_CLIENT=$(detect_http_client)

    case "$HTTP_CLIENT" in
        curl)
            curl -fsSL --retry 3 --retry-delay 2 -o "$OUTPUT" "$URL"
            ;;
        wget)
            wget -q --tries=3 --waitretry=2 -O "$OUTPUT" "$URL"
            ;;
        *)
            error "Neither curl nor wget found. Please install one of them and try again."
            ;;
    esac
}

# Fetch content from URL
# Usage: http_get <url>
http_get() {
    URL="$1"
    HTTP_CLIENT=$(detect_http_client)

    case "$HTTP_CLIENT" in
        curl)
            curl -fsSL --retry 3 --retry-delay 2 "$URL" 2>/dev/null
            ;;
        wget)
            wget -q --tries=3 --waitretry=2 -O - "$URL" 2>/dev/null
            ;;
        *)
            error "Neither curl nor wget found. Please install one of them and try again."
            ;;
    esac
}

# ============================================
# Platform Detection
# ============================================

detect_os() {
    OS="$(uname -s)"
    case "$OS" in
        Linux*)
            echo "linux"
            ;;
        Darwin*)
            echo "darwin"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            error "Windows detected via shell. Please use PowerShell or download manually:
  https://github.com/${GITHUB_REPO}/releases/latest

For Windows (amd64):
  Download vesctl_VERSION_windows_amd64.zip

For Windows (arm64):
  Download vesctl_VERSION_windows_arm64.zip"
            ;;
        FreeBSD*)
            error "FreeBSD is not currently supported. Pre-built binaries are available for:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64, arm64)"
            ;;
        *)
            error "Unsupported operating system: $OS

Pre-built binaries are available for:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64, arm64)"
            ;;
    esac
}

detect_arch() {
    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        armv7l|armv6l)
            error "32-bit ARM is not supported. Please use a 64-bit system.

Supported architectures:
  - amd64 (x86_64)
  - arm64 (aarch64)"
            ;;
        i386|i686)
            error "32-bit x86 is not supported. Please use a 64-bit system.

Supported architectures:
  - amd64 (x86_64)
  - arm64 (aarch64)"
            ;;
        *)
            error "Unsupported architecture: $ARCH

Supported architectures:
  - amd64 (x86_64)
  - arm64 (aarch64)"
            ;;
    esac
}

get_os_display_name() {
    OS="$1"
    case "$OS" in
        linux)  echo "Linux" ;;
        darwin) echo "macOS" ;;
        *)      echo "$OS" ;;
    esac
}

get_arch_display_name() {
    ARCH="$1"
    case "$ARCH" in
        amd64) echo "x86_64 (Intel/AMD)" ;;
        arm64) echo "ARM64 (Apple Silicon/ARM)" ;;
        *)     echo "$ARCH" ;;
    esac
}

# ============================================
# Privilege Handling
# ============================================

# Returns: "system:", "system:sudo", "user:", "custom:", or "custom:sudo"
# Format: "<strategy_type>:<sudo_command>"
determine_install_strategy() {
    REQUESTED_DIR="${VES_INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"

    # If user explicitly set install dir, respect it
    if [ -n "${VES_INSTALL_DIR:-}" ]; then
        # Create directory if it doesn't exist and parent is writable
        if [ ! -d "$REQUESTED_DIR" ]; then
            if mkdir -p "$REQUESTED_DIR" 2>/dev/null; then
                echo "custom:"
                return
            fi
        elif [ -w "$REQUESTED_DIR" ]; then
            echo "custom:"
            return
        fi
        if [ "$(id -u)" -eq 0 ]; then
            echo "custom:"
            return
        fi
        if [ -z "${VES_NO_SUDO:-}" ] && command_exists sudo; then
            echo "custom:sudo"
            return
        fi
        error "Cannot write to $REQUESTED_DIR and sudo is not available.

Try one of:
  - Set VES_INSTALL_DIR to a writable location:
    VES_INSTALL_DIR=\$HOME/.local/bin sh install.sh

  - Run as root:
    sudo sh install.sh"
    fi

    # Default behavior: try /usr/local/bin with sudo, fall back to ~/.local/bin
    if [ -w "$DEFAULT_INSTALL_DIR" ] || [ "$(id -u)" -eq 0 ]; then
        echo "system:"
        return
    fi

    if [ -z "${VES_NO_SUDO:-}" ] && command_exists sudo; then
        echo "system:sudo"
        return
    fi

    # Fall back to user directory (no sudo available)
    echo "user:"
}

# ============================================
# Shell RC File Management
# ============================================

# Get the appropriate RC file for the current shell
get_shell_rc_file() {
    CURRENT_SHELL=$(basename "${SHELL:-/bin/sh}")
    case "$CURRENT_SHELL" in
        bash)
            if [ -f "$HOME/.bashrc" ]; then
                echo "$HOME/.bashrc"
            elif [ -f "$HOME/.bash_profile" ]; then
                echo "$HOME/.bash_profile"
            else
                echo "$HOME/.bashrc"
            fi
            ;;
        zsh)
            echo "${ZDOTDIR:-$HOME}/.zshrc"
            ;;
        fish)
            echo "${XDG_CONFIG_HOME:-$HOME/.config}/fish/config.fish"
            ;;
        *)
            echo "$HOME/.profile"
            ;;
    esac
}

# Check if a line exists in a file (exact match or pattern)
line_exists_in_file() {
    FILE="$1"
    PATTERN="$2"
    [ -f "$FILE" ] && grep -qF "$PATTERN" "$FILE" 2>/dev/null
}

# Add PATH to shell RC file if not already present
add_path_to_rc() {
    DIR_TO_ADD="$1"
    RC_FILE=$(get_shell_rc_file)
    CURRENT_SHELL=$(basename "${SHELL:-/bin/sh}")

    # Check if already in PATH
    case ":$PATH:" in
        *":$DIR_TO_ADD:"*)
            return 0
            ;;
    esac

    # Check if already in RC file
    if line_exists_in_file "$RC_FILE" "$DIR_TO_ADD"; then
        return 0
    fi

    status "Adding $DIR_TO_ADD to PATH in $RC_FILE..."

    # Create RC file directory if needed
    mkdir -p "$(dirname "$RC_FILE")" 2>/dev/null

    # Add appropriate export statement
    if [ "$CURRENT_SHELL" = "fish" ]; then
        printf '\n# Added by vesctl installer\nfish_add_path %s\n' "$DIR_TO_ADD" >> "$RC_FILE"
    else
        printf '\n# Added by vesctl installer\nexport PATH="%s:$PATH"\n' "$DIR_TO_ADD" >> "$RC_FILE"
    fi

    success "Updated $RC_FILE"
    return 1  # Return 1 to indicate RC file was modified
}

# Add zsh completion configuration if needed
add_zsh_completion_config() {
    COMPLETION_DIR="$1"
    RC_FILE="${ZDOTDIR:-$HOME}/.zshrc"

    # Check if fpath already includes our completion directory
    if line_exists_in_file "$RC_FILE" "$COMPLETION_DIR"; then
        return 0
    fi

    # Check for existing fpath and compinit configuration
    HAS_COMPINIT=false
    if line_exists_in_file "$RC_FILE" "compinit"; then
        HAS_COMPINIT=true
    fi

    status "Adding zsh completion configuration to $RC_FILE..."

    # Create RC file if needed
    mkdir -p "$(dirname "$RC_FILE")" 2>/dev/null

    # Add the completion configuration
    {
        printf '\n# vesctl shell completions\n'
        printf 'fpath=(%s $fpath)\n' "$COMPLETION_DIR"
        if [ "$HAS_COMPINIT" = "false" ]; then
            printf 'autoload -Uz compinit && compinit\n'
        fi
    } >> "$RC_FILE"

    success "Updated $RC_FILE with completion configuration"
    return 1  # Return 1 to indicate RC file was modified
}

# Add bash completion sourcing if needed
add_bash_completion_config() {
    COMPLETION_FILE="$1"
    RC_FILE=$(get_shell_rc_file)

    # Check if already configured
    if line_exists_in_file "$RC_FILE" "vesctl"; then
        # Check if it's specifically our completion file
        if line_exists_in_file "$RC_FILE" "$COMPLETION_FILE"; then
            return 0
        fi
    fi

    # Check if bash-completion is being sourced (common setup)
    if line_exists_in_file "$RC_FILE" "bash_completion"; then
        # bash-completion will auto-load from standard directories
        return 0
    fi

    # For user-installed completions, we may need to source manually
    if [ -f "$COMPLETION_FILE" ]; then
        status "Adding bash completion configuration to $RC_FILE..."

        {
            printf '\n# vesctl shell completions\n'
            printf '[ -f "%s" ] && source "%s"\n' "$COMPLETION_FILE" "$COMPLETION_FILE"
        } >> "$RC_FILE"

        success "Updated $RC_FILE with completion configuration"
        return 1
    fi

    return 0
}

# Check if install directory is in PATH and automatically add to RC file if not
configure_path() {
    CHECK_DIR="$1"
    RC_MODIFIED=false

    # Check if install dir is already in PATH
    case ":$PATH:" in
        *":$CHECK_DIR:"*)
            return 0  # Already in PATH
            ;;
    esac

    # Automatically add to RC file
    if add_path_to_rc "$CHECK_DIR"; then
        : # Already configured
    else
        RC_MODIFIED=true
    fi

    # Provide reload instructions if we modified anything
    if [ "$RC_MODIFIED" = "true" ]; then
        RC_FILE=$(get_shell_rc_file)
        printf "\n"
        info "To use vesctl immediately, run:"
        printf "  ${CYAN}source %s${NC}\n" "$RC_FILE"
        printf "\n"
        info "Or start a new terminal session."
    fi
}

# ============================================
# Version Detection
# ============================================

get_latest_version() {
    status "Fetching latest version from GitHub..." >&2

    RESPONSE=$(http_get "$GITHUB_API")

    if [ -z "$RESPONSE" ]; then
        error "Failed to fetch latest version from GitHub.

Please check your internet connection or specify a version:
  VES_VERSION=1.1.0 sh install.sh

Or download manually from:
  https://github.com/${GITHUB_REPO}/releases/latest"
    fi

    # Try to parse with jq if available, fall back to sed if jq fails or isn't installed
    VERSION=""
    if command_exists jq; then
        VERSION=$(echo "$RESPONSE" | jq -r '.tag_name' 2>/dev/null | sed 's/^v//')
    fi

    # Fall back to sed if jq failed or wasn't available
    if [ -z "$VERSION" ] || [ "$VERSION" = "null" ]; then
        VERSION=$(echo "$RESPONSE" | grep '"tag_name":' | sed -E 's/.*"v?([^"]+)".*/\1/')
    fi

    if [ -z "$VERSION" ] || [ "$VERSION" = "null" ]; then
        error "Failed to parse version from GitHub API response.

Please specify a version manually:
  VES_VERSION=1.1.0 sh install.sh

Or download manually from:
  https://github.com/${GITHUB_REPO}/releases/latest"
    fi

    echo "$VERSION"
}

# ============================================
# Checksum Verification
# ============================================

verify_checksum() {
    ARCHIVE="$1"
    CHECKSUMS_FILE="$2"
    ARCHIVE_NAME="$3"

    if [ -n "${VES_NO_VERIFY:-}" ]; then
        warning "Skipping checksum verification (VES_NO_VERIFY is set)"
        return 0
    fi

    if [ ! -f "$CHECKSUMS_FILE" ]; then
        warning "Checksums file not found, skipping verification"
        return 0
    fi

    EXPECTED=$(grep "$ARCHIVE_NAME" "$CHECKSUMS_FILE" | awk '{print $1}')

    if [ -z "$EXPECTED" ]; then
        warning "Checksum not found for $ARCHIVE_NAME, skipping verification"
        return 0
    fi

    status "Verifying SHA256 checksum..."

    if command_exists sha256sum; then
        ACTUAL=$(sha256sum "$ARCHIVE" | awk '{print $1}')
    elif command_exists shasum; then
        ACTUAL=$(shasum -a 256 "$ARCHIVE" | awk '{print $1}')
    elif command_exists openssl; then
        ACTUAL=$(openssl dgst -sha256 "$ARCHIVE" | awk '{print $NF}')
    else
        warning "No SHA256 tool found (sha256sum, shasum, or openssl), skipping verification"
        return 0
    fi

    if [ "$EXPECTED" != "$ACTUAL" ]; then
        error "Checksum verification failed!

Expected: $EXPECTED
Actual:   $ACTUAL

This could indicate a corrupted download or a security issue.
Please try again or download manually from:
  https://github.com/${GITHUB_REPO}/releases"
    fi

    success "Checksum verified"
}

# ============================================
# Download and Install
# ============================================

download_and_install() {
    VERSION="$1"
    OS="$2"
    ARCH="$3"
    INSTALL_DIR="$4"
    SUDO_CMD="$5"

    ARCHIVE_NAME="${BINARY_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="${GITHUB_RELEASES}/v${VERSION}/${ARCHIVE_NAME}"
    CHECKSUMS_URL="${GITHUB_RELEASES}/v${VERSION}/checksums.txt"

    # Create temporary directory
    TEMP_DIR=$(mktemp -d)

    OS_DISPLAY=$(get_os_display_name "$OS")
    ARCH_DISPLAY=$(get_arch_display_name "$ARCH")

    status "Downloading vesctl v${VERSION}..."
    info "Platform: ${OS_DISPLAY} ${ARCH_DISPLAY}"
    info "URL: ${DOWNLOAD_URL}"

    # Download archive
    if ! http_download "$DOWNLOAD_URL" "${TEMP_DIR}/${ARCHIVE_NAME}"; then
        error "Failed to download vesctl.

URL: $DOWNLOAD_URL

Please check:
  - Your internet connection
  - The version exists: https://github.com/${GITHUB_REPO}/releases
  - Your platform is supported (${OS}/${ARCH})"
    fi

    # Download checksums
    if http_download "$CHECKSUMS_URL" "${TEMP_DIR}/checksums.txt" 2>/dev/null; then
        verify_checksum "${TEMP_DIR}/${ARCHIVE_NAME}" "${TEMP_DIR}/checksums.txt" "$ARCHIVE_NAME"
    else
        warning "Could not download checksums file, skipping verification"
    fi

    status "Extracting archive..."

    # Extract archive
    if ! tar -xzf "${TEMP_DIR}/${ARCHIVE_NAME}" -C "$TEMP_DIR"; then
        error "Failed to extract archive. The download may be corrupted.
Please try again or download manually."
    fi

    # Find the binary
    if [ ! -f "${TEMP_DIR}/${BINARY_NAME}" ]; then
        error "Binary not found in archive. This may indicate a packaging issue.
Please report this at: https://github.com/${GITHUB_REPO}/issues"
    fi

    status "Installing to ${INSTALL_DIR}..."

    # Create install directory if needed
    if [ ! -d "$INSTALL_DIR" ]; then
        $SUDO_CMD mkdir -p "$INSTALL_DIR" || error "Failed to create directory: $INSTALL_DIR"
    fi

    # Install binary
    $SUDO_CMD mv "${TEMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}" || error "Failed to install binary"
    $SUDO_CMD chmod +x "${INSTALL_DIR}/${BINARY_NAME}" || error "Failed to set executable permissions"

    success "Installed vesctl to ${INSTALL_DIR}/${BINARY_NAME}"
}

# ============================================
# Shell Completion Setup
# ============================================

setup_completion() {
    INSTALL_DIR="$1"
    SUDO_CMD="$2"

    VESCTL_BIN="${INSTALL_DIR}/${BINARY_NAME}"

    if [ ! -x "$VESCTL_BIN" ]; then
        warning "Cannot set up shell completion: vesctl binary not found"
        return
    fi

    # Detect current shell
    CURRENT_SHELL=$(basename "${SHELL:-sh}")

    case "$CURRENT_SHELL" in
        bash)
            setup_bash_completion "$VESCTL_BIN" "$SUDO_CMD"
            ;;
        zsh)
            setup_zsh_completion "$VESCTL_BIN" "$SUDO_CMD"
            ;;
        fish)
            setup_fish_completion "$VESCTL_BIN" "$SUDO_CMD"
            ;;
        *)
            info "Shell completion available: vesctl completion --help"
            ;;
    esac
}

setup_bash_completion() {
    VESCTL_BIN="$1"
    SUDO_CMD="$2"

    # Try system-wide location first
    if [ -d "/etc/bash_completion.d" ] && [ -w "/etc/bash_completion.d" -o -n "$SUDO_CMD" ]; then
        status "Setting up bash completion (system-wide)..."
        if $SUDO_CMD sh -c "\"$VESCTL_BIN\" completion bash > /etc/bash_completion.d/vesctl" 2>/dev/null; then
            success "Bash completion installed to /etc/bash_completion.d/vesctl"
            return
        fi
    fi

    # Fall back to user location
    BASH_COMPLETION_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/bash-completion/completions"
    if mkdir -p "$BASH_COMPLETION_DIR" 2>/dev/null; then
        status "Setting up bash completion (user)..."
        if "$VESCTL_BIN" completion bash > "${BASH_COMPLETION_DIR}/vesctl" 2>/dev/null; then
            success "Bash completion installed to ${BASH_COMPLETION_DIR}/vesctl"
            # Automatically configure RC file if needed
            add_bash_completion_config "${BASH_COMPLETION_DIR}/vesctl" || true
            return 0
        fi
    fi
}

setup_zsh_completion() {
    VESCTL_BIN="$1"
    SUDO_CMD="$2"

    # Try user location
    ZSH_COMPLETION_DIR="${ZDOTDIR:-$HOME}/.zsh/completions"
    if mkdir -p "$ZSH_COMPLETION_DIR" 2>/dev/null; then
        status "Setting up zsh completion..."
        if "$VESCTL_BIN" completion zsh > "${ZSH_COMPLETION_DIR}/_vesctl" 2>/dev/null; then
            success "Zsh completion installed to ${ZSH_COMPLETION_DIR}/_vesctl"
            # Automatically configure RC file
            add_zsh_completion_config "$ZSH_COMPLETION_DIR" || true
            return 0
        fi
    fi
}

setup_fish_completion() {
    VESCTL_BIN="$1"
    SUDO_CMD="$2"

    FISH_COMPLETION_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/fish/completions"
    if mkdir -p "$FISH_COMPLETION_DIR" 2>/dev/null; then
        status "Setting up fish completion..."
        if "$VESCTL_BIN" completion fish > "${FISH_COMPLETION_DIR}/vesctl.fish" 2>/dev/null; then
            success "Fish completion installed to ${FISH_COMPLETION_DIR}/vesctl.fish"
            # Fish auto-loads from this directory, no RC modification needed
            return
        fi
    fi
}

# ============================================
# Uninstall
# ============================================

uninstall() {
    INSTALL_DIR="${VES_INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"
    VESCTL_PATH="${INSTALL_DIR}/${BINARY_NAME}"

    # Determine if sudo is needed for uninstall
    SUDO_CMD=""
    if [ ! -w "$INSTALL_DIR" ] && [ "$(id -u)" -ne 0 ]; then
        if [ -z "${VES_NO_SUDO:-}" ] && command_exists sudo; then
            SUDO_CMD="sudo"
        fi
    fi

    status "Uninstalling vesctl..."

    if [ ! -f "$VESCTL_PATH" ]; then
        error "vesctl not found at $VESCTL_PATH"
    fi

    $SUDO_CMD rm -f "$VESCTL_PATH" || error "Failed to remove $VESCTL_PATH"
    success "Removed $VESCTL_PATH"

    # Clean up completion files
    status "Cleaning up shell completions..."

    # Bash completions
    if [ -f "/etc/bash_completion.d/vesctl" ]; then
        $SUDO_CMD rm -f "/etc/bash_completion.d/vesctl" 2>/dev/null
    fi
    BASH_COMPLETION_USER="${XDG_DATA_HOME:-$HOME/.local/share}/bash-completion/completions/vesctl"
    if [ -f "$BASH_COMPLETION_USER" ]; then
        rm -f "$BASH_COMPLETION_USER" 2>/dev/null
    fi

    # Zsh completions
    ZSH_COMPLETION="${ZDOTDIR:-$HOME}/.zsh/completions/_vesctl"
    if [ -f "$ZSH_COMPLETION" ]; then
        rm -f "$ZSH_COMPLETION" 2>/dev/null
    fi

    # Fish completions
    FISH_COMPLETION="${XDG_CONFIG_HOME:-$HOME/.config}/fish/completions/vesctl.fish"
    if [ -f "$FISH_COMPLETION" ]; then
        rm -f "$FISH_COMPLETION" 2>/dev/null
    fi

    success "vesctl has been uninstalled"
}

# ============================================
# Help
# ============================================

show_help() {
    cat << 'EOF'
vesctl installer

Automatically detects your platform and installs the appropriate binary
from GitHub releases.

USAGE
    curl -fsSL https://robinmordasiewicz.github.io/vesctl/install.sh | sh
    wget -qO- https://robinmordasiewicz.github.io/vesctl/install.sh | sh

OPTIONS
    --uninstall     Remove vesctl and shell completions
    --help, -h      Show this help message

ENVIRONMENT VARIABLES
    VES_VERSION      Specific version to install (default: latest)
    VES_INSTALL_DIR  Installation directory (default: /usr/local/bin)
    VES_NO_SUDO      Skip sudo even if needed (for custom install dirs)
    VES_NO_VERIFY    Skip checksum verification

SUPPORTED PLATFORMS
    Linux       amd64 (x86_64), arm64 (aarch64)
    macOS       amd64 (Intel), arm64 (Apple Silicon)
    Windows     amd64, arm64 (manual download required)

EXAMPLES
    # Install latest version
    curl -fsSL https://robinmordasiewicz.github.io/vesctl/install.sh | sh

    # Install specific version
    curl -fsSL https://robinmordasiewicz.github.io/vesctl/install.sh | VES_VERSION=1.1.0 sh

    # Install to custom directory (no sudo required)
    curl -fsSL https://robinmordasiewicz.github.io/vesctl/install.sh | VES_INSTALL_DIR=$HOME/.local/bin sh

    # Install using wget instead of curl
    wget -qO- https://robinmordasiewicz.github.io/vesctl/install.sh | sh

    # Uninstall
    curl -fsSL https://robinmordasiewicz.github.io/vesctl/install.sh | sh -s -- --uninstall

WINDOWS INSTALLATION
    Download the appropriate zip file from GitHub releases:
    https://github.com/robinmordasiewicz/vesctl/releases/latest

    - Windows (Intel/AMD): vesctl_VERSION_windows_amd64.zip
    - Windows (ARM):       vesctl_VERSION_windows_arm64.zip

EOF
}

# ============================================
# Main
# ============================================

main() {
    # Parse arguments
    while [ $# -gt 0 ]; do
        case "$1" in
            --uninstall)
                uninstall
                exit 0
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                error "Unknown option: $1

Use --help for usage information."
                ;;
        esac
        shift
    done

    # Print banner
    printf "\n"
    printf "${CYAN}vesctl installer${NC}\n"
    printf "${CYAN}=================${NC}\n"
    printf "\n"

    # Check for HTTP client
    HTTP_CLIENT=$(detect_http_client)
    if [ "$HTTP_CLIENT" = "none" ]; then
        error "Neither curl nor wget found. Please install one of them:

On Debian/Ubuntu: sudo apt-get install curl
On RHEL/CentOS:   sudo yum install curl
On macOS:         curl should be pre-installed
On Alpine:        apk add curl"
    fi

    if ! command_exists tar; then
        error "tar is required but not installed. Please install tar and try again."
    fi

    # Detect platform
    OS=$(detect_os)
    ARCH=$(detect_arch)

    OS_DISPLAY=$(get_os_display_name "$OS")
    ARCH_DISPLAY=$(get_arch_display_name "$ARCH")

    status "Detected platform: ${OS_DISPLAY} ${ARCH_DISPLAY}"

    # Get version
    if [ -n "${VES_VERSION:-}" ]; then
        VERSION="$VES_VERSION"
        status "Using specified version: v${VERSION}"
    else
        VERSION=$(get_latest_version)
        status "Latest version: v${VERSION}"
    fi

    # Determine installation strategy
    STRATEGY=$(determine_install_strategy)
    STRATEGY_TYPE="${STRATEGY%%:*}"
    SUDO_CMD="${STRATEGY#*:}"

    case "$STRATEGY_TYPE" in
        system)
            INSTALL_DIR="$DEFAULT_INSTALL_DIR"
            ;;
        user)
            INSTALL_DIR="$USER_INSTALL_DIR"
            # Create ~/.local/bin if needed
            if [ ! -d "$INSTALL_DIR" ]; then
                mkdir -p "$INSTALL_DIR"
            fi
            SUDO_CMD=""
            status "Installing to $INSTALL_DIR (no sudo required)"
            ;;
        custom)
            INSTALL_DIR="$VES_INSTALL_DIR"
            ;;
    esac

    # Check for existing installation
    if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        EXISTING_VERSION=$("${INSTALL_DIR}/${BINARY_NAME}" version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo "unknown")
        if [ "$EXISTING_VERSION" = "$VERSION" ]; then
            success "vesctl v${VERSION} is already installed"
            printf "\n"
            exit 0
        fi
        warning "Upgrading from v${EXISTING_VERSION} to v${VERSION}"
    fi

    # Download and install
    download_and_install "$VERSION" "$OS" "$ARCH" "$INSTALL_DIR" "$SUDO_CMD"

    # Verify installation
    if ! "${INSTALL_DIR}/${BINARY_NAME}" version >/dev/null 2>&1; then
        error "Installation verification failed. Please check the binary at ${INSTALL_DIR}/${BINARY_NAME}"
    fi

    # Set up shell completion
    setup_completion "$INSTALL_DIR" "$SUDO_CMD"

    # Configure PATH in shell RC file if needed
    configure_path "$INSTALL_DIR"

    # Success message
    printf "\n"
    success "vesctl v${VERSION} installed successfully!"
    printf "\n"
    printf "%s\n" "Get started:"
    printf "  ${CYAN}vesctl --help${NC}              # Show help\n"
    printf "  ${CYAN}vesctl configure${NC}           # Configure credentials\n"
    printf "  ${CYAN}vesctl version${NC}             # Show version info\n"
    printf "\n"
    printf "%s\n" "Documentation: https://github.com/${GITHUB_REPO}"
    printf "\n"
}

main "$@"
