#!/bin/bash
# run-all-tests.sh - Master test runner for vesctl compatibility testing

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
cd "$SCRIPT_DIR"

# Default configuration
PHASE=""
VERBOSE=false
WITH_API=false
HELP=false

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
  --phase N      Run only phase N (1-4)
  --verbose, -v  Show verbose output
  --with-api     Run API-dependent tests (requires VES_P12_* credentials)
  --help, -h     Show this help message

Phases:
  1: Configure/Authentication commands
  2: Simple commands (version, completion)
  3: Configuration CRUD operations
  4: Other commands (request, site, api-endpoint)

Environment Variables:
  ORIGINAL_VESCTL   Path to original vesctl (default: ./vesctl-0.2.47-original)
  OUR_VESCTL        Path to our vesctl (default: ./vesctl.darwin-arm64)
  VES_P12_PASSWORD  Password for P12 credentials
  VES_P12_FILE      Path to P12 credentials file

Examples:
  ./run-all-tests.sh                    # Run all phases (help tests only)
  ./run-all-tests.sh --phase 1          # Run phase 1 only
  ./run-all-tests.sh --with-api         # Run all phases including API tests
  ./run-all-tests.sh --phase 1 -v       # Run phase 1 with verbose output

EOF
    exit 0
fi

# Header
echo "=========================================="
echo "vesctl Compatibility Test Suite"
echo "=========================================="
echo ""
echo "Date: $(date)"
echo "Original: ${ORIGINAL_VESCTL:-${PROJECT_ROOT}/vesctl-0.2.47-original}"
echo "Ours: ${OUR_VESCTL:-${PROJECT_ROOT}/vesctl.darwin-arm64}"
echo ""

# Check if binaries exist (default to project root paths)
ORIGINAL_VESCTL="${ORIGINAL_VESCTL:-${PROJECT_ROOT}/vesctl-0.2.47-original}"
OUR_VESCTL="${OUR_VESCTL:-${PROJECT_ROOT}/vesctl.darwin-arm64}"

if [[ ! -x "$ORIGINAL_VESCTL" ]]; then
    echo "ERROR: Original vesctl not found at $ORIGINAL_VESCTL"
    echo "Download with: ./claudedocs/compatibility/scripts/download-original.sh"
    exit 1
fi

if [[ ! -x "$OUR_VESCTL" ]]; then
    echo "ERROR: Our vesctl not found at $OUR_VESCTL"
    echo "Build with: go build -o vesctl.darwin-arm64"
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

# Run test phases
run_phase() {
    local phase=$1
    local script=$2

    if [[ -n "$PHASE" ]] && [[ "$PHASE" != "$phase" ]]; then
        return 0
    fi

    if [[ -x "$script" ]]; then
        echo ""
        echo "=========================================="
        echo "Running Phase $phase"
        echo "=========================================="

        if [[ "$VERBOSE" == "true" ]]; then
            bash "$script"
        else
            bash "$script" 2>&1
        fi
    else
        echo "Phase $phase script not found or not executable: $script"
    fi
}

# Execute phases
run_phase 1 "${SCRIPT_DIR}/tests/phase1-configure/test-configure.sh"

# Phase 2-4 can be added later
if [[ -x "${SCRIPT_DIR}/tests/phase2-simple/test-simple.sh" ]]; then
    run_phase 2 "${SCRIPT_DIR}/tests/phase2-simple/test-simple.sh"
fi

if [[ -x "${SCRIPT_DIR}/tests/phase3-crud/test-crud.sh" ]]; then
    run_phase 3 "${SCRIPT_DIR}/tests/phase3-crud/test-crud.sh"
fi

if [[ -x "${SCRIPT_DIR}/tests/phase4-other/test-other.sh" ]]; then
    run_phase 4 "${SCRIPT_DIR}/tests/phase4-other/test-other.sh"
fi

echo ""
echo "=========================================="
echo "All requested phases complete"
echo "=========================================="
echo ""
echo "Results: ${SCRIPT_DIR}/results/latest/"
echo ""
echo "View summary:"
echo "  cat ${SCRIPT_DIR}/results/latest/summary.md"
