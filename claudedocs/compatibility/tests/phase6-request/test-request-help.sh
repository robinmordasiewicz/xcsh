#!/bin/bash
# test-request-help.sh - Test request command help text compatibility
# Phase 6: Request Commands Validation
#
# Tests:
# - request --help
# - request rpc --help (and RPC command list comparison)
# - request secrets --help (and subcommand help)
# - request command-sequence --help

set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# Change to project root so binary paths work correctly
cd "$PROJECT_ROOT"

source "${SCRIPT_DIR}/../../lib/common.sh"

echo "=========================================="
echo "Phase 6: Request Commands Tests"
echo "=========================================="
echo ""

# Initialize results directory
init_results_dir
PHASE_DIR="${RESULTS_DIR}/phase6-request"
mkdir -p "$PHASE_DIR"

log_info "Original vesctl: ${ORIGINAL_VESCTL}"
log_info "Our vesctl: ${OUR_VESCTL}"
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
    $OUR_VESCTL "${cmd[@]}" --help > "${test_dir}/ours.txt" 2>&1
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
# Helper: Compare flags only (ignore description text)
# ============================================
test_flags_identical() {
    local test_name="$1"
    shift
    local cmd=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Extract just the flags (lines starting with whitespace and -)
    $ORIGINAL_VESCTL "${cmd[@]}" --help 2>&1 | grep -E '^\s+(-|--)[a-z]' > "${test_dir}/original_flags.txt"
    local orig_exit=$?
    $OUR_VESCTL "${cmd[@]}" --help 2>&1 | grep -E '^\s+(-|--)[a-z]' > "${test_dir}/ours_flags.txt"
    local our_exit=$?

    if [[ ! -s "${test_dir}/original_flags.txt" ]]; then
        log_skip "${test_name} (no flags to compare)"
        echo "SKIP" > "${test_dir}/result.txt"
        return 0
    fi

    if diff -q "${test_dir}/original_flags.txt" "${test_dir}/ours_flags.txt" > /dev/null 2>&1; then
        log_pass "${test_name} (flags identical)"
        echo "PASS" > "${test_dir}/result.txt"
        return 0
    else
        diff -u "${test_dir}/original_flags.txt" "${test_dir}/ours_flags.txt" > "${test_dir}/flags_diff.txt" 2>&1 || true
        log_fail "${test_name} (flags differ)"
        echo "FAIL" > "${test_dir}/result.txt"
        return 1
    fi
}

# ============================================
# Helper: Compare RPC commands list
# ============================================
test_rpc_commands() {
    local test_name="$1"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Extract RPC command names (Available Commands section)
    $ORIGINAL_VESCTL request rpc --help 2>&1 | \
        sed -n '/^Available Commands:/,/^$/p' | \
        grep -E '^\s+[a-z]' | \
        awk '{print $1}' | sort > "${test_dir}/original_rpcs.txt"

    $OUR_VESCTL request rpc --help 2>&1 | \
        sed -n '/^Available Commands:/,/^$/p' | \
        grep -E '^\s+[a-z]' | \
        awk '{print $1}' | sort > "${test_dir}/ours_rpcs.txt"

    # Count common, original-only, and ours-only
    local common=$(comm -12 "${test_dir}/original_rpcs.txt" "${test_dir}/ours_rpcs.txt" | wc -l)
    local orig_only=$(comm -23 "${test_dir}/original_rpcs.txt" "${test_dir}/ours_rpcs.txt" | wc -l)
    local ours_only=$(comm -13 "${test_dir}/original_rpcs.txt" "${test_dir}/ours_rpcs.txt" | wc -l)

    echo "common=$common" > "${test_dir}/stats.txt"
    echo "original_only=$orig_only" >> "${test_dir}/stats.txt"
    echo "ours_only=$ours_only" >> "${test_dir}/stats.txt"

    # Save differences
    comm -23 "${test_dir}/original_rpcs.txt" "${test_dir}/ours_rpcs.txt" > "${test_dir}/deprecated_rpcs.txt"
    comm -13 "${test_dir}/original_rpcs.txt" "${test_dir}/ours_rpcs.txt" > "${test_dir}/new_rpcs.txt"

    if [[ $orig_only -eq 0 ]] && [[ $ours_only -eq 0 ]]; then
        log_pass "${test_name} (${common} RPCs identical)"
        echo "PASS" > "${test_dir}/result.txt"
    else
        log_warn "${test_name} (${common} common, ${orig_only} deprecated, ${ours_only} new)"
        echo "WARN_RPC_DIFF" > "${test_dir}/result.txt"
    fi
}

# ============================================
# Test request main help
# ============================================
echo "--- Request Main Command ---"
echo ""

test_help_identical "request-main" request || true
echo ""

# ============================================
# Test request rpc
# ============================================
echo "--- Request RPC ---"
echo ""

test_flags_identical "request-rpc-flags" request rpc || true
test_rpc_commands "request-rpc-commands" || true
echo ""

# ============================================
# Test request secrets
# ============================================
echo "--- Request Secrets ---"
echo ""

test_help_identical "request-secrets" request secrets || true

# Test secrets subcommands
SECRETS_CMDS=(
    "encrypt"
    "get-public-key"
    "get-policy-document"
    "secret-info"
    "build-blindfold-bundle"
)

for subcmd in "${SECRETS_CMDS[@]}"; do
    test_help_identical "request-secrets-${subcmd}" request secrets "$subcmd" || true
done
echo ""

# ============================================
# Test request command-sequence
# ============================================
echo "--- Request Command-Sequence ---"
echo ""

test_flags_identical "request-cmdseq-flags" request command-sequence || true
echo ""

# ============================================
# Summary
# ============================================
echo ""
print_summary

echo ""
echo "Phase 6 (Request Commands) test results: ${RESULTS_DIR}"
