# xcsh API Reference - Complete Function Documentation

**Version**: 1.0
**Last Updated**: 2025-12-24
**Scope**: All public validation and formatting functions

---

## Table of Contents

1. [Domain Information Functions](#domain-information-functions)
2. [Tier Validation Functions](#tier-validation-functions)
3. [Preview Feature Functions](#preview-feature-functions)
4. [Category Functions](#category-functions)
5. [Use Case Functions](#use-case-functions)
6. [Workflow Functions](#workflow-functions)
7. [Type Definitions](#type-definitions)
8. [Error Codes](#error-codes)

---

## Domain Information Functions

### GetDomainInfo

**Location**: `pkg/types/domains_generated.go`

**Signature**:
```go
func GetDomainInfo(domainName string) (*DomainInfo, bool)
```

**Description**:
Retrieves complete information about a domain from the domain registry.

**Parameters**:
- `domainName` (string) - The name of the domain (e.g., "api", "kubernetes")

**Returns**:
- `*DomainInfo` - Domain information struct (nil if not found)
- `bool` - True if domain found, false otherwise

**Example**:
```go
info, found := types.GetDomainInfo("api")
if found {
    fmt.Println(info.DisplayName)        // "Api"
    fmt.Println(info.Description)        // Domain description
    fmt.Println(info.Category)           // "Security"
    fmt.Println(info.RequiresTier)       // "Standard"
    fmt.Println(len(info.UseCases))      // Number of use cases
}
```

**Error Handling**:
- Returns (nil, false) if domain doesn't exist
- Always safe to call with any string

---

## Tier Validation Functions

### ValidateTierAccess

**Location**: `pkg/validation/tier.go`

**Signature**:
```go
func ValidateTierAccess(userTier, requiredTier string) bool
```

**Description**:
Validates whether a user's subscription tier meets the requirements to access a domain.

**Parameters**:
- `userTier` (string) - User's subscription tier ("Standard", "Professional", "Enterprise")
- `requiredTier` (string) - Domain's required tier

**Returns**:
- `bool` - True if user has sufficient tier, false otherwise

**Tier Hierarchy**:
```
Standard < Professional < Enterprise
```

**Example**:
```go
// Professional user accessing Standard domain - allowed
ok := validation.ValidateTierAccess("Professional", "Standard")  // true

// Standard user accessing Professional domain - denied
ok := validation.ValidateTierAccess("Standard", "Professional")  // false
```

**Error Handling**:
- Unknown tiers default to compatible (conservative approach)
- Returns false only when user tier is explicitly lower

---

### GetDomainsByTier

**Location**: `pkg/validation/tier.go`

**Signature**:
```go
func GetDomainsByTier(tier string) []*types.DomainInfo
```

**Description**:
Returns all domains accessible at a specific subscription tier.

**Parameters**:
- `tier` (string) - Subscription tier ("Standard", "Professional", "Enterprise")

**Returns**:
- `[]*types.DomainInfo` - Slice of accessible domains

**Example**:
```go
// Get all Standard tier domains
standardDomains := validation.GetDomainsByTier("Standard")
fmt.Printf("Standard tier has %d domains\n", len(standardDomains))  // 25

// Get all Professional tier domains (includes Standard)
proDomains := validation.GetDomainsByTier("Professional")
fmt.Printf("Professional tier has %d domains\n", len(proDomains))  // 42
```

---

## Preview Feature Functions

### GetPreviewDomains

**Location**: `pkg/validation/preview.go`

**Signature**:
```go
func GetPreviewDomains() []*types.DomainInfo
```

**Description**:
Returns all domains marked as preview/beta features.

**Returns**:
- `[]*types.DomainInfo` - Slice of preview domains

**Example**:
```go
previewDomains := validation.GetPreviewDomains()
for _, domain := range previewDomains {
    fmt.Printf("Preview: %s - %s\n", domain.Name, domain.Description)
}
```

---

### IsPreviewDomain

**Location**: `pkg/validation/preview.go`

**Signature**:
```go
func IsPreviewDomain(domainName string) bool
```

**Description**:
Checks if a specific domain is a preview feature.

**Parameters**:
- `domainName` (string) - Domain name to check

**Returns**:
- `bool` - True if domain is preview, false otherwise

**Example**:
```go
if validation.IsPreviewDomain("generative_ai") {
    fmt.Println("This feature is preview/beta")
}
```

---

## Category Functions

### GetDomainsByCategory

**Location**: `pkg/validation/categories.go`

**Signature**:
```go
func GetDomainsByCategory(category string) []*types.DomainInfo
```

**Description**:
Returns all domains in a specific category.

**Parameters**:
- `category` (string) - Category name ("Security", "Networking", "Platform", "Infrastructure", "Operations", "Other", "AI")

**Returns**:
- `[]*types.DomainInfo` - Slice of domains in that category

**Example**:
```go
// Get all Security domains
securityDomains := validation.GetDomainsByCategory("Security")
fmt.Printf("Security category has %d domains\n", len(securityDomains))  // 9

for _, domain := range securityDomains {
    fmt.Printf("- %s: %s\n", domain.Name, domain.Description)
}
```

---

### GetAllCategories

**Location**: `pkg/validation/categories.go`

**Signature**:
```go
func GetAllCategories() []string
```

**Description**:
Returns all available domain categories.

**Returns**:
- `[]string` - Slice of category names (sorted)

**Example**:
```go
categories := validation.GetAllCategories()
for _, cat := range categories {
    fmt.Println(cat)
}
// Output:
// AI
// Infrastructure
// Networking
// Operations
// Other
// Platform
// Security
```

---

## Use Case Functions

### GetDomainUseCases

**Location**: `pkg/validation/use_cases.go`

**Signature**:
```go
func GetDomainUseCases(domain string) (string, error)
```

**Description**:
Retrieves formatted use cases for a specific domain.

**Parameters**:
- `domain` (string) - Domain name

**Returns**:
- `string` - Formatted use cases with "USE CASES:" header and bullet points
- `error` - Error if domain not found

**Example**:
```go
useCases, err := validation.GetDomainUseCases("api")
if err != nil {
    fmt.Println("Domain not found")
    return
}

fmt.Println(useCases)
// Output:
// USE CASES:
//   • Discover and catalog APIs
//   • Test API security and behavior
//   • Manage API credentials
//   • Define API groups and testing policies
```

---

### SearchUseCases

**Location**: `pkg/validation/use_cases.go`

**Signature**:
```go
func SearchUseCases(keyword string) []UseCase
```

**Description**:
Searches all use cases across all domains for a keyword.

**Parameters**:
- `keyword` (string) - Search keyword (case-insensitive)

**Returns**:
- `[]UseCase` - Slice of matching use cases

**Example**:
```go
results := validation.SearchUseCases("configure")
fmt.Printf("Found %d use cases with 'configure'\n", len(results))

for _, uc := range results {
    fmt.Printf("Domain: %s\n", uc.Domain)
    fmt.Printf("  %s\n", uc.Description)
}
```

**Notes**:
- Case-insensitive search
- Searches use case descriptions only
- Empty keyword returns all use cases

---

### GetAllUseCases

**Location**: `pkg/validation/use_cases.go`

**Signature**:
```go
func GetAllUseCases() []UseCase
```

**Description**:
Returns all use cases across all domains.

**Returns**:
- `[]UseCase` - Slice of all use cases (73+ items)

**Example**:
```go
allUseCases := validation.GetAllUseCases()
fmt.Printf("Total use cases: %d\n", len(allUseCases))

// Group by domain
byDomain := make(map[string]int)
for _, uc := range allUseCases {
    byDomain[uc.Domain]++
}

for domain, count := range byDomain {
    fmt.Printf("%s: %d use cases\n", domain, count)
}
```

---

### FormatUseCases

**Location**: `pkg/validation/use_cases.go`

**Signature**:
```go
func FormatUseCases(useCases []string) string
```

**Description**:
Formats use cases for display with bullet points.

**Parameters**:
- `useCases` ([]string) - Slice of use case descriptions

**Returns**:
- `string` - Formatted string with "USE CASES:" header

**Example**:
```go
cases := []string{
    "Configure authentication",
    "Manage users",
    "Set up SSO",
}

formatted := validation.FormatUseCases(cases)
fmt.Println(formatted)
// Output:
// USE CASES:
//   • Configure authentication
//   • Manage users
//   • Set up SSO
```

---

## Workflow Functions

### GetWorkflowSuggestions

**Location**: `pkg/validation/workflows.go`

**Signature**:
```go
func GetWorkflowSuggestions(domain string) []WorkflowSuggestion
```

**Description**:
Returns recommended workflow suggestions for a domain based on its category.

**Parameters**:
- `domain` (string) - Domain name

**Returns**:
- `[]WorkflowSuggestion` - Slice of workflow suggestions

**Example**:
```go
workflows := validation.GetWorkflowSuggestions("api")

for _, wf := range workflows {
    fmt.Printf("Workflow: %s\n", wf.Name)
    fmt.Printf("  Description: %s\n", wf.Description)
    fmt.Printf("  Category: %s\n", wf.Category)
    fmt.Printf("  Domains: %v\n", wf.Domains)
}
```

---

### GetRelatedDomains

**Location**: `pkg/validation/workflows.go`

**Signature**:
```go
func GetRelatedDomains(domain string) []*types.DomainInfo
```

**Description**:
Returns domains related to a specific domain using multi-strategy scoring.

**Parameters**:
- `domain` (string) - Domain name

**Returns**:
- `[]*types.DomainInfo` - Slice of up to 5 related domains, sorted by relevance

**Scoring Strategy**:
- Category compatibility: Strength 4 (highest priority)
- Use case overlap: Strength 3 (complementary)
- Tier compatibility: Strength 2 (can work together)

**Example**:
```go
relatedDomains := validation.GetRelatedDomains("api")

for _, domain := range relatedDomains {
    fmt.Printf("- %s (%s)\n", domain.Name, domain.Category)
}
// Output:
// - application_firewall (Security)
// - authentication (Security)
// - threat_campaign (Security)
// - network_security (Security)
// - cdn (Networking)
```

---

### FormatRelatedDomains

**Location**: `pkg/validation/workflows.go`

**Signature**:
```go
func FormatRelatedDomains(domains []*types.DomainInfo) string
```

**Description**:
Formats related domains for display with descriptions.

**Parameters**:
- `domains` ([]*types.DomainInfo) - Slice of domains to format

**Returns**:
- `string` - Formatted string with "RELATED DOMAINS:" header

**Example**:
```go
domains := validation.GetRelatedDomains("api")
formatted := validation.FormatRelatedDomains(domains)

fmt.Println(formatted)
// Output:
// RELATED DOMAINS:
//   • application_firewall - F5 Distributed Cloud Application Firewall...
//   • authentication - F5 Distributed Cloud Authentication...
//   [... up to 5 domains ...]
```

---

### FormatWorkflowSuggestions

**Location**: `pkg/validation/workflows.go`

**Signature**:
```go
func FormatWorkflowSuggestions(workflows []WorkflowSuggestion) string
```

**Description**:
Formats workflow suggestions for display.

**Parameters**:
- `workflows` ([]WorkflowSuggestion) - Slice of workflows to format

**Returns**:
- `string` - Formatted string with "SUGGESTED WORKFLOWS:" header

**Example**:
```go
workflows := validation.GetWorkflowSuggestions("api")
formatted := validation.FormatWorkflowSuggestions(workflows)

fmt.Println(formatted)
// Output:
// SUGGESTED WORKFLOWS:
//   • API Security Workflow
//     Secure APIs with firewall and threat detection
//     Involves: api, application_firewall, threat_campaign
//   [... up to 3 workflows ...]
```

---

### GetWorkflowsByCategory

**Location**: `pkg/validation/workflows.go`

**Signature**:
```go
func GetWorkflowsByCategory(category string) []WorkflowSuggestion
```

**Description**:
Returns all workflow suggestions for a specific category.

**Parameters**:
- `category` (string) - Category name ("Security", "Platform", "Infrastructure", etc.)

**Returns**:
- `[]WorkflowSuggestion` - Slice of workflows for that category

**Example**:
```go
securityWorkflows := validation.GetWorkflowsByCategory("Security")
fmt.Printf("Security category has %d workflows\n", len(securityWorkflows))

for _, wf := range securityWorkflows {
    fmt.Printf("- %s\n", wf.Name)
}
```

---

## Type Definitions

### DomainInfo

**Location**: `pkg/types/domains.go`

**Definition**:
```go
type DomainInfo struct {
    Name           string                 // Domain name (e.g., "api")
    DisplayName    string                 // Display name (e.g., "Api")
    Description    string                 // Domain description
    Aliases        []string               // Alternative names
    Complexity     string                 // "basic" or "advanced"
    IsPreview      bool                   // True if preview/beta
    RequiresTier   string                 // "Standard", "Professional", "Enterprise"
    Category       string                 // Domain category
    UseCases       []string               // Practical use cases
    RelatedDomains []string               // Related domain names
    CLIMetadata    map[string]interface{} // Additional metadata
}
```

**Usage**:
```go
info, _ := types.GetDomainInfo("api")

fmt.Println(info.Name)          // "api"
fmt.Println(info.DisplayName)   // "Api"
fmt.Println(info.Category)      // "Security"
fmt.Println(info.RequiresTier)  // "Standard"
fmt.Println(len(info.UseCases)) // 4
```

---

### WorkflowSuggestion

**Location**: `pkg/validation/workflows.go`

**Definition**:
```go
type WorkflowSuggestion struct {
    Name        string   // Workflow name (e.g., "API Security Workflow")
    Description string   // What this workflow accomplishes
    Domains     []string // Domains involved in this workflow
    Category    string   // Category (e.g., "Security")
}
```

**Usage**:
```go
workflows := validation.GetWorkflowSuggestions("api")

for _, wf := range workflows {
    fmt.Printf("%s\n", wf.Name)
    fmt.Printf("  %s\n", wf.Description)
    fmt.Printf("  Domains: %v\n", wf.Domains)
}
```

---

### UseCase

**Location**: `pkg/validation/use_cases.go`

**Definition**:
```go
type UseCase struct {
    Domain      string // Domain name
    Description string // Use case description
    Category    string // Domain category
}
```

**Usage**:
```go
results := validation.SearchUseCases("configure")

for _, uc := range results {
    fmt.Printf("Domain: %s\n", uc.Domain)
    fmt.Printf("Category: %s\n", uc.Category)
    fmt.Printf("Use case: %s\n", uc.Description)
}
```

---

## Error Codes

### Exit Codes

| Code | Meaning | Example |
|------|---------|---------|
| 0 | Success | Command completed successfully |
| 1 | General error | Unknown error occurred |
| 2 | Invalid argument | Wrong command syntax |
| 8 | Quota exceeded | Tier doesn't allow operation |
| 9 | Feature unavailable | Feature not in your tier |

### Common Errors

**Domain not found**
```
Error: domain not found: "invalid_domain"
Solution: Use "xcsh domains list" to find valid domain names
```

**Tier requirement not met**
```
Error: tier requirement not met (requires Professional)
Solution: Upgrade subscription or use domain available at your tier
```

**Authentication failed**
```
Error: authentication failed
Solution: Check API credentials and token expiration
```

---

## Summary

The xcsh API provides comprehensive functions for:

✅ **Domain Discovery** - Find and retrieve domain information
✅ **Tier Validation** - Check subscription access
✅ **Preview Features** - Identify beta/preview domains
✅ **Categories** - Organize domains by function
✅ **Use Cases** - Find practical examples
✅ **Workflows** - Discover domain combinations
✅ **Formatting** - Display information consistently

**Key Pattern**: All functions are safe to call with any input and handle errors gracefully.

---

*Last updated: 2025-12-24*
*Version: 1.0*
*Part of the xcsh data-driven CLI architecture*
