package cloudstatus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Client provides methods to interact with the F5 Cloud Status API.
type Client struct {
	baseURL      string
	httpClient   *http.Client
	cache        *Cache
	cacheEnabled bool
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the client.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(url, "/")
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithCache enables caching with the specified TTL.
func WithCache(ttl time.Duration) ClientOption {
	return func(c *Client) {
		c.cache = NewCache(ttl)
		c.cacheEnabled = true
	}
}

// WithoutCache disables caching.
func WithoutCache() ClientOption {
	return func(c *Client) {
		c.cacheEnabled = false
	}
}

// NewClient creates a new F5 Cloud Status API client.
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		baseURL: BaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache:        NewCache(60 * time.Second),
		cacheEnabled: true,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// get performs an HTTP GET request and returns the response body.
func (c *Client) get(endpoint string) ([]byte, error) {
	url := c.baseURL + endpoint

	// Check cache first
	if c.cacheEnabled {
		if cached, ok := c.cache.Get(url); ok {
			return cached, nil
		}
	}

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Store in cache
	if c.cacheEnabled {
		c.cache.Set(url, body)
	}

	return body, nil
}

// GetStatus retrieves the overall status indicator.
func (c *Client) GetStatus() (*StatusResponse, error) {
	body, err := c.get("/status.json")
	if err != nil {
		return nil, err
	}

	var resp StatusResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse status response: %w", err)
	}

	return &resp, nil
}

// GetSummary retrieves the complete status summary.
func (c *Client) GetSummary() (*SummaryResponse, error) {
	body, err := c.get("/summary.json")
	if err != nil {
		return nil, err
	}

	var resp SummaryResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse summary response: %w", err)
	}

	return &resp, nil
}

// GetComponents retrieves all components.
func (c *Client) GetComponents() (*ComponentsResponse, error) {
	body, err := c.get("/components.json")
	if err != nil {
		return nil, err
	}

	var resp ComponentsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse components response: %w", err)
	}

	return &resp, nil
}

// GetComponent retrieves a single component by ID.
func (c *Client) GetComponent(id string) (*ComponentResponse, error) {
	body, err := c.get(fmt.Sprintf("/components/%s.json", id))
	if err != nil {
		return nil, err
	}

	var resp ComponentResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse component response: %w", err)
	}

	return &resp, nil
}

// GetIncidents retrieves all incidents.
func (c *Client) GetIncidents() (*IncidentsResponse, error) {
	body, err := c.get("/incidents.json")
	if err != nil {
		return nil, err
	}

	var resp IncidentsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse incidents response: %w", err)
	}

	return &resp, nil
}

// GetUnresolvedIncidents retrieves only unresolved incidents.
func (c *Client) GetUnresolvedIncidents() (*IncidentsResponse, error) {
	body, err := c.get("/incidents/unresolved.json")
	if err != nil {
		return nil, err
	}

	var resp IncidentsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse unresolved incidents response: %w", err)
	}

	return &resp, nil
}

// GetMaintenances retrieves all scheduled maintenances.
func (c *Client) GetMaintenances() (*MaintenancesResponse, error) {
	body, err := c.get("/scheduled-maintenances.json")
	if err != nil {
		return nil, err
	}

	var resp MaintenancesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse maintenances response: %w", err)
	}

	return &resp, nil
}

// GetUpcomingMaintenances retrieves only upcoming maintenances.
func (c *Client) GetUpcomingMaintenances() (*MaintenancesResponse, error) {
	body, err := c.get("/scheduled-maintenances/upcoming.json")
	if err != nil {
		return nil, err
	}

	var resp MaintenancesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse upcoming maintenances response: %w", err)
	}

	return &resp, nil
}

// ClearCache clears the response cache.
func (c *Client) ClearCache() {
	if c.cache != nil {
		c.cache.Clear()
	}
}

// GetComponentGroups extracts component groups from all components.
func (c *Client) GetComponentGroups() ([]ComponentGroup, error) {
	resp, err := c.GetComponents()
	if err != nil {
		return nil, err
	}

	return ExtractComponentGroups(resp.Components), nil
}

// GetPoPs extracts PoP components from all components.
func (c *Client) GetPoPs() ([]Component, error) {
	resp, err := c.GetComponents()
	if err != nil {
		return nil, err
	}

	return FilterPoPs(resp.Components), nil
}

// GetRegionalStatus returns aggregated status by region.
func (c *Client) GetRegionalStatus() ([]RegionalStatus, error) {
	resp, err := c.GetComponents()
	if err != nil {
		return nil, err
	}

	return CalculateRegionalStatus(resp.Components), nil
}

// Helper functions for component filtering and grouping

// ExtractComponentGroups extracts groups and their child components.
func ExtractComponentGroups(components []Component) []ComponentGroup {
	// Map group IDs to groups
	groupMap := make(map[string]*ComponentGroup)

	// First pass: identify groups
	for _, comp := range components {
		if comp.Group {
			group := ComponentGroup{
				ID:          comp.ID,
				Name:        comp.Name,
				Description: comp.Description,
				Components:  []Component{},
			}
			groupMap[comp.ID] = &group
		}
	}

	// Second pass: assign components to groups
	for _, comp := range components {
		if !comp.Group && comp.GroupID != nil {
			if group, ok := groupMap[*comp.GroupID]; ok {
				group.Components = append(group.Components, comp)
				group.ComponentCount++
			}
		}
	}

	// Convert map to slice
	groups := make([]ComponentGroup, 0, len(groupMap))
	for _, group := range groupMap {
		groups = append(groups, *group)
	}

	return groups
}

// FilterPoPs returns only components that are PoPs (Point of Presence).
func FilterPoPs(components []Component) []Component {
	pops := []Component{}
	popRegex := regexp.MustCompile(`(?i)\bpop\b|edge\s*pop|point\s*of\s*presence`)

	for _, comp := range components {
		if !comp.Group && popRegex.MatchString(comp.Description) {
			pops = append(pops, comp)
		}
	}

	return pops
}

// FilterByStatus returns components matching the given status.
func FilterByStatus(components []Component, status string) []Component {
	filtered := []Component{}
	for _, comp := range components {
		if comp.Status == status {
			filtered = append(filtered, comp)
		}
	}
	return filtered
}

// FilterDegraded returns only non-operational components.
func FilterDegraded(components []Component) []Component {
	degraded := []Component{}
	for _, comp := range components {
		if comp.IsDegraded() && !comp.Group {
			degraded = append(degraded, comp)
		}
	}
	return degraded
}

// FilterByGroup returns components belonging to a specific group.
func FilterByGroup(components []Component, groupID string) []Component {
	filtered := []Component{}
	for _, comp := range components {
		if comp.GroupID != nil && *comp.GroupID == groupID {
			filtered = append(filtered, comp)
		}
	}
	return filtered
}

// ExtractSiteCode extracts the site code from a component name or description.
// Site codes are typically in the format (site-code) within the name.
func ExtractSiteCode(name string) string {
	re := regexp.MustCompile(`\(([a-z0-9-]+)\)`)
	matches := re.FindStringSubmatch(strings.ToLower(name))
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// DetectRegion attempts to determine the region from a component's group or name.
func DetectRegion(comp Component, groups []ComponentGroup) string {
	if comp.GroupID == nil {
		return ""
	}

	for _, group := range groups {
		if group.ID == *comp.GroupID {
			groupNameLower := strings.ToLower(group.Name)
			switch {
			case strings.Contains(groupNameLower, "north america"):
				return "north-america"
			case strings.Contains(groupNameLower, "south america"):
				return "south-america"
			case strings.Contains(groupNameLower, "europe"):
				return "europe"
			case strings.Contains(groupNameLower, "asia"):
				return "asia"
			case strings.Contains(groupNameLower, "oceania"):
				return "oceania"
			case strings.Contains(groupNameLower, "middle east"):
				return "middle-east"
			}
		}
	}

	return ""
}

// CalculateRegionalStatus calculates aggregated status for each region.
func CalculateRegionalStatus(components []Component) []RegionalStatus {
	groups := ExtractComponentGroups(components)
	pops := FilterPoPs(components)

	regionMap := make(map[string]*RegionalStatus)

	// Initialize regions
	for _, region := range PredefinedRegions {
		regionMap[region.ID] = &RegionalStatus{
			Region:        region,
			OverallStatus: StatusNone,
			Components:    []Component{},
		}
	}

	// Assign PoPs to regions
	for _, pop := range pops {
		regionID := DetectRegion(pop, groups)
		if regionID == "" {
			continue
		}

		if regional, ok := regionMap[regionID]; ok {
			regional.Components = append(regional.Components, pop)
			regional.TotalCount++
			if pop.IsOperational() {
				regional.OperationalCount++
			} else {
				regional.DegradedCount++
			}
		}
	}

	// Calculate overall status for each region
	for _, regional := range regionMap {
		if regional.DegradedCount == 0 {
			regional.OverallStatus = StatusNone
		} else if float64(regional.DegradedCount)/float64(regional.TotalCount) < 0.25 {
			regional.OverallStatus = StatusMinor
		} else if float64(regional.DegradedCount)/float64(regional.TotalCount) < 0.5 {
			regional.OverallStatus = StatusMajor
		} else {
			regional.OverallStatus = StatusCritical
		}
	}

	// Convert map to slice
	statuses := make([]RegionalStatus, 0, len(regionMap))
	for _, regional := range regionMap {
		statuses = append(statuses, *regional)
	}

	return statuses
}

// FilterIncidentsByStatus filters incidents by their status.
func FilterIncidentsByStatus(incidents []Incident, status string) []Incident {
	filtered := []Incident{}
	for _, inc := range incidents {
		if inc.Status == status {
			filtered = append(filtered, inc)
		}
	}
	return filtered
}

// FilterIncidentsByImpact filters incidents by their impact level.
func FilterIncidentsByImpact(incidents []Incident, impact string) []Incident {
	filtered := []Incident{}
	for _, inc := range incidents {
		if inc.Impact == impact {
			filtered = append(filtered, inc)
		}
	}
	return filtered
}

// FilterIncidentsSince filters incidents created after a given time.
func FilterIncidentsSince(incidents []Incident, since time.Time) []Incident {
	filtered := []Incident{}
	for _, inc := range incidents {
		if inc.CreatedAt.After(since) {
			filtered = append(filtered, inc)
		}
	}
	return filtered
}

// FilterMaintenancesByStatus filters maintenances by their status.
func FilterMaintenancesByStatus(maintenances []ScheduledMaintenance, status string) []ScheduledMaintenance {
	filtered := []ScheduledMaintenance{}
	for _, maint := range maintenances {
		if maint.Status == status {
			filtered = append(filtered, maint)
		}
	}
	return filtered
}

// GetActiveMaintenances returns maintenances that are currently in progress.
func GetActiveMaintenances(maintenances []ScheduledMaintenance) []ScheduledMaintenance {
	active := []ScheduledMaintenance{}
	for _, maint := range maintenances {
		if maint.IsActive() {
			active = append(active, maint)
		}
	}
	return active
}
