// Package cloudstatus provides types and functions for interacting with the F5 Cloud Status API.
// The API is based on Atlassian Statuspage v2 and provides real-time status information
// for F5 Distributed Cloud services.
package cloudstatus

import "time"

// BaseURL is the base URL for the F5 Cloud Status API.
const BaseURL = "https://www.f5cloudstatus.com/api/v2"

// Status indicator values for overall system status.
const (
	StatusNone        = "none"
	StatusMinor       = "minor"
	StatusMajor       = "major"
	StatusCritical    = "critical"
	StatusMaintenance = "maintenance"
)

// Component status values.
const (
	ComponentOperational         = "operational"
	ComponentDegradedPerformance = "degraded_performance"
	ComponentPartialOutage       = "partial_outage"
	ComponentMajorOutage         = "major_outage"
	ComponentUnderMaintenance    = "under_maintenance"
)

// Incident status values.
const (
	IncidentInvestigating = "investigating"
	IncidentIdentified    = "identified"
	IncidentMonitoring    = "monitoring"
	IncidentResolved      = "resolved"
	IncidentPostmortem    = "postmortem"
)

// Incident impact values.
const (
	ImpactNone     = "none"
	ImpactMinor    = "minor"
	ImpactMajor    = "major"
	ImpactCritical = "critical"
)

// Maintenance status values.
const (
	MaintenanceScheduled  = "scheduled"
	MaintenanceInProgress = "in_progress"
	MaintenanceVerifying  = "verifying"
	MaintenanceCompleted  = "completed"
)

// Exit codes for CI/CD integration.
const (
	ExitCodeHealthy     = 0
	ExitCodeMinor       = 1
	ExitCodeMajor       = 2
	ExitCodeCritical    = 3
	ExitCodeMaintenance = 4
	ExitCodeAPIError    = 10
	ExitCodeParseError  = 11
)

// StatusResponse represents the response from /status.json endpoint.
type StatusResponse struct {
	Page   PageInfo `json:"page"`
	Status Status   `json:"status"`
}

// SummaryResponse represents the response from /summary.json endpoint.
type SummaryResponse struct {
	Page                  PageInfo               `json:"page"`
	Components            []Component            `json:"components"`
	Incidents             []Incident             `json:"incidents"`
	ScheduledMaintenances []ScheduledMaintenance `json:"scheduled_maintenances"`
	Status                Status                 `json:"status"`
}

// ComponentsResponse represents the response from /components.json endpoint.
type ComponentsResponse struct {
	Page       PageInfo    `json:"page"`
	Components []Component `json:"components"`
}

// ComponentResponse represents the response from /components/{id}.json endpoint.
type ComponentResponse struct {
	Page      PageInfo  `json:"page"`
	Component Component `json:"component"`
}

// IncidentsResponse represents the response from /incidents.json endpoint.
type IncidentsResponse struct {
	Page      PageInfo   `json:"page"`
	Incidents []Incident `json:"incidents"`
}

// MaintenancesResponse represents the response from /scheduled-maintenances.json endpoint.
type MaintenancesResponse struct {
	Page                  PageInfo               `json:"page"`
	ScheduledMaintenances []ScheduledMaintenance `json:"scheduled_maintenances"`
}

// PageInfo contains metadata about the status page.
type PageInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	TimeZone  string    `json:"time_zone"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Status represents the overall status indicator.
type Status struct {
	Indicator   string `json:"indicator"`
	Description string `json:"description"`
}

// Component represents a service component on the status page.
type Component struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Status             string    `json:"status"`
	Description        string    `json:"description"`
	GroupID            *string   `json:"group_id"`
	Group              bool      `json:"group"`
	Components         []string  `json:"components,omitempty"`
	Position           int       `json:"position"`
	Showcase           bool      `json:"showcase"`
	OnlyShowIfDegraded bool      `json:"only_show_if_degraded"`
	PageID             string    `json:"page_id"`
	StartDate          *string   `json:"start_date"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// IsGroup returns true if this component is a group container.
func (c *Component) IsGroup() bool {
	return c.Group
}

// IsOperational returns true if the component is fully operational.
func (c *Component) IsOperational() bool {
	return c.Status == ComponentOperational
}

// IsDegraded returns true if the component has any non-operational status.
func (c *Component) IsDegraded() bool {
	return c.Status != ComponentOperational
}

// Incident represents an incident on the status page.
type Incident struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Status          string           `json:"status"`
	Impact          string           `json:"impact"`
	Shortlink       string           `json:"shortlink"`
	PageID          string           `json:"page_id"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	StartedAt       time.Time        `json:"started_at"`
	MonitoringAt    *time.Time       `json:"monitoring_at"`
	ResolvedAt      *time.Time       `json:"resolved_at"`
	IncidentUpdates []IncidentUpdate `json:"incident_updates"`
	Components      []Component      `json:"components"`
}

// IsResolved returns true if the incident has been resolved.
func (i *Incident) IsResolved() bool {
	return i.Status == IncidentResolved || i.Status == IncidentPostmortem
}

// IsActive returns true if the incident is still active.
func (i *Incident) IsActive() bool {
	return !i.IsResolved()
}

// IncidentUpdate represents an update to an incident.
type IncidentUpdate struct {
	ID                   string              `json:"id"`
	Status               string              `json:"status"`
	Body                 string              `json:"body"`
	IncidentID           string              `json:"incident_id"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
	DisplayAt            time.Time           `json:"display_at"`
	AffectedComponents   []AffectedComponent `json:"affected_components"`
	DeliverNotifications bool                `json:"deliver_notifications"`
}

// AffectedComponent represents a component affected by an incident update.
type AffectedComponent struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
}

// ScheduledMaintenance represents a scheduled maintenance window.
type ScheduledMaintenance struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Status          string           `json:"status"`
	Impact          string           `json:"impact"`
	Shortlink       string           `json:"shortlink"`
	PageID          string           `json:"page_id"`
	ScheduledFor    time.Time        `json:"scheduled_for"`
	ScheduledUntil  time.Time        `json:"scheduled_until"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	StartedAt       *time.Time       `json:"started_at"`
	MonitoringAt    *time.Time       `json:"monitoring_at"`
	ResolvedAt      *time.Time       `json:"resolved_at"`
	IncidentUpdates []IncidentUpdate `json:"incident_updates"`
	Components      []Component      `json:"components"`
}

// IsUpcoming returns true if the maintenance is scheduled but not yet started.
func (m *ScheduledMaintenance) IsUpcoming() bool {
	return m.Status == MaintenanceScheduled
}

// IsActive returns true if the maintenance is currently in progress.
func (m *ScheduledMaintenance) IsActive() bool {
	return m.Status == MaintenanceInProgress || m.Status == MaintenanceVerifying
}

// IsCompleted returns true if the maintenance has been completed.
func (m *ScheduledMaintenance) IsCompleted() bool {
	return m.Status == MaintenanceCompleted
}

// Region represents a geographic region for PoP grouping.
type Region struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// PredefinedRegions contains the known geographic regions for F5 XC PoPs.
var PredefinedRegions = []Region{
	{ID: "north-america", Name: "north-america", DisplayName: "North America"},
	{ID: "south-america", Name: "south-america", DisplayName: "South America"},
	{ID: "europe", Name: "europe", DisplayName: "Europe"},
	{ID: "asia", Name: "asia", DisplayName: "Asia"},
	{ID: "oceania", Name: "oceania", DisplayName: "Oceania"},
	{ID: "middle-east", Name: "middle-east", DisplayName: "Middle East"},
}

// ComponentGroup represents a group of related components.
type ComponentGroup struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Components     []Component `json:"components"`
	ComponentCount int         `json:"component_count"`
}

// RegionalStatus represents aggregated status for a geographic region.
type RegionalStatus struct {
	Region           Region      `json:"region"`
	OverallStatus    string      `json:"overall_status"`
	OperationalCount int         `json:"operational_count"`
	DegradedCount    int         `json:"degraded_count"`
	TotalCount       int         `json:"total_count"`
	Components       []Component `json:"components"`
}

// StatusIndicatorToExitCode converts a status indicator to an exit code.
func StatusIndicatorToExitCode(indicator string) int {
	switch indicator {
	case StatusNone:
		return ExitCodeHealthy
	case StatusMinor:
		return ExitCodeMinor
	case StatusMajor:
		return ExitCodeMajor
	case StatusCritical:
		return ExitCodeCritical
	case StatusMaintenance:
		return ExitCodeMaintenance
	default:
		return ExitCodeHealthy
	}
}

// StatusIndicatorDescription returns a human-readable description for a status indicator.
func StatusIndicatorDescription(indicator string) string {
	switch indicator {
	case StatusNone:
		return "All Systems Operational"
	case StatusMinor:
		return "Minor System Issue"
	case StatusMajor:
		return "Major System Issue"
	case StatusCritical:
		return "Critical System Outage"
	case StatusMaintenance:
		return "System Under Maintenance"
	default:
		return "Unknown Status"
	}
}

// ComponentStatusDescription returns a human-readable description for a component status.
func ComponentStatusDescription(status string) string {
	switch status {
	case ComponentOperational:
		return "Operational"
	case ComponentDegradedPerformance:
		return "Degraded Performance"
	case ComponentPartialOutage:
		return "Partial Outage"
	case ComponentMajorOutage:
		return "Major Outage"
	case ComponentUnderMaintenance:
		return "Under Maintenance"
	default:
		return "Unknown"
	}
}
