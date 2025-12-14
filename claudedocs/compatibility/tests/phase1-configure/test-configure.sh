#!/bin/bash
# test-configure.sh - Phase 1: Configure/Configuration command help tests
# Focuses on help text consistency between original and our f5xcctl

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/compare.sh"
source "${SCRIPT_DIR}/../../lib/report.sh"

echo "=========================================="
echo "Phase 1: Configure/Configuration Help Tests"
echo "=========================================="
echo ""

# Initialize results directory
init_results_dir
PHASE_DIR="${RESULTS_DIR}/phase1-configure"
mkdir -p "$PHASE_DIR"

# Check binaries exist
if ! check_original_vesctl; then
    echo "Original vesctl not found. Exiting."
    exit 1
fi

if ! check_our_vesctl; then
    echo "Our f5xcctl not found. Exiting."
    exit 1
fi

log_info "Original vesctl: ${ORIGINAL_VESCTL}"
log_info "Our f5xcctl: ${OUR_F5XCCTL}"
log_info "Results directory: ${RESULTS_DIR}"
echo ""

# Override RESULTS_DIR for test_help function
export RESULTS_DIR="$PHASE_DIR"

# Test 1: Main help
log_info "Testing: vesctl --help"
test_help "vesctl-help" --help || true

# Test 2: Configuration top-level help
log_info "Testing: vesctl configuration --help"
test_help "configuration-help" configuration --help || true

# Test 3: Configuration subcommand help
SUBCOMMANDS=(
    "list"
    "get"
    "create"
    "delete"
    "replace"
    "apply"
    "status"
    "add-labels"
    "remove-labels"
    "patch"
)

for subcmd in "${SUBCOMMANDS[@]}"; do
    log_info "Testing: vesctl configuration ${subcmd} --help"
    test_help "configuration-${subcmd}-help" configuration "$subcmd" --help || true
done

# Test 4: Test a few key resource types with list help
KEY_RESOURCES=(
    "namespace"
    "http_loadbalancer"
    "origin_pool"
    "virtual_host"
    "service_policy"
    "app_firewall"
    "certificate"
    "healthcheck"
)

for resource in "${KEY_RESOURCES[@]}"; do
    log_info "Testing: vesctl configuration list ${resource} --help"
    test_help "configuration-list-${resource}-help" configuration list "$resource" --help || true
done

# Test 5: Test get help for same resources
for resource in "${KEY_RESOURCES[@]}"; do
    log_info "Testing: vesctl configuration get ${resource} --help"
    test_help "configuration-get-${resource}-help" configuration get "$resource" --help || true
done

# Reset RESULTS_DIR for report generation
export RESULTS_DIR=$(dirname "$PHASE_DIR")

echo ""
echo "=========================================="
echo "Phase 1 Complete"
echo "=========================================="

# Print summary
print_summary

# Generate reports
generate_json_report "$RESULTS_DIR"
generate_markdown_report "$RESULTS_DIR"

echo ""
echo "Reports generated:"
echo "  - ${RESULTS_DIR}/summary.json"
echo "  - ${RESULTS_DIR}/summary.md"
