// Package subscription provides pre-execution validation against subscription capabilities.
package subscription

import (
	"context"
	"fmt"
	"strings"
)

// Validator provides pre-execution validation against subscription tier and addon availability.
type Validator struct {
	client   *Client
	registry *FeatureRegistry
}

// NewValidator creates a new subscription validator.
func NewValidator(client *Client) *Validator {
	return &Validator{
		client:   client,
		registry: NewFeatureRegistry(),
	}
}

// DetectAndCacheTier detects the current subscription tier and caches it.
// If the tier is already cached, returns the cached value.
// Returns the detected tier (Standard or Advanced).
func (v *Validator) DetectAndCacheTier(ctx context.Context) (string, error) {
	// Check if already cached
	cachedTier := GetCachedTier()
	if cachedTier != "" {
		return cachedTier, nil
	}

	// Detect from API
	tier, err := v.detectTierFromAPI(ctx)
	if err != nil {
		// Fall back to Standard if detection fails
		_ = SetCachedTierWithSource("Standard", "default")
		return "Standard", nil
	}

	// Cache the detected tier
	if err := SetCachedTierWithSource(tier, "api"); err != nil {
		return tier, fmt.Errorf("failed to cache tier: %w", err)
	}

	return tier, nil
}

// detectTierFromAPI queries the subscription API to determine the current tier.
func (v *Validator) detectTierFromAPI(ctx context.Context) (string, error) {
	if v.client == nil {
		return "Standard", fmt.Errorf("subscription client not available")
	}

	// Get subscription info which determines tier from plans and addons
	info, err := v.client.GetSubscriptionInfo(ctx)
	if err != nil {
		return "Standard", fmt.Errorf("failed to get subscription info: %w", err)
	}

	return info.Tier, nil
}

// GetCurrentTier returns the current subscription tier (cached or detected).
func (v *Validator) GetCurrentTier(ctx context.Context) (string, error) {
	return v.DetectAndCacheTier(ctx)
}

// GetRegistry returns the feature registry.
func (v *Validator) GetRegistry() *FeatureRegistry {
	return v.registry
}

// ResourceValidationResult contains the result of validating a resource type.
type ResourceValidationResult struct {
	ResourceType   string              `json:"resource_type" yaml:"resource_type"`
	IsAllowed      bool                `json:"is_allowed" yaml:"is_allowed"`
	CurrentTier    string              `json:"current_tier" yaml:"current_tier"`
	RequiredTier   string              `json:"required_tier,omitempty" yaml:"required_tier,omitempty"`
	RequiredAddons []string            `json:"required_addons,omitempty" yaml:"required_addons,omitempty"`
	MissingAddons  []string            `json:"missing_addons,omitempty" yaml:"missing_addons,omitempty"`
	ErrorMessage   string              `json:"error_message,omitempty" yaml:"error_message,omitempty"`
	HelpAnnotation string              `json:"help_annotation,omitempty" yaml:"help_annotation,omitempty"`
	Recommendation string              `json:"recommendation,omitempty" yaml:"recommendation,omitempty"`
	FeatureDetails *FeatureRequirement `json:"feature_details,omitempty" yaml:"feature_details,omitempty"`
}

// ValidateResourceAccess validates if a resource type can be created/modified
// with the current subscription tier and addon state.
func (v *Validator) ValidateResourceAccess(ctx context.Context, resourceType string) (*ResourceValidationResult, error) {
	result := &ResourceValidationResult{
		ResourceType: resourceType,
		IsAllowed:    true,
	}

	// Get current tier
	currentTier, err := v.GetCurrentTier(ctx)
	if err != nil {
		// If we can't determine tier, allow but warn
		result.CurrentTier = "Unknown"
		return result, nil
	}
	result.CurrentTier = currentTier

	// Check feature requirements for this resource type
	features := v.registry.GetFeaturesForResource(resourceType)
	if len(features) == 0 {
		// No restrictions, resource is allowed
		return result, nil
	}

	// Check each feature requirement
	for _, feature := range features {
		result.RequiredTier = feature.MinimumTier
		result.RequiredAddons = feature.RequiredAddons
		result.HelpAnnotation = feature.HelpAnnotation
		result.FeatureDetails = feature

		// Check tier requirement
		if !isTierSufficient(currentTier, feature.MinimumTier) {
			result.IsAllowed = false
			result.ErrorMessage = fmt.Sprintf(
				"Resource type '%s' requires %s tier (current: %s)",
				resourceType, feature.MinimumTier, currentTier,
			)
			result.Recommendation = fmt.Sprintf(
				"Upgrade to %s tier via F5 XC Console or contact F5 sales",
				feature.MinimumTier,
			)
			return result, nil
		}

		// Check addon requirements if we have a client
		if len(feature.RequiredAddons) > 0 && v.client != nil {
			addons, err := v.client.GetAddonServices(ctx, "system")
			if err == nil {
				missingAddons := v.checkMissingAddons(feature.RequiredAddons, addons)
				if len(missingAddons) > 0 {
					result.IsAllowed = false
					result.MissingAddons = missingAddons
					result.ErrorMessage = fmt.Sprintf(
						"Resource type '%s' requires addon service(s): %s",
						resourceType, strings.Join(missingAddons, ", "),
					)
					result.Recommendation = fmt.Sprintf(
						"Subscribe to %s addon via F5 XC Console",
						strings.Join(missingAddons, ", "),
					)
					return result, nil
				}
			}
		}
	}

	return result, nil
}

// checkMissingAddons returns the list of required addons that are not subscribed.
func (v *Validator) checkMissingAddons(required []string, available []AddonServiceInfo) []string {
	var missing []string

	for _, req := range required {
		found := false
		for _, addon := range available {
			if strings.EqualFold(addon.Name, req) && addon.IsActive() {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, req)
		}
	}

	return missing
}

// FeatureValidationResult contains the result of validating a feature.
type FeatureValidationResult struct {
	Feature        string              `json:"feature" yaml:"feature"`
	IsAvailable    bool                `json:"is_available" yaml:"is_available"`
	CurrentTier    string              `json:"current_tier" yaml:"current_tier"`
	RequiredTier   string              `json:"required_tier,omitempty" yaml:"required_tier,omitempty"`
	AddonState     string              `json:"addon_state,omitempty" yaml:"addon_state,omitempty"`
	AddonAccess    string              `json:"addon_access,omitempty" yaml:"addon_access,omitempty"`
	ErrorMessage   string              `json:"error_message,omitempty" yaml:"error_message,omitempty"`
	Recommendation string              `json:"recommendation,omitempty" yaml:"recommendation,omitempty"`
	FeatureDetails *FeatureRequirement `json:"feature_details,omitempty" yaml:"feature_details,omitempty"`
}

// ValidateFeatureAccess validates if a feature is available with the current subscription.
func (v *Validator) ValidateFeatureAccess(ctx context.Context, featureName string) (*FeatureValidationResult, error) {
	result := &FeatureValidationResult{
		Feature:     featureName,
		IsAvailable: true,
	}

	// Get current tier
	currentTier, err := v.GetCurrentTier(ctx)
	if err != nil {
		result.CurrentTier = "Unknown"
		return result, nil
	}
	result.CurrentTier = currentTier

	// Check feature registry
	feature := v.registry.GetFeature(featureName)
	if feature == nil {
		// Feature not in registry, check addons directly
		return v.validateAddonFeature(ctx, featureName, currentTier)
	}

	result.RequiredTier = feature.MinimumTier
	result.FeatureDetails = feature

	// Check tier requirement
	if !isTierSufficient(currentTier, feature.MinimumTier) {
		result.IsAvailable = false
		result.ErrorMessage = fmt.Sprintf(
			"Feature '%s' requires %s tier (current: %s)",
			feature.DisplayName, feature.MinimumTier, currentTier,
		)
		result.Recommendation = fmt.Sprintf(
			"Upgrade to %s tier to access %s",
			feature.MinimumTier, feature.DisplayName,
		)
		return result, nil
	}

	// Check addon requirements
	if len(feature.RequiredAddons) > 0 && v.client != nil {
		addons, err := v.client.GetAddonServices(ctx, "system")
		if err == nil {
			for _, reqAddon := range feature.RequiredAddons {
				for _, addon := range addons {
					if strings.EqualFold(addon.Name, reqAddon) {
						result.AddonState = addon.State
						result.AddonAccess = addon.AccessStatus

						if !addon.IsActive() {
							result.IsAvailable = false
							if addon.NeedsUpgrade() {
								result.ErrorMessage = fmt.Sprintf(
									"Feature '%s' requires a plan upgrade",
									feature.DisplayName,
								)
								result.Recommendation = "Upgrade your subscription plan via F5 XC Console"
							} else if addon.NeedsContactSales() {
								result.ErrorMessage = fmt.Sprintf(
									"Feature '%s' requires contacting F5 sales",
									feature.DisplayName,
								)
								result.Recommendation = "Contact F5 sales to enable this feature"
							} else if addon.IsAvailable() {
								result.ErrorMessage = fmt.Sprintf(
									"Feature '%s' addon is available but not subscribed",
									feature.DisplayName,
								)
								result.Recommendation = fmt.Sprintf(
									"Subscribe to %s addon via F5 XC Console",
									reqAddon,
								)
							} else {
								result.ErrorMessage = fmt.Sprintf(
									"Feature '%s' addon is not available",
									feature.DisplayName,
								)
							}
							return result, nil
						}
						break
					}
				}
			}
		}
	}

	return result, nil
}

// validateAddonFeature validates a feature by checking addon services directly.
func (v *Validator) validateAddonFeature(ctx context.Context, featureName, currentTier string) (*FeatureValidationResult, error) {
	result := &FeatureValidationResult{
		Feature:     featureName,
		CurrentTier: currentTier,
		IsAvailable: true,
	}

	if v.client == nil {
		return result, nil
	}

	addons, err := v.client.GetAddonServices(ctx, "system")
	if err != nil {
		return result, nil
	}

	// Find matching addon
	for _, addon := range addons {
		if strings.EqualFold(addon.Name, featureName) ||
			strings.Contains(strings.ToLower(addon.Name), strings.ToLower(featureName)) {

			result.AddonState = addon.State
			result.AddonAccess = addon.AccessStatus
			result.RequiredTier = addon.Tier

			if !addon.IsActive() {
				result.IsAvailable = false
				if addon.NeedsUpgrade() {
					result.ErrorMessage = fmt.Sprintf("Feature '%s' requires a plan upgrade", featureName)
					result.Recommendation = "Upgrade your subscription plan via F5 XC Console"
				} else if addon.NeedsContactSales() {
					result.ErrorMessage = fmt.Sprintf("Feature '%s' requires contacting F5 sales", featureName)
					result.Recommendation = "Contact F5 sales to enable this feature"
				} else if addon.IsAvailable() {
					result.ErrorMessage = fmt.Sprintf("Feature '%s' is available but not subscribed", featureName)
					result.Recommendation = fmt.Sprintf("Subscribe to %s addon via F5 XC Console", addon.Name)
				} else {
					result.ErrorMessage = fmt.Sprintf("Feature '%s' is not available (access: %s)", featureName, addon.AccessStatus)
				}
			}
			return result, nil
		}
	}

	// Feature not found in addons, assume available
	return result, nil
}

// ValidateFieldAccess validates if a specific field in a resource requires special subscription.
func (v *Validator) ValidateFieldAccess(ctx context.Context, fieldPath string) (*FeatureValidationResult, error) {
	result := &FeatureValidationResult{
		Feature:     fieldPath,
		IsAvailable: true,
	}

	currentTier, err := v.GetCurrentTier(ctx)
	if err != nil {
		result.CurrentTier = "Unknown"
		return result, nil
	}
	result.CurrentTier = currentTier

	// Check if this field requires specific features
	features := v.registry.GetFeaturesForField(fieldPath)
	if len(features) == 0 {
		return result, nil
	}

	for _, feature := range features {
		if !isTierSufficient(currentTier, feature.MinimumTier) {
			result.IsAvailable = false
			result.RequiredTier = feature.MinimumTier
			result.FeatureDetails = feature
			result.ErrorMessage = fmt.Sprintf(
				"Field '%s' requires %s tier (current: %s) - %s",
				fieldPath, feature.MinimumTier, currentTier, feature.DisplayName,
			)
			result.Recommendation = fmt.Sprintf(
				"Remove '%s' from configuration or upgrade to %s tier",
				fieldPath, feature.MinimumTier,
			)
			return result, nil
		}
	}

	return result, nil
}

// SubscriptionContextSpec represents subscription context for spec output
type SubscriptionContextSpec struct {
	CurrentTier        string                  `json:"current_tier" yaml:"current_tier"`
	TierSource         string                  `json:"tier_source" yaml:"tier_source"`
	AvailableFeatures  []string                `json:"available_features" yaml:"available_features"`
	RestrictedFeatures []RestrictedFeatureSpec `json:"restricted_features" yaml:"restricted_features"`
}

// RestrictedFeatureSpec describes a feature not available in current tier
type RestrictedFeatureSpec struct {
	Feature        string   `json:"feature" yaml:"feature"`
	DisplayName    string   `json:"display_name" yaml:"display_name"`
	RequiredTier   string   `json:"required_tier" yaml:"required_tier"`
	RequiredAddons []string `json:"required_addons,omitempty" yaml:"required_addons,omitempty"`
	HelpAnnotation string   `json:"help_annotation" yaml:"help_annotation"`
}

// GenerateSubscriptionContext generates the subscription context for spec output.
func (v *Validator) GenerateSubscriptionContext(ctx context.Context) (*SubscriptionContextSpec, error) {
	cacheInfo := GetTierCacheInfo()

	// Determine tier source description
	tierSource := "Unknown"
	switch {
	case cacheInfo.IsFromEnv:
		tierSource = "F5XC_SUBSCRIPTION_TIER environment variable"
	case cacheInfo.IsFromAPI:
		tierSource = "Detected from F5 XC API"
	case cacheInfo.IsDefault:
		tierSource = "Default (API detection unavailable)"
	}

	currentTier := cacheInfo.Tier
	if currentTier == "" {
		// Try to detect
		if tier, err := v.GetCurrentTier(ctx); err == nil {
			currentTier = tier
			cacheInfo = GetTierCacheInfo()
			if cacheInfo.IsFromAPI {
				tierSource = "Detected from F5 XC API"
			}
		} else {
			currentTier = "Standard"
			tierSource = "Default (API detection unavailable)"
		}
	}

	// Build available and restricted feature lists
	var available []string
	var restricted []RestrictedFeatureSpec

	for _, feature := range v.registry.GetAllFeatures() {
		if isTierSufficient(currentTier, feature.MinimumTier) {
			available = append(available, feature.FeatureName)
		} else {
			restricted = append(restricted, RestrictedFeatureSpec{
				Feature:        feature.FeatureName,
				DisplayName:    feature.DisplayName,
				RequiredTier:   feature.MinimumTier,
				RequiredAddons: feature.RequiredAddons,
				HelpAnnotation: feature.HelpAnnotation,
			})
		}
	}

	return &SubscriptionContextSpec{
		CurrentTier:        currentTier,
		TierSource:         tierSource,
		AvailableFeatures:  available,
		RestrictedFeatures: restricted,
	}, nil
}

// HelpAnnotation returns the help annotation for a resource type based on current tier.
func (v *Validator) HelpAnnotation(resourceType string) string {
	currentTier := GetCachedTier()
	if currentTier == "" {
		currentTier = "Standard" // Default to Standard for help display
	}

	return v.registry.GetResourceRestriction(resourceType, currentTier)
}

// IsResourceAllowedQuick performs a quick check without API calls using cached tier.
func (v *Validator) IsResourceAllowedQuick(resourceType string) bool {
	currentTier := GetCachedTier()
	if currentTier == "" {
		return true // Allow if tier unknown
	}

	return v.registry.IsResourceAvailable(resourceType, currentTier)
}
