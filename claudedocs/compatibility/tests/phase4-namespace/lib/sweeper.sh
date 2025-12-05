#!/bin/bash
# sweeper.sh - Resource cleanup system for vesctl compatibility testing
#
# This module provides functions to track created resources and ensure
# complete cleanup after tests, including orphaned resources from failed runs.

# State file location (set by init_sweeper)
SWEEPER_STATE_FILE=""
SWEEPER_TEST_PREFIX="vesctl-test-"

# Initialize sweeper with state file path
# Usage: init_sweeper <results_dir>
init_sweeper() {
    local results_dir="$1"
    SWEEPER_STATE_FILE="${results_dir}/state.json"

    # Create initial state file
    cat > "$SWEEPER_STATE_FILE" << 'EOF'
{
    "tracked_resources": [],
    "orphaned_cleanup": [],
    "test_run_id": "",
    "started_at": ""
}
EOF

    # Set test run ID and timestamp
    local run_id=$(date +%s)
    local started_at=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    local tmp=$(mktemp)
    jq --arg id "$run_id" --arg time "$started_at" \
       '.test_run_id = $id | .started_at = $time' \
       "$SWEEPER_STATE_FILE" > "$tmp" && mv "$tmp" "$SWEEPER_STATE_FILE"

    log_info "Sweeper initialized with state file: $SWEEPER_STATE_FILE"
}

# Track a created resource
# Usage: track_resource <resource_type> <resource_name> [namespace]
track_resource() {
    local resource_type="$1"
    local resource_name="$2"
    local namespace="${3:-}"
    local created_at=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    if [[ -z "$SWEEPER_STATE_FILE" || ! -f "$SWEEPER_STATE_FILE" ]]; then
        log_warn "Sweeper not initialized, cannot track resource: $resource_type/$resource_name"
        return 1
    fi

    # Add to tracked resources using jq
    local tmp=$(mktemp)
    jq --arg type "$resource_type" \
       --arg name "$resource_name" \
       --arg ns "$namespace" \
       --arg time "$created_at" \
       '.tracked_resources += [{
           "type": $type,
           "name": $name,
           "namespace": $ns,
           "created_at": $time,
           "deleted": false
       }]' "$SWEEPER_STATE_FILE" > "$tmp" && mv "$tmp" "$SWEEPER_STATE_FILE"

    log_info "Tracking resource: $resource_type/$resource_name"
}

# Mark resource as deleted
# Usage: mark_deleted <resource_type> <resource_name>
mark_deleted() {
    local resource_type="$1"
    local resource_name="$2"

    if [[ -z "$SWEEPER_STATE_FILE" || ! -f "$SWEEPER_STATE_FILE" ]]; then
        return 1
    fi

    local tmp=$(mktemp)
    jq --arg type "$resource_type" \
       --arg name "$resource_name" \
       '(.tracked_resources[] | select(.type == $type and .name == $name)) .deleted = true' \
       "$SWEEPER_STATE_FILE" > "$tmp" && mv "$tmp" "$SWEEPER_STATE_FILE"

    log_info "Marked as deleted: $resource_type/$resource_name"
}

# Get list of resources not yet deleted
# Returns: lines of "type name namespace"
get_pending_cleanup() {
    if [[ -z "$SWEEPER_STATE_FILE" || ! -f "$SWEEPER_STATE_FILE" ]]; then
        return
    fi

    jq -r '.tracked_resources[] | select(.deleted == false) | "\(.type) \(.name) \(.namespace)"' "$SWEEPER_STATE_FILE"
}

# Scan for orphaned test resources from previous runs
# Usage: scan_orphaned_resources <binary>
scan_orphaned_resources() {
    local binary="$1"

    # List namespaces and find any with test prefix
    $binary configuration list namespace --outfmt json 2>/dev/null | \
        jq -r '.items[]?.metadata.name // empty' 2>/dev/null | \
        grep "^${SWEEPER_TEST_PREFIX}" || true
}

# Delete a resource using appropriate binary
# Usage: delete_resource <binary> <resource_type> <resource_name> [namespace]
delete_resource() {
    local binary="$1"
    local resource_type="$2"
    local resource_name="$3"
    local namespace="${4:-}"

    log_info "Deleting resource: $resource_type/$resource_name"

    if [[ "$resource_type" == "namespace" ]]; then
        # Namespace doesn't need -n flag
        $binary configuration delete "$resource_type" "$resource_name" 2>&1 || true
    else
        # Namespaced resources need -n flag
        $binary configuration delete "$resource_type" "$resource_name" -n "$namespace" 2>&1 || true
    fi
}

# Run full cleanup sweep
# Usage: run_cleanup_sweep <binary> <results_dir>
run_cleanup_sweep() {
    local binary="$1"
    local results_dir="$2"

    log_info "Running cleanup sweep..."

    # First: clean up any tracked but not deleted resources
    local pending=$(get_pending_cleanup)
    if [[ -n "$pending" ]]; then
        log_info "Cleaning up tracked resources..."
        echo "$pending" | while IFS=' ' read -r type name ns; do
            [[ -z "$type" || -z "$name" ]] && continue
            log_info "  Cleaning up tracked: $type/$name"
            delete_resource "$binary" "$type" "$name" "$ns"
            mark_deleted "$type" "$name"
        done
    fi

    # Second: scan for orphaned test resources
    log_info "Scanning for orphaned test resources..."
    local orphaned=$(scan_orphaned_resources "$binary")
    if [[ -n "$orphaned" ]]; then
        log_warn "Found orphaned test resources from previous runs:"
        echo "$orphaned" | while read -r name; do
            [[ -z "$name" ]] && continue
            log_info "  Deleting orphaned: namespace/$name"
            delete_resource "$binary" "namespace" "$name" ""

            # Track in state file
            if [[ -n "$SWEEPER_STATE_FILE" && -f "$SWEEPER_STATE_FILE" ]]; then
                local tmp=$(mktemp)
                jq --arg name "$name" \
                   '.orphaned_cleanup += [$name]' "$SWEEPER_STATE_FILE" > "$tmp" && mv "$tmp" "$SWEEPER_STATE_FILE"
            fi
        done
    else
        log_info "No orphaned test resources found"
    fi
}

# Verify cleanup was successful
# Usage: verify_cleanup <binary>
# Returns: 0 if clean, 1 if test resources remain
verify_cleanup() {
    local binary="$1"

    log_info "Verifying cleanup..."

    local remaining=$(scan_orphaned_resources "$binary")
    if [[ -n "$remaining" ]]; then
        log_fail "Cleanup verification failed - remaining test resources:"
        echo "$remaining" | while read -r name; do
            echo "  - $name"
        done
        return 1
    fi

    # Also check tracked resources
    local pending=$(get_pending_cleanup)
    if [[ -n "$pending" ]]; then
        log_warn "Some tracked resources may not have been deleted:"
        echo "$pending"
    fi

    log_pass "Cleanup verification successful - no test resources remaining"
    return 0
}

# Get cleanup summary from state file
# Usage: get_cleanup_summary
get_cleanup_summary() {
    if [[ -z "$SWEEPER_STATE_FILE" || ! -f "$SWEEPER_STATE_FILE" ]]; then
        echo "No state file available"
        return
    fi

    local tracked=$(jq '.tracked_resources | length' "$SWEEPER_STATE_FILE")
    local deleted=$(jq '[.tracked_resources[] | select(.deleted == true)] | length' "$SWEEPER_STATE_FILE")
    local orphaned=$(jq '.orphaned_cleanup | length' "$SWEEPER_STATE_FILE")

    echo "Cleanup Summary:"
    echo "  Tracked resources: $tracked"
    echo "  Deleted: $deleted"
    echo "  Orphaned cleaned: $orphaned"
}
