#!/bin/bash
# common.sh - Shared utilities for vesctl compatibility testing

# Binary paths (relative to project root)
# Original F5 proprietary binary - download with: ./claudedocs/compatibility/scripts/download-original.sh
# On Linux amd64: use vesctl-0.2.35 (last working version)
# On macOS: use vesctl-0.2.47-original (only for offline tests - has TLS bug)
ORIGINAL_VESCTL="${ORIGINAL_VESCTL:-./vesctl-0.2.35}"
OUR_VESCTL="${OUR_VESCTL:-./vesctl}"

# Timeout for original vesctl commands (in seconds)
# NOTE: vesctl-0.2.35 help commands can take 20-25 seconds to complete
ORIGINAL_VESCTL_TIMEOUT=${ORIGINAL_VESCTL_TIMEOUT:-60}

# Known buggy commands array (currently empty - commands are slow but work)
KNOWN_BUGGY_COMMANDS=()

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_WARNED=0
TESTS_SKIPPED=0

# Results directory
RESULTS_DIR=""

# Initialize results directory with timestamp
init_results_dir() {
    local timestamp=$(date +"%Y-%m-%d_%H-%M-%S")
    RESULTS_DIR="${SCRIPT_DIR}/../results/${timestamp}"
    mkdir -p "$RESULTS_DIR"

    # Create/update latest symlink
    local latest_link="${SCRIPT_DIR}/../results/latest"
    rm -f "$latest_link"
    ln -s "$timestamp" "$latest_link"

    echo "$RESULTS_DIR"
}

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_pass() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((++TESTS_PASSED))
}

log_fail() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((++TESTS_FAILED))
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
    ((++TESTS_WARNED))
}

log_skip() {
    echo -e "${YELLOW}[SKIP]${NC} $1"
    ((++TESTS_SKIPPED))
}

# Check if original vesctl exists
check_original_vesctl() {
    if [[ ! -x "$ORIGINAL_VESCTL" ]]; then
        echo -e "${RED}ERROR: Original vesctl not found at $ORIGINAL_VESCTL${NC}"
        return 1
    fi
    return 0
}

# Check if a command is known to be buggy in vesctl-0.2.35
# Usage: is_buggy_command "configuration" "list"
# Returns: 0 if buggy, 1 if not buggy
is_buggy_command() {
    local cmd_prefix="$*"
    for buggy in "${KNOWN_BUGGY_COMMANDS[@]}"; do
        # Check if the command starts with the buggy prefix
        if [[ "$cmd_prefix" == "$buggy"* ]]; then
            return 0
        fi
    done
    return 1
}

# Run original vesctl with timeout protection
# Usage: run_original_vesctl_safe <args...>
# Returns: exit code (124 if timeout)
run_original_vesctl_safe() {
    timeout "${ORIGINAL_VESCTL_TIMEOUT}s" "$ORIGINAL_VESCTL" "$@"
    return $?
}

# Check if our vesctl exists
check_our_vesctl() {
    if [[ ! -x "$OUR_VESCTL" ]]; then
        echo -e "${RED}ERROR: Our vesctl not found at $OUR_VESCTL${NC}"
        return 1
    fi
    return 0
}

# Check if API credentials are available
# Supports both P12 and key/cert authentication methods
has_api_credentials() {
    # Check for P12 credentials
    if [[ -n "${VES_P12_PASSWORD:-}" && -n "${VES_P12_FILE:-}" && -f "${VES_P12_FILE:-/nonexistent}" ]]; then
        return 0
    fi

    # Check for key/cert credentials from .vesconfig
    local vesconfig="${HOME}/.vesconfig"
    if [[ -f "$vesconfig" ]]; then
        local key_file=$(grep -E "^key:" "$vesconfig" 2>/dev/null | awk '{print $2}')
        local cert_file=$(grep -E "^cert:" "$vesconfig" 2>/dev/null | awk '{print $2}')
        if [[ -n "$key_file" && -f "$key_file" && -n "$cert_file" && -f "$cert_file" ]]; then
            return 0
        fi
    fi

    return 1
}

# Run command and capture stdout, stderr, and exit code
# Usage: run_and_capture <binary> <output_prefix> <args...>
run_and_capture() {
    local binary="$1"
    local output_prefix="$2"
    shift 2
    local args=("$@")

    local stdout_file="${output_prefix}.stdout"
    local stderr_file="${output_prefix}.stderr"
    local exit_file="${output_prefix}.exit"

    # Run the command
    set +e
    "$binary" "${args[@]}" > "$stdout_file" 2> "$stderr_file"
    local exit_code=$?
    set -e

    echo "$exit_code" > "$exit_file"

    echo "$exit_code"
}

# Normalize output for comparison (remove dynamic values)
normalize_output() {
    local input_file="$1"
    local output_file="$2"

    sed -E \
        -e 's/[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}[^ ]*/TIMESTAMP/g' \
        -e 's/[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}/TIMESTAMP/g' \
        -e 's/[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}/UUID/g' \
        -e 's/[a-f0-9]{40}/COMMIT_HASH/g' \
        -e 's/[a-f0-9]{7,8}([^a-f0-9])/COMMIT\1/g' \
        -e 's/v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?/VERSION/g' \
        -e 's/vesctl\.darwin-arm64/vesctl/g' \
        -e 's/f5xc/vesctl/g' \
        "$input_file" > "$output_file"
}

# Compare two files and return result
# Returns: 0 if identical, 1 if different
compare_files() {
    local file1="$1"
    local file2="$2"
    local diff_file="$3"

    if diff -u "$file1" "$file2" > "$diff_file" 2>&1; then
        return 0
    else
        return 1
    fi
}

# Run a help comparison test
# Usage: test_help <test_name> <args...>
test_help() {
    local test_name="$1"
    shift
    local args=("$@")

    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Run original vesctl
    local orig_exit=$(run_and_capture "$ORIGINAL_VESCTL" "${test_dir}/original" "${args[@]}")

    # Run our vesctl
    local our_exit=$(run_and_capture "$OUR_VESCTL" "${test_dir}/ours" "${args[@]}")

    # Normalize outputs
    normalize_output "${test_dir}/original.stdout" "${test_dir}/original.stdout.normalized"
    normalize_output "${test_dir}/ours.stdout" "${test_dir}/ours.stdout.normalized"

    # Compare exit codes
    if [[ "$orig_exit" != "$our_exit" ]]; then
        log_fail "${test_name}: Exit code mismatch (original: $orig_exit, ours: $our_exit)"
        echo "EXIT_MISMATCH: original=$orig_exit ours=$our_exit" > "${test_dir}/result.txt"
        return 0  # Don't fail the script, just log the failure
    fi

    # Compare normalized stdout
    if ! compare_files "${test_dir}/original.stdout.normalized" "${test_dir}/ours.stdout.normalized" "${test_dir}/stdout.diff"; then
        log_fail "${test_name}: stdout differs (see ${test_dir}/stdout.diff)"
        echo "STDOUT_DIFFERS" > "${test_dir}/result.txt"
        return 0  # Don't fail the script, just log the failure
    fi

    # Compare stderr (normalized)
    normalize_output "${test_dir}/original.stderr" "${test_dir}/original.stderr.normalized"
    normalize_output "${test_dir}/ours.stderr" "${test_dir}/ours.stderr.normalized"

    if ! compare_files "${test_dir}/original.stderr.normalized" "${test_dir}/ours.stderr.normalized" "${test_dir}/stderr.diff"; then
        log_warn "${test_name}: stderr differs (see ${test_dir}/stderr.diff)"
        echo "STDERR_DIFFERS" > "${test_dir}/result.txt"
        # Continue - stderr differences are warnings, not failures
    fi

    log_pass "${test_name}"
    echo "PASS" > "${test_dir}/result.txt"
    return 0
}

# Print test summary
print_summary() {
    echo ""
    echo "=========================================="
    echo "Test Summary"
    echo "=========================================="
    echo -e "Passed:  ${GREEN}${TESTS_PASSED}${NC}"
    echo -e "Failed:  ${RED}${TESTS_FAILED}${NC}"
    echo -e "Warned:  ${YELLOW}${TESTS_WARNED}${NC}"
    echo -e "Skipped: ${YELLOW}${TESTS_SKIPPED}${NC}"
    echo "=========================================="

    local total=$((TESTS_PASSED + TESTS_FAILED))
    if [[ $total -gt 0 ]]; then
        local rate=$((TESTS_PASSED * 100 / total))
        echo "Pass Rate: ${rate}%"
    fi

    if [[ $TESTS_FAILED -gt 0 ]]; then
        echo ""
        echo -e "${RED}Some tests failed. Check results in: ${RESULTS_DIR}${NC}"
        return 1
    fi

    return 0
}

# Export counters for subshells
export_counters() {
    export TESTS_PASSED TESTS_FAILED TESTS_WARNED TESTS_SKIPPED
}

# ============================================
# API Response Comparison Functions
# ============================================

# Normalize API response for comparison
# Replaces dynamic fields (timestamps, UUIDs, generations) with placeholders
# Usage: normalize_api_response <input_file> <output_file>
normalize_api_response() {
    local input_file="$1"
    local output_file="$2"

    # Check if file exists and is valid JSON
    if [[ ! -f "$input_file" ]]; then
        cp /dev/null "$output_file"
        return 1
    fi

    # Use jq to normalize dynamic fields
    jq '
        # Normalize ISO timestamps
        walk(if type == "string" and test("^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}") then "TIMESTAMP" else . end) |
        # Normalize UUIDs
        walk(if type == "string" and test("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$") then "UUID" else . end) |
        # Normalize large generation numbers (epoch-like)
        walk(if type == "number" and . > 1000000000 then "GENERATION" else . end) |
        # Sort object keys for consistent comparison
        walk(if type == "object" then to_entries | sort_by(.key) | from_entries else . end) |
        # Sort arrays of objects by name if they have metadata.name
        walk(if type == "array" and length > 0 and (.[0] | type) == "object" and (.[0].metadata?.name // null) != null then sort_by(.metadata.name) else . end)
    ' "$input_file" > "$output_file" 2>/dev/null || cp "$input_file" "$output_file"
}

# Compare API responses with detailed diff
# Usage: compare_api_responses <orig_file> <ours_file> <diff_file> <test_dir>
# Returns: 0 if match, 1 if differ
compare_api_responses() {
    local orig_file="$1"
    local ours_file="$2"
    local diff_file="$3"
    local test_dir="$4"

    # Normalize both files
    normalize_api_response "$orig_file" "${test_dir}/original.normalized.json"
    normalize_api_response "$ours_file" "${test_dir}/ours.normalized.json"

    # Compare normalized versions
    if diff -u "${test_dir}/original.normalized.json" "${test_dir}/ours.normalized.json" > "$diff_file" 2>&1; then
        return 0
    else
        # Generate detailed analysis
        generate_diff_analysis "$diff_file" "${test_dir}/diff_analysis.txt"
        return 1
    fi
}

# Generate human-readable diff analysis
# Usage: generate_diff_analysis <diff_file> <analysis_file>
generate_diff_analysis() {
    local diff_file="$1"
    local analysis_file="$2"

    {
        echo "=== Diff Analysis ==="
        echo ""

        # Count lines added/removed
        local added=$(grep -c "^+" "$diff_file" 2>/dev/null || echo 0)
        local removed=$(grep -c "^-" "$diff_file" 2>/dev/null || echo 0)

        echo "Lines added in ours: $added"
        echo "Lines removed from original: $removed"
        echo ""

        # Extract key differences (JSON field changes)
        echo "=== Key Differences ==="
        grep -E "^[+-].*\"[a-z_]+\":" "$diff_file" 2>/dev/null | head -30 || echo "No field differences found"
        echo ""

        # Check for structural differences
        echo "=== Structural Analysis ==="
        if grep -q "^-.*\[" "$diff_file" 2>/dev/null || grep -q "^+.*\[" "$diff_file" 2>/dev/null; then
            echo "Array structure differences detected"
        fi
        if grep -q "^-.*{" "$diff_file" 2>/dev/null || grep -q "^+.*{" "$diff_file" 2>/dev/null; then
            echo "Object structure differences detected"
        fi

    } > "$analysis_file"
}

# Compare error messages (more lenient)
# Usage: compare_error_messages <orig_stderr> <ours_stderr>
# Returns: 0 if similar enough, 1 if significantly different
compare_error_messages() {
    local orig_file="$1"
    local ours_file="$2"

    # If both empty, they match
    if [[ ! -s "$orig_file" && ! -s "$ours_file" ]]; then
        return 0
    fi

    # If one is empty and the other isn't, they differ
    if [[ ! -s "$orig_file" || ! -s "$ours_file" ]]; then
        return 1
    fi

    # Normalize and compare - extract key error patterns
    local orig_pattern=$(grep -oE "(error|Error|ERROR|not found|denied|failed|invalid)" "$orig_file" 2>/dev/null | sort -u | tr '\n' ' ')
    local ours_pattern=$(grep -oE "(error|Error|ERROR|not found|denied|failed|invalid)" "$ours_file" 2>/dev/null | sort -u | tr '\n' ' ')

    if [[ "$orig_pattern" == "$ours_pattern" ]]; then
        return 0
    fi

    return 1
}

# Run API test with full capture and comparison
# Usage: run_api_test <test_name> <args...>
# Usage: run_api_test --ours-first <test_name> <args...>
# Returns: 0 if pass, 1 if fail
run_api_test() {
    local ours_first=false
    if [[ "${1:-}" == "--ours-first" ]]; then
        ours_first=true
        shift
    fi

    local test_name="$1"
    local test_dir="${RESULTS_DIR}/${test_name}"
    shift
    local args=("$@")

    mkdir -p "$test_dir"

    local orig_exit our_exit
    if [[ "$ours_first" == "true" ]]; then
        # Run ours first (for create tests to avoid race condition)
        our_exit=$(run_and_capture "$OUR_VESCTL" "${test_dir}/ours" "${args[@]}")
        orig_exit=$(run_and_capture "$ORIGINAL_VESCTL" "${test_dir}/original" "${args[@]}")
    else
        # Run original first (default)
        orig_exit=$(run_and_capture "$ORIGINAL_VESCTL" "${test_dir}/original" "${args[@]}")
        our_exit=$(run_and_capture "$OUR_VESCTL" "${test_dir}/ours" "${args[@]}")
    fi

    # Store exit codes
    echo "original_exit=$orig_exit" > "${test_dir}/exits.txt"
    echo "our_exit=$our_exit" >> "${test_dir}/exits.txt"

    # Compare exit codes first
    if [[ "$orig_exit" != "$our_exit" ]]; then
        log_fail "${test_name}: Exit code mismatch (original: $orig_exit, ours: $our_exit)"
        echo "EXIT_MISMATCH" > "${test_dir}/result.txt"
        return 0  # Don't fail the script, just log the failure
    fi

    # If non-zero exit, compare stderr for error consistency
    if [[ "$orig_exit" != "0" ]]; then
        if compare_error_messages "${test_dir}/original.stderr" "${test_dir}/ours.stderr"; then
            log_pass "${test_name} (both errored consistently)"
            echo "PASS_ERROR" > "${test_dir}/result.txt"
            return 0
        else
            log_warn "${test_name}: Error messages differ"
            echo "WARN_ERROR_MSG" > "${test_dir}/result.txt"
            # Store diff for inspection
            diff -u "${test_dir}/original.stderr" "${test_dir}/ours.stderr" > "${test_dir}/stderr.diff" 2>&1 || true
            return 0  # Still pass, just warn
        fi
    fi

    # Compare stdout (API response)
    if compare_api_responses \
        "${test_dir}/original.stdout" \
        "${test_dir}/ours.stdout" \
        "${test_dir}/stdout.diff" \
        "$test_dir"; then
        log_pass "${test_name}"
        echo "PASS" > "${test_dir}/result.txt"
        return 0
    else
        log_fail "${test_name}: Response differs (see ${test_dir}/stdout.diff)"
        echo "FAIL" > "${test_dir}/result.txt"
        return 0  # Don't fail the script, just log the failure
    fi
}

# Test flags section comparison (for --help output)
# Usage: test_flags_section <test_name> <args...>
test_flags_section() {
    local test_name="$1"
    shift
    local args=("$@")

    local test_dir="${RESULTS_DIR}/${test_name}"
    mkdir -p "$test_dir"

    # Capture help output
    local orig_exit=$(run_and_capture "$ORIGINAL_VESCTL" "${test_dir}/original" "${args[@]}" --help)
    local our_exit=$(run_and_capture "$OUR_VESCTL" "${test_dir}/ours" "${args[@]}" --help)

    # Extract just the flags section
    sed -n '/^Flags:/,/^[A-Z]/p' "${test_dir}/original.stdout" | head -n -1 > "${test_dir}/original.flags"
    sed -n '/^Flags:/,/^[A-Z]/p' "${test_dir}/ours.stdout" | head -n -1 > "${test_dir}/ours.flags"

    # Normalize flags (remove extra whitespace, sort)
    sort "${test_dir}/original.flags" > "${test_dir}/original.flags.sorted"
    sort "${test_dir}/ours.flags" > "${test_dir}/ours.flags.sorted"

    if diff -u "${test_dir}/original.flags.sorted" "${test_dir}/ours.flags.sorted" > "${test_dir}/flags.diff" 2>&1; then
        log_pass "${test_name}"
        echo "PASS" > "${test_dir}/result.txt"
        return 0
    else
        log_fail "${test_name}: Flags differ"
        echo "FAIL_FLAGS" > "${test_dir}/result.txt"
        return 0  # Don't fail the script, just log the failure
    fi
}
