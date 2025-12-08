#!/bin/bash
# test-resource-help.sh - Test help text compatibility for multiple resource types
# Phase 5: Multi-Resource Validation
#
# NOTE: vesctl-0.2.35 help commands can take 20-25 seconds to complete.
# The timeout is set to 60 seconds in lib/common.sh to accommodate this.

set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# Change to project root so binary paths work correctly
cd "$PROJECT_ROOT"

source "${SCRIPT_DIR}/../../lib/common.sh"

echo "=========================================="
echo "Phase 5A: Multi-Resource Help Tests"
echo "=========================================="
echo ""

# Initialize results directory
init_results_dir
PHASE_DIR="${RESULTS_DIR}/phase5-resources"
mkdir -p "$PHASE_DIR"

log_info "Original vesctl: ${ORIGINAL_VESCTL}"
log_info "Our vesctl: ${OUR_VESCTL}"
log_info "Results directory: ${RESULTS_DIR}"
echo ""

export RESULTS_DIR="$PHASE_DIR"

# ============================================
# Representative resource types to test
# ============================================
# These resources cover different categories and are commonly used
RESOURCES=(
    "http_loadbalancer"
    "origin_pool"
    "app_firewall"
    "healthcheck"
    "dns_zone"
    "certificate"
    "user"
    "token"
    "virtual_network"
    "cloud_credentials"
)

# Additional resources from the original vesctl
ADDITIONAL_RESOURCES=(
    "namespace"
    "virtual_host"
    "virtual_site"
    "service_policy"
    "network_policy"
    "cluster"
    "endpoint"
    "route"
    "secret_policy"
)

# ============================================
# Helper: Test help structure (excluding resource list)
# ============================================
test_help_structure() {
    local test_name="$1"
    shift
    local cmd=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Capture full help
    $ORIGINAL_VESCTL "${cmd[@]}" --help > "${test_dir}/original_full.txt" 2>&1
    local orig_exit=$?
    $OUR_VESCTL "${cmd[@]}" --help > "${test_dir}/ours_full.txt" 2>&1
    local our_exit=$?

    # Store exit codes
    echo "orig_exit=$orig_exit" > "${test_dir}/exits.txt"
    echo "our_exit=$our_exit" >> "${test_dir}/exits.txt"

    # If original doesn't support this resource, it's expected
    if [[ $orig_exit -ne 0 ]]; then
        if [[ $our_exit -eq 0 ]]; then
            log_warn "${test_name} (resource not in original - our enhancement)"
            echo "WARN_NEW_RESOURCE" > "${test_dir}/result.txt"
            return 0
        else
            log_warn "${test_name} (resource not in either)"
            echo "WARN_MISSING" > "${test_dir}/result.txt"
            return 0
        fi
    fi

    # Check if original shows parent command help (resource not found but exit 0)
    # This happens when original doesn't have the resource - it shows parent help instead
    local orig_first_line=$(head -1 "${test_dir}/original_full.txt")
    local our_first_line=$(head -1 "${test_dir}/ours_full.txt")

    # If original shows generic "List/Get/Create configuration" but ours shows specific resource
    # this means the resource doesn't exist in original
    if [[ "$orig_first_line" =~ ^(List|Get|Create)\ configuration ]] && \
       ! [[ "$our_first_line" =~ ^(List|Get|Create)\ configuration ]]; then
        log_warn "${test_name} (resource not in original - our enhancement)"
        echo "WARN_NEW_RESOURCE" > "${test_dir}/result.txt"
        return 0
    fi

    # Extract structure (Usage, Aliases, Examples, Flags, Global Flags sections)
    # Exclude lines that are resource names (2 spaces + lowercase/underscore)
    grep -v "^  [a-z_]" "${test_dir}/original_full.txt" > "${test_dir}/original_structure.txt"
    grep -v "^  [a-z_]" "${test_dir}/ours_full.txt" > "${test_dir}/ours_structure.txt"

    if diff -q "${test_dir}/original_structure.txt" "${test_dir}/ours_structure.txt" > /dev/null 2>&1; then
        log_pass "${test_name} (structure identical)"
        echo "PASS" > "${test_dir}/result.txt"
        return 0
    else
        diff -u "${test_dir}/original_structure.txt" "${test_dir}/ours_structure.txt" > "${test_dir}/diff.txt" 2>&1 || true
        log_fail "${test_name} (structure differs)"
        echo "FAIL" > "${test_dir}/result.txt"
        return 1
    fi
}

# ============================================
# Helper: Test flags section only
# ============================================
test_flags_section() {
    local test_name="$1"
    shift
    local cmd=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Get full help output to check for parent command pattern
    local orig_full=$($ORIGINAL_VESCTL "${cmd[@]}" --help 2>&1)
    local our_full=$($OUR_VESCTL "${cmd[@]}" --help 2>&1)
    local orig_exit=$?
    local our_exit=$?

    # Check if original shows parent command help (resource not found but exit 0)
    local orig_first_line=$(echo "$orig_full" | head -1)
    local our_first_line=$(echo "$our_full" | head -1)

    if [[ "$orig_first_line" =~ ^(List|Get|Create)\ configuration ]] && \
       ! [[ "$our_first_line" =~ ^(List|Get|Create)\ configuration ]]; then
        log_warn "${test_name} (resource not in original - our enhancement)"
        echo "WARN_NEW_RESOURCE" > "${test_dir}/result.txt"
        return 0
    fi

    # Extract Flags and Global Flags sections
    echo "$orig_full" | sed -n '/^Flags:/,/^$/p' > "${test_dir}/original_flags.txt"
    echo "$our_full" | sed -n '/^Flags:/,/^$/p' > "${test_dir}/ours_flags.txt"

    # If original doesn't support this, skip
    if [[ $orig_exit -ne 0 ]] || [[ ! -s "${test_dir}/original_flags.txt" ]]; then
        log_skip "${test_name} (resource not in original)"
        echo "SKIP" > "${test_dir}/result.txt"
        return 0
    fi

    echo "$orig_full" | sed -n '/^Global Flags:/,/^$/p' > "${test_dir}/original_global.txt"
    echo "$our_full" | sed -n '/^Global Flags:/,/^$/p' > "${test_dir}/ours_global.txt"

    local flags_match=true
    if ! diff -q "${test_dir}/original_flags.txt" "${test_dir}/ours_flags.txt" > /dev/null 2>&1; then
        flags_match=false
    fi
    if ! diff -q "${test_dir}/original_global.txt" "${test_dir}/ours_global.txt" > /dev/null 2>&1; then
        flags_match=false
    fi

    if $flags_match; then
        log_pass "${test_name} (flags identical)"
        echo "PASS" > "${test_dir}/result.txt"
        return 0
    else
        diff -u "${test_dir}/original_flags.txt" "${test_dir}/ours_flags.txt" > "${test_dir}/flags_diff.txt" 2>&1 || true
        diff -u "${test_dir}/original_global.txt" "${test_dir}/ours_global.txt" > "${test_dir}/global_diff.txt" 2>&1 || true
        log_fail "${test_name} (flags differ)"
        echo "FAIL" > "${test_dir}/result.txt"
        return 1
    fi
}

# ============================================
# Test each resource type
# ============================================
echo "--- Primary Resources (10) ---"
echo ""

for resource in "${RESOURCES[@]}"; do
    echo "Testing: $resource"

    # Test list --help - SKIP due to vesctl-0.2.35 bug
    if is_buggy_command "configuration" "list"; then
        log_skip "list-${resource} (vesctl-0.2.35 hangs on configuration list --help)"
        mkdir -p "${RESULTS_DIR}/list-${resource}"
        echo "SKIP" > "${RESULTS_DIR}/list-${resource}/result.txt"
        mkdir -p "${RESULTS_DIR}/list-${resource}-flags"
        echo "SKIP" > "${RESULTS_DIR}/list-${resource}-flags/result.txt"
    else
        test_help_structure "list-${resource}" configuration list "$resource" || true
        test_flags_section "list-${resource}-flags" configuration list "$resource" || true
    fi

    # Test get --help
    test_help_structure "get-${resource}" configuration get "$resource" || true
    test_flags_section "get-${resource}-flags" configuration get "$resource" || true

    # Test create --help
    test_help_structure "create-${resource}" configuration create "$resource" || true

    echo ""
done

# ============================================
# Test additional resources
# ============================================
echo "--- Additional Resources (9) ---"
echo ""

for resource in "${ADDITIONAL_RESOURCES[@]}"; do
    echo "Testing: $resource"

    # Test list --help only - SKIP due to vesctl-0.2.35 bug
    if is_buggy_command "configuration" "list"; then
        log_skip "list-${resource} (vesctl-0.2.35 hangs on configuration list --help)"
        mkdir -p "${RESULTS_DIR}/list-${resource}"
        echo "SKIP" > "${RESULTS_DIR}/list-${resource}/result.txt"
    else
        test_help_structure "list-${resource}" configuration list "$resource" || true
    fi

    echo ""
done

# ============================================
# Summary
# ============================================
echo ""
print_summary

echo ""
echo "Phase 5A (Resource Help) test results: ${RESULTS_DIR}"
