#!/bin/bash
# run-all-tests.sh - Master test runner for vesctl compatibility testing
# Phase 10: Integration & Regression Suite

# Don't use set -e so we can continue running phases even when some fail
# We track failures manually and exit with appropriate code at the end

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
cd "$SCRIPT_DIR"

# Default configuration
PHASE=""
VERBOSE=false
WITH_API=false
HELP=false
GENERATE_REPORT=true
OVERALL_RESULT=0  # Track if any phase failed

# Arrays to track phase results
declare -A PHASE_RESULTS
declare -A PHASE_PASS
declare -A PHASE_FAIL
declare -A PHASE_WARN
declare -A PHASE_TIME

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --phase)
            PHASE="$2"
            shift 2
            ;;
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --with-api)
            WITH_API=true
            shift
            ;;
        --no-report)
            GENERATE_REPORT=false
            shift
            ;;
        --help|-h)
            HELP=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Show help
if [[ "$HELP" == "true" ]]; then
    cat << EOF
vesctl Compatibility Test Runner

Usage: ./run-all-tests.sh [OPTIONS]

Options:
  --phase N      Run only phase N (1-9, or "all")
  --verbose, -v  Show verbose output
  --with-api     Run API-dependent tests (requires VES_* credentials)
  --no-report    Skip generating summary report
  --help, -h     Show this help message

Phases:
  1: Configure/Authentication commands (behavioral)
  2: Simple commands (version, completion)
  3: Configuration CRUD operations (help text)
  4: Namespace CRUD operations (requires API)
  5: Multi-resource validation (help tests + API list)
  6: Request commands (rpc, secrets, command-sequence)
  7: Site management commands (aws_vpc, azure_vnet)
  8: Output format precision (table, json, yaml, tsv)
  9: Error handling & edge cases

Environment Variables:
  ORIGINAL_VESCTL   Path to original vesctl binary
  OUR_VESCTL        Path to our vesctl binary
  VES_API_URL       F5 XC API URL
  VES_P12_FILE      Path to P12 credentials file
  VES_P12_PASSWORD  Password for P12 credentials

Examples:
  ./run-all-tests.sh                    # Run all phases (help tests only)
  ./run-all-tests.sh --phase 1          # Run phase 1 only
  ./run-all-tests.sh --with-api         # Run all phases including API tests
  ./run-all-tests.sh --phase all -v     # Run all phases with verbose output

EOF
    exit 0
fi

# Header
echo "=========================================="
echo "vesctl Compatibility Test Suite"
echo "=========================================="
echo ""
echo "Date: $(date)"
echo "Original: ${ORIGINAL_VESCTL:-${PROJECT_ROOT}/vesctl-0.2.35}"
echo "Ours: ${OUR_VESCTL:-${PROJECT_ROOT}/vesctl}"
echo ""

# Check if binaries exist (default to project root paths)
ORIGINAL_VESCTL="${ORIGINAL_VESCTL:-${PROJECT_ROOT}/vesctl-0.2.35}"
OUR_VESCTL="${OUR_VESCTL:-${PROJECT_ROOT}/vesctl}"

if [[ ! -x "$ORIGINAL_VESCTL" ]]; then
    echo "ERROR: Original vesctl not found at $ORIGINAL_VESCTL"
    echo "Download with: ./claudedocs/compatibility/scripts/download-original.sh"
    exit 1
fi

if [[ ! -x "$OUR_VESCTL" ]]; then
    echo "ERROR: Our vesctl not found at $OUR_VESCTL"
    echo "Build with: go build -o vesctl"
    exit 1
fi

export ORIGINAL_VESCTL OUR_VESCTL

# Check API credentials
if [[ "$WITH_API" == "true" ]]; then
    if [[ -z "$VES_P12_PASSWORD" ]] || [[ -z "$VES_P12_FILE" ]]; then
        echo "WARNING: API tests requested but credentials not set"
        echo "Set VES_P12_PASSWORD and VES_P12_FILE environment variables"
        WITH_API=false
    elif [[ ! -f "$VES_P12_FILE" ]]; then
        echo "WARNING: VES_P12_FILE does not exist: $VES_P12_FILE"
        WITH_API=false
    fi
fi

# Create results directory with timestamp
TIMESTAMP=$(date +%Y-%m-%d_%H-%M-%S)
RESULTS_DIR="${SCRIPT_DIR}/tests/results/${TIMESTAMP}"
mkdir -p "$RESULTS_DIR"

# Update latest symlink
ln -sfn "${TIMESTAMP}" "${SCRIPT_DIR}/tests/results/latest"

export RESULTS_DIR

# Run test phases
run_phase() {
    local phase=$1
    local script=$2
    local description=$3
    local requires_api=${4:-false}
    local phase_result=0
    local start_time end_time duration

    # Skip if specific phase requested and this isn't it
    if [[ -n "$PHASE" ]] && [[ "$PHASE" != "$phase" ]] && [[ "$PHASE" != "all" ]]; then
        return 0
    fi

    # Skip API-required phases if no API credentials
    if [[ "$requires_api" == "true" ]] && [[ "$WITH_API" != "true" ]]; then
        echo ""
        echo "[SKIP] Phase $phase: $description (requires --with-api)"
        PHASE_RESULTS[$phase]="SKIP"
        return 0
    fi

    if [[ -x "$script" ]]; then
        echo ""
        echo "=========================================="
        echo "Phase $phase: $description"
        echo "=========================================="

        start_time=$(date +%s)

        if [[ "$VERBOSE" == "true" ]]; then
            bash "$script" 2>&1 | tee "${RESULTS_DIR}/phase${phase}.log" || phase_result=$?
        else
            bash "$script" > "${RESULTS_DIR}/phase${phase}.log" 2>&1 || phase_result=$?
            # Show summary from log
            tail -20 "${RESULTS_DIR}/phase${phase}.log"
        fi

        end_time=$(date +%s)
        duration=$((end_time - start_time))
        PHASE_TIME[$phase]="${duration}s"

        # Extract pass/fail/warn counts from log (strip ANSI codes first)
        local pass_count=$(sed 's/\x1b\[[0-9;]*m//g' "${RESULTS_DIR}/phase${phase}.log" 2>/dev/null | grep -c '\[PASS\]' || echo 0)
        local fail_count=$(sed 's/\x1b\[[0-9;]*m//g' "${RESULTS_DIR}/phase${phase}.log" 2>/dev/null | grep -c '\[FAIL\]' || echo 0)
        local warn_count=$(sed 's/\x1b\[[0-9;]*m//g' "${RESULTS_DIR}/phase${phase}.log" 2>/dev/null | grep -c '\[WARN\]' || echo 0)

        PHASE_PASS[$phase]=$pass_count
        PHASE_FAIL[$phase]=$fail_count
        PHASE_WARN[$phase]=$warn_count

        if [[ $phase_result -ne 0 ]] || [[ $fail_count -gt 0 ]]; then
            PHASE_RESULTS[$phase]="FAIL"
            OVERALL_RESULT=1
        else
            PHASE_RESULTS[$phase]="PASS"
        fi

        echo ""
        echo "Phase $phase completed in ${duration}s (Pass: $pass_count, Fail: $fail_count, Warn: $warn_count)"
    else
        echo ""
        echo "[SKIP] Phase $phase: $description (script not found: $script)"
        PHASE_RESULTS[$phase]="SKIP"
    fi
}

# Execute all phases
echo ""
echo "Starting test execution..."
SUITE_START=$(date +%s)

# Phase 1: Configure/behavior tests (no API required)
run_phase 1 "${SCRIPT_DIR}/tests/phase1-configure/test-configure.sh" "Configure/Authentication" false

# Phase 2: Version/completion tests (no API required)
run_phase 2 "${SCRIPT_DIR}/tests/phase2-simple/test-simple.sh" "Simple Commands" false

# Phase 3: Auth structure tests (no API required)
run_phase 3 "${SCRIPT_DIR}/tests/phase3-auth-crud/test-no-api.sh" "Configuration CRUD" false

# Phase 4: Namespace CRUD tests (requires API credentials)
run_phase 4 "${SCRIPT_DIR}/tests/phase4-namespace/test-namespace-crud.sh" "Namespace CRUD" true

# Phase 5: Multi-resource validation (help tests, optionally API)
run_phase 5 "${SCRIPT_DIR}/tests/phase5-resources/test-resource-help.sh" "Multi-Resource Validation" false

# Phase 6: Request commands
run_phase 6 "${SCRIPT_DIR}/tests/phase6-request/test-request-help.sh" "Request Commands" false

# Phase 7: Site management commands
run_phase 7 "${SCRIPT_DIR}/tests/phase7-site/test-site-help.sh" "Site Management" false

# Phase 8: Output format precision (requires API)
run_phase 8 "${SCRIPT_DIR}/tests/phase8-output/test-output-formats.sh" "Output Formats" true

# Phase 9: Error handling & edge cases (requires API for some tests)
run_phase 9 "${SCRIPT_DIR}/tests/phase9-errors/test-error-handling.sh" "Error Handling" true

SUITE_END=$(date +%s)
SUITE_DURATION=$((SUITE_END - SUITE_START))

# Generate summary report
if [[ "$GENERATE_REPORT" == "true" ]]; then
    echo ""
    echo "Generating summary report..."

    REPORT_FILE="${RESULTS_DIR}/summary.md"

    cat > "$REPORT_FILE" << EOF
# vesctl Compatibility Test Report

**Date:** $(date)
**Duration:** ${SUITE_DURATION}s
**Original:** ${ORIGINAL_VESCTL}
**Ours:** ${OUR_VESCTL}

## Summary

| Phase | Description | Status | Pass | Fail | Warn | Time |
|-------|-------------|--------|------|------|------|------|
EOF

    # Add each phase to report
    total_pass=0
    total_fail=0
    total_warn=0
    phases_run=0

    for phase in 1 2 3 4 5 6 7 8 9; do
        status="${PHASE_RESULTS[$phase]:-N/A}"
        pass="${PHASE_PASS[$phase]:-0}"
        fail="${PHASE_FAIL[$phase]:-0}"
        warn="${PHASE_WARN[$phase]:-0}"
        time="${PHASE_TIME[$phase]:-N/A}"

        case $phase in
            1) desc="Configure/Authentication" ;;
            2) desc="Simple Commands" ;;
            3) desc="Configuration CRUD" ;;
            4) desc="Namespace CRUD" ;;
            5) desc="Multi-Resource Validation" ;;
            6) desc="Request Commands" ;;
            7) desc="Site Management" ;;
            8) desc="Output Formats" ;;
            9) desc="Error Handling" ;;
        esac

        if [[ "$status" != "SKIP" ]] && [[ "$status" != "N/A" ]]; then
            total_pass=$((total_pass + pass))
            total_fail=$((total_fail + fail))
            total_warn=$((total_warn + warn))
            phases_run=$((phases_run + 1))
        fi

        echo "| $phase | $desc | $status | $pass | $fail | $warn | $time |" >> "$REPORT_FILE"
    done

    cat >> "$REPORT_FILE" << EOF
| **Total** | | | **$total_pass** | **$total_fail** | **$total_warn** | **${SUITE_DURATION}s** |

## Overall Result

EOF

    if [[ $OVERALL_RESULT -eq 0 ]]; then
        echo "**PASS** - All phases completed successfully" >> "$REPORT_FILE"
    else
        echo "**FAIL** - Some phases had failures" >> "$REPORT_FILE"
    fi

    # Calculate pass rate
    total_tests=$((total_pass + total_fail))
    if [[ $total_tests -gt 0 ]]; then
        pass_rate=$((total_pass * 100 / total_tests))
        echo "" >> "$REPORT_FILE"
        echo "**Pass Rate:** ${pass_rate}% ($total_pass/$total_tests tests)" >> "$REPORT_FILE"
    fi

    cat >> "$REPORT_FILE" << EOF

## Test Environment

- **API Tests:** $(if [[ "$WITH_API" == "true" ]]; then echo "Enabled"; else echo "Disabled"; fi)
- **Phases Run:** $phases_run/9
- **Results Directory:** $RESULTS_DIR

## Log Files

EOF

    for phase in 1 2 3 4 5 6 7 8 9; do
        if [[ -f "${RESULTS_DIR}/phase${phase}.log" ]]; then
            echo "- [Phase $phase Log](phase${phase}.log)" >> "$REPORT_FILE"
        fi
    done

    cat >> "$REPORT_FILE" << EOF

---
*Generated by vesctl Compatibility Test Suite*
EOF

    echo "Report saved to: $REPORT_FILE"
fi

# Final summary
echo ""
echo "=========================================="
echo "Test Suite Complete"
echo "=========================================="
echo ""
echo "Duration: ${SUITE_DURATION}s"
echo "Results:  ${RESULTS_DIR}"
echo ""

# Print phase summary table
printf "%-7s %-25s %-8s %-6s %-6s %-6s\n" "Phase" "Description" "Status" "Pass" "Fail" "Warn"
printf "%-7s %-25s %-8s %-6s %-6s %-6s\n" "-----" "-------------------------" "------" "----" "----" "----"

for phase in 1 2 3 4 5 6 7 8 9; do
    status="${PHASE_RESULTS[$phase]:-N/A}"
    pass="${PHASE_PASS[$phase]:-0}"
    fail="${PHASE_FAIL[$phase]:-0}"
    warn="${PHASE_WARN[$phase]:-0}"

    case $phase in
        1) desc="Configure/Auth" ;;
        2) desc="Simple Commands" ;;
        3) desc="Config CRUD" ;;
        4) desc="Namespace CRUD" ;;
        5) desc="Multi-Resource" ;;
        6) desc="Request Commands" ;;
        7) desc="Site Management" ;;
        8) desc="Output Formats" ;;
        9) desc="Error Handling" ;;
    esac

    printf "%-7s %-25s %-8s %-6s %-6s %-6s\n" "$phase" "$desc" "$status" "$pass" "$fail" "$warn"
done

echo ""

if [[ $OVERALL_RESULT -eq 0 ]]; then
    echo "Overall: PASS"
else
    echo "Overall: FAIL"
fi

echo ""
echo "View full report: cat ${RESULTS_DIR}/summary.md"

# Exit with overall result (0 = all passed, 1 = some failed)
exit $OVERALL_RESULT
