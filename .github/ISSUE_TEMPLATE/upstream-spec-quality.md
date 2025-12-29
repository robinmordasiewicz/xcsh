---
name: Upstream API Spec Quality Issue
about: Report specification organization or completeness issues to upstream repository
title: '[SPEC QUALITY] '
labels: upstream, spec-quality
---

## Summary
<!-- Brief description of the specification issue -->

## Affected Specification

- **Domain**: <!-- e.g., "api_security", "load_balancer" -->
- **Resource(s)**: <!-- e.g., "http_loadbalancer", "api_definition" -->
- **Spec Version**: <!-- e.g., "v1.0.30" -->
- **File(s)**: <!-- e.g., ".specs/domains/api_security.json" -->

## Current State

**Current Organization**:
<!-- Describe how the resource is currently organized -->

**Current Issue**:
<!-- Explain the problem with current organization -->

## Expected State

**Suggested Organization**:
<!-- How it should be organized -->

**Reasoning**:
<!-- Why this change would be beneficial -->

## Impact on xcsh

**User Experience Impact**:
<!-- How this affects xcsh CLI users -->

**CLI Example**:

```bash
# Current (problematic):
xcsh api_security list api_definition

# Expected (cleaner):
xcsh api_security list api_definition
```

## Additional Context

- [ ] This affects CLI command naming/organization
- [ ] This is a resource classification issue
- [ ] This is a missing metadata issue
- [ ] This is a naming consistency issue
- [ ] This is a documentation issue

## Related Issues
<!-- Link to any related issues in this repository or upstream -->
- Related upstream issue: [link]
- xcsh issue: [link]

## Workaround
<!-- Describe any temporary workaround for xcsh -->
We have documented this in `.specs/domain_config.yaml` with the following override:

```yaml
# If applicable, show any workarounds configured
```

---

**Note**: Issues of this type are automatically detected during daily upstream syncs. This issue may have been created automatically by the [xcsh sync-upstream-specs](https://github.com/robinmordasiewicz/f5xc-xcsh/blob/main/.github/workflows/sync-upstream-specs.yml) workflow.
