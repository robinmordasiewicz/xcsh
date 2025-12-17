// Package subscription provides subscription tier feature mapping and validation.
package subscription

// FeatureRequirement defines what subscription tier and addons are needed for a feature
type FeatureRequirement struct {
	FeatureName    string   // Unique identifier for the feature
	DisplayName    string   // Human-readable name
	MinimumTier    string   // Minimum subscription tier (STANDARD, ADVANCED)
	RequiredAddons []string // Addon services that must be subscribed
	AddonTier      string   // For addons with sub-tiers (e.g., "standard", "advanced")
	Description    string   // Description of what this feature enables
	HelpAnnotation string   // Annotation shown in --help output (e.g., "[Requires Advanced]")
	ResourceTypes  []string // Resource types that require this feature
	FieldPatterns  []string // Field paths in resources that require this feature
}

// FeatureRegistry maintains the mapping of features to subscription requirements
type FeatureRegistry struct {
	features    map[string]*FeatureRequirement
	resourceMap map[string][]*FeatureRequirement // Resource type -> required features
	fieldMap    map[string][]*FeatureRequirement // Field pattern -> required features
}

// NewFeatureRegistry creates a new feature registry with default mappings
func NewFeatureRegistry() *FeatureRegistry {
	registry := &FeatureRegistry{
		features:    make(map[string]*FeatureRequirement),
		resourceMap: make(map[string][]*FeatureRequirement),
		fieldMap:    make(map[string][]*FeatureRequirement),
	}

	// Register all default features
	for _, feature := range defaultFeatures {
		registry.Register(feature)
	}

	return registry
}

// Register adds a feature requirement to the registry
func (r *FeatureRegistry) Register(req *FeatureRequirement) {
	r.features[req.FeatureName] = req

	// Index by resource type
	for _, rt := range req.ResourceTypes {
		r.resourceMap[rt] = append(r.resourceMap[rt], req)
	}

	// Index by field pattern
	for _, fp := range req.FieldPatterns {
		r.fieldMap[fp] = append(r.fieldMap[fp], req)
	}
}

// GetFeature returns the requirement for a specific feature
func (r *FeatureRegistry) GetFeature(name string) *FeatureRequirement {
	return r.features[name]
}

// GetFeaturesForResource returns all features required for a resource type
func (r *FeatureRegistry) GetFeaturesForResource(resourceType string) []*FeatureRequirement {
	return r.resourceMap[resourceType]
}

// GetFeaturesForField returns all features required for a field pattern
func (r *FeatureRegistry) GetFeaturesForField(fieldPath string) []*FeatureRequirement {
	return r.fieldMap[fieldPath]
}

// GetAllFeatures returns all registered features
func (r *FeatureRegistry) GetAllFeatures() []*FeatureRequirement {
	features := make([]*FeatureRequirement, 0, len(r.features))
	for _, f := range r.features {
		features = append(features, f)
	}
	return features
}

// GetAdvancedFeatures returns all features requiring Advanced tier
func (r *FeatureRegistry) GetAdvancedFeatures() []*FeatureRequirement {
	var advanced []*FeatureRequirement
	for _, f := range r.features {
		if f.MinimumTier == TierAdvanced {
			advanced = append(advanced, f)
		}
	}
	return advanced
}

// GetStandardFeatures returns all features available in Standard tier
func (r *FeatureRegistry) GetStandardFeatures() []*FeatureRequirement {
	var standard []*FeatureRequirement
	for _, f := range r.features {
		if f.MinimumTier == TierStandard || f.MinimumTier == "" {
			standard = append(standard, f)
		}
	}
	return standard
}

// IsResourceAvailable checks if a resource type is available for the given tier
func (r *FeatureRegistry) IsResourceAvailable(resourceType, currentTier string) bool {
	features := r.GetFeaturesForResource(resourceType)
	if len(features) == 0 {
		return true // No restrictions
	}

	for _, f := range features {
		if !isTierSufficient(currentTier, f.MinimumTier) {
			return false
		}
	}
	return true
}

// GetResourceRestriction returns the restriction message for a resource type
func (r *FeatureRegistry) GetResourceRestriction(resourceType, currentTier string) string {
	features := r.GetFeaturesForResource(resourceType)
	if len(features) == 0 {
		return ""
	}

	for _, f := range features {
		if !isTierSufficient(currentTier, f.MinimumTier) {
			return f.HelpAnnotation
		}
	}
	return ""
}

// isTierSufficient checks if currentTier meets or exceeds requiredTier
// Note: Only Standard and Advanced are active tiers.
// Basic maps to Standard, Premium maps to Advanced.
func isTierSufficient(currentTier, requiredTier string) bool {
	tierOrder := map[string]int{
		"":           0,
		TierNoTier:   0,
		TierBasic:    1, // Discontinued, maps to Standard
		TierStandard: 1,
		"Standard":   1,
		"Basic":      1,
		TierAdvanced: 2,
		TierPremium:  2, // Discontinued, maps to Advanced
		"Advanced":   2,
		"Premium":    2,
	}

	currentOrder := tierOrder[currentTier]
	requiredOrder := tierOrder[requiredTier]

	return currentOrder >= requiredOrder
}

// defaultFeatures defines the hardcoded feature-to-tier mappings
// This is enriched at runtime with live addon service data
var defaultFeatures = []*FeatureRequirement{
	// ============================================================
	// STANDARD TIER FEATURES
	// ============================================================

	// Core Load Balancing - Standard
	{
		FeatureName:    "http-loadbalancer",
		DisplayName:    "HTTP Load Balancer",
		MinimumTier:    TierStandard,
		RequiredAddons: nil,
		Description:    "HTTP and HTTPS load balancing with WAF capabilities",
		HelpAnnotation: "",
		ResourceTypes:  []string{"http_loadbalancer"},
	},
	{
		FeatureName:    "tcp-loadbalancer",
		DisplayName:    "TCP Load Balancer",
		MinimumTier:    TierStandard,
		RequiredAddons: nil,
		Description:    "TCP load balancing for non-HTTP protocols",
		HelpAnnotation: "",
		ResourceTypes:  []string{"tcp_loadbalancer"},
	},
	{
		FeatureName:    "origin-pool",
		DisplayName:    "Origin Pool",
		MinimumTier:    TierStandard,
		RequiredAddons: nil,
		Description:    "Backend server pools for load balancing",
		HelpAnnotation: "",
		ResourceTypes:  []string{"origin_pool"},
	},
	{
		FeatureName:    "healthcheck",
		DisplayName:    "Health Check",
		MinimumTier:    TierStandard,
		RequiredAddons: nil,
		Description:    "Health monitoring for origin servers",
		HelpAnnotation: "",
		ResourceTypes:  []string{"healthcheck"},
	},

	// WAF/App Firewall - Standard
	{
		FeatureName:    "app-firewall",
		DisplayName:    "Application Firewall",
		MinimumTier:    TierStandard,
		RequiredAddons: nil,
		Description:    "Web Application Firewall protection",
		HelpAnnotation: "",
		ResourceTypes:  []string{"app_firewall"},
	},
	{
		FeatureName:    "service-policy",
		DisplayName:    "Service Policy",
		MinimumTier:    TierStandard,
		RequiredAddons: nil,
		Description:    "Traffic management and access control policies",
		HelpAnnotation: "",
		ResourceTypes:  []string{"service_policy"},
	},

	// DNS - Standard
	{
		FeatureName:    "dns-zone",
		DisplayName:    "DNS Zone",
		MinimumTier:    TierStandard,
		RequiredAddons: nil,
		Description:    "DNS zone management",
		HelpAnnotation: "",
		ResourceTypes:  []string{"dns_zone"},
	},

	// Bot Defense Standard - Included with Standard subscription
	{
		FeatureName:    "bot-defense-standard",
		DisplayName:    "Bot Defense Standard",
		MinimumTier:    TierStandard,
		RequiredAddons: []string{"bot-defense"},
		AddonTier:      "standard",
		Description:    "Basic bot detection and mitigation capabilities",
		HelpAnnotation: "[Requires Bot Defense addon]",
		ResourceTypes:  []string{"bot_defense_config"},
		FieldPatterns:  []string{"spec.bot_defense"},
	},

	// ============================================================
	// ADVANCED TIER FEATURES
	// ============================================================

	// Bot Defense Advanced - Requires Advanced subscription
	{
		FeatureName:    "bot-defense-advanced",
		DisplayName:    "Bot Defense Advanced",
		MinimumTier:    TierAdvanced,
		RequiredAddons: []string{"bot-defense"},
		AddonTier:      "advanced",
		Description:    "Advanced bot protection with behavioral analysis and client-side defense",
		HelpAnnotation: "[Requires Advanced]",
		ResourceTypes:  []string{"bot_defense_advanced_policy"},
		FieldPatterns:  []string{"spec.bot_defense_advanced", "spec.client_side_defense"},
	},

	// Client-Side Defense - Advanced
	{
		FeatureName:    "client-side-defense",
		DisplayName:    "Client-Side Defense",
		MinimumTier:    TierAdvanced,
		RequiredAddons: []string{"client-side-defense"},
		Description:    "JavaScript-based client-side protection against Magecart and formjacking",
		HelpAnnotation: "[Requires Advanced]",
		ResourceTypes:  []string{"client_side_defense_policy"},
		FieldPatterns:  []string{"spec.client_side_defense"},
	},

	// API Security - Advanced only
	{
		FeatureName:    "api-security",
		DisplayName:    "API Security",
		MinimumTier:    TierAdvanced,
		RequiredAddons: []string{"api-security"},
		Description:    "API discovery, schema validation, and security enforcement",
		HelpAnnotation: "[Requires Advanced]",
		ResourceTypes:  []string{"api_definition", "api_group", "api_rate_limiter"},
		FieldPatterns:  []string{"spec.api_definition", "spec.api_specification"},
	},

	// API Discovery - Advanced only
	{
		FeatureName:    "api-discovery",
		DisplayName:    "API Discovery",
		MinimumTier:    TierAdvanced,
		RequiredAddons: []string{"api-security"},
		Description:    "Automatic API endpoint discovery and documentation",
		HelpAnnotation: "[Requires Advanced]",
		ResourceTypes:  []string{"api_discovery_config"},
		FieldPatterns:  []string{"spec.enable_api_discovery"},
	},

	// Malicious User Detection - Advanced only
	{
		FeatureName:    "malicious-user",
		DisplayName:    "Malicious User Detection",
		MinimumTier:    TierAdvanced,
		RequiredAddons: []string{"malicious-user"},
		Description:    "Detect and mitigate malicious user behavior patterns",
		HelpAnnotation: "[Requires Advanced]",
		ResourceTypes:  []string{"malicious_user_mitigation"},
		FieldPatterns:  []string{"spec.malicious_user_mitigation"},
	},

	// DDoS Protection - Advanced features
	{
		FeatureName:    "ddos-protection-advanced",
		DisplayName:    "DDoS Protection Advanced",
		MinimumTier:    TierAdvanced,
		RequiredAddons: []string{"ddos-protection"},
		AddonTier:      "advanced",
		Description:    "Advanced DDoS protection with auto-mitigation",
		HelpAnnotation: "[Requires Advanced]",
		ResourceTypes:  []string{"ddos_protection_policy"},
		FieldPatterns:  []string{"spec.ddos_mitigation"},
	},

	// Network Firewall - Advanced
	{
		FeatureName:    "network-firewall",
		DisplayName:    "Network Firewall",
		MinimumTier:    TierAdvanced,
		RequiredAddons: []string{"network-firewall"},
		Description:    "Layer 3-4 network firewall with stateful inspection",
		HelpAnnotation: "[Requires Advanced]",
		ResourceTypes:  []string{"network_firewall", "forward_proxy_policy"},
	},

	// Synthetic Monitoring - Advanced
	{
		FeatureName:    "synthetic-monitoring",
		DisplayName:    "Synthetic Monitoring",
		MinimumTier:    TierAdvanced,
		RequiredAddons: []string{"synthetic-monitoring"},
		Description:    "Proactive endpoint monitoring and alerting",
		HelpAnnotation: "[Requires Advanced]",
		ResourceTypes:  []string{"synthetic_monitor"},
	},
}

// FeatureTierMapSpec represents the tier mapping for spec output
type FeatureTierMapSpec struct {
	Description string                      `json:"description" yaml:"description"`
	Tiers       []TierSpec                  `json:"tiers" yaml:"tiers"`
	Resources   map[string]ResourceTierSpec `json:"resources" yaml:"resources"`
	AIGuidance  []string                    `json:"ai_guidance" yaml:"ai_guidance"`
}

// TierSpec describes a subscription tier
type TierSpec struct {
	Name        string   `json:"name" yaml:"name"`
	DisplayName string   `json:"display_name" yaml:"display_name"`
	Description string   `json:"description" yaml:"description"`
	Features    []string `json:"features" yaml:"features"`
}

// ResourceTierSpec describes tier requirements for a resource type
type ResourceTierSpec struct {
	RequiredTier   string   `json:"required_tier" yaml:"required_tier"`
	RequiredAddons []string `json:"required_addons,omitempty" yaml:"required_addons,omitempty"`
	HelpAnnotation string   `json:"help_annotation,omitempty" yaml:"help_annotation,omitempty"`
}

// GenerateFeatureTierMapSpec generates the feature tier map for spec output
func GenerateFeatureTierMapSpec() FeatureTierMapSpec {
	registry := NewFeatureRegistry()

	// Build resource map
	resources := make(map[string]ResourceTierSpec)
	for _, feature := range registry.GetAllFeatures() {
		for _, rt := range feature.ResourceTypes {
			resources[rt] = ResourceTierSpec{
				RequiredTier:   feature.MinimumTier,
				RequiredAddons: feature.RequiredAddons,
				HelpAnnotation: feature.HelpAnnotation,
			}
		}
	}

	// Collect Standard and Advanced features
	var standardFeatures, advancedFeatures []string
	for _, f := range registry.GetStandardFeatures() {
		standardFeatures = append(standardFeatures, f.FeatureName)
	}
	for _, f := range registry.GetAdvancedFeatures() {
		advancedFeatures = append(advancedFeatures, f.FeatureName)
	}

	return FeatureTierMapSpec{
		Description: "F5 XC subscription tiers determine feature and resource availability. Standard tier includes core load balancing, WAF, and basic bot protection. Advanced tier adds API security, advanced bot defense, malicious user detection, and more.",
		Tiers: []TierSpec{
			{
				Name:        TierStandard,
				DisplayName: "Standard",
				Description: "Base subscription tier with core F5 XC functionality including load balancing, WAF, and basic bot protection",
				Features:    standardFeatures,
			},
			{
				Name:        TierAdvanced,
				DisplayName: "Advanced",
				Description: "Enhanced tier with API security, advanced bot defense, malicious user detection, and premium features",
				Features:    advancedFeatures,
			},
		},
		Resources: resources,
		AIGuidance: []string{
			"Check subscription_context.current_tier before recommending tier-restricted features",
			"Use 'f5xcctl subscription validate --feature <name>' to verify feature availability before deployment",
			"Bot Defense has two tiers: bot-defense-standard (Standard subscription) and bot-defense-advanced (Advanced subscription)",
			"Resources marked with [Requires Advanced] will fail to deploy on Standard tier subscriptions",
			"Use 'f5xcctl subscription show --output-format json' to discover current tenant capabilities",
			"Pre-validate resource creation with 'f5xcctl subscription validate --resource-type <type>'",
		},
	}
}
