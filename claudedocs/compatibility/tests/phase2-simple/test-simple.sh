#!/bin/bash
# test-simple.sh - Test version and completion commands (no credentials needed)
# Focuses on behavioral consistency, not exact content match

# Don't exit on errors - we want to run all tests
set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/common.sh"

echo "=========================================="
echo "Phase 2: Version & Completion Tests"
echo "=========================================="
echo ""

# Initialize results directory
init_results_dir
PHASE_DIR="${RESULTS_DIR}/phase2-simple"
mkdir -p "$PHASE_DIR"

log_info "Original vesctl: ${ORIGINAL_VESCTL}"
log_info "Our vesctl: ${OUR_VESCTL}"
log_info "Results directory: ${RESULTS_DIR}"
echo ""

# Override RESULTS_DIR for test functions
export RESULTS_DIR="$PHASE_DIR"

# ============================================
# Test 1: Version Command Help
# ============================================
echo "--- Test Category: Version Command ---"

test_help "version-help" version --help || true

# ============================================
# Test 2: Version Output Style
# ============================================
echo ""
echo "--- Test Category: Version Output Style ---"

test_version_output() {
    local test_name="version-output"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Capture version output
    $ORIGINAL_VESCTL version > "${test_dir}/original.txt" 2>&1 || true
    $OUR_VESCTL version > "${test_dir}/ours.txt" 2>&1 || true

    # Check if both produce output
    local orig_has_output=$(wc -l < "${test_dir}/original.txt")
    local our_has_output=$(wc -l < "${test_dir}/ours.txt")

    if [[ $orig_has_output -gt 0 ]] && [[ $our_has_output -gt 0 ]]; then
        # Both produce output - check if ours mentions "vesctl" (not f5xc)
        if grep -q "vesctl" "${test_dir}/ours.txt"; then
            log_pass "${test_name} (both produce output, ours shows 'vesctl')"
            echo "PASS" > "${test_dir}/result.txt"
        elif grep -q "f5xc" "${test_dir}/ours.txt"; then
            log_fail "${test_name} (ours shows 'f5xc' instead of 'vesctl')"
            echo "FAIL" > "${test_dir}/result.txt"
        else
            log_pass "${test_name} (both produce output)"
            echo "PASS" > "${test_dir}/result.txt"
        fi
    else
        log_fail "${test_name} (missing output)"
        echo "FAIL" > "${test_dir}/result.txt"
    fi

    # Save comparison info
    echo "Original output:" > "${test_dir}/comparison.txt"
    cat "${test_dir}/original.txt" >> "${test_dir}/comparison.txt"
    echo "" >> "${test_dir}/comparison.txt"
    echo "Our output:" >> "${test_dir}/comparison.txt"
    cat "${test_dir}/ours.txt" >> "${test_dir}/comparison.txt"
}

test_version_output

# ============================================
# Test 3: Completion Command Help
# ============================================
echo ""
echo "--- Test Category: Completion Command ---"

test_help "completion-help" completion --help || true

# ============================================
# Test 4: Bash Completion
# ============================================
echo ""
echo "--- Test Category: Bash Completion ---"

test_completion_bash() {
    local test_name="completion-bash"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Capture completion output
    $ORIGINAL_VESCTL completion bash > "${test_dir}/original.txt" 2>&1 || true
    $OUR_VESCTL completion bash > "${test_dir}/ours.txt" 2>&1 || true

    # Compare structure (first 50 lines to check format)
    head -50 "${test_dir}/original.txt" | grep -v "^#" | grep -v "^$" > "${test_dir}/original_structure.txt"
    head -50 "${test_dir}/ours.txt" | grep -v "^#" | grep -v "^$" > "${test_dir}/ours_structure.txt"

    # Check if both have bash completion patterns
    local orig_has_bash=$(grep -c "_vesctl" "${test_dir}/original.txt" 2>/dev/null || echo "0")
    local our_has_bash=$(grep -c "_vesctl" "${test_dir}/ours.txt" 2>/dev/null || echo "0")

    if [[ $orig_has_bash -gt 0 ]] && [[ $our_has_bash -gt 0 ]]; then
        log_pass "${test_name} (both produce bash completion with _vesctl pattern)"
        echo "PASS" > "${test_dir}/result.txt"
    elif diff -q "${test_dir}/original.txt" "${test_dir}/ours.txt" > /dev/null 2>&1; then
        log_pass "${test_name} (identical)"
        echo "PASS" > "${test_dir}/result.txt"
    else
        # Check if both at least produce bash completion syntax
        local orig_has_func=$(grep -c "function" "${test_dir}/original.txt" 2>/dev/null || grep -c "__" "${test_dir}/original.txt" 2>/dev/null || echo "0")
        local our_has_func=$(grep -c "function" "${test_dir}/ours.txt" 2>/dev/null || grep -c "__" "${test_dir}/ours.txt" 2>/dev/null || echo "0")

        if [[ $orig_has_func -gt 0 ]] && [[ $our_has_func -gt 0 ]]; then
            log_pass "${test_name} (both produce bash completion)"
            echo "PASS" > "${test_dir}/result.txt"
        else
            diff -u "${test_dir}/original.txt" "${test_dir}/ours.txt" > "${test_dir}/diff.txt" 2>&1 || true
            log_fail "${test_name}"
            echo "FAIL" > "${test_dir}/result.txt"
        fi
    fi
}

test_completion_bash

# ============================================
# Test 5: Zsh Completion
# ============================================
echo ""
echo "--- Test Category: Zsh Completion ---"

test_completion_zsh() {
    local test_name="completion-zsh"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Capture completion output
    $ORIGINAL_VESCTL completion zsh > "${test_dir}/original.txt" 2>&1 || true
    $OUR_VESCTL completion zsh > "${test_dir}/ours.txt" 2>&1 || true

    # Check if both have zsh completion patterns
    local orig_has_zsh=$(grep -c "compdef" "${test_dir}/original.txt" 2>/dev/null || echo "0")
    local our_has_zsh=$(grep -c "compdef" "${test_dir}/ours.txt" 2>/dev/null || echo "0")

    if [[ $orig_has_zsh -gt 0 ]] && [[ $our_has_zsh -gt 0 ]]; then
        log_pass "${test_name} (both produce zsh completion with compdef)"
        echo "PASS" > "${test_dir}/result.txt"
    elif diff -q "${test_dir}/original.txt" "${test_dir}/ours.txt" > /dev/null 2>&1; then
        log_pass "${test_name} (identical)"
        echo "PASS" > "${test_dir}/result.txt"
    else
        # Check if both at least produce zsh completion syntax
        local orig_has_pattern=$(grep -c "_vesctl" "${test_dir}/original.txt" 2>/dev/null || echo "0")
        local our_has_pattern=$(grep -c "_vesctl" "${test_dir}/ours.txt" 2>/dev/null || echo "0")

        if [[ $orig_has_pattern -gt 0 ]] && [[ $our_has_pattern -gt 0 ]]; then
            log_pass "${test_name} (both produce zsh completion)"
            echo "PASS" > "${test_dir}/result.txt"
        else
            diff -u "${test_dir}/original.txt" "${test_dir}/ours.txt" > "${test_dir}/diff.txt" 2>&1 || true
            log_fail "${test_name}"
            echo "FAIL" > "${test_dir}/result.txt"
        fi
    fi
}

test_completion_zsh

# ============================================
# Test 6: Configure Command Help
# ============================================
echo ""
echo "--- Test Category: Configure Command ---"

# Note: Configure command may not exist in original vesctl
test_configure_exists() {
    local test_name="configure-exists"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Check if configure command exists in original
    $ORIGINAL_VESCTL configure --help > "${test_dir}/original.txt" 2>&1
    local orig_exit=$?

    $OUR_VESCTL configure --help > "${test_dir}/ours.txt" 2>&1
    local our_exit=$?

    if [[ $our_exit -eq 0 ]]; then
        if [[ $orig_exit -eq 0 ]]; then
            # Both have configure command - compare help
            if diff -q "${test_dir}/original.txt" "${test_dir}/ours.txt" > /dev/null 2>&1; then
                log_pass "${test_name} (identical)"
                echo "PASS" > "${test_dir}/result.txt"
            else
                log_warn "${test_name} (help differs - may be enhanced)"
                echo "WARN" > "${test_dir}/result.txt"
            fi
        else
            log_warn "${test_name} (new command - not in original)"
            echo "WARN" > "${test_dir}/result.txt"
        fi
    else
        log_fail "${test_name} (configure command failed)"
        echo "FAIL" > "${test_dir}/result.txt"
    fi
}

test_configure_exists

# ============================================
# Test 7: Login Command Help
# ============================================
echo ""
echo "--- Test Category: Login Command ---"

test_login_exists() {
    local test_name="login-exists"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Check if login command exists in original
    $ORIGINAL_VESCTL login --help > "${test_dir}/original.txt" 2>&1
    local orig_exit=$?

    $OUR_VESCTL login --help > "${test_dir}/ours.txt" 2>&1
    local our_exit=$?

    if [[ $our_exit -eq 0 ]]; then
        if [[ $orig_exit -eq 0 ]]; then
            # Both have login command - compare help
            if diff -q "${test_dir}/original.txt" "${test_dir}/ours.txt" > /dev/null 2>&1; then
                log_pass "${test_name} (identical)"
                echo "PASS" > "${test_dir}/result.txt"
            else
                log_warn "${test_name} (help differs - may be enhanced)"
                echo "WARN" > "${test_dir}/result.txt"
            fi
        else
            log_warn "${test_name} (new command - not in original)"
            echo "WARN" > "${test_dir}/result.txt"
        fi
    else
        log_fail "${test_name} (login command failed)"
        echo "FAIL" > "${test_dir}/result.txt"
    fi
}

test_login_exists

# ============================================
# Test 8: Global Flags in Root Help
# ============================================
echo ""
echo "--- Test Category: Global Flags Consistency ---"

test_global_flags_root() {
    local test_name="global-flags-root"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Extract Global Flags section from root help
    $ORIGINAL_VESCTL --help 2>&1 | sed -n '/^Flags:/,/^$/p' > "${test_dir}/original.txt"
    $OUR_VESCTL --help 2>&1 | sed -n '/^Flags:/,/^$/p' > "${test_dir}/ours.txt"

    if diff -q "${test_dir}/original.txt" "${test_dir}/ours.txt" > /dev/null 2>&1; then
        log_pass "${test_name}"
        echo "PASS" > "${test_dir}/result.txt"
    else
        diff -u "${test_dir}/original.txt" "${test_dir}/ours.txt" > "${test_dir}/diff.txt" 2>&1 || true
        log_fail "${test_name}"
        echo "FAIL" > "${test_dir}/result.txt"
    fi
}

test_global_flags_root

# ============================================
# Test 9: Output Format Flag
# ============================================
echo ""
echo "--- Test Category: Output Format Flag ---"

test_outfmt_flag() {
    local test_name="outfmt-flag"
    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Check if --outfmt flag exists and accepts json/yaml/table
    local formats=("json" "yaml" "table")
    local all_pass=true

    for fmt in "${formats[@]}"; do
        # Test with version command (doesn't need API)
        $ORIGINAL_VESCTL version --outfmt "$fmt" > "${test_dir}/original_${fmt}.txt" 2>&1 || true
        $OUR_VESCTL version --outfmt "$fmt" > "${test_dir}/ours_${fmt}.txt" 2>&1 || true

        # Both should either succeed or fail similarly
        local orig_lines=$(wc -l < "${test_dir}/original_${fmt}.txt")
        local our_lines=$(wc -l < "${test_dir}/ours_${fmt}.txt")

        if [[ $orig_lines -eq 0 ]] && [[ $our_lines -eq 0 ]]; then
            # Both failed - that's ok for some formats
            continue
        fi

        if [[ $orig_lines -gt 0 ]] && [[ $our_lines -gt 0 ]]; then
            # Both succeeded
            continue
        fi

        # One succeeded, one failed
        all_pass=false
    done

    if $all_pass; then
        log_pass "${test_name} (outfmt flag behaves consistently)"
        echo "PASS" > "${test_dir}/result.txt"
    else
        log_warn "${test_name} (outfmt behavior differs)"
        echo "WARN" > "${test_dir}/result.txt"
    fi
}

test_outfmt_flag

# ============================================
# Summary
# ============================================
echo ""
print_summary

# Generate reports (function not yet implemented)
export RESULTS_DIR=$(dirname "$PHASE_DIR")
# generate_markdown_report "$RESULTS_DIR"

echo ""
echo "Phase 2 test results: ${RESULTS_DIR}/summary.md"
