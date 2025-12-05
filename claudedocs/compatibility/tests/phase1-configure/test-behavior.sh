#!/bin/bash
# test-behavior.sh - Test behavioral consistency (structure, style, shared resources)
# Ignores content differences due to API spec evolution

# Don't exit on errors - we want to run all tests
set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/common.sh"

echo "=========================================="
echo "Behavioral Consistency Tests"
echo "=========================================="
echo ""

# Initialize results directory
init_results_dir
PHASE_DIR="${RESULTS_DIR}/behavior"
mkdir -p "$PHASE_DIR"

log_info "Original vesctl: ${ORIGINAL_VESCTL}"
log_info "Our vesctl: ${OUR_VESCTL}"
log_info "Results directory: ${RESULTS_DIR}"
echo ""

# Override RESULTS_DIR for test functions
export RESULTS_DIR="$PHASE_DIR"

# ============================================
# Test 1: Global Flags Consistency
# ============================================
echo "--- Test Category: Global Flags ---"

test_global_flags() {
    local cmd="$1"
    local test_name="global-flags-${cmd//[ ]/-}"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Extract Global Flags section only (use eval for multi-word commands)
    eval "$ORIGINAL_VESCTL $cmd --help" 2>&1 | sed -n '/^Global Flags:/,/^$/p' > "${test_dir}/original.txt"
    eval "$OUR_VESCTL $cmd --help" 2>&1 | sed -n '/^Global Flags:/,/^$/p' > "${test_dir}/ours.txt"

    if diff -q "${test_dir}/original.txt" "${test_dir}/ours.txt" > /dev/null 2>&1; then
        log_pass "${test_name}"
        echo "PASS" > "${test_dir}/result.txt"
    else
        diff -u "${test_dir}/original.txt" "${test_dir}/ours.txt" > "${test_dir}/diff.txt" 2>&1 || true
        log_fail "${test_name}"
        echo "FAIL" > "${test_dir}/result.txt"
    fi
}

test_global_flags "configuration"
test_global_flags "configuration list"
test_global_flags "configuration get"

# ============================================
# Test 2: Shared Resource Help (exact match expected)
# ============================================
echo ""
echo "--- Test Category: Shared Resource Help ---"

# Resources that exist in both binaries
SHARED_RESOURCES=(
    "namespace"
    "http_loadbalancer"
    "origin_pool"
    "virtual_host"
    "service_policy"
    "app_firewall"
    "certificate"
    "healthcheck"
    "cloud_credentials"
    "aws_vpc_site"
    "azure_vnet_site"
    "gcp_vpc_site"
)

for resource in "${SHARED_RESOURCES[@]}"; do
    test_help "shared-list-${resource}" configuration list "$resource" --help || true
    test_help "shared-get-${resource}" configuration get "$resource" --help || true
done

# ============================================
# Test 3: Help Text Structure
# ============================================
echo ""
echo "--- Test Category: Help Structure ---"

test_help_structure() {
    local cmd="$1"
    local test_name="structure-${cmd//[ ]/-}"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Extract structural elements only (section headers)
    eval "$ORIGINAL_VESCTL $cmd --help" 2>&1 | grep -E "^(Usage:|Aliases:|Examples:|Available Commands:|Flags:|Global Flags:|Use )" > "${test_dir}/original.txt" 2>/dev/null || true
    eval "$OUR_VESCTL $cmd --help" 2>&1 | grep -E "^(Usage:|Aliases:|Examples:|Available Commands:|Flags:|Global Flags:|Use )" > "${test_dir}/ours.txt" 2>/dev/null || true

    if diff -q "${test_dir}/original.txt" "${test_dir}/ours.txt" > /dev/null 2>&1; then
        log_pass "${test_name}"
        echo "PASS" > "${test_dir}/result.txt"
    else
        diff -u "${test_dir}/original.txt" "${test_dir}/ours.txt" > "${test_dir}/diff.txt" 2>&1 || true
        log_fail "${test_name}"
        echo "FAIL" > "${test_dir}/result.txt"
    fi
}

test_help_structure "configuration"
test_help_structure "configuration list"
test_help_structure "configuration get"
test_help_structure "configuration create"
test_help_structure "configuration delete"

# ============================================
# Test 4: Error Message Style
# ============================================
echo ""
echo "--- Test Category: Error Message Style ---"

test_error_style() {
    local test_name="$1"
    shift
    local args=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Capture error output
    $ORIGINAL_VESCTL "${args[@]}" > "${test_dir}/original.stdout" 2> "${test_dir}/original.stderr" || true
    $OUR_VESCTL "${args[@]}" > "${test_dir}/ours.stdout" 2> "${test_dir}/ours.stderr" || true

    # Extract error pattern (first line of stderr or stdout containing "Error" or "error")
    local orig_error=$(grep -i "error" "${test_dir}/original.stderr" "${test_dir}/original.stdout" 2>/dev/null | head -1 | sed 's/vesctl\.[a-z0-9-]*/vesctl/g')
    local our_error=$(grep -i "error" "${test_dir}/ours.stderr" "${test_dir}/ours.stdout" 2>/dev/null | head -1 | sed 's/vesctl\.[a-z0-9-]*/vesctl/g')

    echo "$orig_error" > "${test_dir}/original_error.txt"
    echo "$our_error" > "${test_dir}/ours_error.txt"

    if [[ "$orig_error" == "$our_error" ]]; then
        log_pass "${test_name}"
        echo "PASS" > "${test_dir}/result.txt"
    else
        # Check if error style is similar (contains same key phrases)
        if [[ "$orig_error" == *"unknown command"* ]] && [[ "$our_error" == *"unknown command"* ]]; then
            log_pass "${test_name} (style match)"
            echo "PASS" > "${test_dir}/result.txt"
        elif [[ "$orig_error" == *"required"* ]] && [[ "$our_error" == *"required"* ]]; then
            log_pass "${test_name} (style match)"
            echo "PASS" > "${test_dir}/result.txt"
        else
            log_warn "${test_name}: error style differs"
            echo "orig: $orig_error" >> "${test_dir}/diff.txt"
            echo "ours: $our_error" >> "${test_dir}/diff.txt"
            echo "WARN" > "${test_dir}/result.txt"
        fi
    fi
}

test_error_style "error-unknown-cmd" invalid_command
test_error_style "error-missing-resource" configuration list
test_error_style "error-invalid-flag" configuration --invalid-flag

# ============================================
# Test 5: Command Aliases
# ============================================
echo ""
echo "--- Test Category: Command Aliases ---"

test_alias() {
    local test_name="$1"
    local alias_cmd="$2"
    local full_cmd="$3"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Compare alias output with full command (for our binary)
    $OUR_VESCTL $alias_cmd --help 2>&1 > "${test_dir}/alias.txt"
    $OUR_VESCTL $full_cmd --help 2>&1 > "${test_dir}/full.txt"

    # Normalize (replace command names in Usage line)
    sed -i.bak "s/$alias_cmd/$full_cmd/g" "${test_dir}/alias.txt" 2>/dev/null || \
        sed "s/$alias_cmd/$full_cmd/g" "${test_dir}/alias.txt" > "${test_dir}/alias_normalized.txt"

    if diff -q "${test_dir}/alias.txt" "${test_dir}/full.txt" > /dev/null 2>&1 || \
       diff -q "${test_dir}/alias_normalized.txt" "${test_dir}/full.txt" > /dev/null 2>&1; then
        log_pass "${test_name}"
        echo "PASS" > "${test_dir}/result.txt"
    else
        log_fail "${test_name}"
        echo "FAIL" > "${test_dir}/result.txt"
    fi
}

test_alias "alias-cfg" "cfg" "configuration"
test_alias "alias-c" "c" "configuration"

# ============================================
# Summary
# ============================================
echo ""
print_summary

# Generate reports
export RESULTS_DIR=$(dirname "$PHASE_DIR")
generate_markdown_report "$RESULTS_DIR"

echo ""
echo "Behavioral test results: ${RESULTS_DIR}/summary.md"
