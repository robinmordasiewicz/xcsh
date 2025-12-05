# vesctl Compatibility Report

**Generated**: $(date)

## Summary

| Metric | Count |
|--------|-------|
| Original vesctl resources | 140 |
| Our vesctl resources | 269 |
| Missing from ours | 22 |
| Additional in ours | 151 |

## Test Results

- **Passed**: 18 (64%)
- **Failed**: 10 (36%)

### Passing Tests
- `vesctl --help` - Identical
- `vesctl configuration --help` - Identical
- All specific resource help (namespace, http_loadbalancer, origin_pool, etc.)

### Failing Tests
Tests fail where the "Available Commands" lists differ due to resource additions/removals.

## Missing Resources (22)

These resources exist in the original vesctl but not in our implementation.
**None have current API spec files** - they appear to be deprecated/removed from the F5 XC API.

| Resource | Status |
|----------|--------|
| addon_service | Deprecated - no API spec |
| allowed_domain | Deprecated - no API spec |
| allowed_tenant | Deprecated - no API spec |
| apm | Deprecated - no API spec |
| child_tenant | Deprecated - no API spec |
| child_tenant_manager | Deprecated - no API spec |
| deployment | Deprecated - no API spec |
| dns_proxy | Deprecated - no API spec |
| infraprotect_ta_server | Deprecated - no API spec |
| kms_key | Deprecated - no API spec |
| kms_policy | Deprecated - no API spec |
| kms_policy_rule | Deprecated - no API spec |
| managed_tenant | Deprecated - no API spec |
| mitigated_domain | Deprecated - no API spec |
| mobile_base_config | Deprecated - no API spec |
| navigation_tile | Deprecated - no API spec |
| plan | Deprecated - no API spec |
| protected_application | Deprecated - no API spec |
| protected_domain | Deprecated - no API spec |
| tenant_profile | Deprecated - no API spec |
| v1_dns_monitor | Deprecated - no API spec |
| v1_http_monitor | Deprecated - no API spec |

## Additional Resources (151)

Our implementation includes 151 resources that don't exist in the original vesctl.
These are from newer F5 XC API specifications.

Key additions include:
- AI/ML services (ai_assistant, ai_data_bfdp, etc.)
- API Security (api_sec_*, api_credential)
- BIG-IP integration (bigip_*, bigcne_*)
- Bot Defense (bot_defense_*, shape_*)
- CDN (cdn_cache_rule)
- DNS enhancements (dns_zone_rrset, dns_zone_subscription)
- Network security (infraprotect_*, ike_*)
- Observability (graph_*, flow_*, synthetic_monitor_*)
- And many more...

## Recommendations

1. **Help text differences are expected** - Our implementation is based on current API specs
2. **No action required for missing resources** - They're deprecated in the current API
3. **Focus on behavioral consistency** for resources that exist in both implementations
4. **Consider adding deprecation notices** if backwards compatibility is critical

## Next Steps

1. Run API tests with credentials to verify actual API behavior matches
2. Test common operations (list, get, create, delete) on shared resources
3. Document any behavioral differences in shared resources
