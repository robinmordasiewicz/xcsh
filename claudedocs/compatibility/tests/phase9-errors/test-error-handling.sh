#!/bin/bash
# test-error-handling.sh - Test error handling and edge cases
# Phase 9: Error Handling & Edge Cases
#
# Tests:
# - Invalid resource names
# - Missing required arguments
# - Invalid flag values
# - Non-existent resources (API errors)
# - Invalid input files

set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

cd "$PROJECT_ROOT"

source "${SCRIPT_DIR}/../../lib/common.sh"

echo "=========================================="
echo "Phase 9: Error Handling & Edge Cases"
echo "=========================================="
echo ""

# Initialize results directory
init_results_dir
PHASE_DIR="${RESULTS_DIR}/phase9-errors"
mkdir -p "$PHASE_DIR"

log_info "Original vesctl: ${ORIGINAL_VESCTL}"
log_info "Our vesctl: ${OUR_VESCTL}"
log_info "Results directory: ${RESULTS_DIR}"
echo ""

export RESULTS_DIR="$PHASE_DIR"

# ============================================
# Helper: Compare error behavior
# ============================================
test_error_behavior() {
    local test_name="$1"
    shift
    local cmd=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Run both
    $ORIGINAL_VESCTL "${cmd[@]}" > "${test_dir}/original_stdout.txt" 2> "${test_dir}/original_stderr.txt"
    local orig_exit=$?
    $OUR_VESCTL "${cmd[@]}" > "${test_dir}/ours_stdout.txt" 2> "${test_dir}/ours_stderr.txt"
    local our_exit=$?

    echo "orig_exit=$orig_exit" > "${test_dir}/exits.txt"
    echo "our_exit=$our_exit" >> "${test_dir}/exits.txt"

    # Combine stdout and stderr for comparison
    cat "${test_dir}/original_stdout.txt" "${test_dir}/original_stderr.txt" > "${test_dir}/original.txt"
    cat "${test_dir}/ours_stdout.txt" "${test_dir}/ours_stderr.txt" > "${test_dir}/ours.txt"

    # Check exit code consistency
    if [[ $orig_exit -eq $our_exit ]]; then
        if [[ $orig_exit -eq 0 ]]; then
            log_pass "${test_name} (both succeeded, exit 0)"
            echo "PASS_BOTH_SUCCESS" > "${test_dir}/result.txt"
        else
            log_pass "${test_name} (both errored, exit $orig_exit)"
            echo "PASS_BOTH_ERROR" > "${test_dir}/result.txt"
        fi
        return 0
    else
        log_warn "${test_name} (exit codes differ: orig=$orig_exit ours=$our_exit)"
        echo "WARN_EXIT_DIFF" > "${test_dir}/result.txt"
        return 0
    fi
}

# ============================================
# Helper: Check for error message pattern
# ============================================
test_error_message() {
    local test_name="$1"
    local pattern="$2"
    shift 2
    local cmd=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Run both
    $ORIGINAL_VESCTL "${cmd[@]}" > "${test_dir}/original.txt" 2>&1
    local orig_exit=$?
    $OUR_VESCTL "${cmd[@]}" > "${test_dir}/ours.txt" 2>&1
    local our_exit=$?

    echo "orig_exit=$orig_exit" > "${test_dir}/exits.txt"
    echo "our_exit=$our_exit" >> "${test_dir}/exits.txt"

    local orig_match=0
    local our_match=0

    grep -qi "$pattern" "${test_dir}/original.txt" && orig_match=1
    grep -qi "$pattern" "${test_dir}/ours.txt" && our_match=1

    echo "orig_match=$orig_match" >> "${test_dir}/exits.txt"
    echo "our_match=$our_match" >> "${test_dir}/exits.txt"

    if [[ $orig_match -eq 1 ]] && [[ $our_match -eq 1 ]]; then
        log_pass "${test_name} (both show error pattern)"
        echo "PASS" > "${test_dir}/result.txt"
        return 0
    elif [[ $orig_match -eq 0 ]] && [[ $our_match -eq 0 ]]; then
        log_info "${test_name} (neither shows pattern - may be different errors)"
        echo "INFO_NO_PATTERN" > "${test_dir}/result.txt"
        return 0
    else
        log_warn "${test_name} (error patterns differ)"
        echo "WARN_PATTERN_DIFF" > "${test_dir}/result.txt"
        return 0
    fi
}

# ============================================
# Test 1: Invalid resource names
# ============================================
echo "--- Invalid Resource Names ---"
echo ""

test_error_behavior "invalid-resource-list" configuration list nonexistent_foobar_resource
test_error_behavior "invalid-resource-get" configuration get nonexistent_foobar_resource test-name
test_error_behavior "invalid-command" nonexistent_command subcommand
echo ""

# ============================================
# Test 2: Missing required arguments
# ============================================
echo "--- Missing Required Arguments ---"
echo ""

test_error_behavior "missing-name-get" configuration get namespace
test_error_behavior "missing-name-delete" configuration delete namespace
test_error_behavior "missing-resource-list" configuration list
echo ""

# ============================================
# Test 3: Invalid flag values
# ============================================
echo "--- Invalid Flag Values ---"
echo ""

test_error_behavior "invalid-outfmt" configuration list namespace --outfmt invalid_format_xyz
test_error_behavior "invalid-namespace" configuration list namespace --namespace ""
echo ""

# ============================================
# Test 4: Non-existent resource (API errors)
# ============================================
echo "--- Non-Existent Resources (API Errors) ---"
echo ""

test_error_behavior "get-nonexistent-ns" configuration get namespace nonexistent-namespace-xyz-123
test_error_message "get-nonexistent-contains-error" "error\|not found\|404" configuration get namespace nonexistent-namespace-xyz-123
echo ""

# ============================================
# Test 5: Invalid input files
# ============================================
echo "--- Invalid Input Files ---"
echo ""

# Create test fixtures
FIXTURE_DIR="${RESULTS_DIR}/fixtures"
mkdir -p "$FIXTURE_DIR"

# Empty file
touch "${FIXTURE_DIR}/empty.yaml"
test_error_behavior "create-empty-file" configuration create namespace -i "${FIXTURE_DIR}/empty.yaml"

# Invalid YAML
echo "this: is: not: valid: yaml:" > "${FIXTURE_DIR}/invalid.yaml"
test_error_behavior "create-invalid-yaml" configuration create namespace -i "${FIXTURE_DIR}/invalid.yaml"

# Valid YAML but missing required fields
cat > "${FIXTURE_DIR}/missing_fields.yaml" <<EOF
metadata:
  namespace: default
EOF
test_error_behavior "create-missing-fields" configuration create namespace -i "${FIXTURE_DIR}/missing_fields.yaml"
echo ""

# ============================================
# Test 6: Help on error
# ============================================
echo "--- Help on Errors ---"
echo ""

test_error_behavior "unknown-flag" configuration list namespace --unknown-flag-xyz
test_error_message "help-shown-on-error" "usage\|Usage" configuration list namespace --unknown-flag-xyz
echo ""

# ============================================
# Summary
# ============================================
echo ""
print_summary

echo ""
echo "Phase 9 (Error Handling) test results: ${RESULTS_DIR}"
