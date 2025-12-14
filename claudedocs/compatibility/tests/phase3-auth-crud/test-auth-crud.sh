#!/bin/bash
# test-auth-crud.sh - Test authentication commands and configuration CRUD operations
# Part A: Auth commands (no credentials needed)
# Part B: CRUD operations (requires API credentials)

set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/common.sh"

echo "=========================================="
echo "Phase 3: Authentication & CRUD Tests"
echo "=========================================="
echo ""

# Initialize results directory
init_results_dir
PHASE_DIR="${RESULTS_DIR}/phase3-auth-crud"
mkdir -p "$PHASE_DIR"

log_info "Original vesctl: ${ORIGINAL_VESCTL}"
log_info "Our f5xcctl: ${OUR_F5XCCTL}"
log_info "Results directory: ${RESULTS_DIR}"
echo ""

# Override RESULTS_DIR for test functions
export RESULTS_DIR="$PHASE_DIR"

# ============================================
# Part A: Authentication Commands
# ============================================
echo "=========================================="
echo "Part A: Authentication Commands"
echo "=========================================="

# ============================================
# Test 1: Check which auth commands exist
# ============================================
echo ""
echo "--- Command Availability ---"

check_command_exists() {
    local cmd="$1"
    local test_name="exists-${cmd}"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Check original
    $ORIGINAL_VESCTL $cmd --help > "${test_dir}/original.txt" 2>&1
    local orig_exit=$?

    # Check ours
    $OUR_F5XCCTL $cmd --help > "${test_dir}/ours.txt" 2>&1
    local our_exit=$?

    echo "orig_exit=$orig_exit" > "${test_dir}/exits.txt"
    echo "our_exit=$our_exit" >> "${test_dir}/exits.txt"

    if [[ $orig_exit -eq 0 ]] && [[ $our_exit -eq 0 ]]; then
        log_pass "${test_name} (both have command)"
        echo "BOTH" > "${test_dir}/result.txt"
        return 0
    elif [[ $orig_exit -ne 0 ]] && [[ $our_exit -eq 0 ]]; then
        log_warn "${test_name} (ours only - enhancement)"
        echo "OURS_ONLY" > "${test_dir}/result.txt"
        return 1
    elif [[ $orig_exit -eq 0 ]] && [[ $our_exit -ne 0 ]]; then
        log_fail "${test_name} (missing from ours)"
        echo "ORIG_ONLY" > "${test_dir}/result.txt"
        return 2
    else
        log_warn "${test_name} (neither has command)"
        echo "NEITHER" > "${test_dir}/result.txt"
        return 3
    fi
}

# Check authentication-related commands
AUTH_COMMANDS=("configure" "login" "logout" "whoami")
BOTH_HAVE_COMMANDS=()

for cmd in "${AUTH_COMMANDS[@]}"; do
    if check_command_exists "$cmd"; then
        BOTH_HAVE_COMMANDS+=("$cmd")
    fi
done

# ============================================
# Test 2: Compare help for shared auth commands
# ============================================
echo ""
echo "--- Auth Command Help Comparison ---"

for cmd in "${BOTH_HAVE_COMMANDS[@]}"; do
    test_help "auth-help-${cmd}" "$cmd" --help || true
done

# ============================================
# Test 3: Configure subcommands
# ============================================
echo ""
echo "--- Configure Subcommands ---"

# Check if configure command exists in original
$ORIGINAL_VESCTL configure --help > /dev/null 2>&1
ORIG_HAS_CONFIGURE=$?

if [[ $ORIG_HAS_CONFIGURE -eq 0 ]]; then
    # Test configure subcommands
    CONFIGURE_SUBCMDS=("show" "set")
    for subcmd in "${CONFIGURE_SUBCMDS[@]}"; do
        check_command_exists "configure $subcmd" || true
    done
else
    log_info "Original vesctl does not have 'configure' command - skipping subcommand tests"
fi

# ============================================
# Test 4: Error handling without credentials
# ============================================
echo ""
echo "--- Error Handling (No Credentials) ---"

test_no_creds_error() {
    local cmd="$1"
    local test_name="nocreds-${cmd// /-}"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Run without credentials (should fail gracefully)
    # Temporarily unset credential env vars
    (
        unset F5XC_P12_PASSWORD F5XC_P12_FILE
        $ORIGINAL_VESCTL $cmd > "${test_dir}/original.stdout" 2> "${test_dir}/original.stderr"
        echo $? > "${test_dir}/original.exit"
    ) || true

    (
        unset F5XC_P12_PASSWORD F5XC_P12_FILE
        $OUR_F5XCCTL $cmd > "${test_dir}/ours.stdout" 2> "${test_dir}/ours.stderr"
        echo $? > "${test_dir}/ours.exit"
    ) || true

    local orig_exit=$(cat "${test_dir}/original.exit" 2>/dev/null || echo "999")
    local our_exit=$(cat "${test_dir}/ours.exit" 2>/dev/null || echo "999")

    # Both should fail (non-zero exit) when no credentials
    if [[ $orig_exit -ne 0 ]] && [[ $our_exit -ne 0 ]]; then
        log_pass "${test_name} (both fail without credentials)"
        echo "PASS" > "${test_dir}/result.txt"
    elif [[ $orig_exit -eq 0 ]] && [[ $our_exit -eq 0 ]]; then
        log_pass "${test_name} (both succeed - command doesn't need creds)"
        echo "PASS" > "${test_dir}/result.txt"
    else
        log_warn "${test_name} (different behavior)"
        echo "WARN" > "${test_dir}/result.txt"
    fi
}

# Test commands that might need credentials
test_no_creds_error "configuration list namespace" || true
test_no_creds_error "whoami" || true

# ============================================
# Part B: Configuration CRUD (Requires API)
# ============================================
echo ""
echo "=========================================="
echo "Part B: Configuration CRUD Operations"
echo "=========================================="

if has_api_credentials; then
    log_info "API credentials detected - running CRUD tests"
    echo ""

    # ============================================
    # Test 5: List operations
    # ============================================
    echo "--- List Operations ---"

    test_list_operation() {
        local resource="$1"
        local namespace="${2:-system}"
        local test_name="list-${resource}"
        local test_dir="${RESULTS_DIR}/${test_name}"
        mkdir -p "$test_dir"

        # Run list command
        $ORIGINAL_VESCTL configuration list "$resource" -n "$namespace" --outfmt json > "${test_dir}/original.json" 2> "${test_dir}/original.stderr"
        local orig_exit=$?

        $OUR_F5XCCTL configuration list "$resource" -n "$namespace" --outfmt json > "${test_dir}/ours.json" 2> "${test_dir}/ours.stderr"
        local our_exit=$?

        echo "orig_exit=$orig_exit" > "${test_dir}/exits.txt"
        echo "our_exit=$our_exit" >> "${test_dir}/exits.txt"

        # Compare exit codes
        if [[ $orig_exit -ne $our_exit ]]; then
            log_fail "${test_name} (exit code: orig=$orig_exit ours=$our_exit)"
            echo "FAIL" > "${test_dir}/result.txt"
            return 1
        fi

        if [[ $orig_exit -ne 0 ]]; then
            log_warn "${test_name} (both failed with exit $orig_exit)"
            echo "WARN" > "${test_dir}/result.txt"
            return 0
        fi

        # Compare JSON structure (not content)
        # Extract top-level keys
        jq -r 'keys[]' "${test_dir}/original.json" 2>/dev/null | sort > "${test_dir}/original_keys.txt" || true
        jq -r 'keys[]' "${test_dir}/ours.json" 2>/dev/null | sort > "${test_dir}/ours_keys.txt" || true

        if diff -q "${test_dir}/original_keys.txt" "${test_dir}/ours_keys.txt" > /dev/null 2>&1; then
            log_pass "${test_name} (JSON structure matches)"
            echo "PASS" > "${test_dir}/result.txt"
        else
            log_warn "${test_name} (JSON structure differs)"
            diff -u "${test_dir}/original_keys.txt" "${test_dir}/ours_keys.txt" > "${test_dir}/keys_diff.txt" 2>&1 || true
            echo "WARN" > "${test_dir}/result.txt"
        fi
    }

    # Test list for common resources
    CRUD_RESOURCES=("namespace" "http_loadbalancer" "origin_pool")
    for resource in "${CRUD_RESOURCES[@]}"; do
        test_list_operation "$resource" || true
    done

    # ============================================
    # Test 6: Get operations
    # ============================================
    echo ""
    echo "--- Get Operations ---"

    test_get_operation() {
        local resource="$1"
        local name="$2"
        local namespace="${3:-system}"
        local test_name="get-${resource}-${name}"
        local test_dir="${RESULTS_DIR}/${test_name}"
        mkdir -p "$test_dir"

        # Run get command
        $ORIGINAL_VESCTL configuration get "$resource" "$name" -n "$namespace" --outfmt json > "${test_dir}/original.json" 2> "${test_dir}/original.stderr"
        local orig_exit=$?

        $OUR_F5XCCTL configuration get "$resource" "$name" -n "$namespace" --outfmt json > "${test_dir}/ours.json" 2> "${test_dir}/ours.stderr"
        local our_exit=$?

        echo "orig_exit=$orig_exit" > "${test_dir}/exits.txt"
        echo "our_exit=$our_exit" >> "${test_dir}/exits.txt"

        # Compare exit codes
        if [[ $orig_exit -ne $our_exit ]]; then
            log_fail "${test_name} (exit code: orig=$orig_exit ours=$our_exit)"
            echo "FAIL" > "${test_dir}/result.txt"
            return 1
        fi

        if [[ $orig_exit -ne 0 ]]; then
            log_warn "${test_name} (both failed - resource may not exist)"
            echo "WARN" > "${test_dir}/result.txt"
            return 0
        fi

        # Compare JSON structure
        jq -r 'keys[]' "${test_dir}/original.json" 2>/dev/null | sort > "${test_dir}/original_keys.txt" || true
        jq -r 'keys[]' "${test_dir}/ours.json" 2>/dev/null | sort > "${test_dir}/ours_keys.txt" || true

        if diff -q "${test_dir}/original_keys.txt" "${test_dir}/ours_keys.txt" > /dev/null 2>&1; then
            log_pass "${test_name} (JSON structure matches)"
            echo "PASS" > "${test_dir}/result.txt"
        else
            log_warn "${test_name} (JSON structure differs)"
            echo "WARN" > "${test_dir}/result.txt"
        fi
    }

    # Test get for system namespace (should always exist)
    test_get_operation "namespace" "system" "system" || true

    # ============================================
    # Test 7: Output format consistency
    # ============================================
    echo ""
    echo "--- Output Format Tests ---"

    test_output_format() {
        local format="$1"
        local test_name="outfmt-${format}"
        local test_dir="${RESULTS_DIR}/${test_name}"
        mkdir -p "$test_dir"

        # Use namespace list as test command
        $ORIGINAL_VESCTL configuration list namespace -n system --outfmt "$format" > "${test_dir}/original.txt" 2>&1
        local orig_exit=$?

        $OUR_F5XCCTL configuration list namespace -n system --outfmt "$format" > "${test_dir}/ours.txt" 2>&1
        local our_exit=$?

        if [[ $orig_exit -eq $our_exit ]]; then
            if [[ $orig_exit -eq 0 ]]; then
                # Both succeeded - check output is non-empty
                if [[ -s "${test_dir}/original.txt" ]] && [[ -s "${test_dir}/ours.txt" ]]; then
                    log_pass "${test_name} (both produce output)"
                    echo "PASS" > "${test_dir}/result.txt"
                else
                    log_warn "${test_name} (empty output)"
                    echo "WARN" > "${test_dir}/result.txt"
                fi
            else
                log_warn "${test_name} (both failed)"
                echo "WARN" > "${test_dir}/result.txt"
            fi
        else
            log_fail "${test_name} (exit code mismatch)"
            echo "FAIL" > "${test_dir}/result.txt"
        fi
    }

    for fmt in "json" "yaml" "table"; do
        test_output_format "$fmt" || true
    done

else
    log_warn "No API credentials - skipping CRUD tests"
    echo ""
    echo "To run CRUD tests, set:"
    echo "  export F5XC_P12_PASSWORD='your-password'"
    echo "  export F5XC_P12_FILE='/path/to/credentials.p12'"
    echo ""

    # Mark as skipped
    echo "SKIPPED" > "${RESULTS_DIR}/crud-tests-skipped.txt"
fi

# ============================================
# Summary
# ============================================
echo ""
print_summary

echo ""
echo "Phase 3 test results: ${RESULTS_DIR}"
