# Branding Configuration
# Single source of truth for CLI branding
# Update this file to rebrand the entire project

# CLI branding
CLI_NAME := xcsh
CLI_FULL_NAME := F5 Distributed Cloud Shell
CLI_DESCRIPTION := Command-line interface for F5 Distributed Cloud services
CLI_SHORT_DESCRIPTION := F5 Distributed Cloud Shell

# Repository information
REPO_OWNER := robinmordasiewicz
REPO_NAME := $(CLI_NAME)
MODULE_PATH := github.com/$(REPO_OWNER)/$(CLI_NAME)

# Documentation
DOCS_SITE_NAME := $(CLI_NAME) Documentation
DOCS_SITE_URL := https://$(REPO_OWNER).github.io/$(CLI_NAME)/
DOCS_REPO_URL := https://github.com/$(REPO_OWNER)/$(CLI_NAME)

# Build artifacts
BINARY_NAME := $(CLI_NAME)
PROJECT_NAME := $(CLI_NAME)
