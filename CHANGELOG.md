# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [5.0.0] - 2025-12-21

### ⚠️ BREAKING CHANGES

This is a major version release with **significant breaking changes** to the CLI structure. All users must update their command references.

- **Removed**: `f5xcctl configuration` command entirely (no deprecation period)
- **New structure**: `f5xcctl <domain> <operation> <resource>` instead of `f5xcctl configuration <operation> <resource>`
- All commands must use the new domain-based structure

### Migration Guide

**Old Format → New Format:**

```bash
# List resources
f5xcctl configuration list http_loadbalancer
  → f5xcctl load_balancer list http_loadbalancer
  → f5xcctl lb list http_loadbalancer (using alias)

# Get a resource
f5xcctl configuration get origin_pool mypool -n production
  → f5xcctl infrastructure get origin_pool mypool -n production
  → f5xcctl infra get origin_pool mypool -n production (using alias)

# Create a resource
f5xcctl configuration create certificate -i cert.yaml
  → f5xcctl identity create certificate -i cert.yaml

# Delete a resource
f5xcctl configuration delete http_loadbalancer my-lb --yes
  → f5xcctl load_balancer delete http_loadbalancer my-lb --yes
  → f5xcctl lb delete http_loadbalancer my-lb --yes
```

### Added

- **22 Domain-based command groups** organized by functional area:
  - `load_balancer` (alias: `lb`) - HTTP, TCP, UDP load balancing
  - `infrastructure` (alias: `infra`) - Core infrastructure resources
  - `security` (alias: `sec`) - WAF, bot defense, API security
  - `networking` (alias: `net`) - Network configuration and routing
  - `observability` (aliases: `obs`, `o11y`) - Monitoring and observability
  - `application` (aliases: `app`, `apps`) - Application configuration
  - `api_security` (alias: `apisec`) - API protection and security
  - `service_mesh` (alias: `mesh`) - Service mesh management
  - `identity` (alias: `iam`) - Identity and access management
  - `cloud_infrastructure` (alias: `infraprot`) - Cloud infrastructure protection
  - `integration` (alias: `int`) - Integration services
  - `connectivity` (alias: `connect`) - Connectivity services
  - `operations` (alias: `ops`) - Operations and management
  - `subscriptions` (alias: `sub`) - Subscription management
  - `tenancy` (alias: `tenant`) - Tenancy configuration
  - `ai_automation` (alias: `ai`) - AI and automation features
  - And 6 additional domains

- **Domain aliases** for convenience (e.g., `f5xcctl lb` instead of `f5xcctl load_balancer`)
- **Cross-domain resource access**: 132 resources now accessible from all applicable domains
  - Example: `http_loadbalancer` accessible from `load_balancer`, `security`, `api_security`, `applications`, `cdn`, `service_mesh`
- **Enhanced resource registry** with O(1) domain lookups via lazy-initialized indexing
- **Domain discovery** - List all resources available in a domain using `f5xcctl <domain> list <resource-type>`

### Changed

- **Command structure**: All operations must now specify the domain first
  - `<domain> <operation> <resource-type> [name]`
  - Operations: `list`, `get`, `create`, `replace`, `apply`, `delete`, `status`, `patch`, `add-labels`, `remove-labels`
- **Help text**: Updated to reflect new domain-based organization
- **Resource organization**: Resources are now grouped by domain rather than flat structure

### Removed

- **`f5xcctl configuration` command** and all related subcommands
  - All functionality is now available through domain-based commands
  - No equivalent direct command - must use appropriate domain

### Fixed

- Improved thread safety in registry with RWMutex lazy initialization
- Better resource organization for discoverability

### Technical Details

- Added `PrimaryDomain` and `Domains` fields to ResourceType
- Enhanced Registry with domain indexing for efficient lookups
- Updated schema generator to extract domain mappings from enriched OpenAPI specs
- Single domain command factory generates all 22 domain commands (avoiding duplication)
- 268 resources regenerated with domain mappings
- All existing tests pass without modification

### Infrastructure

- Bumped version from 4.45.0 to 5.0.0
- Updated Go version to 1.25.5
- All CI/CD checks passing

### Documentation

- Updated command examples throughout
- Migration guide available in this CHANGELOG
- Help text for all new domain commands included

---

## [4.45.0] - Previous Release

See git history for details on v4.45.0 and earlier releases.

---

## Upgrading to v5.0.0

### Required Changes

1. **Update all `f5xcctl configuration` commands** to use the new domain-based structure
2. **Scripts and automation** must be updated to use new command format
3. **Documentation** referencing old commands needs updating

### Benefits

- ✅ Better organization matching API specification structure
- ✅ Easier resource discovery through domain grouping
- ✅ Convenience aliases for common domains
- ✅ Cross-domain access for frequently multi-domain resources
- ✅ Clearer command hierarchy

### Support

For migration assistance or questions about the new structure:

- Check `f5xcctl <domain> --help` for available operations
- Use `f5xcctl <domain> <operation> --help` for operation-specific help
- See examples in domain help text

---

[5.0.0]: https://github.com/robinmordasiewicz/f5xcctl/releases/tag/v5.0.0
[4.45.0]: https://github.com/robinmordasiewicz/f5xcctl/releases/tag/v4.45.0
