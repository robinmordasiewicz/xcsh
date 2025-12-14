#!/bin/bash
# test-namespace-crud.sh - Comprehensive namespace CRUD compatibility testing
#
# This script performs pedantic comparison of namespace operations between
# the original vesctl and our implementation. It serves as a template for
# testing all other resources.
#
# Usage:
#   ./test-namespace-crud.sh              # Full test run
#   ./test-namespace-crud.sh --cleanup-only   # Cleanup only (recover from failed run)
#   ./test-namespace-crud.sh --skip-cleanup   # Skip cleanup (debugging)
#
# Required environment variables:
#   F5XC_P12_PASSWORD - Password for the P12 certificate
#   F5XC_P12_FILE     - Path to the P12 certificate file

set -euo pipefail

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source common utilities
source "${SCRIPT_DIR}/../../lib/common.sh"

# Source sweeper
source "${SCRIPT_DIR}/lib/sweeper.sh"

# Test configuration
TEST_PREFIX="vesctl-test-"
TEST_NAMESPACE_1="${TEST_PREFIX}ns-$(date +%s)-1"
TEST_NAMESPACE_2="${TEST_PREFIX}ns-$(date +%s)-2"

# Command line options
CLEANUP_ONLY=false
SKIP_CLEANUP=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --cleanup-only)
            CLEANUP_ONLY=true
            shift
            ;;
        --skip-cleanup)
            SKIP_CLEANUP=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# ============================================
# Pre-flight Checks
# ============================================

preflight_checks() {
    log_info "Running pre-flight checks..."

    # Check binaries
    if ! check_original_vesctl; then
        echo "Please download the original vesctl binary first:"
        echo "  ./claudedocs/compatibility/scripts/download-original.sh"
        exit 1
    fi

    if ! check_our_vesctl; then
        echo "Please build our f5xcctl binary first:"
        echo "  go build -o vesctl.darwin-arm64 ."
        exit 1
    fi

    # Check API credentials
    if ! has_api_credentials; then
        echo -e "${RED}ERROR: API credentials not configured${NC}"
        echo "Configure either:"
        echo "  1. P12 credentials via environment variables:"
        echo "     F5XC_P12_PASSWORD - Password for the P12 certificate"
        echo "     F5XC_P12_FILE     - Path to the P12 certificate file"
        echo "  2. Key/cert credentials in ~/.vesconfig:"
        echo "     key: /path/to/key.key"
        echo "     cert: /path/to/cert.cert"
        exit 1
    fi

    log_pass "Pre-flight checks passed"
}

# ============================================
# Help Comparison Tests
# ============================================

run_help_tests() {
    log_info "=========================================="
    log_info "Phase: Help Comparison Tests"
    log_info "=========================================="

    # Test: configuration namespace --help
    test_help "help-configuration-namespace" configuration namespace --help

    # Test: configuration create namespace --help
    test_help "help-create-namespace" configuration create namespace --help

    # Test: configuration get namespace --help
    test_help "help-get-namespace" configuration get namespace --help

    # Test: configuration list namespace --help
    test_help "help-list-namespace" configuration list namespace --help

    # Test: configuration replace namespace --help
    test_help "help-replace-namespace" configuration replace namespace --help

    # Test: configuration delete namespace --help
    test_help "help-delete-namespace" configuration delete namespace --help

    # Test: flags section comparison for create
    test_flags_section "flags-create-namespace" configuration create namespace
}

# ============================================
# List Tests (Pre-Create Baseline)
# ============================================

run_list_baseline_tests() {
    log_info "=========================================="
    log_info "Phase: List Tests (Pre-Create Baseline)"
    log_info "=========================================="

    # Test: list namespaces (default format)
    run_api_test "list-baseline-default" configuration list namespace

    # Test: list namespaces (JSON format)
    run_api_test "list-baseline-json" configuration list namespace --outfmt json

    # Test: list namespaces (YAML format)
    run_api_test "list-baseline-yaml" configuration list namespace --outfmt yaml

    # Test: list namespaces (with label selector - should return empty or filtered)
    run_api_test "list-baseline-label-selector" configuration list namespace -l "test-label=nonexistent"
}

# ============================================
# Create Tests
# ============================================

run_create_tests() {
    log_info "=========================================="
    log_info "Phase: Create Tests"
    log_info "=========================================="

    # Create fixture files
    local fixtures_dir="${SCRIPT_DIR}/fixtures"
    mkdir -p "$fixtures_dir"

    # Create namespace YAML fixture
    cat > "${fixtures_dir}/namespace-create-1.yaml" << EOF
metadata:
  name: ${TEST_NAMESPACE_1}
  labels:
    test-type: compatibility
    created-by: vesctl-test
  annotations:
    description: "Test namespace for vesctl compatibility testing"
spec: {}
EOF

    # Create second namespace YAML fixture
    cat > "${fixtures_dir}/namespace-create-2.yaml" << EOF
metadata:
  name: ${TEST_NAMESPACE_2}
  labels:
    test-type: compatibility
    created-by: vesctl-test
  annotations:
    description: "Second test namespace for vesctl compatibility testing"
spec: {}
EOF

    # Test: create namespace from file
    # Note: Use --ours-first so our version creates the namespace first,
    # and original gets 409 conflict (same as if original ran first but reversed)
    log_info "Creating test namespace 1: ${TEST_NAMESPACE_1}"
    run_api_test --ours-first "create-from-file" configuration create namespace -i "${fixtures_dir}/namespace-create-1.yaml"
    track_resource "namespace" "${TEST_NAMESPACE_1}" ""

    # Wait for namespace to be ready
    sleep 2

    # Test: create second namespace (for testing list with multiple items)
    log_info "Creating test namespace 2: ${TEST_NAMESPACE_2}"
    run_api_test --ours-first "create-second-namespace" configuration create namespace -i "${fixtures_dir}/namespace-create-2.yaml"
    track_resource "namespace" "${TEST_NAMESPACE_2}" ""

    # Wait for namespace to be ready
    sleep 2

    # Test: create duplicate namespace (should fail)
    run_api_test "create-duplicate-error" configuration create namespace -i "${fixtures_dir}/namespace-create-1.yaml"

    # Test: create with JSON data (if supported)
    local json_data="{\"metadata\":{\"name\":\"${TEST_PREFIX}json-test\"},\"spec\":{}}"
    # Note: This may or may not be supported - we compare behavior
    run_api_test "create-json-data" configuration create namespace --json-data "$json_data" || true
}

# ============================================
# Get Tests
# ============================================

run_get_tests() {
    log_info "=========================================="
    log_info "Phase: Get Tests"
    log_info "=========================================="

    # Test: get namespace (default format)
    run_api_test "get-default" configuration get namespace "${TEST_NAMESPACE_1}"

    # Test: get namespace (JSON format)
    run_api_test "get-json" configuration get namespace "${TEST_NAMESPACE_1}" --outfmt json

    # Test: get namespace (YAML format)
    run_api_test "get-yaml" configuration get namespace "${TEST_NAMESPACE_1}" --outfmt yaml

    # Test: get namespace with response-format proto
    run_api_test "get-response-proto" configuration get namespace "${TEST_NAMESPACE_1}" --response-format GET_RSP_FORMAT_READ

    # Test: get second namespace
    run_api_test "get-second-namespace" configuration get namespace "${TEST_NAMESPACE_2}" --outfmt json

    # Test: get nonexistent namespace (should fail)
    run_api_test "get-nonexistent" configuration get namespace "${TEST_PREFIX}nonexistent"

    # Test: get system namespace (should succeed)
    run_api_test "get-system" configuration get namespace "system" --outfmt json
}

# ============================================
# List Tests (Post-Create)
# ============================================

run_list_post_create_tests() {
    log_info "=========================================="
    log_info "Phase: List Tests (Post-Create)"
    log_info "=========================================="

    # Test: list namespaces (should include our test namespaces)
    run_api_test "list-post-create-json" configuration list namespace --outfmt json

    # Test: list with label selector (should find our test namespaces)
    run_api_test "list-post-create-label" configuration list namespace -l "test-type=compatibility" --outfmt json

    # Test: list with specific label value
    run_api_test "list-post-create-created-by" configuration list namespace -l "created-by=vesctl-test" --outfmt json
}

# ============================================
# Replace Tests
# ============================================

run_replace_tests() {
    log_info "=========================================="
    log_info "Phase: Replace Tests"
    log_info "=========================================="

    local fixtures_dir="${SCRIPT_DIR}/fixtures"

    # Create updated namespace YAML fixture
    cat > "${fixtures_dir}/namespace-replace.yaml" << EOF
metadata:
  name: ${TEST_NAMESPACE_1}
  labels:
    test-type: compatibility
    created-by: vesctl-test
    updated: "true"
  annotations:
    description: "Updated test namespace for vesctl compatibility testing"
spec: {}
EOF

    # Test: replace namespace from file
    run_api_test "replace-from-file" configuration replace namespace -i "${fixtures_dir}/namespace-replace.yaml"

    # Wait for update to propagate
    sleep 2

    # Test: verify replacement
    run_api_test "replace-verify" configuration get namespace "${TEST_NAMESPACE_1}" --outfmt json

    # Test: replace nonexistent namespace (should fail)
    cat > "${fixtures_dir}/namespace-replace-nonexistent.yaml" << EOF
metadata:
  name: ${TEST_PREFIX}nonexistent-replace
  labels:
    test-type: compatibility
spec: {}
EOF
    run_api_test "replace-nonexistent" configuration replace namespace -i "${fixtures_dir}/namespace-replace-nonexistent.yaml"
}

# ============================================
# Delete Tests
# ============================================

run_delete_tests() {
    log_info "=========================================="
    log_info "Phase: Delete Tests"
    log_info "=========================================="

    # Test: delete first namespace
    log_info "Deleting test namespace 1: ${TEST_NAMESPACE_1}"
    run_api_test "delete-first" configuration delete namespace "${TEST_NAMESPACE_1}"
    mark_deleted "namespace" "${TEST_NAMESPACE_1}"

    # Wait for deletion
    sleep 3

    # Test: verify deletion (should fail - namespace gone)
    run_api_test "delete-verify-gone" configuration get namespace "${TEST_NAMESPACE_1}"

    # Test: delete second namespace
    log_info "Deleting test namespace 2: ${TEST_NAMESPACE_2}"
    run_api_test "delete-second" configuration delete namespace "${TEST_NAMESPACE_2}"
    mark_deleted "namespace" "${TEST_NAMESPACE_2}"

    # Wait for deletion
    sleep 3

    # Test: delete nonexistent namespace (should fail gracefully)
    run_api_test "delete-nonexistent" configuration delete namespace "${TEST_PREFIX}nonexistent-delete"
}

# ============================================
# Cleanup Phase
# ============================================

run_cleanup_phase() {
    log_info "=========================================="
    log_info "Phase: Cleanup"
    log_info "=========================================="

    # Run full cleanup sweep using our binary (it should work like original)
    run_cleanup_sweep "$OUR_F5XCCTL" "$RESULTS_DIR"

    # Verify cleanup
    verify_cleanup "$OUR_F5XCCTL"

    # Print cleanup summary
    get_cleanup_summary
}

# ============================================
# Generate Report
# ============================================

generate_report() {
    log_info "=========================================="
    log_info "Generating Report"
    log_info "=========================================="

    local report_file="${RESULTS_DIR}/REPORT.md"

    cat > "$report_file" << EOF
# Phase 4: Namespace CRUD Compatibility Test Report

**Date**: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
**Results Directory**: ${RESULTS_DIR}

## Test Summary

| Metric | Value |
|--------|-------|
| Passed | ${TESTS_PASSED} |
| Failed | ${TESTS_FAILED} |
| Warned | ${TESTS_WARNED} |
| Skipped | ${TESTS_SKIPPED} |
| Pass Rate | $((TESTS_PASSED * 100 / (TESTS_PASSED + TESTS_FAILED + 1)))% |

## Binary Information

| Binary | Path |
|--------|------|
| Original | ${ORIGINAL_VESCTL} |
| Ours | ${OUR_F5XCCTL} |

## Test Namespaces

| Namespace | Purpose |
|-----------|---------|
| ${TEST_NAMESPACE_1} | Primary test namespace |
| ${TEST_NAMESPACE_2} | Secondary test namespace |

## Test Phases

### Help Comparison
Tests that compare --help output for namespace subcommands.

### List Tests (Baseline)
Tests that compare list output before creating test namespaces.

### Create Tests
Tests that compare namespace creation behavior.

### Get Tests
Tests that compare namespace retrieval in various formats.

### List Tests (Post-Create)
Tests that compare list output after creating test namespaces.

### Replace Tests
Tests that compare namespace update behavior.

### Delete Tests
Tests that compare namespace deletion behavior.

### Cleanup
Sweeper validation to ensure all test resources are removed.

## Failed Tests

EOF

    # List failed tests
    if [[ -d "$RESULTS_DIR" ]]; then
        find "$RESULTS_DIR" -name "result.txt" -exec sh -c '
            result=$(cat "$1")
            if [[ "$result" == "FAIL"* || "$result" == "EXIT_MISMATCH"* || "$result" == "STDOUT_DIFFERS"* ]]; then
                test_name=$(dirname "$1" | xargs basename)
                echo "- **${test_name}**: ${result}"
            fi
        ' _ {} \; >> "$report_file"
    fi

    cat >> "$report_file" << EOF

## Diff Analysis

For detailed diff analysis of each test, see the individual test directories:
\`${RESULTS_DIR}/<test-name>/\`

Each test directory contains:
- \`original.stdout\` / \`ours.stdout\` - Raw stdout output
- \`original.stderr\` / \`ours.stderr\` - Raw stderr output
- \`*.normalized.json\` - Normalized JSON for comparison
- \`stdout.diff\` - Unified diff of stdout
- \`diff_analysis.txt\` - Human-readable diff analysis
- \`result.txt\` - Test result (PASS/FAIL/etc.)

## Cleanup Summary

$(get_cleanup_summary)

---
Generated by vesctl compatibility testing framework
EOF

    log_info "Report generated: ${report_file}"
}

# ============================================
# Main Execution
# ============================================

main() {
    echo ""
    echo "=========================================="
    echo "Phase 4: Namespace CRUD Compatibility Test"
    echo "=========================================="
    echo ""

    # Handle cleanup-only mode
    if [[ "$CLEANUP_ONLY" == "true" ]]; then
        log_info "Running in cleanup-only mode"
        RESULTS_DIR="${SCRIPT_DIR}/results/cleanup-$(date +%s)"
        mkdir -p "$RESULTS_DIR"
        init_sweeper "$RESULTS_DIR"
        run_cleanup_phase
        print_summary
        exit $?
    fi

    # Pre-flight checks
    preflight_checks

    # Initialize results directory
    RESULTS_DIR=$(init_results_dir)
    log_info "Results directory: ${RESULTS_DIR}"

    # Initialize sweeper
    init_sweeper "$RESULTS_DIR"

    # Run test phases
    run_help_tests
    run_list_baseline_tests
    run_create_tests
    run_get_tests
    run_list_post_create_tests
    run_replace_tests
    run_delete_tests

    # Cleanup phase (unless skipped)
    if [[ "$SKIP_CLEANUP" != "true" ]]; then
        run_cleanup_phase
    else
        log_warn "Cleanup skipped (--skip-cleanup flag)"
    fi

    # Generate report
    generate_report

    # Print summary
    echo ""
    print_summary

    # Return appropriate exit code
    if [[ $TESTS_FAILED -gt 0 ]]; then
        exit 1
    fi
    exit 0
}

# Run main
main "$@"
