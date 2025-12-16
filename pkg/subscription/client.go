package subscription

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/robinmordasiewicz/f5xcctl/pkg/client"
)

// Client provides methods to interact with F5 XC subscription APIs
type Client struct {
	apiClient *client.Client
}

// NewClient creates a new subscription client
func NewClient(apiClient *client.Client) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

// API response types - internal use for parsing API responses

type listPlansResponse struct {
	Items []planItem `json:"items"`
}

type planItem struct {
	Name     string   `json:"name"`
	Metadata metadata `json:"metadata"`
	Spec     planSpec `json:"spec"`
}

type metadata struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type planSpec struct {
	DisplayName      string             `json:"display_name,omitempty"`
	Description      string             `json:"description,omitempty"`
	IncludedServices []serviceReference `json:"included_services,omitempty"`
	AllowedServices  []serviceReference `json:"allowed_services,omitempty"`
}

type serviceReference struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type listAddonServicesResponse struct {
	Items []addonServiceItem `json:"items"`
}

// addonServiceItem matches the actual F5 XC API response structure for addon_services
// The API returns a flat structure with fields directly on the item
type addonServiceItem struct {
	Tenant      string `json:"tenant,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	Name        string `json:"name"`
	UID         string `json:"uid,omitempty"`
	Description string `json:"description,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
}

type activationStatusResponse struct {
	ActivationStatus string `json:"activation_status,omitempty"`
	State            string `json:"state,omitempty"`
	AccessStatus     string `json:"access_status,omitempty"`
	Tier             string `json:"tier,omitempty"`
}

type allActivationStatusResponse struct {
	Items []namespaceActivationStatus `json:"items"`
}

type namespaceActivationStatus struct {
	Namespace        string `json:"namespace"`
	ActivationStatus string `json:"activation_status,omitempty"`
	State            string `json:"state,omitempty"`
	AccessStatus     string `json:"access_status,omitempty"`
}

// quotaAPIResponse matches the actual F5 XC API response structure for quota endpoints
// Quotas are TENANT-level - the namespace parameter in the API is actually the tenant context
type quotaAPIResponse struct {
	// Primary fields (new API structure)
	Objects   map[string]quotaEntry `json:"objects,omitempty"`
	Resources map[string]quotaEntry `json:"resources,omitempty"`
	APIs      map[string]quotaEntry `json:"apis,omitempty"`
	// Deprecated fields (still returned by API, fallback for older responses)
	QuotaUsage map[string]quotaEntry `json:"quota_usage,omitempty"`
}

type quotaEntry struct {
	Limit       quotaLimit  `json:"limit,omitempty"`
	Usage       quotaUsage  `json:"usage,omitempty"`
	DisplayName string      `json:"display_name,omitempty"`
	Description string      `json:"description,omitempty"`
}

type quotaLimit struct {
	Maximum float64 `json:"maximum"`
}

type quotaUsage struct {
	Current float64 `json:"current"`
}

// GetPlans retrieves subscription plans for a namespace
func (c *Client) GetPlans(ctx context.Context, namespace string) ([]PlanInfo, error) {
	if namespace == "" {
		namespace = "system"
	}

	path := fmt.Sprintf("/api/web/namespaces/%s/plans", namespace)
	resp, err := c.apiClient.Get(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get plans: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(resp.Body))
	}

	var plansResp listPlansResponse
	if err := json.Unmarshal(resp.Body, &plansResp); err != nil {
		return nil, fmt.Errorf("failed to parse plans response: %w", err)
	}

	var plans []PlanInfo
	for _, item := range plansResp.Items {
		plan := PlanInfo{
			Name:        item.Metadata.Name,
			DisplayName: item.Spec.DisplayName,
			Description: item.Spec.Description,
		}

		// Extract included services
		for _, svc := range item.Spec.IncludedServices {
			plan.IncludedServices = append(plan.IncludedServices, svc.Name)
		}

		// Extract allowed services
		for _, svc := range item.Spec.AllowedServices {
			plan.AllowedServices = append(plan.AllowedServices, svc.Name)
		}

		plans = append(plans, plan)
	}

	return plans, nil
}

// GetAddonServices retrieves addon services for a namespace
// The API returns a flat structure with basic info, and activation status is fetched separately
func (c *Client) GetAddonServices(ctx context.Context, namespace string) ([]AddonServiceInfo, error) {
	if namespace == "" {
		namespace = "system"
	}

	path := fmt.Sprintf("/api/web/namespaces/%s/addon_services", namespace)
	resp, err := c.apiClient.Get(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get addon services: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(resp.Body))
	}

	var addonsResp listAddonServicesResponse
	if err := json.Unmarshal(resp.Body, &addonsResp); err != nil {
		return nil, fmt.Errorf("failed to parse addon services response: %w", err)
	}

	var addons []AddonServiceInfo
	for _, item := range addonsResp.Items {
		// Start with basic info from the flat structure
		addon := AddonServiceInfo{
			Name:        item.Name,
			DisplayName: formatDisplayName(item.Name), // Generate display name from service name
			Description: item.Description,
			Namespace:   item.Namespace,
			// State defaults to None, will be updated by activation status
			State:        StateNone,
			AccessStatus: AccessAllowed,
		}

		// Skip disabled services
		if item.Disabled {
			addon.State = StateNone
			addon.AccessStatus = AccessDenied
		}

		// Fetch activation status for each addon to get state and access info
		activationStatus, err := c.GetAddonServiceActivationStatus(ctx, item.Name)
		if err == nil && activationStatus != nil {
			addon.State = activationStatus.State
			addon.AccessStatus = activationStatus.AccessStatus
			if activationStatus.Tier != "" {
				addon.Tier = activationStatus.Tier
			}
		}

		addons = append(addons, addon)
	}

	return addons, nil
}

// formatDisplayName converts a service name to a display-friendly format
// e.g., "client-side-defense" -> "Client Side Defense"
func formatDisplayName(name string) string {
	words := strings.Split(name, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

// GetAddonServiceActivationStatus gets the activation status for a specific addon
func (c *Client) GetAddonServiceActivationStatus(ctx context.Context, addonName string) (*AddonServiceInfo, error) {
	path := fmt.Sprintf("/api/web/namespaces/system/addon_services/%s/activation-status", addonName)
	resp, err := c.apiClient.Get(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get addon activation status: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(resp.Body))
	}

	var statusResp activationStatusResponse
	if err := json.Unmarshal(resp.Body, &statusResp); err != nil {
		return nil, fmt.Errorf("failed to parse activation status response: %w", err)
	}

	return &AddonServiceInfo{
		Name:         addonName,
		State:        normalizeState(statusResp.State),
		AccessStatus: normalizeAccessStatus(statusResp.AccessStatus),
		Tier:         normalizeTier(statusResp.Tier),
	}, nil
}

// GetAllAddonServiceActivationStatus gets activation status across all namespaces
func (c *Client) GetAllAddonServiceActivationStatus(ctx context.Context, addonName string) ([]AddonServiceInfo, error) {
	path := fmt.Sprintf("/api/web/namespaces/system/addon_services/%s/all-activation-status", addonName)
	resp, err := c.apiClient.Get(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get all addon activation status: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(resp.Body))
	}

	var statusResp allActivationStatusResponse
	if err := json.Unmarshal(resp.Body, &statusResp); err != nil {
		return nil, fmt.Errorf("failed to parse all activation status response: %w", err)
	}

	var addons []AddonServiceInfo
	for _, item := range statusResp.Items {
		addon := AddonServiceInfo{
			Name:         addonName,
			State:        normalizeState(item.State),
			AccessStatus: normalizeAccessStatus(item.AccessStatus),
			Namespace:    item.Namespace,
		}
		addons = append(addons, addon)
	}

	return addons, nil
}

// GetQuotaInfo retrieves tenant-level quota limits and usage.
// Note: Quotas are enforced at the TENANT level, not per-namespace.
// The "system" namespace is used as the tenant context - querying any namespace
// returns the same tenant-wide quota data.
func (c *Client) GetQuotaInfo(ctx context.Context) (*QuotaUsageInfo, error) {
	// Use /quota/usage endpoint which returns both limits AND current usage
	// The namespace "system" is used as tenant context
	path := "/api/web/namespaces/system/quota/usage"
	resp, err := c.apiClient.Get(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get quota info: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(resp.Body))
	}

	var quotaResp quotaAPIResponse
	if err := json.Unmarshal(resp.Body, &quotaResp); err != nil {
		return nil, fmt.Errorf("failed to parse quota response: %w", err)
	}

	// Convert map response to slice for easier iteration
	var objects []QuotaItem

	// Use "objects" field if present, fall back to deprecated "quota_usage"
	quotaMap := quotaResp.Objects
	if len(quotaMap) == 0 {
		quotaMap = quotaResp.QuotaUsage
	}

	for name, entry := range quotaMap {
		limit := entry.Limit.Maximum
		usage := entry.Usage.Current

		// Skip unlimited quotas (maximum = -1) and items with negative usage
		if limit < 0 || usage < 0 {
			continue
		}

		percentage := float64(0)
		if limit > 0 {
			percentage = (usage / limit) * 100
		}

		displayName := entry.DisplayName
		if displayName == "" {
			displayName = name // Use name as display name if not provided
		}

		objects = append(objects, QuotaItem{
			Name:        name,
			DisplayName: displayName,
			Description: entry.Description,
			ObjectType:  name,
			Limit:       limit,
			Usage:       usage,
			Percentage:  percentage,
			Status:      QuotaStatusFromPercentage(percentage),
		})
	}

	return &QuotaUsageInfo{
		Namespace: "tenant", // Quotas are tenant-level, not namespace-specific
		Objects:   objects,
	}, nil
}

// GetSubscriptionInfo retrieves complete subscription information
func (c *Client) GetSubscriptionInfo(ctx context.Context) (*SubscriptionInfo, error) {
	// Get plans
	plans, err := c.GetPlans(ctx, "system")
	if err != nil {
		return nil, fmt.Errorf("failed to get plans: %w", err)
	}

	// Get addon services
	addons, err := c.GetAddonServices(ctx, "system")
	if err != nil {
		return nil, fmt.Errorf("failed to get addon services: %w", err)
	}

	// Get quota summary (tenant-level)
	quotaInfo, err := c.GetQuotaInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get quota info: %w", err)
	}

	// Determine tier from plan or addons
	tier := determineTier(plans, addons)

	// Separate active and available addons
	var activeAddons, availableAddons []AddonServiceInfo
	for _, addon := range addons {
		if addon.IsActive() {
			activeAddons = append(activeAddons, addon)
		} else if addon.IsAvailable() {
			availableAddons = append(availableAddons, addon)
		}
	}

	// Build quota summary
	var atRisk, exceeded int
	for _, q := range quotaInfo.Objects {
		if q.IsExceeded() {
			exceeded++
		} else if q.IsAtRisk() {
			atRisk++
		}
	}

	info := &SubscriptionInfo{
		Tier:            tier,
		ActiveAddons:    activeAddons,
		AvailableAddons: availableAddons,
		QuotaSummary: QuotaSummary{
			TotalLimits:    len(quotaInfo.Objects),
			LimitsAtRisk:   atRisk,
			LimitsExceeded: exceeded,
			Objects:        quotaInfo.Objects,
		},
	}

	// Set plan info if available
	if len(plans) > 0 {
		info.Plan = plans[0]
	}

	return info, nil
}

// ValidateResource validates if a resource can be deployed
func (c *Client) ValidateResource(ctx context.Context, req ValidationRequest) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:  true,
		Checks: []ValidationCheck{},
	}

	namespace := req.Namespace
	if namespace == "" {
		namespace = "system"
	}

	// Validate quota if resource type and count specified
	// Note: Quotas are tenant-level, not namespace-level
	if req.ResourceType != "" && req.Count > 0 {
		quotaInfo, err := c.GetQuotaInfo(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get quota info: %w", err)
		}

		// Find the quota for this resource type
		var found bool
		for _, q := range quotaInfo.Objects {
			if strings.EqualFold(q.ObjectType, req.ResourceType) || strings.EqualFold(q.Name, req.ResourceType) {
				found = true
				remaining := q.RemainingCapacity()
				check := ValidationCheck{
					Type:      "quota",
					Resource:  req.ResourceType,
					Current:   int(q.Usage),
					Requested: req.Count,
					Limit:     int(q.Limit),
				}

				if float64(req.Count) > remaining {
					check.Result = ValidationFail
					check.Message = fmt.Sprintf("Quota exceeded: %s would have %d/%d (requesting %d, only %d available)",
						req.ResourceType, int(q.Usage)+req.Count, int(q.Limit), req.Count, int(remaining))
				} else if (q.Usage+float64(req.Count))/q.Limit >= 0.8 {
					check.Result = ValidationWarning
					check.Message = fmt.Sprintf("Quota warning: %s will be at %.0f%% after deployment",
						req.ResourceType, ((q.Usage+float64(req.Count))/q.Limit)*100)
				} else {
					check.Result = ValidationPass
					check.Message = fmt.Sprintf("Quota OK: %s has sufficient capacity (%d available)",
						req.ResourceType, int(remaining))
				}

				result.AddCheck(check)
				break
			}
		}

		if !found {
			check := ValidationCheck{
				Type:      "quota",
				Resource:  req.ResourceType,
				Requested: req.Count,
				Result:    ValidationWarning,
				Message:   fmt.Sprintf("No quota limit found for resource type: %s", req.ResourceType),
			}
			result.AddCheck(check)
		}
	}

	// Validate feature/addon if specified
	if req.Feature != "" {
		addons, err := c.GetAddonServices(ctx, "system")
		if err != nil {
			return nil, fmt.Errorf("failed to get addon services: %w", err)
		}

		var found bool
		for _, addon := range addons {
			if strings.EqualFold(addon.Name, req.Feature) {
				found = true
				check := ValidationCheck{
					Type:        "feature",
					Feature:     req.Feature,
					CurrentTier: addon.Tier,
					Status:      StateDescription(addon.State),
				}

				if addon.IsActive() {
					check.Result = ValidationPass
					check.Message = fmt.Sprintf("Feature '%s' is active (tier: %s)", req.Feature, addon.Tier)
				} else if addon.NeedsUpgrade() {
					check.Result = ValidationFail
					check.Message = fmt.Sprintf("Feature '%s' requires a plan upgrade", req.Feature)
				} else if addon.NeedsContactSales() {
					check.Result = ValidationFail
					check.Message = fmt.Sprintf("Feature '%s' requires contacting F5 sales", req.Feature)
				} else if addon.IsAvailable() {
					check.Result = ValidationWarning
					check.Message = fmt.Sprintf("Feature '%s' is available but not subscribed", req.Feature)
				} else {
					check.Result = ValidationFail
					check.Message = fmt.Sprintf("Feature '%s' is not available (access: %s)", req.Feature, AccessStatusDescription(addon.AccessStatus))
				}

				result.AddCheck(check)
				break
			}
		}

		if !found {
			check := ValidationCheck{
				Type:    "feature",
				Feature: req.Feature,
				Result:  ValidationWarning,
				Message: fmt.Sprintf("Feature '%s' not found in addon services", req.Feature),
			}
			result.AddCheck(check)
		}
	}

	return result, nil
}

// FilterAddons filters addon services based on criteria
func FilterAddons(addons []AddonServiceInfo, filter string) []AddonServiceInfo {
	if filter == "" {
		return addons
	}

	var filtered []AddonServiceInfo
	for _, addon := range addons {
		switch strings.ToLower(filter) {
		case "active":
			if addon.IsActive() {
				filtered = append(filtered, addon)
			}
		case "available":
			if addon.IsAvailable() {
				filtered = append(filtered, addon)
			}
		case "denied":
			if addon.IsDenied() {
				filtered = append(filtered, addon)
			}
		}
	}

	return filtered
}

// Helper functions

func normalizeTier(tier string) string {
	tier = strings.ToUpper(strings.TrimSpace(tier))
	switch tier {
	case "NO_TIER", "NOTIER", "":
		return TierNoTier
	case "BASIC":
		return TierBasic
	case "STANDARD":
		return TierStandard
	case "ADVANCED":
		return TierAdvanced
	case "PREMIUM":
		return TierPremium
	default:
		return tier
	}
}

func normalizeState(state string) string {
	state = strings.ToUpper(strings.TrimSpace(state))
	switch state {
	case "AS_NONE", "NONE", "":
		return StateNone
	case "AS_PENDING", "PENDING":
		return StatePending
	case "AS_SUBSCRIBED", "SUBSCRIBED":
		return StateSubscribed
	case "AS_ERROR", "ERROR":
		return StateError
	default:
		return state
	}
}

func normalizeAccessStatus(status string) string {
	status = strings.ToUpper(strings.TrimSpace(status))
	switch status {
	case "AS_AC_ALLOWED", "ALLOWED", "":
		return AccessAllowed
	case "AS_AC_PBAC_DENY", "DENIED", "DENY":
		return AccessDenied
	case "AS_AC_PBAC_DENY_UPGRADE_PLAN", "UPGRADE_REQUIRED", "UPGRADE_PLAN":
		return AccessUpgradeRequired
	case "AS_AC_PBAC_DENY_CONTACT_SALES", "CONTACT_SALES":
		return AccessContactSales
	case "AS_AC_PBAC_DENY_INTERNAL_SVC", "INTERNAL_SERVICE":
		return AccessInternalService
	default:
		return status
	}
}

func determineTier(plans []PlanInfo, addons []AddonServiceInfo) string {
	// Check plans first
	for _, plan := range plans {
		nameLower := strings.ToLower(plan.Name)
		displayLower := strings.ToLower(plan.DisplayName)

		if strings.Contains(nameLower, "advanced") || strings.Contains(displayLower, "advanced") {
			return "Advanced"
		}
		if strings.Contains(nameLower, "enterprise") || strings.Contains(displayLower, "enterprise") {
			return "Advanced"
		}
	}

	// Check if any advanced-tier addons are active
	for _, addon := range addons {
		if addon.IsActive() && (addon.Tier == TierAdvanced || addon.Tier == TierPremium) {
			return "Advanced"
		}
	}

	return "Standard"
}

// BuildCurlCommand builds a curl command for debugging
func BuildCurlCommand(baseURL, path string, query url.Values) string {
	fullURL := strings.TrimRight(baseURL, "/") + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}
	return fmt.Sprintf("curl -X GET '%s' -H 'Content-Type: application/json'", fullURL)
}
