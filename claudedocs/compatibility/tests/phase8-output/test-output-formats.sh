#!/bin/bash
# test-output-formats.sh - Test output format precision
# Phase 8: Output Format Compatibility Testing
#
# Tests:
# - Table format (list operations)
# - JSON format (get operations)
# - YAML format (get operations)
# - TSV format (list operations)

set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

cd "$PROJECT_ROOT"

source "${SCRIPT_DIR}/../../lib/common.sh"

echo "=========================================="
echo "Phase 8: Output Format Precision Tests"
echo "=========================================="
echo ""

# Initialize results directory
init_results_dir
PHASE_DIR="${RESULTS_DIR}/phase8-output"
mkdir -p "$PHASE_DIR"

log_info "Original vesctl: ${ORIGINAL_VESCTL}"
log_info "Our f5xcctl: ${OUR_F5XCCTL}"
log_info "Results directory: ${RESULTS_DIR}"
echo ""

export RESULTS_DIR="$PHASE_DIR"

# ============================================
# Helper: Compare byte-for-byte
# ============================================
test_output_identical() {
    local test_name="$1"
    shift
    local cmd_orig=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Run original
    $ORIGINAL_VESCTL "${cmd_orig[@]}" > "${test_dir}/original.txt" 2>&1
    local orig_exit=$?

    # Run ours with same command
    $OUR_F5XCCTL "${cmd_orig[@]}" > "${test_dir}/ours.txt" 2>&1
    local our_exit=$?

    echo "orig_exit=$orig_exit" > "${test_dir}/exits.txt"
    echo "our_exit=$our_exit" >> "${test_dir}/exits.txt"

    if [[ $orig_exit -ne 0 ]]; then
        log_warn "${test_name} (original failed with exit $orig_exit)"
        echo "WARN_ORIG_FAIL" > "${test_dir}/result.txt"
        return 0
    fi

    if [[ $our_exit -ne 0 ]]; then
        log_fail "${test_name} (ours failed with exit $our_exit)"
        echo "FAIL_OUR_EXIT" > "${test_dir}/result.txt"
        return 1
    fi

    if diff -q "${test_dir}/original.txt" "${test_dir}/ours.txt" > /dev/null 2>&1; then
        log_pass "${test_name} (byte-identical)"
        echo "PASS" > "${test_dir}/result.txt"
        return 0
    else
        diff -u "${test_dir}/original.txt" "${test_dir}/ours.txt" > "${test_dir}/diff.txt" 2>&1 || true

        # Check if it's just key ordering difference (JSON)
        if command -v jq &> /dev/null && [[ "$test_name" == *json* ]]; then
            # Sort JSON keys and compare
            jq -S '.' "${test_dir}/original.txt" > "${test_dir}/original_sorted.json" 2>/dev/null || true
            jq -S '.' "${test_dir}/ours.txt" > "${test_dir}/ours_sorted.json" 2>/dev/null || true

            if diff -q "${test_dir}/original_sorted.json" "${test_dir}/ours_sorted.json" > /dev/null 2>&1; then
                log_info "${test_name} (semantically equivalent - key order differs)"
                echo "PASS_KEY_ORDER" > "${test_dir}/result.txt"
                return 0
            fi
        fi

        log_fail "${test_name} (differs)"
        echo "FAIL" > "${test_dir}/result.txt"
        return 1
    fi
}

# ============================================
# Helper: Semantic JSON comparison
# ============================================
test_json_semantic() {
    local test_name="$1"
    shift
    local cmd_orig=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Run both
    $ORIGINAL_VESCTL "${cmd_orig[@]}" > "${test_dir}/original.json" 2>&1
    local orig_exit=$?
    $OUR_F5XCCTL "${cmd_orig[@]}" > "${test_dir}/ours.json" 2>&1
    local our_exit=$?

    echo "orig_exit=$orig_exit" > "${test_dir}/exits.txt"
    echo "our_exit=$our_exit" >> "${test_dir}/exits.txt"

    if [[ $orig_exit -ne 0 ]]; then
        log_warn "${test_name} (original failed)"
        echo "WARN" > "${test_dir}/result.txt"
        return 0
    fi

    if ! command -v jq &> /dev/null; then
        log_warn "${test_name} (jq not installed, skipping semantic check)"
        echo "SKIP_NO_JQ" > "${test_dir}/result.txt"
        return 0
    fi

    # Sort and compare
    jq -S '.' "${test_dir}/original.json" > "${test_dir}/original_sorted.json" 2>/dev/null
    jq -S '.' "${test_dir}/ours.json" > "${test_dir}/ours_sorted.json" 2>/dev/null

    if diff -q "${test_dir}/original_sorted.json" "${test_dir}/ours_sorted.json" > /dev/null 2>&1; then
        log_pass "${test_name} (semantically identical)"
        echo "PASS" > "${test_dir}/result.txt"
        return 0
    else
        diff -u "${test_dir}/original_sorted.json" "${test_dir}/ours_sorted.json" > "${test_dir}/diff.txt" 2>&1 || true
        log_fail "${test_name} (semantic diff)"
        echo "FAIL" > "${test_dir}/result.txt"
        return 1
    fi
}

# ============================================
# Test Table format
# ============================================
echo "--- Table Format Tests ---"
echo ""

test_output_identical "table-list-namespace" configuration list namespace
echo ""

# ============================================
# Test JSON format
# ============================================
echo "--- JSON Format Tests ---"
echo ""

test_output_identical "json-get-namespace" configuration get namespace default --outfmt json
test_json_semantic "json-get-namespace-semantic" configuration get namespace default --outfmt json
echo ""

# ============================================
# Test YAML format
# ============================================
echo "--- YAML Format Tests ---"
echo ""

test_output_identical "yaml-get-namespace" configuration get namespace default --outfmt yaml
echo ""

# ============================================
# Test TSV format
# ============================================
echo "--- TSV Format Tests ---"
echo ""

test_output_identical "tsv-list-namespace" configuration list namespace --outfmt tsv
echo ""

# ============================================
# Test multiple resources
# ============================================
echo "--- Multiple Resource Tests ---"
echo ""

# Get a specific namespace to compare more fields
test_output_identical "json-get-system" configuration get namespace system --outfmt json
test_json_semantic "json-get-system-semantic" configuration get namespace system --outfmt json
test_output_identical "yaml-get-system" configuration get namespace system --outfmt yaml
echo ""

# ============================================
# Summary
# ============================================
echo ""
print_summary

echo ""
echo "Phase 8 (Output Format) test results: ${RESULTS_DIR}"
