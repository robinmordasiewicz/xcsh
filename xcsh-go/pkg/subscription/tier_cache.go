// Package subscription provides subscription tier caching via environment variables.
package subscription

import (
	"os"
	"strings"
)

// Environment variable names for subscription caching
const (
	// EnvSubscriptionTier caches the detected subscription tier to avoid repeated API calls.
	// Valid values: "Standard", "Advanced"
	EnvSubscriptionTier = "F5XC_SUBSCRIPTION_TIER"

	// EnvSubscriptionTierSource indicates how the tier was determined.
	// Values: "api" (detected from API), "env" (from environment variable), "default" (fallback)
	EnvSubscriptionTierSource = "F5XC_SUBSCRIPTION_TIER_SOURCE"
)

// GetCachedTier returns the cached subscription tier from the environment variable.
// Returns empty string if not set.
func GetCachedTier() string {
	tier := os.Getenv(EnvSubscriptionTier)
	return normalizeTierDisplay(tier)
}

// SetCachedTier stores the subscription tier in an environment variable.
// This persists for the duration of the process and can be inherited by child processes.
func SetCachedTier(tier string) error {
	normalizedTier := normalizeTierDisplay(tier)
	return os.Setenv(EnvSubscriptionTier, normalizedTier)
}

// SetCachedTierWithSource stores both the tier and its detection source.
func SetCachedTierWithSource(tier, source string) error {
	if err := SetCachedTier(tier); err != nil {
		return err
	}
	return os.Setenv(EnvSubscriptionTierSource, source)
}

// GetCachedTierSource returns how the cached tier was determined.
// Returns "env" if F5XC_SUBSCRIPTION_TIER was set externally,
// "api" if detected from API, or "default" if using fallback.
func GetCachedTierSource() string {
	source := os.Getenv(EnvSubscriptionTierSource)
	if source == "" {
		// If tier is set but source is not, assume it was set externally
		if GetCachedTier() != "" {
			return "env"
		}
		return ""
	}
	return source
}

// ClearCachedTier removes the cached subscription tier.
func ClearCachedTier() {
	_ = os.Unsetenv(EnvSubscriptionTier)
	_ = os.Unsetenv(EnvSubscriptionTierSource)
}

// IsTierCached returns true if the subscription tier is already cached.
func IsTierCached() bool {
	return GetCachedTier() != ""
}

// TierCacheInfo provides information about the cached tier.
type TierCacheInfo struct {
	Tier      string `json:"tier" yaml:"tier"`
	Source    string `json:"source" yaml:"source"`
	IsCached  bool   `json:"is_cached" yaml:"is_cached"`
	IsFromAPI bool   `json:"is_from_api" yaml:"is_from_api"`
	IsFromEnv bool   `json:"is_from_env" yaml:"is_from_env"`
	IsDefault bool   `json:"is_default" yaml:"is_default"`
}

// GetTierCacheInfo returns detailed information about the cached tier.
func GetTierCacheInfo() *TierCacheInfo {
	tier := GetCachedTier()
	source := GetCachedTierSource()

	return &TierCacheInfo{
		Tier:      tier,
		Source:    source,
		IsCached:  tier != "",
		IsFromAPI: source == "api",
		IsFromEnv: source == "env" || (tier != "" && source == ""),
		IsDefault: source == "default",
	}
}

// normalizeTierDisplay converts tier constants to display names.
// STANDARD -> Standard, ADVANCED -> Advanced
// Discontinued tiers (Basic, Premium) are mapped to active tiers.
func normalizeTierDisplay(tier string) string {
	upper := strings.ToUpper(strings.TrimSpace(tier))
	switch upper {
	case TierStandard, TierBasic:
		return "Standard" // Basic discontinued, maps to Standard
	case TierAdvanced, TierPremium:
		return "Advanced" // Premium discontinued, maps to Advanced
	default:
		// Return as-is if already in display format
		if tier == "Standard" || tier == "Basic" {
			return "Standard"
		}
		if tier == "Advanced" || tier == "Premium" {
			return "Advanced"
		}
		return tier
	}
}

// IsStandardTier returns true if the given tier is Standard (or discontinued Basic).
func IsStandardTier(tier string) bool {
	normalized := strings.ToUpper(strings.TrimSpace(tier))
	return normalized == TierStandard || normalized == TierBasic ||
		tier == "Standard" || tier == "Basic"
}

// IsAdvancedTier returns true if the given tier is Advanced (or discontinued Premium).
func IsAdvancedTier(tier string) bool {
	normalized := strings.ToUpper(strings.TrimSpace(tier))
	return normalized == TierAdvanced || normalized == TierPremium ||
		tier == "Advanced" || tier == "Premium"
}

// CompareTiers compares two tiers and returns:
// -1 if tier1 < tier2
//
//	0 if tier1 == tier2
//	1 if tier1 > tier2
//
// Note: Only Standard and Advanced are active tiers.
// Basic maps to Standard, Premium maps to Advanced.
func CompareTiers(tier1, tier2 string) int {
	order := func(tier string) int {
		// Handle display names first
		switch tier {
		case "Standard", "Basic":
			return 1
		case "Advanced", "Premium":
			return 2
		}
		// Then handle constants
		switch strings.ToUpper(strings.TrimSpace(tier)) {
		case "", TierNoTier:
			return 0
		case TierBasic, TierStandard:
			return 1 // Basic discontinued, maps to Standard
		case TierAdvanced, TierPremium:
			return 2 // Premium discontinued, maps to Advanced
		default:
			return 0
		}
	}

	o1, o2 := order(tier1), order(tier2)
	if o1 < o2 {
		return -1
	} else if o1 > o2 {
		return 1
	}
	return 0
}
