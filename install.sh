#!/bin/sh
# vesctl install script
# Usage: curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | sh
#
# Environment variables:
#   VESCTL_VERSION      - Specific version to install (default: latest)
#   VESCTL_INSTALL_DIR  - Installation directory (default: /usr/local/bin)
#   VESCTL_NO_SUDO      - Skip sudo if set to any value
#   VESCTL_NO_VERIFY    - Skip checksum verification if set

set -e

# Configuration
GITHUB_REPO="robinmordasiewicz/vesctl"
GITHUB_API="https://api.github.com/repos/${GITHUB_REPO}/releases/latest"
GITHUB_RELEASES="https://github.com/${GITHUB_REPO}/releases/download"
DEFAULT_INSTALL_DIR="/usr/local/bin"
BINARY_NAME="vesctl"

# Colors for output (check if terminal supports colors)
if [ -t 1 ] && [ -n "$(tput colors 2>/dev/null)" ] && [ "$(tput colors)" -ge 8 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    BLUE='\033[0;34m'
    NC='\033[0m' # No Color
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
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
# Platform Detection
# ============================================

detect_os() {
    OS="$(uname -s)"
    case "$OS" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*)
            error "Windows detected. Please use PowerShell or download manually from GitHub releases."
            ;;
        *)       error "Unsupported operating system: $OS" ;;
    esac
}

detect_arch() {
    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64)  echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        armv7l)        error "32-bit ARM is not supported. Please use a 64-bit system." ;;
        i386|i686)     error "32-bit x86 is not supported. Please use a 64-bit system." ;;
        *)             error "Unsupported architecture: $ARCH" ;;
    esac
}

# ============================================
# Privilege Handling
# ============================================

need_sudo() {
    INSTALL_DIR="${VESCTL_INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"

    # Check if we can write to the install directory
    if [ -w "$INSTALL_DIR" ] || [ -w "$(dirname "$INSTALL_DIR")" ]; then
        echo ""
        return
    fi

    # Check if running as root
    if [ "$(id -u)" -eq 0 ]; then
        echo ""
        return
    fi

    # Check if sudo is available
    if [ -n "$VESCTL_NO_SUDO" ]; then
        error "Cannot write to $INSTALL_DIR and VESCTL_NO_SUDO is set. Try setting VESCTL_INSTALL_DIR to a writable location."
    fi

    if ! command_exists sudo; then
        error "Cannot write to $INSTALL_DIR and sudo is not available. Try running as root or set VESCTL_INSTALL_DIR to a writable location."
    fi

    echo "sudo"
}

# ============================================
# Version Detection
# ============================================

get_latest_version() {
    if ! command_exists curl; then
        error "curl is required but not installed. Please install curl and try again."
    fi

    VERSION=$(curl -fsSL "$GITHUB_API" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/^v//')

    if [ -z "$VERSION" ]; then
        error "Failed to fetch latest version from GitHub. Please check your internet connection or specify VESCTL_VERSION."
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

    if [ -n "$VESCTL_NO_VERIFY" ]; then
        warning "Skipping checksum verification (VESCTL_NO_VERIFY is set)"
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

    status "Verifying checksum..."

    if command_exists sha256sum; then
        ACTUAL=$(sha256sum "$ARCHIVE" | awk '{print $1}')
    elif command_exists shasum; then
        ACTUAL=$(shasum -a 256 "$ARCHIVE" | awk '{print $1}')
    else
        warning "Neither sha256sum nor shasum found, skipping verification"
        return 0
    fi

    if [ "$EXPECTED" != "$ACTUAL" ]; then
        error "Checksum verification failed!
Expected: $EXPECTED
Actual:   $ACTUAL

This could indicate a corrupted download or a security issue.
Please try again or download manually from GitHub."
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

    status "Downloading vesctl v${VERSION} for ${OS}/${ARCH}..."

    # Download archive
    if ! curl -fsSL -o "${TEMP_DIR}/${ARCHIVE_NAME}" "$DOWNLOAD_URL"; then
        error "Failed to download vesctl. Please check your internet connection or try again later.
URL: $DOWNLOAD_URL"
    fi

    # Download checksums
    if curl -fsSL -o "${TEMP_DIR}/checksums.txt" "$CHECKSUMS_URL" 2>/dev/null; then
        verify_checksum "${TEMP_DIR}/${ARCHIVE_NAME}" "${TEMP_DIR}/checksums.txt" "$ARCHIVE_NAME"
    else
        warning "Could not download checksums file, skipping verification"
    fi

    status "Extracting..."

    # Extract archive
    if ! tar -xzf "${TEMP_DIR}/${ARCHIVE_NAME}" -C "$TEMP_DIR"; then
        error "Failed to extract archive"
    fi

    # Find the binary
    if [ ! -f "${TEMP_DIR}/${BINARY_NAME}" ]; then
        error "Binary not found in archive"
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
            status "Shell completion is available. Run 'vesctl completion --help' for instructions."
            ;;
    esac
}

setup_bash_completion() {
    VESCTL_BIN="$1"
    SUDO_CMD="$2"

    # Try system-wide location first
    if [ -d "/etc/bash_completion.d" ] && [ -w "/etc/bash_completion.d" -o -n "$SUDO_CMD" ]; then
        status "Setting up bash completion (system-wide)..."
        $SUDO_CMD sh -c "\"$VESCTL_BIN\" completion bash > /etc/bash_completion.d/vesctl" 2>/dev/null && \
            success "Bash completion installed to /etc/bash_completion.d/vesctl" && return
    fi

    # Fall back to user location
    BASH_COMPLETION_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/bash-completion/completions"
    if mkdir -p "$BASH_COMPLETION_DIR" 2>/dev/null; then
        status "Setting up bash completion (user)..."
        "$VESCTL_BIN" completion bash > "${BASH_COMPLETION_DIR}/vesctl" 2>/dev/null && \
            success "Bash completion installed to ${BASH_COMPLETION_DIR}/vesctl" && return
    fi

    # Provide manual instructions
    printf "\n%s\n" "To enable bash completion, add this to your ~/.bashrc:"
    printf "  %s\n" 'eval "$(vesctl completion bash)"'
}

setup_zsh_completion() {
    VESCTL_BIN="$1"
    SUDO_CMD="$2"

    # Try user location
    ZSH_COMPLETION_DIR="${ZDOTDIR:-$HOME}/.zsh/completions"
    if mkdir -p "$ZSH_COMPLETION_DIR" 2>/dev/null; then
        status "Setting up zsh completion..."
        "$VESCTL_BIN" completion zsh > "${ZSH_COMPLETION_DIR}/_vesctl" 2>/dev/null
        success "Zsh completion installed to ${ZSH_COMPLETION_DIR}/_vesctl"
        printf "\n%s\n" "Add this to your ~/.zshrc if not already present:"
        printf "  %s\n" "fpath=(${ZSH_COMPLETION_DIR} \$fpath)"
        printf "  %s\n" "autoload -Uz compinit && compinit"
        return
    fi

    # Provide manual instructions
    printf "\n%s\n" "To enable zsh completion, add this to your ~/.zshrc:"
    printf "  %s\n" 'eval "$(vesctl completion zsh)"'
}

setup_fish_completion() {
    VESCTL_BIN="$1"
    SUDO_CMD="$2"

    FISH_COMPLETION_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/fish/completions"
    if mkdir -p "$FISH_COMPLETION_DIR" 2>/dev/null; then
        status "Setting up fish completion..."
        "$VESCTL_BIN" completion fish > "${FISH_COMPLETION_DIR}/vesctl.fish" 2>/dev/null && \
            success "Fish completion installed to ${FISH_COMPLETION_DIR}/vesctl.fish" && return
    fi

    # Provide manual instructions
    printf "\n%s\n" "To enable fish completion, run:"
    printf "  %s\n" "vesctl completion fish > ~/.config/fish/completions/vesctl.fish"
}

# ============================================
# Uninstall
# ============================================

uninstall() {
    INSTALL_DIR="${VESCTL_INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"
    VESCTL_PATH="${INSTALL_DIR}/${BINARY_NAME}"
    SUDO_CMD=$(need_sudo)

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
    cat << EOF
vesctl installer

Usage: $0 [OPTIONS]

Options:
    --uninstall     Remove vesctl and shell completions
    --help, -h      Show this help message

Environment Variables:
    VESCTL_VERSION      Specific version to install (default: latest)
    VESCTL_INSTALL_DIR  Installation directory (default: /usr/local/bin)
    VESCTL_NO_SUDO      Skip sudo even if needed (for custom install dirs)
    VESCTL_NO_VERIFY    Skip checksum verification

Examples:
    # Install latest version
    curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | sh

    # Install specific version
    curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | VESCTL_VERSION=0.1.0 sh

    # Install to custom directory
    curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | VESCTL_INSTALL_DIR=\$HOME/.local/bin sh

    # Uninstall
    curl -fsSL https://raw.githubusercontent.com/robinmordasiewicz/vesctl/main/install.sh | sh -s -- --uninstall

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
                error "Unknown option: $1. Use --help for usage."
                ;;
        esac
        shift
    done

    # Check dependencies
    if ! command_exists curl; then
        error "curl is required but not installed. Please install curl and try again.

On Debian/Ubuntu: sudo apt-get install curl
On RHEL/CentOS:   sudo yum install curl
On macOS:         curl should be pre-installed"
    fi

    if ! command_exists tar; then
        error "tar is required but not installed. Please install tar and try again."
    fi

    # Detect platform
    OS=$(detect_os)
    ARCH=$(detect_arch)

    status "Detected platform: ${OS}/${ARCH}"

    # Get version
    if [ -n "$VESCTL_VERSION" ]; then
        VERSION="$VESCTL_VERSION"
        status "Using specified version: v${VERSION}"
    else
        status "Fetching latest version..."
        VERSION=$(get_latest_version)
        status "Latest version: v${VERSION}"
    fi

    # Check for existing installation
    INSTALL_DIR="${VESCTL_INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"
    if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        EXISTING_VERSION=$("${INSTALL_DIR}/${BINARY_NAME}" version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo "unknown")
        if [ "$EXISTING_VERSION" = "$VERSION" ]; then
            success "vesctl v${VERSION} is already installed"
            exit 0
        fi
        warning "Existing installation found (v${EXISTING_VERSION}). Will upgrade to v${VERSION}."
    fi

    # Determine if sudo is needed
    SUDO_CMD=$(need_sudo)

    # Download and install
    download_and_install "$VERSION" "$OS" "$ARCH" "$INSTALL_DIR" "$SUDO_CMD"

    # Verify installation
    if ! "${INSTALL_DIR}/${BINARY_NAME}" version >/dev/null 2>&1; then
        error "Installation verification failed. Please check the binary at ${INSTALL_DIR}/${BINARY_NAME}"
    fi

    # Set up shell completion
    setup_completion "$INSTALL_DIR" "$SUDO_CMD"

    # Success message
    printf "\n"
    success "vesctl v${VERSION} installed successfully!"
    printf "\n"
    printf "%s\n" "Get started:"
    printf "  %s\n" "vesctl --help              # Show help"
    printf "  %s\n" "vesctl configure           # Configure credentials"
    printf "  %s\n" "vesctl version             # Show version info"
    printf "\n"
    printf "%s\n" "Documentation: https://github.com/${GITHUB_REPO}"
    printf "\n"
}

main "$@"
