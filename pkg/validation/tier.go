// Package validation provides subscription tier validation and access control.
package validation

import (
	"fmt"
)

// Tier constants - subscription tier levels for domain access
const (
	TierStandard     = "Standard"
	TierProfessional = "Professional"
	TierEnterprise   = "Enterprise"
)

// TierLevel returns the numeric level of a tier for comparison purposes.
// Higher numbers indicate more privileged tiers.
// Returns 0 for unknown tiers (treated as lowest level).
func TierLevel(tier string) int {
	switch tier {
	case TierStandard:
		return 1
	case TierProfessional:
		return 2
	case TierEnterprise:
		return 3
	default:
		return 0 // Unknown tier, treated as no access
	}
}

// IsSufficientTier checks if the current tier is sufficient to access a domain requiring a specific tier.
// Tier hierarchy: Standard (1) < Professional (2) < Enterprise (3)
//
// Examples:
// - IsSufficientTier("Professional", "Standard") → true (Professional >= Standard)
// - IsSufficientTier("Professional", "Professional") → true (exact match)
// - IsSufficientTier("Standard", "Professional") → false (Standard < Professional)
// - IsSufficientTier("Enterprise", "Any") → true (Enterprise is highest)
// - IsSufficientTier("Standard", "") → true (empty requirement defaults to accessible)
func IsSufficientTier(currentTier, requiredTier string) bool {
	// Empty required tier means no restriction
	if requiredTier == "" {
		return true
	}

	currentLevel := TierLevel(currentTier)
	requiredLevel := TierLevel(requiredTier)

	// If either tier is unknown (level 0), deny access (fail secure)
	if currentLevel == 0 || requiredLevel == 0 {
		// Special case: if required tier is unknown but current is valid, deny
		// If current tier is unknown but required is valid, deny
		return false
	}

	return currentLevel >= requiredLevel
}

// TierName returns a user-friendly name for a tier.
func TierName(tier string) string {
	switch tier {
	case TierStandard:
		return "Standard"
	case TierProfessional:
		return "Professional"
	case TierEnterprise:
		return "Enterprise"
	default:
		return tier
	}
}

// GetNextTier returns the next tier in the hierarchy after the given tier.
// Returns empty string if already at highest tier.
func GetNextTier(currentTier string) string {
	switch currentTier {
	case TierStandard:
		return TierProfessional
	case TierProfessional:
		return TierEnterprise
	case TierEnterprise:
		return "" // Already at highest tier
	default:
		return TierProfessional // Default next tier
	}
}

// GetUpgradePath returns a formatted string describing the upgrade needed.
// Returns empty string if tier is already sufficient.
func GetUpgradePath(currentTier, requiredTier string) string {
	if IsSufficientTier(currentTier, requiredTier) {
		return ""
	}

	nextTier := GetNextTier(currentTier)
	if nextTier == "" {
		return fmt.Sprintf("You already have %s tier", currentTier)
	}

	return fmt.Sprintf("Upgrade from %s to %s tier", TierName(currentTier), TierName(requiredTier))
}

// TierAccessError represents an error when a user lacks sufficient tier for a resource.
type TierAccessError struct {
	Domain        string // The domain being accessed
	CurrentTier   string // User's current tier
	RequiredTier  string // Required tier for the domain
	DomainDisplay string // Human-readable domain name
}

// Error implements the error interface and returns a formatted error message.
func (e *TierAccessError) Error() string {
	upgradePath := GetUpgradePath(e.CurrentTier, e.RequiredTier)
	if upgradePath == "" {
		// This shouldn't happen - would mean sufficient tier already
		return fmt.Sprintf("domain '%s' is not accessible", e.Domain)
	}

	return fmt.Sprintf(
		"Domain '%s' requires %s tier\nYour subscription: %s\n%s\n\n"+
			"To upgrade your subscription, visit: https://console.volterra.io/account/upgrade\n"+
			"For assistance, contact F5 support: support@f5.com",
		e.DomainDisplay,
		TierName(e.RequiredTier),
		TierName(e.CurrentTier),
		upgradePath,
	)
}

// NewTierAccessError creates a new TierAccessError with the given parameters.
func NewTierAccessError(domain, domainDisplay, currentTier, requiredTier string) *TierAccessError {
	return &TierAccessError{
		Domain:        domain,
		CurrentTier:   currentTier,
		RequiredTier:  requiredTier,
		DomainDisplay: domainDisplay,
	}
}
