#!/bin/bash
# compare.sh - Advanced comparison functions for vesctl compatibility testing

# Note: common.sh must be sourced before this file

# Compare help text with smart normalization
# Handles differences in binary name, version strings, etc.
compare_help_text() {
    local orig_file="$1"
    local our_file="$2"
    local diff_file="$3"

    # Create temp files for normalized content
    local orig_normalized=$(mktemp)
    local our_normalized=$(mktemp)

    # Normalize help text:
    # 1. Replace binary names (vesctl.darwin-arm64 -> vesctl, f5xc -> vesctl)
    # 2. Normalize whitespace
    # 3. Remove version strings
    sed -E \
        -e 's/vesctl\.darwin-arm64/vesctl/g' \
        -e 's/f5xc/vesctl/g' \
        -e 's/[[:space:]]+$//g' \
        -e 's/v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?/VERSION/g' \
        "$orig_file" > "$orig_normalized"

    sed -E \
        -e 's/vesctl\.darwin-arm64/vesctl/g' \
        -e 's/f5xc/vesctl/g' \
        -e 's/[[:space:]]+$//g' \
        -e 's/v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?/VERSION/g' \
        "$our_file" > "$our_normalized"

    # Compare
    local result=0
    if ! diff -u "$orig_normalized" "$our_normalized" > "$diff_file" 2>&1; then
        result=1
    fi

    # Cleanup
    rm -f "$orig_normalized" "$our_normalized"

    return $result
}

# Compare JSON output (structure only, ignore values)
compare_json_structure() {
    local orig_file="$1"
    local our_file="$2"
    local diff_file="$3"

    # Extract JSON keys/structure using jq
    if command -v jq &> /dev/null; then
        local orig_keys=$(mktemp)
        local our_keys=$(mktemp)

        jq -r 'paths | map(tostring) | join(".")' "$orig_file" 2>/dev/null | sort > "$orig_keys"
        jq -r 'paths | map(tostring) | join(".")' "$our_file" 2>/dev/null | sort > "$our_keys"

        local result=0
        if ! diff -u "$orig_keys" "$our_keys" > "$diff_file" 2>&1; then
            result=1
        fi

        rm -f "$orig_keys" "$our_keys"
        return $result
    else
        # Fallback: direct comparison with normalization
        compare_files "$orig_file" "$our_file" "$diff_file"
        return $?
    fi
}

# Compare YAML output (structure only)
compare_yaml_structure() {
    local orig_file="$1"
    local our_file="$2"
    local diff_file="$3"

    # If yq is available, use it for structure comparison
    if command -v yq &> /dev/null; then
        local orig_keys=$(mktemp)
        local our_keys=$(mktemp)

        yq -r '.. | path | map(tostring) | join(".")' "$orig_file" 2>/dev/null | sort -u > "$orig_keys"
        yq -r '.. | path | map(tostring) | join(".")' "$our_file" 2>/dev/null | sort -u > "$our_keys"

        local result=0
        if ! diff -u "$orig_keys" "$our_keys" > "$diff_file" 2>&1; then
            result=1
        fi

        rm -f "$orig_keys" "$our_keys"
        return $result
    else
        # Fallback: direct comparison with normalization
        normalize_output "$orig_file" "${orig_file}.normalized"
        normalize_output "$our_file" "${our_file}.normalized"
        compare_files "${orig_file}.normalized" "${our_file}.normalized" "$diff_file"
        return $?
    fi
}

# Compare error messages (more lenient)
# Allows for different wording but checks for same error type
compare_error_messages() {
    local orig_file="$1"
    local our_file="$2"

    # Extract error patterns
    local orig_has_error=$(grep -qi "error\|fail\|invalid\|missing\|required" "$orig_file" && echo "yes" || echo "no")
    local our_has_error=$(grep -qi "error\|fail\|invalid\|missing\|required" "$our_file" && echo "yes" || echo "no")

    if [[ "$orig_has_error" == "$our_has_error" ]]; then
        return 0
    else
        return 1
    fi
}

# Get line-by-line diff with context
detailed_diff() {
    local file1="$1"
    local file2="$2"
    local output="$3"

    diff -y --suppress-common-lines -W 200 "$file1" "$file2" > "$output" 2>&1 || true
}

# Count differences between files
count_differences() {
    local file1="$1"
    local file2="$2"

    diff "$file1" "$file2" 2>/dev/null | grep -c "^[<>]" || echo "0"
}

# Check if outputs are semantically equivalent
# (same information, possibly different formatting)
are_semantically_equivalent() {
    local file1="$1"
    local file2="$2"

    # Normalize both files aggressively
    local norm1=$(mktemp)
    local norm2=$(mktemp)

    # Remove all whitespace variations, normalize case for non-code content
    tr -s '[:space:]' ' ' < "$file1" | sed 's/^ //; s/ $//' > "$norm1"
    tr -s '[:space:]' ' ' < "$file2" | sed 's/^ //; s/ $//' > "$norm2"

    local result=0
    if ! diff -q "$norm1" "$norm2" > /dev/null 2>&1; then
        result=1
    fi

    rm -f "$norm1" "$norm2"
    return $result
}
