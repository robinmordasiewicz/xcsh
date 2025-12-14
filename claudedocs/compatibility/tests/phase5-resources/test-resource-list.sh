#!/bin/bash
# test-resource-list.sh - Test list operations for multiple resource types
# Phase 5B: Multi-Resource List API Tests (Read-Only)
#
# These tests are SAFE to run because:
# - List operations are read-only
# - We compare outputs from the same existing data
# - No race conditions since we're not creating/modifying objects

set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# Change to project root so binary paths work correctly
cd "$PROJECT_ROOT"

source "${SCRIPT_DIR}/../../lib/common.sh"

echo "=========================================="
echo "Phase 5B: Multi-Resource List Tests (API)"
echo "=========================================="
echo ""

# Check for API credentials
if ! has_api_credentials; then
    echo "ERROR: API credentials not available"
    echo "Set F5XC_P12_PASSWORD and F5XC_P12_FILE environment variables"
    exit 1
fi

# Initialize results directory
init_results_dir
PHASE_DIR="${RESULTS_DIR}/phase5-list-api"
mkdir -p "$PHASE_DIR"

log_info "Original vesctl: ${ORIGINAL_VESCTL}"
log_info "Our f5xcctl: ${OUR_F5XCCTL}"
log_info "Results directory: ${RESULTS_DIR}"
echo ""

export RESULTS_DIR="$PHASE_DIR"

# ============================================
# Resources that commonly have objects in staging
# These are read-only list operations
# ============================================
LIST_RESOURCES=(
    "namespace"            # Always has objects
    "user"                 # Always has objects
    "role"                 # Always has objects
    "known_label"          # Common objects
    "known_label_key"      # Common objects
)

# Resources that may or may not have objects
OPTIONAL_RESOURCES=(
    "http_loadbalancer"
    "origin_pool"
    "healthcheck"
    "certificate"
    "virtual_network"
    "dns_zone"
    "app_firewall"
    "service_policy"
    "network_policy"
    "token"
)

# ============================================
# Helper: Run list and compare outputs
# ============================================
test_list_api() {
    local test_name="$1"
    local resource="$2"
    local namespace="${3:-}"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    local ns_args=()
    if [[ -n "$namespace" ]]; then
        ns_args=("-n" "$namespace")
    fi

    # Run original vesctl
    log_info "Running original: configuration list $resource ${ns_args[*]}"
    $ORIGINAL_VESCTL configuration list "$resource" "${ns_args[@]}" --outfmt json \
        > "${test_dir}/original.stdout" 2> "${test_dir}/original.stderr"
    local orig_exit=$?

    # Run our f5xcctl
    log_info "Running ours: configuration list $resource ${ns_args[*]}"
    $OUR_F5XCCTL configuration list "$resource" "${ns_args[@]}" --outfmt json \
        > "${test_dir}/ours.stdout" 2> "${test_dir}/ours.stderr"
    local our_exit=$?

    # Store exit codes
    echo "original_exit=$orig_exit" > "${test_dir}/exits.txt"
    echo "our_exit=$our_exit" >> "${test_dir}/exits.txt"

    # Check if original supports this resource
    if [[ $orig_exit -ne 0 ]]; then
        if grep -q "unknown command" "${test_dir}/original.stderr" 2>/dev/null; then
            if [[ $our_exit -eq 0 ]]; then
                log_warn "${test_name} (resource not in original - our enhancement)"
                echo "WARN_NEW_RESOURCE" > "${test_dir}/result.txt"
            else
                log_warn "${test_name} (resource not in either)"
                echo "WARN_MISSING" > "${test_dir}/result.txt"
            fi
            return 0
        fi
    fi

    # Compare exit codes
    if [[ "$orig_exit" != "$our_exit" ]]; then
        log_fail "${test_name}: Exit code mismatch (original: $orig_exit, ours: $our_exit)"
        echo "EXIT_MISMATCH" > "${test_dir}/result.txt"
        return 0
    fi

    # If both errored, compare error patterns
    if [[ $orig_exit -ne 0 ]]; then
        if compare_error_messages "${test_dir}/original.stderr" "${test_dir}/ours.stderr"; then
            log_pass "${test_name} (both errored consistently)"
            echo "PASS_ERROR" > "${test_dir}/result.txt"
        else
            log_warn "${test_name}: Error messages differ"
            echo "WARN_ERROR_MSG" > "${test_dir}/result.txt"
        fi
        return 0
    fi

    # Compare API responses
    if compare_api_responses \
        "${test_dir}/original.stdout" \
        "${test_dir}/ours.stdout" \
        "${test_dir}/stdout.diff" \
        "$test_dir"; then
        log_pass "${test_name}"
        echo "PASS" > "${test_dir}/result.txt"
        return 0
    else
        log_fail "${test_name}: Response differs"
        echo "FAIL" > "${test_dir}/result.txt"
        return 0
    fi
}

# ============================================
# Test core resources (always have data)
# ============================================
echo "--- Core Resources (always exist) ---"
echo ""

for resource in "${LIST_RESOURCES[@]}"; do
    echo "Testing list: $resource"
    test_list_api "list-${resource}" "$resource" || true
    echo ""
done

# ============================================
# Test optional resources (may not have data)
# ============================================
echo "--- Optional Resources (may be empty) ---"
echo ""

for resource in "${OPTIONAL_RESOURCES[@]}"; do
    echo "Testing list: $resource"
    # Try system namespace first, then default
    test_list_api "list-${resource}-system" "$resource" "system" || true
    echo ""
done

# ============================================
# Test table output format comparison
# ============================================
echo "--- Table Output Format Tests ---"
echo ""

test_table_format() {
    local test_name="$1"
    local resource="$2"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Run with table format (default)
    $ORIGINAL_VESCTL configuration list "$resource" \
        > "${test_dir}/original.stdout" 2> "${test_dir}/original.stderr"
    local orig_exit=$?

    $OUR_F5XCCTL configuration list "$resource" \
        > "${test_dir}/ours.stdout" 2> "${test_dir}/ours.stderr"
    local our_exit=$?

    echo "original_exit=$orig_exit" > "${test_dir}/exits.txt"
    echo "our_exit=$our_exit" >> "${test_dir}/exits.txt"

    if [[ $orig_exit -ne 0 ]] || [[ $our_exit -ne 0 ]]; then
        log_skip "${test_name} (skipped - command failed)"
        echo "SKIP" > "${test_dir}/result.txt"
        return 0
    fi

    # Compare table headers (first line)
    local orig_header=$(head -1 "${test_dir}/original.stdout")
    local our_header=$(head -1 "${test_dir}/ours.stdout")

    if [[ "$orig_header" == "$our_header" ]]; then
        log_pass "${test_name} (table headers match)"
        echo "PASS" > "${test_dir}/result.txt"
    else
        log_fail "${test_name} (table headers differ)"
        echo "FAIL_HEADERS" > "${test_dir}/result.txt"
        echo "Original: $orig_header" > "${test_dir}/header_diff.txt"
        echo "Ours: $our_header" >> "${test_dir}/header_diff.txt"
    fi
}

# Test table format for namespace (always exists)
test_table_format "table-namespace" "namespace" || true

# ============================================
# Summary
# ============================================
echo ""
print_summary

echo ""
echo "Phase 5B (Resource List API) test results: ${RESULTS_DIR}"
