package cloudstatus

// Spec represents the complete cloudstatus specification for AI agents
type Spec struct {
	Name               string               `json:"name" yaml:"name"`
	Description        string               `json:"description" yaml:"description"`
	BaseURL            string               `json:"base_url" yaml:"base_url"`
	Authentication     string               `json:"authentication" yaml:"authentication"`
	AIHints            AIHints              `json:"ai_hints" yaml:"ai_hints"`
	DataModel          DataModel            `json:"data_model" yaml:"data_model"`
	StatusIndicators   StatusIndicators     `json:"status_indicators" yaml:"status_indicators"`
	Endpoints          []EndpointSpec       `json:"endpoints" yaml:"endpoints"`
	Commands           []CommandSpec        `json:"commands" yaml:"commands"`
	Workflows          []WorkflowSpec       `json:"workflows" yaml:"workflows"`
	ExitCodes          []ExitCodeSpec       `json:"exit_codes" yaml:"exit_codes"`
	FilterCapabilities FilterCapabilities   `json:"filter_capabilities" yaml:"filter_capabilities"`
	OutputFormats      []string             `json:"output_formats" yaml:"output_formats"`
	CachingBehavior    CachingSpec          `json:"caching_behavior" yaml:"caching_behavior"`
	ComponentGroups    []ComponentGroupSpec `json:"component_groups" yaml:"component_groups"`
	Regions            []RegionSpec         `json:"regions" yaml:"regions"`
}

// AIHints provides guidance for AI agents on how to use the cloudstatus commands
type AIHints struct {
	DiscoveryCommand       string        `json:"discovery_command" yaml:"discovery_command"`
	QuickStatusCheck       string        `json:"quick_status_check" yaml:"quick_status_check"`
	ComprehensiveStatus    string        `json:"comprehensive_status" yaml:"comprehensive_status"`
	MonitoringSetup        string        `json:"monitoring_setup" yaml:"monitoring_setup"`
	AuthenticationRequired bool          `json:"authentication_required" yaml:"authentication_required"`
	RecommendedPolling     int           `json:"recommended_polling_seconds" yaml:"recommended_polling_seconds"`
	BestPractices          []string      `json:"best_practices" yaml:"best_practices"`
	UseCases               []UseCaseSpec `json:"use_cases" yaml:"use_cases"`
}

// UseCaseSpec describes a specific use case for the cloudstatus API
type UseCaseSpec struct {
	Scenario    string `json:"scenario" yaml:"scenario"`
	Command     string `json:"command" yaml:"command"`
	Description string `json:"description" yaml:"description"`
}

// DataModel describes the structure of API responses
type DataModel struct {
	Components  []FieldSpec `json:"components" yaml:"components"`
	Incidents   []FieldSpec `json:"incidents" yaml:"incidents"`
	Maintenance []FieldSpec `json:"maintenance" yaml:"maintenance"`
	Status      []FieldSpec `json:"status" yaml:"status"`
}

// FieldSpec describes a field in the data model
type FieldSpec struct {
	Field       string `json:"field" yaml:"field"`
	Type        string `json:"type" yaml:"type"`
	Description string `json:"description" yaml:"description"`
	Example     string `json:"example" yaml:"example"`
}

// StatusIndicators describes the possible status values
type StatusIndicators struct {
	Overall           []StatusValueSpec `json:"overall" yaml:"overall"`
	Component         []StatusValueSpec `json:"component" yaml:"component"`
	IncidentImpact    []StatusValueSpec `json:"incident_impact" yaml:"incident_impact"`
	IncidentStatus    []StatusValueSpec `json:"incident_status" yaml:"incident_status"`
	MaintenanceStatus []StatusValueSpec `json:"maintenance_status" yaml:"maintenance_status"`
}

// StatusValueSpec describes a status value
type StatusValueSpec struct {
	Value       string `json:"value" yaml:"value"`
	Description string `json:"description" yaml:"description"`
	Severity    int    `json:"severity" yaml:"severity"`
	ExitCode    int    `json:"exit_code,omitempty" yaml:"exit_code,omitempty"`
	Color       string `json:"color,omitempty" yaml:"color,omitempty"`
}

// EndpointSpec describes an API endpoint
type EndpointSpec struct {
	Path         string `json:"path" yaml:"path"`
	Method       string `json:"method" yaml:"method"`
	Description  string `json:"description" yaml:"description"`
	ResponseType string `json:"response_type" yaml:"response_type"`
}

// CommandSpec describes a CLI command
type CommandSpec struct {
	Name        string     `json:"name" yaml:"name"`
	Path        []string   `json:"path" yaml:"path"`
	Description string     `json:"description" yaml:"description"`
	Flags       []FlagSpec `json:"flags,omitempty" yaml:"flags,omitempty"`
	Args        []string   `json:"args,omitempty" yaml:"args,omitempty"`
	ExitCodes   []int      `json:"exit_codes,omitempty" yaml:"exit_codes,omitempty"`
	Examples    []string   `json:"examples,omitempty" yaml:"examples,omitempty"`
}

// FlagSpec describes a command flag
type FlagSpec struct {
	Name        string   `json:"name" yaml:"name"`
	Shorthand   string   `json:"shorthand,omitempty" yaml:"shorthand,omitempty"`
	Type        string   `json:"type" yaml:"type"`
	Default     string   `json:"default,omitempty" yaml:"default,omitempty"`
	Description string   `json:"description" yaml:"description"`
	ValidValues []string `json:"valid_values,omitempty" yaml:"valid_values,omitempty"`
}

// WorkflowSpec describes a multi-step workflow
type WorkflowSpec struct {
	Name        string         `json:"name" yaml:"name"`
	Description string         `json:"description" yaml:"description"`
	Steps       []WorkflowStep `json:"steps" yaml:"steps"`
}

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	Step        int    `json:"step" yaml:"step"`
	Command     string `json:"command" yaml:"command"`
	Description string `json:"description" yaml:"description"`
}

// ExitCodeSpec describes an exit code
type ExitCodeSpec struct {
	Code        int    `json:"code" yaml:"code"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

// FilterCapabilities describes filtering options
type FilterCapabilities struct {
	ComponentFilters   []FilterSpec `json:"component_filters" yaml:"component_filters"`
	IncidentFilters    []FilterSpec `json:"incident_filters" yaml:"incident_filters"`
	MaintenanceFilters []FilterSpec `json:"maintenance_filters" yaml:"maintenance_filters"`
	PopFilters         []FilterSpec `json:"pop_filters" yaml:"pop_filters"`
	TimeFilters        []FilterSpec `json:"time_filters" yaml:"time_filters"`
}

// FilterSpec describes a filter option
type FilterSpec struct {
	Name        string   `json:"name" yaml:"name"`
	Flag        string   `json:"flag" yaml:"flag"`
	Type        string   `json:"type" yaml:"type"`
	Description string   `json:"description" yaml:"description"`
	ValidValues []string `json:"valid_values,omitempty" yaml:"valid_values,omitempty"`
	Example     string   `json:"example" yaml:"example"`
}

// CachingSpec describes caching behavior
type CachingSpec struct {
	DefaultTTLSeconds int    `json:"default_ttl_seconds" yaml:"default_ttl_seconds"`
	BypassFlag        string `json:"bypass_flag" yaml:"bypass_flag"`
	TTLOverrideFlag   string `json:"ttl_override_flag" yaml:"ttl_override_flag"`
	CacheLocation     string `json:"cache_location" yaml:"cache_location"`
}

// ComponentGroupSpec describes a component group
type ComponentGroupSpec struct {
	ID             string `json:"id" yaml:"id"`
	Name           string `json:"name" yaml:"name"`
	Description    string `json:"description" yaml:"description"`
	ComponentCount int    `json:"component_count" yaml:"component_count"`
}

// RegionSpec describes a geographic region
type RegionSpec struct {
	Name        string `json:"name" yaml:"name"`
	DisplayName string `json:"display_name" yaml:"display_name"`
	PopCount    int    `json:"pop_count" yaml:"pop_count"`
}

// GenerateSpec generates the complete cloudstatus specification
func GenerateSpec() *Spec {
	return &Spec{
		Name:               "cloudstatus",
		Description:        "Monitor F5 Distributed Cloud service status, incidents, and maintenance",
		BaseURL:            BaseURL + "/",
		Authentication:     "none (public API)",
		AIHints:            getAIHints(),
		DataModel:          getDataModel(),
		StatusIndicators:   getStatusIndicators(),
		Endpoints:          getEndpoints(),
		Commands:           getCommands(),
		Workflows:          getWorkflows(),
		ExitCodes:          getExitCodes(),
		FilterCapabilities: getFilterCapabilities(),
		OutputFormats:      []string{"json", "yaml", "table", "wide"},
		CachingBehavior:    getCachingSpec(),
		ComponentGroups:    getComponentGroups(),
		Regions:            getRegions(),
	}
}

func getAIHints() AIHints {
	return AIHints{
		DiscoveryCommand:       "xcsh cloudstatus --spec --output-format json",
		QuickStatusCheck:       "xcsh cloudstatus status",
		ComprehensiveStatus:    "xcsh cloudstatus summary --output-format json",
		MonitoringSetup:        "xcsh cloudstatus watch --interval 60",
		AuthenticationRequired: false,
		RecommendedPolling:     60,
		BestPractices: []string{
			"Use 'xcsh cloudstatus status --quiet' in CI/CD to get exit codes (0=healthy)",
			"Use 'xcsh cloudstatus summary --output-format json' for complete state",
			"Use 'xcsh cloudstatus watch' for continuous monitoring",
			"Filter by region with 'xcsh cloudstatus pops list --region <region>'",
			"Check incidents first when troubleshooting connectivity issues",
			"Use --no-cache for real-time status during incidents",
		},
		UseCases: []UseCaseSpec{
			{
				Scenario:    "Pre-deployment health check",
				Command:     "xcsh cloudstatus status --quiet || exit 1",
				Description: "Block deployment if F5 XC has active issues",
			},
			{
				Scenario:    "Incident investigation",
				Command:     "xcsh cloudstatus incidents active --output-format json",
				Description: "Get current incidents affecting services",
			},
			{
				Scenario:    "Regional status check",
				Command:     "xcsh cloudstatus pops status --region north-america",
				Description: "Check status of North America PoPs",
			},
			{
				Scenario:    "Maintenance awareness",
				Command:     "xcsh cloudstatus maintenance upcoming",
				Description: "View upcoming maintenance windows",
			},
		},
	}
}

func getDataModel() DataModel {
	return DataModel{
		Components: []FieldSpec{
			{Field: "id", Type: "string", Description: "Unique component identifier", Example: "ybcpdlwcdq67"},
			{Field: "name", Type: "string", Description: "Human-readable component name", Example: "Portal & Customer Login"},
			{Field: "status", Type: "string", Description: "Current operational status", Example: "operational"},
			{Field: "description", Type: "string", Description: "Component description", Example: "Edge PoP (iad3)"},
			{Field: "group_id", Type: "string", Description: "Parent group identifier", Example: "22ynq9vg49qm"},
			{Field: "group", Type: "boolean", Description: "Whether this is a group container", Example: "false"},
		},
		Incidents: []FieldSpec{
			{Field: "id", Type: "string", Description: "Unique incident identifier", Example: "kcxnsw71xmwp"},
			{Field: "name", Type: "string", Description: "Incident title", Example: "API Latency Issues"},
			{Field: "status", Type: "string", Description: "Current incident status", Example: "investigating"},
			{Field: "impact", Type: "string", Description: "Impact severity", Example: "minor"},
			{Field: "started_at", Type: "timestamp", Description: "When incident started", Example: "2024-01-15T10:30:00Z"},
			{Field: "resolved_at", Type: "timestamp", Description: "When incident was resolved", Example: "2024-01-15T12:45:00Z"},
		},
		Maintenance: []FieldSpec{
			{Field: "id", Type: "string", Description: "Unique maintenance identifier", Example: "xp5l86wjjzyy"},
			{Field: "name", Type: "string", Description: "Maintenance title", Example: "Scheduled Network Upgrade"},
			{Field: "status", Type: "string", Description: "Maintenance status", Example: "scheduled"},
			{Field: "scheduled_for", Type: "timestamp", Description: "Planned start time", Example: "2024-02-01T02:00:00Z"},
			{Field: "scheduled_until", Type: "timestamp", Description: "Planned end time", Example: "2024-02-01T06:00:00Z"},
		},
		Status: []FieldSpec{
			{Field: "indicator", Type: "string", Description: "Overall status indicator", Example: "none"},
			{Field: "description", Type: "string", Description: "Human-readable status", Example: "All Systems Operational"},
		},
	}
}

func getStatusIndicators() StatusIndicators {
	return StatusIndicators{
		Overall: []StatusValueSpec{
			{Value: StatusNone, Description: "All systems operational", Severity: 0, ExitCode: ExitCodeHealthy, Color: "green"},
			{Value: StatusMinor, Description: "Minor system issue", Severity: 1, ExitCode: ExitCodeMinor, Color: "yellow"},
			{Value: StatusMajor, Description: "Major system issue", Severity: 2, ExitCode: ExitCodeMajor, Color: "orange"},
			{Value: StatusCritical, Description: "Critical system outage", Severity: 3, ExitCode: ExitCodeCritical, Color: "red"},
			{Value: StatusMaintenance, Description: "System under maintenance", Severity: 1, ExitCode: ExitCodeMaintenance, Color: "blue"},
		},
		Component: []StatusValueSpec{
			{Value: ComponentOperational, Description: "Component fully operational", Severity: 0},
			{Value: ComponentDegradedPerformance, Description: "Component experiencing degraded performance", Severity: 1},
			{Value: ComponentPartialOutage, Description: "Component has partial outage", Severity: 2},
			{Value: ComponentMajorOutage, Description: "Component has major outage", Severity: 3},
			{Value: ComponentUnderMaintenance, Description: "Component under scheduled maintenance", Severity: 1},
		},
		IncidentImpact: []StatusValueSpec{
			{Value: ImpactNone, Description: "No impact", Severity: 0},
			{Value: ImpactMinor, Description: "Minor impact", Severity: 1},
			{Value: ImpactMajor, Description: "Major impact", Severity: 2},
			{Value: ImpactCritical, Description: "Critical impact", Severity: 3},
		},
		IncidentStatus: []StatusValueSpec{
			{Value: IncidentInvestigating, Description: "Issue under investigation", Severity: 2},
			{Value: IncidentIdentified, Description: "Root cause identified", Severity: 2},
			{Value: IncidentMonitoring, Description: "Fix deployed, monitoring", Severity: 1},
			{Value: IncidentResolved, Description: "Issue resolved", Severity: 0},
			{Value: IncidentPostmortem, Description: "Post-incident review", Severity: 0},
		},
		MaintenanceStatus: []StatusValueSpec{
			{Value: MaintenanceScheduled, Description: "Maintenance scheduled", Severity: 0},
			{Value: MaintenanceInProgress, Description: "Maintenance in progress", Severity: 1},
			{Value: MaintenanceVerifying, Description: "Verifying maintenance completion", Severity: 1},
			{Value: MaintenanceCompleted, Description: "Maintenance completed", Severity: 0},
		},
	}
}

func getEndpoints() []EndpointSpec {
	return []EndpointSpec{
		{Path: "/status.json", Method: "GET", Description: "Overall status indicator", ResponseType: "StatusResponse"},
		{Path: "/summary.json", Method: "GET", Description: "Complete summary with components, incidents, maintenance", ResponseType: "SummaryResponse"},
		{Path: "/components.json", Method: "GET", Description: "All components", ResponseType: "ComponentsResponse"},
		{Path: "/components/{id}.json", Method: "GET", Description: "Single component by ID", ResponseType: "ComponentResponse"},
		{Path: "/incidents.json", Method: "GET", Description: "All incidents", ResponseType: "IncidentsResponse"},
		{Path: "/incidents/unresolved.json", Method: "GET", Description: "Unresolved incidents only", ResponseType: "IncidentsResponse"},
		{Path: "/scheduled-maintenances.json", Method: "GET", Description: "All scheduled maintenance", ResponseType: "MaintenancesResponse"},
		{Path: "/scheduled-maintenances/upcoming.json", Method: "GET", Description: "Upcoming maintenance only", ResponseType: "MaintenancesResponse"},
	}
}

func getCommands() []CommandSpec {
	return []CommandSpec{
		{
			Name:        "status",
			Path:        []string{"cloudstatus", "status"},
			Description: "Get overall F5 Cloud status indicator",
			Flags: []FlagSpec{
				{Name: "quiet", Shorthand: "q", Type: "bool", Description: "Suppress output, return exit code only"},
			},
			ExitCodes: []int{0, 1, 2, 3, 4},
			Examples: []string{
				"xcsh cloudstatus status",
				"xcsh cloudstatus status --quiet && echo 'All systems operational'",
			},
		},
		{
			Name:        "summary",
			Path:        []string{"cloudstatus", "summary"},
			Description: "Get complete status summary including components, incidents, and maintenance",
			Flags: []FlagSpec{
				{Name: "brief", Type: "bool", Description: "Condensed one-liner per section"},
			},
		},
		{
			Name:        "components list",
			Path:        []string{"cloudstatus", "components", "list"},
			Description: "List all components with optional filtering",
			Flags: []FlagSpec{
				{Name: "group", Type: "string", Description: "Filter by group name or ID"},
				{Name: "status", Type: "string", Description: "Filter by status", ValidValues: []string{"operational", "degraded_performance", "partial_outage", "major_outage", "under_maintenance"}},
				{Name: "pop", Type: "bool", Description: "Show PoP components only"},
				{Name: "services", Type: "bool", Description: "Show service components only"},
				{Name: "degraded-only", Type: "bool", Description: "Show non-operational only"},
			},
		},
		{
			Name:        "components get",
			Path:        []string{"cloudstatus", "components", "get"},
			Description: "Get details for a specific component",
			Args:        []string{"component_id_or_name"},
		},
		{
			Name:        "components groups",
			Path:        []string{"cloudstatus", "components", "groups"},
			Description: "List component groups",
			Flags: []FlagSpec{
				{Name: "with-components", Type: "bool", Description: "Include component count per group"},
			},
		},
		{
			Name:        "incidents list",
			Path:        []string{"cloudstatus", "incidents", "list"},
			Description: "List all incidents",
			Flags: []FlagSpec{
				{Name: "status", Type: "string", Description: "Filter by status", ValidValues: []string{"investigating", "identified", "monitoring", "resolved", "postmortem"}},
				{Name: "impact", Type: "string", Description: "Filter by impact level", ValidValues: []string{"none", "minor", "major", "critical"}},
				{Name: "since", Type: "duration", Description: "Time filter (1h, 1d, 7d, 30d)"},
				{Name: "limit", Type: "int", Description: "Limit results"},
			},
		},
		{
			Name:        "incidents active",
			Path:        []string{"cloudstatus", "incidents", "active"},
			Description: "List only unresolved incidents",
		},
		{
			Name:        "incidents get",
			Path:        []string{"cloudstatus", "incidents", "get"},
			Description: "Get details for a specific incident",
			Args:        []string{"incident_id"},
		},
		{
			Name:        "incidents updates",
			Path:        []string{"cloudstatus", "incidents", "updates"},
			Description: "Show incident update timeline",
			Args:        []string{"incident_id"},
		},
		{
			Name:        "maintenance list",
			Path:        []string{"cloudstatus", "maintenance", "list"},
			Description: "List all scheduled maintenance",
			Flags: []FlagSpec{
				{Name: "status", Type: "string", Description: "Filter by status", ValidValues: []string{"scheduled", "in_progress", "verifying", "completed"}},
			},
		},
		{
			Name:        "maintenance upcoming",
			Path:        []string{"cloudstatus", "maintenance", "upcoming"},
			Description: "List upcoming maintenance only",
		},
		{
			Name:        "maintenance active",
			Path:        []string{"cloudstatus", "maintenance", "active"},
			Description: "List in-progress maintenance",
		},
		{
			Name:        "maintenance get",
			Path:        []string{"cloudstatus", "maintenance", "get"},
			Description: "Get details for a specific maintenance window",
			Args:        []string{"maintenance_id"},
		},
		{
			Name:        "pops list",
			Path:        []string{"cloudstatus", "pops", "list"},
			Description: "List PoP locations",
			Flags: []FlagSpec{
				{Name: "region", Type: "string", Description: "Filter by region", ValidValues: []string{"north-america", "south-america", "europe", "asia", "oceania", "middle-east"}},
			},
		},
		{
			Name:        "pops status",
			Path:        []string{"cloudstatus", "pops", "status"},
			Description: "Get aggregated status by region",
			Flags: []FlagSpec{
				{Name: "region", Type: "string", Description: "Filter by region", ValidValues: []string{"north-america", "south-america", "europe", "asia", "oceania", "middle-east"}},
			},
		},
		{
			Name:        "watch",
			Path:        []string{"cloudstatus", "watch"},
			Description: "Real-time status monitoring",
			Flags: []FlagSpec{
				{Name: "interval", Type: "int", Default: "60", Description: "Polling interval in seconds"},
				{Name: "components", Type: "string", Description: "Watch specific components (comma-separated)"},
				{Name: "exit-on-change", Type: "bool", Description: "Exit when status changes"},
				{Name: "no-clear", Type: "bool", Description: "Don't clear screen between updates"},
			},
		},
	}
}

func getWorkflows() []WorkflowSpec {
	return []WorkflowSpec{
		{
			Name:        "ci-cd-health-gate",
			Description: "Block deployment when F5 XC has issues",
			Steps: []WorkflowStep{
				{Step: 1, Command: "xcsh cloudstatus status --quiet", Description: "Check overall status"},
				{Step: 2, Command: "[ $? -eq 0 ] || exit 1", Description: "Exit if not healthy"},
				{Step: 3, Command: "xcsh cloudstatus incidents active --output-format json | jq '. | length'", Description: "Count active incidents"},
				{Step: 4, Command: "# Proceed with deployment if healthy", Description: "Continue deployment"},
			},
		},
		{
			Name:        "incident-investigation",
			Description: "Investigate service issues",
			Steps: []WorkflowStep{
				{Step: 1, Command: "xcsh cloudstatus status", Description: "Quick status check"},
				{Step: 2, Command: "xcsh cloudstatus incidents active", Description: "View active incidents"},
				{Step: 3, Command: "xcsh cloudstatus incidents get <id>", Description: "Get incident details"},
				{Step: 4, Command: "xcsh cloudstatus incidents updates <id>", Description: "View incident timeline"},
			},
		},
		{
			Name:        "regional-status-check",
			Description: "Check status for specific regions",
			Steps: []WorkflowStep{
				{Step: 1, Command: "xcsh cloudstatus pops status", Description: "Overview of all regions"},
				{Step: 2, Command: "xcsh cloudstatus pops list --region north-america", Description: "List North America PoPs"},
				{Step: 3, Command: "xcsh cloudstatus components list --pop --degraded-only", Description: "Find degraded PoPs"},
			},
		},
		{
			Name:        "continuous-monitoring",
			Description: "Set up continuous monitoring",
			Steps: []WorkflowStep{
				{Step: 1, Command: "xcsh cloudstatus watch --interval 30", Description: "Start monitoring"},
				{Step: 2, Command: "xcsh cloudstatus watch --exit-on-change && alert", Description: "Alert on change"},
			},
		},
	}
}

func getExitCodes() []ExitCodeSpec {
	return []ExitCodeSpec{
		{Code: ExitCodeHealthy, Name: "Healthy", Description: "All systems operational (or command success)"},
		{Code: ExitCodeMinor, Name: "Minor", Description: "Minor system issue detected"},
		{Code: ExitCodeMajor, Name: "Major", Description: "Major system issue detected"},
		{Code: ExitCodeCritical, Name: "Critical", Description: "Critical system outage"},
		{Code: ExitCodeMaintenance, Name: "Maintenance", Description: "System under maintenance"},
		{Code: ExitCodeAPIError, Name: "APIError", Description: "Failed to reach status API"},
		{Code: ExitCodeParseError, Name: "ParseError", Description: "Failed to parse API response"},
	}
}

func getFilterCapabilities() FilterCapabilities {
	return FilterCapabilities{
		ComponentFilters: []FilterSpec{
			{Name: "group", Flag: "--group", Type: "string", Description: "Filter by component group", Example: "--group Services"},
			{Name: "status", Flag: "--status", Type: "string", Description: "Filter by status", ValidValues: []string{"operational", "degraded_performance", "partial_outage", "major_outage", "under_maintenance"}, Example: "--status operational"},
			{Name: "pop", Flag: "--pop", Type: "bool", Description: "Show only PoP components", Example: "--pop"},
			{Name: "services", Flag: "--services", Type: "bool", Description: "Show only service components", Example: "--services"},
			{Name: "degraded-only", Flag: "--degraded-only", Type: "bool", Description: "Show only non-operational components", Example: "--degraded-only"},
		},
		IncidentFilters: []FilterSpec{
			{Name: "status", Flag: "--status", Type: "string", Description: "Filter by incident status", ValidValues: []string{"investigating", "identified", "monitoring", "resolved", "postmortem"}, Example: "--status monitoring"},
			{Name: "impact", Flag: "--impact", Type: "string", Description: "Filter by impact level", ValidValues: []string{"none", "minor", "major", "critical"}, Example: "--impact major"},
			{Name: "since", Flag: "--since", Type: "duration", Description: "Filter by time", Example: "--since 7d"},
			{Name: "limit", Flag: "--limit", Type: "int", Description: "Limit results", Example: "--limit 10"},
		},
		MaintenanceFilters: []FilterSpec{
			{Name: "status", Flag: "--status", Type: "string", Description: "Filter by maintenance status", ValidValues: []string{"scheduled", "in_progress", "verifying", "completed"}, Example: "--status scheduled"},
		},
		PopFilters: []FilterSpec{
			{Name: "region", Flag: "--region", Type: "string", Description: "Filter by region", ValidValues: []string{"north-america", "south-america", "europe", "asia", "oceania", "middle-east"}, Example: "--region north-america"},
		},
		TimeFilters: []FilterSpec{
			{Name: "since", Flag: "--since", Type: "duration", Description: "Time range filter", Example: "--since 24h"},
		},
	}
}

func getCachingSpec() CachingSpec {
	return CachingSpec{
		DefaultTTLSeconds: 60,
		BypassFlag:        "--no-cache",
		TTLOverrideFlag:   "--cache-ttl",
		CacheLocation:     "in-memory",
	}
}

func getComponentGroups() []ComponentGroupSpec {
	return []ComponentGroupSpec{
		{ID: "22ynq9vg49qm", Name: "Services", Description: "Core F5 XC services", ComponentCount: 16},
		{ID: "kwm9bq7w9v22", Name: "North America PoPs", Description: "Edge PoPs in North America", ComponentCount: 14},
		{ID: "9dyxpcltwdx7", Name: "Europe PoPs", Description: "Edge PoPs in Europe", ComponentCount: 11},
		{ID: "b5yxhx1rxzby", Name: "Asia PoPs", Description: "Edge PoPs in Asia", ComponentCount: 7},
		{ID: "4xcms7l3k7kv", Name: "South America PoPs", Description: "Edge PoPs in South America", ComponentCount: 1},
		{ID: "dvqs5zbrxrrs", Name: "Oceania PoPs", Description: "Edge PoPs in Oceania", ComponentCount: 2},
		{ID: "mfrs7kmq9k8y", Name: "Middle East PoPs", Description: "Edge PoPs in Middle East", ComponentCount: 2},
	}
}

func getRegions() []RegionSpec {
	return []RegionSpec{
		{Name: "north-america", DisplayName: "North America", PopCount: 14},
		{Name: "south-america", DisplayName: "South America", PopCount: 1},
		{Name: "europe", DisplayName: "Europe", PopCount: 11},
		{Name: "asia", DisplayName: "Asia", PopCount: 7},
		{Name: "oceania", DisplayName: "Oceania", PopCount: 2},
		{Name: "middle-east", DisplayName: "Middle East", PopCount: 2},
	}
}
