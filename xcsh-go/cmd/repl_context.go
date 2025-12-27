package cmd

import (
	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

// ContextPath represents the current navigation context in the REPL
type ContextPath struct {
	Domain string // e.g., "load_balancer" - empty at root
	Action string // e.g., "list" - empty if only domain context
}

// IsRoot returns true if at root context (no domain selected)
func (c *ContextPath) IsRoot() bool {
	return c.Domain == ""
}

// IsDomain returns true if in a domain context (but no action)
func (c *ContextPath) IsDomain() bool {
	return c.Domain != "" && c.Action == ""
}

// IsAction returns true if in an action context
func (c *ContextPath) IsAction() bool {
	return c.Domain != "" && c.Action != ""
}

// String returns the path as "domain/action" or "domain" or ""
func (c *ContextPath) String() string {
	if c.Domain == "" {
		return ""
	}
	if c.Action == "" {
		return c.Domain
	}
	return c.Domain + "/" + c.Action
}

// Reset clears the context to root
func (c *ContextPath) Reset() {
	c.Domain = ""
	c.Action = ""
}

// NavigateUp moves up one level in the context hierarchy
// Returns true if navigation occurred, false if already at root
func (c *ContextPath) NavigateUp() bool {
	if c.Action != "" {
		c.Action = ""
		return true
	}
	if c.Domain != "" {
		c.Domain = ""
		return true
	}
	return false // Already at root
}

// SetDomain enters a domain context
func (c *ContextPath) SetDomain(domain string) {
	c.Domain = domain
	c.Action = ""
}

// SetAction enters an action context within current domain
func (c *ContextPath) SetAction(action string) {
	c.Action = action
}

// ValidActions defines all valid action commands for domains
var ValidActions = map[string]bool{
	"list":          true,
	"get":           true,
	"create":        true,
	"delete":        true,
	"replace":       true,
	"apply":         true,
	"status":        true,
	"patch":         true,
	"add-labels":    true,
	"remove-labels": true,
}

// ContextValidator provides validation for navigation commands
type ContextValidator struct {
	domains map[string]bool // Cache of valid domain names and aliases
}

// NewContextValidator creates a validator with cached domain list
func NewContextValidator() *ContextValidator {
	v := &ContextValidator{
		domains: make(map[string]bool),
	}
	// Populate from DomainRegistry
	for domain := range types.DomainRegistry {
		v.domains[domain] = true
	}
	// Also include aliases
	for alias := range types.AliasRegistry {
		v.domains[alias] = true
	}
	return v
}

// IsValidDomain checks if input is a valid domain name or alias
func (v *ContextValidator) IsValidDomain(name string) bool {
	return v.domains[name]
}

// IsValidAction checks if input is a valid action command
func (v *ContextValidator) IsValidAction(name string) bool {
	return ValidActions[name]
}

// ResolveDomain returns the canonical domain name for an alias
func (v *ContextValidator) ResolveDomain(name string) (string, bool) {
	return types.ResolveDomain(name)
}
