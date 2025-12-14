#!/bin/bash
# report.sh - Report generation for vesctl compatibility testing

# Note: common.sh must be sourced before this file

# Generate JSON summary report
generate_json_report() {
    local results_dir="$1"
    local output_file="${results_dir}/summary.json"

    local date=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    local total=$((TESTS_PASSED + TESTS_FAILED))
    local pass_rate=0
    if [[ $total -gt 0 ]]; then
        pass_rate=$((TESTS_PASSED * 100 / total))
    fi

    # Collect failed tests
    local failed_tests=()
    for result_file in "${results_dir}"/*/result.txt; do
        if [[ -f "$result_file" ]]; then
            local result=$(cat "$result_file")
            if [[ "$result" != "PASS" ]]; then
                local test_name=$(basename "$(dirname "$result_file")")
                failed_tests+=("\"${test_name}\"")
            fi
        fi
    done

    # Build JSON
    cat > "$output_file" << EOF
{
  "date": "${date}",
  "summary": {
    "passed": ${TESTS_PASSED},
    "failed": ${TESTS_FAILED},
    "warned": ${TESTS_WARNED},
    "skipped": ${TESTS_SKIPPED},
    "total": ${total},
    "pass_rate": ${pass_rate}
  },
  "failed_tests": [$(IFS=,; echo "${failed_tests[*]}")],
  "original_binary": "${ORIGINAL_VESCTL}",
  "our_binary": "${OUR_F5XCCTL}"
}
EOF

    echo "$output_file"
}

# Generate Markdown summary report
generate_markdown_report() {
    local results_dir="$1"
    local output_file="${results_dir}/summary.md"

    local date=$(date +"%Y-%m-%d %H:%M:%S")
    local total=$((TESTS_PASSED + TESTS_FAILED))
    local pass_rate=0
    if [[ $total -gt 0 ]]; then
        pass_rate=$((TESTS_PASSED * 100 / total))
    fi

    cat > "$output_file" << EOF
# vesctl Compatibility Test Results

**Date**: ${date}
**Pass Rate**: ${pass_rate}% (${TESTS_PASSED}/${total})

## Summary

| Metric | Count |
|--------|-------|
| Passed | ${TESTS_PASSED} |
| Failed | ${TESTS_FAILED} |
| Warned | ${TESTS_WARNED} |
| Skipped | ${TESTS_SKIPPED} |
| **Total** | **${total}** |

## Binaries Tested

| Binary | Path |
|--------|------|
| Original | \`${ORIGINAL_VESCTL}\` |
| Ours | \`${OUR_F5XCCTL}\` |

EOF

    # Add failed tests section
    local has_failures=false
    for result_file in "${results_dir}"/*/result.txt; do
        if [[ -f "$result_file" ]]; then
            local result=$(cat "$result_file")
            if [[ "$result" != "PASS" ]]; then
                if [[ "$has_failures" == "false" ]]; then
                    echo "## Failed Tests" >> "$output_file"
                    echo "" >> "$output_file"
                    has_failures=true
                fi

                local test_name=$(basename "$(dirname "$result_file")")
                local test_dir=$(dirname "$result_file")

                echo "### ${test_name}" >> "$output_file"
                echo "" >> "$output_file"
                echo "**Result**: ${result}" >> "$output_file"
                echo "" >> "$output_file"

                # Add diff if exists
                if [[ -f "${test_dir}/stdout.diff" ]] && [[ -s "${test_dir}/stdout.diff" ]]; then
                    echo "<details>" >> "$output_file"
                    echo "<summary>stdout diff</summary>" >> "$output_file"
                    echo "" >> "$output_file"
                    echo '```diff' >> "$output_file"
                    head -100 "${test_dir}/stdout.diff" >> "$output_file"
                    echo '```' >> "$output_file"
                    echo "</details>" >> "$output_file"
                    echo "" >> "$output_file"
                fi
            fi
        fi
    done

    if [[ "$has_failures" == "false" ]]; then
        echo "## All Tests Passed!" >> "$output_file"
        echo "" >> "$output_file"
        echo "No compatibility issues detected." >> "$output_file"
    fi

    echo "$output_file"
}

# Generate detailed per-test report
generate_detailed_report() {
    local results_dir="$1"
    local output_file="${results_dir}/detailed.md"

    cat > "$output_file" << EOF
# Detailed Test Results

EOF

    for test_dir in "${results_dir}"/*/; do
        if [[ -d "$test_dir" ]] && [[ -f "${test_dir}/result.txt" ]]; then
            local test_name=$(basename "$test_dir")
            local result=$(cat "${test_dir}/result.txt")

            echo "## ${test_name}" >> "$output_file"
            echo "" >> "$output_file"
            echo "**Result**: ${result}" >> "$output_file"
            echo "" >> "$output_file"

            # Exit codes
            if [[ -f "${test_dir}/original.exit" ]] && [[ -f "${test_dir}/ours.exit" ]]; then
                local orig_exit=$(cat "${test_dir}/original.exit")
                local our_exit=$(cat "${test_dir}/ours.exit")
                echo "**Exit Codes**: Original=${orig_exit}, Ours=${our_exit}" >> "$output_file"
                echo "" >> "$output_file"
            fi

            # Original stdout preview
            if [[ -f "${test_dir}/original.stdout" ]]; then
                echo "<details>" >> "$output_file"
                echo "<summary>Original stdout (first 50 lines)</summary>" >> "$output_file"
                echo "" >> "$output_file"
                echo '```' >> "$output_file"
                head -50 "${test_dir}/original.stdout" >> "$output_file"
                echo '```' >> "$output_file"
                echo "</details>" >> "$output_file"
                echo "" >> "$output_file"
            fi

            # Our stdout preview
            if [[ -f "${test_dir}/ours.stdout" ]]; then
                echo "<details>" >> "$output_file"
                echo "<summary>Our stdout (first 50 lines)</summary>" >> "$output_file"
                echo "" >> "$output_file"
                echo '```' >> "$output_file"
                head -50 "${test_dir}/ours.stdout" >> "$output_file"
                echo '```' >> "$output_file"
                echo "</details>" >> "$output_file"
                echo "" >> "$output_file"
            fi

            echo "---" >> "$output_file"
            echo "" >> "$output_file"
        fi
    done

    echo "$output_file"
}

# Print quick summary to terminal
print_quick_summary() {
    local results_dir="$1"

    echo ""
    echo "=========================================="
    echo "Quick Summary"
    echo "=========================================="

    local pass_count=0
    local fail_count=0

    for result_file in "${results_dir}"/*/result.txt; do
        if [[ -f "$result_file" ]]; then
            local result=$(cat "$result_file")
            if [[ "$result" == "PASS" ]]; then
                ((pass_count++))
            else
                ((fail_count++))
                local test_name=$(basename "$(dirname "$result_file")")
                echo -e "${RED}FAIL${NC}: ${test_name} - ${result}"
            fi
        fi
    done

    echo ""
    echo "Passed: ${pass_count}"
    echo "Failed: ${fail_count}"

    if [[ $fail_count -eq 0 ]]; then
        echo -e "${GREEN}All tests passed!${NC}"
    fi
}
