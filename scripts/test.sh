#!/bin/bash
# vesctl CLI Test Runner
#
# This script runs all tests for the vesctl CLI tool.
#
# Usage:
#   ./scripts/test.sh              # Run all tests
#   ./scripts/test.sh unit         # Run only unit tests
#   ./scripts/test.sh integration  # Run only integration tests
#   ./scripts/test.sh -v           # Run with verbose output
#
# Environment Variables (required for integration tests):
#   VES_API_URL      - API URL (e.g., https://tenant.staging.volterra.us)
#   VES_API_P12_FILE - Path to P12 certificate bundle
#   VES_P12_PASSWORD - Password for P12 bundle

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default settings
VERBOSE=""
TEST_TYPE="all"
COVERAGE=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        unit)
            TEST_TYPE="unit"
            shift
            ;;
        integration)
            TEST_TYPE="integration"
            shift
            ;;
        -v|--verbose)
            VERBOSE="-v"
            shift
            ;;
        --coverage)
            COVERAGE="-coverprofile=coverage.out"
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [unit|integration] [-v|--verbose] [--coverage]"
            echo ""
            echo "Options:"
            echo "  unit          Run only unit tests"
            echo "  integration   Run only integration tests"
            echo "  -v, --verbose Enable verbose output"
            echo "  --coverage    Generate coverage report"
            echo ""
            echo "Environment Variables (for integration tests):"
            echo "  VES_API_URL      API URL"
            echo "  VES_API_P12_FILE Path to P12 certificate bundle"
            echo "  VES_P12_PASSWORD Password for P12 bundle"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

echo -e "${GREEN}vesctl CLI Test Runner${NC}"
echo "================================"

# Build the binary first
echo -e "\n${YELLOW}Building vesctl binary...${NC}"
go build -o vesctl .
echo -e "${GREEN}✓ Build successful${NC}"

# Check if integration test environment is configured
check_integration_env() {
    if [[ -z "$VES_API_URL" ]]; then
        echo -e "${YELLOW}Warning: VES_API_URL not set${NC}"
        return 1
    fi
    if [[ -z "$VES_API_P12_FILE" ]]; then
        echo -e "${YELLOW}Warning: VES_API_P12_FILE not set${NC}"
        return 1
    fi
    if [[ -z "$VES_P12_PASSWORD" ]]; then
        echo -e "${YELLOW}Warning: VES_P12_PASSWORD not set${NC}"
        return 1
    fi
    if [[ ! -f "$VES_API_P12_FILE" ]]; then
        echo -e "${YELLOW}Warning: P12 file not found at $VES_API_P12_FILE${NC}"
        return 1
    fi
    return 0
}

# Run unit tests
run_unit_tests() {
    echo -e "\n${YELLOW}Running unit tests...${NC}"
    echo "--------------------------------"

    # VES_P12_PASSWORD should already be set in environment

    go test $VERBOSE $COVERAGE ./pkg/... 2>&1 | while IFS= read -r line; do
        if [[ $line == *"PASS"* ]]; then
            echo -e "${GREEN}$line${NC}"
        elif [[ $line == *"FAIL"* ]]; then
            echo -e "${RED}$line${NC}"
        elif [[ $line == *"SKIP"* ]]; then
            echo -e "${YELLOW}$line${NC}"
        else
            echo "$line"
        fi
    done

    if [[ ${PIPESTATUS[0]} -eq 0 ]]; then
        echo -e "\n${GREEN}✓ Unit tests passed${NC}"
        return 0
    else
        echo -e "\n${RED}✗ Unit tests failed${NC}"
        return 1
    fi
}

# Run integration tests
run_integration_tests() {
    echo -e "\n${YELLOW}Running integration tests...${NC}"
    echo "--------------------------------"

    if ! check_integration_env; then
        echo -e "${YELLOW}Skipping integration tests - environment not configured${NC}"
        echo ""
        echo "To run integration tests, set these environment variables:"
        echo "  export VES_API_URL=\"https://your-tenant.staging.volterra.us\""
        echo "  export VES_API_P12_FILE=\"/path/to/cert.p12\""
        echo "  export VES_P12_PASSWORD=\"your-password\""
        return 0
    fi

    echo -e "Using API URL: ${GREEN}$VES_API_URL${NC}"
    echo -e "Using P12 file: ${GREEN}$VES_API_P12_FILE${NC}"

    go test $VERBOSE ./tests/integration/... 2>&1 | while IFS= read -r line; do
        if [[ $line == *"PASS"* ]]; then
            echo -e "${GREEN}$line${NC}"
        elif [[ $line == *"FAIL"* ]]; then
            echo -e "${RED}$line${NC}"
        elif [[ $line == *"SKIP"* ]]; then
            echo -e "${YELLOW}$line${NC}"
        else
            echo "$line"
        fi
    done

    if [[ ${PIPESTATUS[0]} -eq 0 ]]; then
        echo -e "\n${GREEN}✓ Integration tests passed${NC}"
        return 0
    else
        echo -e "\n${RED}✗ Integration tests failed${NC}"
        return 1
    fi
}

# Run tests based on type
UNIT_RESULT=0
INTEGRATION_RESULT=0

case $TEST_TYPE in
    unit)
        run_unit_tests
        UNIT_RESULT=$?
        ;;
    integration)
        run_integration_tests
        INTEGRATION_RESULT=$?
        ;;
    all)
        run_unit_tests
        UNIT_RESULT=$?
        run_integration_tests
        INTEGRATION_RESULT=$?
        ;;
esac

# Generate coverage report if requested
if [[ -n "$COVERAGE" && -f "coverage.out" ]]; then
    echo -e "\n${YELLOW}Generating coverage report...${NC}"
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}✓ Coverage report: coverage.html${NC}"
fi

# Print summary
echo ""
echo "================================"
echo -e "${GREEN}Test Summary${NC}"
echo "================================"

if [[ "$TEST_TYPE" == "all" || "$TEST_TYPE" == "unit" ]]; then
    if [[ $UNIT_RESULT -eq 0 ]]; then
        echo -e "Unit Tests:        ${GREEN}PASSED${NC}"
    else
        echo -e "Unit Tests:        ${RED}FAILED${NC}"
    fi
fi

if [[ "$TEST_TYPE" == "all" || "$TEST_TYPE" == "integration" ]]; then
    if [[ $INTEGRATION_RESULT -eq 0 ]]; then
        echo -e "Integration Tests: ${GREEN}PASSED${NC}"
    else
        echo -e "Integration Tests: ${RED}FAILED${NC}"
    fi
fi

# Exit with appropriate code
if [[ $UNIT_RESULT -ne 0 || $INTEGRATION_RESULT -ne 0 ]]; then
    exit 1
fi

exit 0
