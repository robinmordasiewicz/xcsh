// Package subscription provides types and functions for interacting with F5 XC subscription,
// addon services, and quota management APIs. This package enables AI assistants and Terraform
// users to validate deployment feasibility before execution.
package subscription

// Addon service tier values - determines feature capabilities
// Note: Only Standard and Advanced are active subscription tiers.
// Basic and Premium are discontinued but kept for API compatibility.
const (
	TierNoTier   = "NO_TIER"
	TierBasic    = "BASIC"    // Discontinued - maps to Standard
	TierStandard = "STANDARD" // Active tier
	TierAdvanced = "ADVANCED" // Active tier
	TierPremium  = "PREMIUM"  // Discontinued - maps to Advanced
)

// Addon service state values - subscription status
const (
	StateNone       = "AS_NONE"       // Not subscribed
	StatePending    = "AS_PENDING"    // Subscription pending
	StateSubscribed = "AS_SUBSCRIBED" // Actively subscribed
	StateError      = "AS_ERROR"      // Subscription error
)

// Addon service access status values - availability based on plan
const (
	AccessAllowed         = "AS_AC_ALLOWED"                 // Can subscribe
	AccessDenied          = "AS_AC_PBAC_DENY"               // Access denied by policy
	AccessUpgradeRequired = "AS_AC_PBAC_DENY_UPGRADE_PLAN"  // Requires plan upgrade
	AccessContactSales    = "AS_AC_PBAC_DENY_CONTACT_SALES" // Requires sales contact
	AccessInternalService = "AS_AC_PBAC_DENY_INTERNAL_SVC"  // Internal service only
	AccessUnknown         = "AS_AC_UNKNOWN"                 // Unknown status
)

// Activation type values - how the addon is managed
const (
	ActivationSelf             = "self"              // User can self-activate
	ActivationPartiallyManaged = "partially_managed" // Some features managed
	ActivationManaged          = "managed"           // Fully managed by F5
)

// Subscription state values - lifecycle of addon subscriptions
const (
	SubscriptionPending        = "SUBSCRIPTION_PENDING"         // Request submitted, awaiting activation
	SubscriptionEnabled        = "SUBSCRIPTION_ENABLED"         // Addon is active and enabled
	SubscriptionDisablePending = "SUBSCRIPTION_DISABLE_PENDING" // Deactivation in progress
	SubscriptionDisabled       = "SUBSCRIPTION_DISABLED"        // Addon is disabled
)

// Access status for end-of-life addons
const (
	AccessEOL = "AS_AC_EOL" // End of life, cannot activate
)

// Validation result status
const (
	ValidationPass    = "PASS"
	ValidationFail    = "FAIL"
	ValidationWarning = "WARNING"
)

// Exit codes for subscription operations
const (
	ExitCodeValid   = 0 // Validation passed
	ExitCodeInvalid = 2 // Validation failed
)

// SubscriptionInfo represents the complete subscription state for a tenant
type SubscriptionInfo struct {
	Tier            string             `json:"tier" yaml:"tier"`
	TenantName      string             `json:"tenant_name" yaml:"tenant_name"`
	Plan            PlanInfo           `json:"plan" yaml:"plan"`
	ActiveAddons    []AddonServiceInfo `json:"active_addons" yaml:"active_addons"`
	AvailableAddons []AddonServiceInfo `json:"available_addons" yaml:"available_addons"`
	QuotaSummary    QuotaSummary       `json:"quota_summary" yaml:"quota_summary"`
}

// PlanInfo represents the subscription plan details
type PlanInfo struct {
	Name             string   `json:"name" yaml:"name"`
	DisplayName      string   `json:"display_name" yaml:"display_name"`
	Description      string   `json:"description,omitempty" yaml:"description,omitempty"`
	IncludedServices []string `json:"included_services,omitempty" yaml:"included_services,omitempty"`
	AllowedServices  []string `json:"allowed_services,omitempty" yaml:"allowed_services,omitempty"`
}

// AddonServiceInfo represents an addon service and its current status
type AddonServiceInfo struct {
	Name           string `json:"name" yaml:"name"`
	DisplayName    string `json:"display_name" yaml:"display_name"`
	Description    string `json:"description,omitempty" yaml:"description,omitempty"`
	Tier           string `json:"tier" yaml:"tier"`
	State          string `json:"state" yaml:"state"`
	AccessStatus   string `json:"access_status" yaml:"access_status"`
	ActivationType string `json:"activation_type,omitempty" yaml:"activation_type,omitempty"`
	Namespace      string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

// IsActive returns true if the addon service is actively subscribed
func (a *AddonServiceInfo) IsActive() bool {
	return a.State == StateSubscribed
}

// IsAvailable returns true if the addon service can be subscribed to
func (a *AddonServiceInfo) IsAvailable() bool {
	return a.AccessStatus == AccessAllowed && a.State != StateSubscribed
}

// IsDenied returns true if access to the addon service is denied
func (a *AddonServiceInfo) IsDenied() bool {
	return a.AccessStatus == AccessDenied ||
		a.AccessStatus == AccessUpgradeRequired ||
		a.AccessStatus == AccessContactSales ||
		a.AccessStatus == AccessInternalService
}

// NeedsUpgrade returns true if the addon requires a plan upgrade
func (a *AddonServiceInfo) NeedsUpgrade() bool {
	return a.AccessStatus == AccessUpgradeRequired
}

// NeedsContactSales returns true if the addon requires contacting sales
func (a *AddonServiceInfo) NeedsContactSales() bool {
	return a.AccessStatus == AccessContactSales
}

// IsSelfActivation returns true if the addon can be self-activated
func (a *AddonServiceInfo) IsSelfActivation() bool {
	return a.ActivationType == ActivationSelf
}

// IsManagedActivation returns true if the addon requires backend intervention
func (a *AddonServiceInfo) IsManagedActivation() bool {
	return a.ActivationType == ActivationPartiallyManaged || a.ActivationType == ActivationManaged
}

// IsEndOfLife returns true if the addon is end-of-life
func (a *AddonServiceInfo) IsEndOfLife() bool {
	return a.AccessStatus == AccessEOL
}

// CanActivate returns true if the addon can be activated (allowed and not already active)
func (a *AddonServiceInfo) CanActivate() bool {
	return a.AccessStatus == AccessAllowed &&
		a.State != StateSubscribed &&
		!a.IsEndOfLife()
}

// IsPending returns true if the addon has a pending activation
func (a *AddonServiceInfo) IsPending() bool {
	return a.State == StatePending
}

// QuotaSummary provides an overview of quota usage
type QuotaSummary struct {
	TotalLimits    int         `json:"total_limits" yaml:"total_limits"`
	LimitsAtRisk   int         `json:"limits_at_risk" yaml:"limits_at_risk"`
	LimitsExceeded int         `json:"limits_exceeded" yaml:"limits_exceeded"`
	Objects        []QuotaItem `json:"objects,omitempty" yaml:"objects,omitempty"`
	Resources      []QuotaItem `json:"resources,omitempty" yaml:"resources,omitempty"`
}

// QuotaUsageInfo represents detailed quota limits and current usage
type QuotaUsageInfo struct {
	Namespace string      `json:"namespace" yaml:"namespace"`
	Objects   []QuotaItem `json:"objects" yaml:"objects"`
	Resources []QuotaItem `json:"resources" yaml:"resources"`
	APIs      []RateLimit `json:"apis,omitempty" yaml:"apis,omitempty"`
}

// QuotaItem represents a single quota with limit and usage
type QuotaItem struct {
	Name        string  `json:"name" yaml:"name"`
	DisplayName string  `json:"display_name" yaml:"display_name"`
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
	ObjectType  string  `json:"object_type,omitempty" yaml:"object_type,omitempty"`
	Limit       float64 `json:"limit" yaml:"limit"`
	Usage       float64 `json:"usage" yaml:"usage"`
	Percentage  float64 `json:"percentage" yaml:"percentage"`
	Status      string  `json:"status" yaml:"status"` // OK, WARNING, EXCEEDED
}

// IsExceeded returns true if usage equals or exceeds the limit
func (q *QuotaItem) IsExceeded() bool {
	return q.Usage >= q.Limit
}

// IsAtRisk returns true if usage is above 80% of the limit
func (q *QuotaItem) IsAtRisk() bool {
	return q.Percentage >= 80.0 && q.Percentage < 100.0
}

// RemainingCapacity returns how much capacity remains
func (q *QuotaItem) RemainingCapacity() float64 {
	remaining := q.Limit - q.Usage
	if remaining < 0 {
		return 0
	}
	return remaining
}

// RateLimit represents API rate limiting configuration
type RateLimit struct {
	Name  string `json:"name" yaml:"name"`
	Rate  int    `json:"rate" yaml:"rate"`
	Burst int    `json:"burst" yaml:"burst"`
	Unit  string `json:"unit" yaml:"unit"` // per-minute, per-second, etc.
}

// ValidationRequest represents a request to validate deployment feasibility
type ValidationRequest struct {
	ResourceType  string `json:"resource_type,omitempty" yaml:"resource_type,omitempty"`
	Count         int    `json:"count,omitempty" yaml:"count,omitempty"`
	Feature       string `json:"feature,omitempty" yaml:"feature,omitempty"`
	TerraformPlan string `json:"terraform_plan,omitempty" yaml:"terraform_plan,omitempty"`
	Namespace     string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

// ValidationResult represents the result of a deployment validation
type ValidationResult struct {
	Valid    bool              `json:"valid" yaml:"valid"`
	Checks   []ValidationCheck `json:"checks" yaml:"checks"`
	Warnings []string          `json:"warnings,omitempty" yaml:"warnings,omitempty"`
	Errors   []string          `json:"errors,omitempty" yaml:"errors,omitempty"`
}

// ValidationCheck represents a single validation check
type ValidationCheck struct {
	Type         string `json:"type" yaml:"type"`                               // quota, feature, addon
	Resource     string `json:"resource,omitempty" yaml:"resource,omitempty"`   // Resource type being checked
	Feature      string `json:"feature,omitempty" yaml:"feature,omitempty"`     // Feature being checked
	Current      int    `json:"current,omitempty" yaml:"current,omitempty"`     // Current usage
	Requested    int    `json:"requested,omitempty" yaml:"requested,omitempty"` // Requested count
	Limit        int    `json:"limit,omitempty" yaml:"limit,omitempty"`         // Quota limit
	RequiredTier string `json:"required_tier,omitempty" yaml:"required_tier,omitempty"`
	CurrentTier  string `json:"current_tier,omitempty" yaml:"current_tier,omitempty"`
	Status       string `json:"status,omitempty" yaml:"status,omitempty"` // Addon subscription status
	Result       string `json:"result" yaml:"result"`                     // PASS, FAIL, WARNING
	Message      string `json:"message,omitempty" yaml:"message,omitempty"`
}

// IsPassed returns true if the check passed
func (v *ValidationCheck) IsPassed() bool {
	return v.Result == ValidationPass
}

// IsFailed returns true if the check failed
func (v *ValidationCheck) IsFailed() bool {
	return v.Result == ValidationFail
}

// IsWarning returns true if the check has a warning
func (v *ValidationCheck) IsWarning() bool {
	return v.Result == ValidationWarning
}

// AddCheck adds a validation check to the result
func (v *ValidationResult) AddCheck(check ValidationCheck) {
	v.Checks = append(v.Checks, check)
	if check.IsFailed() {
		v.Valid = false
		if check.Message != "" {
			v.Errors = append(v.Errors, check.Message)
		}
	} else if check.IsWarning() && check.Message != "" {
		v.Warnings = append(v.Warnings, check.Message)
	}
}

// PassedCount returns the number of passed checks
func (v *ValidationResult) PassedCount() int {
	count := 0
	for _, check := range v.Checks {
		if check.IsPassed() {
			count++
		}
	}
	return count
}

// FailedCount returns the number of failed checks
func (v *ValidationResult) FailedCount() int {
	count := 0
	for _, check := range v.Checks {
		if check.IsFailed() {
			count++
		}
	}
	return count
}

// ActivationResponse represents the result of an addon activation attempt
type ActivationResponse struct {
	AddonService      string `json:"addon_service" yaml:"addon_service"`
	Namespace         string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	SubscriptionState string `json:"subscription_state,omitempty" yaml:"subscription_state,omitempty"`
	ActivationType    string `json:"activation_type,omitempty" yaml:"activation_type,omitempty"`
	AccessStatus      string `json:"access_status,omitempty" yaml:"access_status,omitempty"`
	RequestID         string `json:"request_id,omitempty" yaml:"request_id,omitempty"`
	Message           string `json:"message" yaml:"message"`
	NextSteps         string `json:"next_steps,omitempty" yaml:"next_steps,omitempty"`
	IsPending         bool   `json:"is_pending" yaml:"is_pending"`
	IsImmediate       bool   `json:"is_immediate" yaml:"is_immediate"`
}

// PendingActivation represents a pending addon activation request
type PendingActivation struct {
	AddonService      string `json:"addon_service" yaml:"addon_service"`
	Namespace         string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	SubscriptionState string `json:"subscription_state" yaml:"subscription_state"`
	ActivationType    string `json:"activation_type,omitempty" yaml:"activation_type,omitempty"`
	Message           string `json:"message,omitempty" yaml:"message,omitempty"`
}

// ActivationStatusResult contains pending activations and active addons
type ActivationStatusResult struct {
	PendingActivations []PendingActivation `json:"pending_activations" yaml:"pending_activations"`
	ActiveAddons       []string            `json:"active_addons" yaml:"active_addons"`
	TotalPending       int                 `json:"total_pending" yaml:"total_pending"`
}

// Spec represents the subscription command specification for AI assistants
type Spec struct {
	CommandGroup      string         `json:"command_group" yaml:"command_group"`
	Description       string         `json:"description" yaml:"description"`
	Discovery         DiscoverySpec  `json:"discovery" yaml:"discovery"`
	ValidationCommand string         `json:"validation_command" yaml:"validation_command"`
	AddonTiers        []string       `json:"addon_tiers" yaml:"addon_tiers"`
	AddonStates       []string       `json:"addon_states" yaml:"addon_states"`
	AccessStatuses    []string       `json:"access_statuses" yaml:"access_statuses"`
	QuotaTypes        []string       `json:"quota_types" yaml:"quota_types"`
	AIHints           []string       `json:"ai_hints" yaml:"ai_hints"`
	ExitCodes         []ExitCodeSpec `json:"exit_codes" yaml:"exit_codes"`
	Workflows         []WorkflowSpec `json:"workflows" yaml:"workflows"`
}

// DiscoverySpec describes how AI assistants can discover subscription information
type DiscoverySpec struct {
	Commands    []string `json:"commands" yaml:"commands"`
	Description string   `json:"description" yaml:"description"`
}

// ExitCodeSpec describes an exit code
type ExitCodeSpec struct {
	Code        int    `json:"code" yaml:"code"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

// WorkflowSpec describes a subscription-related workflow
type WorkflowSpec struct {
	Name        string         `json:"name" yaml:"name"`
	Description string         `json:"description" yaml:"description"`
	Steps       []WorkflowStep `json:"steps" yaml:"steps"`
}

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	Step        int    `json:"step" yaml:"step"`
	Description string `json:"description" yaml:"description"`
	Command     string `json:"command" yaml:"command"`
}

// GenerateSpec generates the subscription command specification for AI assistants
func GenerateSpec() *Spec {
	return &Spec{
		CommandGroup: "subscription",
		Description:  "Subscription management for F5 Distributed Cloud tenant capabilities, addon services, and quota limits.",
		Discovery: DiscoverySpec{
			Commands: []string{
				"xcsh subscription show --output-format json",
				"xcsh subscription addons --output-format json",
				"xcsh subscription quota --output-format json",
				"xcsh subscription activation-status --output-format json",
			},
			Description: "Use these commands to discover tenant subscription capabilities before deployment.",
		},
		ValidationCommand: "xcsh subscription validate",
		AddonTiers:        []string{TierStandard, TierAdvanced}, // Only active tiers
		AddonStates:       []string{StateNone, StatePending, StateSubscribed, StateError},
		AccessStatuses:    []string{AccessAllowed, AccessDenied, AccessUpgradeRequired, AccessContactSales, AccessInternalService},
		QuotaTypes:        []string{"objects", "resources", "apis"},
		AIHints: []string{
			"Before deploying resources, check subscription tier with 'xcsh subscription show'",
			"Validate features like bot-defense are subscribed before using them in Terraform",
			"Check quota limits before creating multiple resources of the same type",
			"Use 'xcsh subscription validate' to pre-flight check Terraform plans",
			"Exit code 2 from validate indicates deployment will fail due to quota/feature issues",
			"Addon services with AccessUpgradeRequired need a plan upgrade to enable",
			"Addon services with AccessContactSales require contacting F5 sales",
			"Use --namespace flag to check quotas in specific namespaces",
			"Use 'xcsh subscription activate --addon <name>' to activate available addons",
			"Exit code 9 from activate indicates feature not available (denied, upgrade required, contact sales)",
			"Use 'xcsh subscription activation-status' to check pending activation requests",
			"Self-activation addons activate immediately; managed addons require SRE approval",
		},
		ExitCodes: []ExitCodeSpec{
			{Code: 0, Name: "ExitSuccess", Description: "All validations passed"},
			{Code: 1, Name: "ExitGenericError", Description: "Generic error occurred"},
			{Code: 2, Name: "ExitValidationFailed", Description: "Validation failed - quota exceeded or feature unavailable"},
			{Code: 8, Name: "ExitQuotaExceeded", Description: "Quota would be exceeded by the operation"},
			{Code: 9, Name: "ExitFeatureNotAvailable", Description: "Feature not available in current subscription"},
		},
		Workflows: []WorkflowSpec{
			{
				Name:        "pre-deployment-validation",
				Description: "Validate subscription capabilities before Terraform apply",
				Steps: []WorkflowStep{
					{Step: 1, Description: "Check subscription tier", Command: "xcsh subscription show --output-format json"},
					{Step: 2, Description: "Verify required addons are active", Command: "xcsh subscription addons --filter active --output-format json"},
					{Step: 3, Description: "Check quota availability", Command: "xcsh subscription quota -n <namespace> --output-format json"},
					{Step: 4, Description: "Validate specific resources", Command: "xcsh subscription validate --resource-type http_loadbalancer --count 5"},
				},
			},
			{
				Name:        "addon-activation-check",
				Description: "Check and understand addon service availability",
				Steps: []WorkflowStep{
					{Step: 1, Description: "List all addon services", Command: "xcsh subscription addons --all --output-format json"},
					{Step: 2, Description: "Check specific addon status", Command: "xcsh subscription validate --feature bot-defense"},
				},
			},
			{
				Name:        "addon-activation",
				Description: "Activate addon services and monitor status",
				Steps: []WorkflowStep{
					{Step: 1, Description: "List available addons", Command: "xcsh subscription addons --filter available --output-format json"},
					{Step: 2, Description: "Activate desired addon", Command: "xcsh subscription activate --addon <name>"},
					{Step: 3, Description: "Check activation status", Command: "xcsh subscription activation-status --addon <name>"},
				},
			},
		},
	}
}

// TierDescription returns a human-readable description for an addon tier
func TierDescription(tier string) string {
	switch tier {
	case TierNoTier:
		return "No Tier"
	case TierBasic:
		return "Standard" // Basic discontinued, maps to Standard
	case TierStandard:
		return "Standard"
	case TierAdvanced:
		return "Advanced"
	case TierPremium:
		return "Advanced" // Premium discontinued, maps to Advanced
	default:
		return "Unknown"
	}
}

// StateDescription returns a human-readable description for an addon state
func StateDescription(state string) string {
	switch state {
	case StateNone:
		return "Not Subscribed"
	case StatePending:
		return "Pending"
	case StateSubscribed:
		return "Subscribed"
	case StateError:
		return "Error"
	default:
		return "Unknown"
	}
}

// AccessStatusDescription returns a human-readable description for an access status
func AccessStatusDescription(status string) string {
	switch status {
	case AccessAllowed:
		return "Allowed"
	case AccessDenied:
		return "Denied"
	case AccessUpgradeRequired:
		return "Upgrade Required"
	case AccessContactSales:
		return "Contact Sales"
	case AccessInternalService:
		return "Internal Service"
	case AccessUnknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// QuotaStatusFromPercentage returns the status based on usage percentage
func QuotaStatusFromPercentage(percentage float64) string {
	switch {
	case percentage >= 100:
		return "EXCEEDED"
	case percentage >= 80:
		return "WARNING"
	default:
		return "OK"
	}
}

// ActivationTypeDescription returns a human-readable description for an activation type
func ActivationTypeDescription(activationType string) string {
	switch activationType {
	case ActivationSelf:
		return "Self-Activation (immediate)"
	case ActivationPartiallyManaged:
		return "Partially Managed (requires approval)"
	case ActivationManaged:
		return "Fully Managed (SRE intervention required)"
	default:
		return "Unknown"
	}
}

// SubscriptionStateDescription returns a human-readable description for a subscription state
func SubscriptionStateDescription(state string) string {
	switch state {
	case SubscriptionPending:
		return "Pending Activation"
	case SubscriptionEnabled:
		return "Enabled"
	case SubscriptionDisablePending:
		return "Pending Deactivation"
	case SubscriptionDisabled:
		return "Disabled"
	default:
		return state
	}
}
