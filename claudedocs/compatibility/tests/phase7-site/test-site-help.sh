#!/bin/bash
# test-site-help.sh - Test site command help text compatibility
# Phase 7: Site Management Commands Validation
#
# Tests:
# - site --help
# - site aws_vpc --help (and subcommands)
# - site azure_vnet --help (and subcommands)

set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# Change to project root so binary paths work correctly
cd "$PROJECT_ROOT"

source "${SCRIPT_DIR}/../../lib/common.sh"

echo "=========================================="
echo "Phase 7: Site Management Commands Tests"
echo "=========================================="
echo ""

# Initialize results directory
init_results_dir
PHASE_DIR="${RESULTS_DIR}/phase7-site"
mkdir -p "$PHASE_DIR"

log_info "Original vesctl: ${ORIGINAL_VESCTL}"
log_info "Our f5xcctl: ${OUR_F5XCCTL}"
log_info "Results directory: ${RESULTS_DIR}"
echo ""

export RESULTS_DIR="$PHASE_DIR"

# ============================================
# Helper: Compare full help text
# ============================================
test_help_identical() {
    local test_name="$1"
    shift
    local cmd=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    $ORIGINAL_VESCTL "${cmd[@]}" --help > "${test_dir}/original.txt" 2>&1
    local orig_exit=$?
    $OUR_F5XCCTL "${cmd[@]}" --help > "${test_dir}/ours.txt" 2>&1
    local our_exit=$?

    echo "orig_exit=$orig_exit" > "${test_dir}/exits.txt"
    echo "our_exit=$our_exit" >> "${test_dir}/exits.txt"

    if [[ $orig_exit -ne 0 ]]; then
        if [[ $our_exit -eq 0 ]]; then
            log_warn "${test_name} (not in original - our enhancement)"
            echo "WARN_NEW_COMMAND" > "${test_dir}/result.txt"
            return 0
        else
            log_warn "${test_name} (not in either)"
            echo "WARN_MISSING" > "${test_dir}/result.txt"
            return 0
        fi
    fi

    if diff -q "${test_dir}/original.txt" "${test_dir}/ours.txt" > /dev/null 2>&1; then
        log_pass "${test_name} (identical)"
        echo "PASS" > "${test_dir}/result.txt"
        return 0
    else
        diff -u "${test_dir}/original.txt" "${test_dir}/ours.txt" > "${test_dir}/diff.txt" 2>&1 || true
        log_fail "${test_name} (differs)"
        echo "FAIL" > "${test_dir}/result.txt"
        return 1
    fi
}

# ============================================
# Helper: Compare flags only
# ============================================
test_flags_count() {
    local test_name="$1"
    shift
    local cmd=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Count flags in each version
    local orig_flags=$($ORIGINAL_VESCTL "${cmd[@]}" --help 2>&1 | grep -E '^\s+(-|--)[a-z]' | wc -l)
    local our_flags=$($OUR_F5XCCTL "${cmd[@]}" --help 2>&1 | grep -E '^\s+(-|--)[a-z]' | wc -l)

    echo "original_flags=$orig_flags" > "${test_dir}/flag_count.txt"
    echo "our_flags=$our_flags" >> "${test_dir}/flag_count.txt"

    # Extract flag names only
    $ORIGINAL_VESCTL "${cmd[@]}" --help 2>&1 | grep -E '^\s+(-|--)[a-z]' | sed 's/^\s*//' | awk '{print $1}' | sort > "${test_dir}/original_flag_names.txt"
    $OUR_F5XCCTL "${cmd[@]}" --help 2>&1 | grep -E '^\s+(-|--)[a-z]' | sed 's/^\s*//' | awk '{print $1}' | sort > "${test_dir}/our_flag_names.txt"

    # Compare
    local common=$(comm -12 "${test_dir}/original_flag_names.txt" "${test_dir}/our_flag_names.txt" | wc -l)
    local orig_only=$(comm -23 "${test_dir}/original_flag_names.txt" "${test_dir}/our_flag_names.txt" | wc -l)
    local ours_only=$(comm -13 "${test_dir}/original_flag_names.txt" "${test_dir}/our_flag_names.txt" | wc -l)

    echo "common=$common" >> "${test_dir}/flag_count.txt"
    echo "original_only=$orig_only" >> "${test_dir}/flag_count.txt"
    echo "ours_only=$ours_only" >> "${test_dir}/flag_count.txt"

    # Save the differences
    comm -23 "${test_dir}/original_flag_names.txt" "${test_dir}/our_flag_names.txt" > "${test_dir}/missing_flags.txt"
    comm -13 "${test_dir}/original_flag_names.txt" "${test_dir}/our_flag_names.txt" > "${test_dir}/new_flags.txt"

    if [[ $orig_only -eq 0 ]] && [[ $ours_only -eq 0 ]]; then
        log_pass "${test_name} (${common} flags identical)"
        echo "PASS" > "${test_dir}/result.txt"
    elif [[ $orig_only -gt 0 ]]; then
        log_warn "${test_name} (${common} common, ${orig_only} missing, ${ours_only} new)"
        echo "WARN_MISSING_FLAGS" > "${test_dir}/result.txt"
    else
        log_info "${test_name} (${common} common, ${ours_only} new flags)"
        echo "PASS_NEW_FLAGS" > "${test_dir}/result.txt"
    fi
}

# ============================================
# Test site main help
# ============================================
echo "--- Site Main Command ---"
echo ""

test_help_identical "site-main" site || true
echo ""

# ============================================
# Test site aws_vpc
# ============================================
echo "--- Site AWS VPC ---"
echo ""

test_help_identical "site-aws-vpc" site aws_vpc || true
test_help_identical "site-aws-vpc-create" site aws_vpc create || true
test_flags_count "site-aws-vpc-create-flags" site aws_vpc create || true
test_help_identical "site-aws-vpc-delete" site aws_vpc delete || true
test_help_identical "site-aws-vpc-replace" site aws_vpc replace || true
test_help_identical "site-aws-vpc-run" site aws_vpc run || true
test_flags_count "site-aws-vpc-run-flags" site aws_vpc run || true
echo ""

# ============================================
# Test site azure_vnet
# ============================================
echo "--- Site Azure VNet ---"
echo ""

test_help_identical "site-azure-vnet" site azure_vnet || true
test_help_identical "site-azure-vnet-create" site azure_vnet create || true
test_flags_count "site-azure-vnet-create-flags" site azure_vnet create || true
test_help_identical "site-azure-vnet-delete" site azure_vnet delete || true
test_help_identical "site-azure-vnet-replace" site azure_vnet replace || true
test_help_identical "site-azure-vnet-run" site azure_vnet run || true
test_flags_count "site-azure-vnet-run-flags" site azure_vnet run || true
echo ""

# ============================================
# Summary
# ============================================
echo ""
print_summary

echo ""
echo "Phase 7 (Site Commands) test results: ${RESULTS_DIR}"
