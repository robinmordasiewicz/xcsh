#!/bin/bash
# test-no-api.sh - Comprehensive test of all commands that don't require API credentials
# Tests help text, flags, and structure consistency
#
# NOTE: vesctl-0.2.35 help commands can take 20-25 seconds to complete.
# The timeout is set to 60 seconds in lib/common.sh to accommodate this.

set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/common.sh"

echo "=========================================="
echo "Phase 3A: No-API Command Tests"
echo "=========================================="
echo ""

# Initialize results directory
init_results_dir
PHASE_DIR="${RESULTS_DIR}/phase3-no-api"
mkdir -p "$PHASE_DIR"

log_info "Original vesctl: ${ORIGINAL_VESCTL}"
log_info "Our f5xcctl: ${OUR_F5XCCTL}"
log_info "Results directory: ${RESULTS_DIR}"
echo ""

export RESULTS_DIR="$PHASE_DIR"

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
    $OUR_F5XCCTL "${cmd[@]}" --help > "${test_dir}/ours_full.txt" 2>&1

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
# Helper: Test exact help match
# ============================================
test_help_exact() {
    local test_name="$1"
    shift
    local cmd=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    $ORIGINAL_VESCTL "${cmd[@]}" --help > "${test_dir}/original.txt" 2>&1
    $OUR_F5XCCTL "${cmd[@]}" --help > "${test_dir}/ours.txt" 2>&1

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
# Helper: Test flags section only
# ============================================
test_flags_section() {
    local test_name="$1"
    shift
    local cmd=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Extract Flags and Global Flags sections
    $ORIGINAL_VESCTL "${cmd[@]}" --help 2>&1 | sed -n '/^Flags:/,/^$/p' > "${test_dir}/original_flags.txt"
    $OUR_F5XCCTL "${cmd[@]}" --help 2>&1 | sed -n '/^Flags:/,/^$/p' > "${test_dir}/ours_flags.txt"

    $ORIGINAL_VESCTL "${cmd[@]}" --help 2>&1 | sed -n '/^Global Flags:/,/^$/p' > "${test_dir}/original_global.txt"
    $OUR_F5XCCTL "${cmd[@]}" --help 2>&1 | sed -n '/^Global Flags:/,/^$/p' > "${test_dir}/ours_global.txt"

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
# Test 1: Root command
# ============================================
echo "--- Root Command ---"
test_help_structure "root" || true
test_flags_section "root-flags" || true

# ============================================
# Test 2: Version command
# ============================================
echo ""
echo "--- Version Command ---"
test_help_exact "version-help" version || true

# ============================================
# Test 3: Completion command
# ============================================
echo ""
echo "--- Completion Command ---"
test_help_exact "completion-help" completion || true
test_help_exact "completion-bash-help" completion bash || true
test_help_exact "completion-zsh-help" completion zsh || true

# ============================================
# Test 4: Configuration command structure
# ============================================
echo ""
echo "--- Configuration Command ---"
test_help_structure "configuration" configuration || true
test_flags_section "configuration-flags" configuration || true

# ============================================
# Test 5: All configuration subcommands
# ============================================
echo ""
echo "--- Configuration Subcommands ---"
CONFIG_SUBCMDS=("add-labels" "apply" "create" "delete" "get" "list" "patch" "remove-labels" "replace" "status")

for subcmd in "${CONFIG_SUBCMDS[@]}"; do
    # Use the is_buggy_command function from common.sh
    if is_buggy_command "configuration" "$subcmd"; then
        log_skip "config-${subcmd} (vesctl-0.2.35 hangs on this help command)"
        mkdir -p "${RESULTS_DIR}/config-${subcmd}"
        echo "SKIP" > "${RESULTS_DIR}/config-${subcmd}/result.txt"
        mkdir -p "${RESULTS_DIR}/config-${subcmd}-flags"
        echo "SKIP" > "${RESULTS_DIR}/config-${subcmd}-flags/result.txt"
        continue
    fi
    test_help_structure "config-${subcmd}" configuration "$subcmd" || true
    test_flags_section "config-${subcmd}-flags" configuration "$subcmd" || true
done

# ============================================
# Test 6: Request command (if exists in both)
# ============================================
echo ""
echo "--- Request Command ---"
if $ORIGINAL_VESCTL request --help > /dev/null 2>&1 && $OUR_F5XCCTL request --help > /dev/null 2>&1; then
    test_help_structure "request" request || true
    test_flags_section "request-flags" request || true

    # Request subcommands
    REQUEST_SUBCMDS=("http" "rpc")
    for subcmd in "${REQUEST_SUBCMDS[@]}"; do
        if $ORIGINAL_VESCTL request "$subcmd" --help > /dev/null 2>&1; then
            test_help_structure "request-${subcmd}" request "$subcmd" || true
        fi
    done
else
    log_info "Request command - checking availability..."
    $ORIGINAL_VESCTL request --help > /dev/null 2>&1 && log_info "  Original: has request" || log_info "  Original: no request"
    $OUR_F5XCCTL request --help > /dev/null 2>&1 && log_info "  Ours: has request" || log_info "  Ours: no request"
fi

# ============================================
# Test 7: Site command (if exists in both)
# ============================================
echo ""
echo "--- Site Command ---"
if $ORIGINAL_VESCTL site --help > /dev/null 2>&1 && $OUR_F5XCCTL site --help > /dev/null 2>&1; then
    test_help_structure "site" site || true
    test_flags_section "site-flags" site || true
else
    log_info "Site command - checking availability..."
    $ORIGINAL_VESCTL site --help > /dev/null 2>&1 && log_info "  Original: has site" || log_info "  Original: no site"
    $OUR_F5XCCTL site --help > /dev/null 2>&1 && log_info "  Ours: has site" || log_info "  Ours: no site"
fi

# ============================================
# Test 8: API-endpoint command (if exists in both)
# ============================================
echo ""
echo "--- API-Endpoint Command ---"
if $ORIGINAL_VESCTL api-endpoint --help > /dev/null 2>&1 && $OUR_F5XCCTL api-endpoint --help > /dev/null 2>&1; then
    test_help_structure "api-endpoint" api-endpoint || true
    test_flags_section "api-endpoint-flags" api-endpoint || true
else
    log_info "API-endpoint command - checking availability..."
    $ORIGINAL_VESCTL api-endpoint --help > /dev/null 2>&1 && log_info "  Original: has api-endpoint" || log_info "  Original: no api-endpoint"
    $OUR_F5XCCTL api-endpoint --help > /dev/null 2>&1 && log_info "  Ours: has api-endpoint" || log_info "  Ours: no api-endpoint"
fi

# ============================================
# Test 9: Error handling
# ============================================
echo ""
echo "--- Error Handling ---"

test_error() {
    local test_name="$1"
    shift
    local cmd=("$@")
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    $ORIGINAL_VESCTL "${cmd[@]}" > "${test_dir}/original.stdout" 2> "${test_dir}/original.stderr"
    local orig_exit=$?

    $OUR_F5XCCTL "${cmd[@]}" > "${test_dir}/ours.stdout" 2> "${test_dir}/ours.stderr"
    local our_exit=$?

    echo "orig_exit=$orig_exit our_exit=$our_exit" > "${test_dir}/exits.txt"

    if [[ $orig_exit -eq $our_exit ]]; then
        log_pass "${test_name} (exit code: $orig_exit)"
        echo "PASS" > "${test_dir}/result.txt"
    else
        log_fail "${test_name} (orig=$orig_exit ours=$our_exit)"
        echo "FAIL" > "${test_dir}/result.txt"
    fi
}

test_error "error-unknown-cmd" unknowncommand || true
test_error "error-invalid-flag" --invalidflag || true
test_error "error-config-missing-resource" configuration list || true

# ============================================
# Summary
# ============================================
echo ""
print_summary

echo ""
echo "Phase 3A (No-API) test results: ${RESULTS_DIR}"
