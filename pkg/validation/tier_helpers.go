package validation

import (
	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

// ValidateTierAccess checks if a user's subscription tier meets domain requirements
func ValidateTierAccess(userTier, requiredTier string) bool {
	tierOrder := map[string]int{
		"Standard":     1,
		"Professional": 2,
		"Enterprise":   3,
	}

	userLevel := tierOrder[userTier]
	requiredLevel := tierOrder[requiredTier]

	if userLevel == 0 || requiredLevel == 0 {
		return true // Unknown tiers default to compatible
	}

	return userLevel >= requiredLevel
}

// GetDomainsByTier returns all domains accessible at a specific subscription tier
func GetDomainsByTier(tier string) []*types.DomainInfo {
	var result []*types.DomainInfo

	for _, domain := range types.DomainRegistry {
		if ValidateTierAccess(tier, domain.RequiresTier) {
			result = append(result, domain)
		}
	}

	return result
}

// GetPreviewDomains returns all domains marked as preview/beta features
func GetPreviewDomains() []*types.DomainInfo {
	var result []*types.DomainInfo

	for _, domain := range types.DomainRegistry {
		if domain.IsPreview {
			result = append(result, domain)
		}
	}

	return result
}

// IsPreviewDomain checks if a specific domain is a preview feature
func IsPreviewDomain(domainName string) bool {
	domain, found := types.GetDomainInfo(domainName)
	if !found {
		return false
	}
	return domain.IsPreview
}

// GetNonPreviewDomains returns all domains that are not preview features
func GetNonPreviewDomains() []*types.DomainInfo {
	var result []*types.DomainInfo

	for _, domain := range types.DomainRegistry {
		if !domain.IsPreview {
			result = append(result, domain)
		}
	}

	return result
}

// GetPreviewDomainsInTier returns preview domains accessible at a specific tier
func GetPreviewDomainsInTier(tier string) []*types.DomainInfo {
	var result []*types.DomainInfo

	for _, domain := range types.DomainRegistry {
		if domain.IsPreview && ValidateTierAccess(tier, domain.RequiresTier) {
			result = append(result, domain)
		}
	}

	return result
}
