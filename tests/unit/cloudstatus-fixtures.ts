/**
 * Cloudstatus Test Fixtures
 *
 * Mock API response data matching the Atlassian Statuspage v2 API format
 * used by F5 Cloud Status (https://www.f5cloudstatus.com/api/v2)
 */

import type {
	StatusResponse,
	SummaryResponse,
	ComponentsResponse,
	IncidentsResponse,
	MaintenancesResponse,
	Component,
	Incident,
	ScheduledMaintenance,
	PageInfo,
	Status,
} from "../../src/cloudstatus/types.js";

// Common page info for all responses
export const mockPageInfo: PageInfo = {
	id: "test-page-id",
	name: "F5 Distributed Cloud Status",
	url: "https://www.f5cloudstatus.com",
	time_zone: "UTC",
	updated_at: "2024-01-15T12:00:00Z",
};

// Status indicators for different scenarios
export const mockStatusHealthy: Status = {
	indicator: "none",
	description: "All Systems Operational",
};

export const mockStatusMinor: Status = {
	indicator: "minor",
	description: "Minor System Issue",
};

export const mockStatusMajor: Status = {
	indicator: "major",
	description: "Major System Issue",
};

export const mockStatusCritical: Status = {
	indicator: "critical",
	description: "Critical System Outage",
};

export const mockStatusMaintenance: Status = {
	indicator: "maintenance",
	description: "System Under Maintenance",
};

// StatusResponse fixtures
export const mockStatusResponse: StatusResponse = {
	page: mockPageInfo,
	status: mockStatusHealthy,
};

export const mockStatusResponseMinor: StatusResponse = {
	page: mockPageInfo,
	status: mockStatusMinor,
};

export const mockStatusResponseMajor: StatusResponse = {
	page: mockPageInfo,
	status: mockStatusMajor,
};

export const mockStatusResponseCritical: StatusResponse = {
	page: mockPageInfo,
	status: mockStatusCritical,
};

export const mockStatusResponseMaintenance: StatusResponse = {
	page: mockPageInfo,
	status: mockStatusMaintenance,
};

// Component fixtures
export const mockOperationalComponent: Component = {
	id: "comp-1",
	name: "F5 XC Console",
	status: "operational",
	description: "Main console application",
	group_id: null,
	group: false,
	position: 1,
	showcase: true,
	only_show_if_degraded: false,
	page_id: "test-page-id",
	start_date: null,
	created_at: "2023-01-01T00:00:00Z",
	updated_at: "2024-01-15T12:00:00Z",
};

export const mockDegradedComponent: Component = {
	id: "comp-2",
	name: "API Gateway",
	status: "degraded_performance",
	description: "API gateway service",
	group_id: null,
	group: false,
	position: 2,
	showcase: true,
	only_show_if_degraded: false,
	page_id: "test-page-id",
	start_date: null,
	created_at: "2023-01-01T00:00:00Z",
	updated_at: "2024-01-15T12:00:00Z",
};

export const mockPartialOutageComponent: Component = {
	id: "comp-3",
	name: "Load Balancer",
	status: "partial_outage",
	description: "Load balancing service",
	group_id: null,
	group: false,
	position: 3,
	showcase: true,
	only_show_if_degraded: false,
	page_id: "test-page-id",
	start_date: null,
	created_at: "2023-01-01T00:00:00Z",
	updated_at: "2024-01-15T12:00:00Z",
};

export const mockGroupComponent: Component = {
	id: "group-1",
	name: "Core Services",
	status: "operational",
	description: "Group of core services",
	group_id: null,
	group: true,
	components: ["comp-1", "comp-2"],
	position: 0,
	showcase: true,
	only_show_if_degraded: false,
	page_id: "test-page-id",
	start_date: null,
	created_at: "2023-01-01T00:00:00Z",
	updated_at: "2024-01-15T12:00:00Z",
};

// ComponentsResponse fixture
export const mockComponentsResponse: ComponentsResponse = {
	page: mockPageInfo,
	components: [
		mockGroupComponent,
		mockOperationalComponent,
		mockDegradedComponent,
		mockPartialOutageComponent,
	],
};

// All operational components for testing
export const mockComponentsAllOperational: ComponentsResponse = {
	page: mockPageInfo,
	components: [
		mockGroupComponent,
		mockOperationalComponent,
		{ ...mockOperationalComponent, id: "comp-4", name: "DNS Service" },
	],
};

// Incident fixtures
export const mockActiveIncident: Incident = {
	id: "inc-1",
	name: "API Latency Issues",
	status: "investigating",
	impact: "minor",
	shortlink: "https://stspg.io/abc123",
	page_id: "test-page-id",
	created_at: "2024-01-15T10:00:00Z",
	updated_at: "2024-01-15T11:00:00Z",
	started_at: "2024-01-15T10:00:00Z",
	monitoring_at: null,
	resolved_at: null,
	incident_updates: [
		{
			id: "update-1",
			status: "investigating",
			body: "We are currently investigating increased API response times in the US-East region. Our team is actively working to identify the root cause.",
			incident_id: "inc-1",
			created_at: "2024-01-15T11:00:00Z",
			updated_at: "2024-01-15T11:00:00Z",
			display_at: "2024-01-15T11:00:00Z",
			affected_components: [],
			deliver_notifications: true,
		},
	],
	components: [mockDegradedComponent],
};

export const mockResolvedIncident: Incident = {
	id: "inc-2",
	name: "Console Login Issues",
	status: "resolved",
	impact: "major",
	shortlink: "https://stspg.io/def456",
	page_id: "test-page-id",
	created_at: "2024-01-14T08:00:00Z",
	updated_at: "2024-01-14T12:00:00Z",
	started_at: "2024-01-14T08:00:00Z",
	monitoring_at: "2024-01-14T11:00:00Z",
	resolved_at: "2024-01-14T12:00:00Z",
	incident_updates: [
		{
			id: "update-2",
			status: "resolved",
			body: "This incident has been resolved. Login functionality has been restored.",
			incident_id: "inc-2",
			created_at: "2024-01-14T12:00:00Z",
			updated_at: "2024-01-14T12:00:00Z",
			display_at: "2024-01-14T12:00:00Z",
			affected_components: [],
			deliver_notifications: true,
		},
	],
	components: [mockOperationalComponent],
};

// IncidentsResponse fixtures
export const mockIncidentsResponse: IncidentsResponse = {
	page: mockPageInfo,
	incidents: [mockActiveIncident, mockResolvedIncident],
};

export const mockIncidentsActiveOnly: IncidentsResponse = {
	page: mockPageInfo,
	incidents: [mockActiveIncident],
};

export const mockIncidentsEmpty: IncidentsResponse = {
	page: mockPageInfo,
	incidents: [],
};

// Maintenance fixtures
export const mockScheduledMaintenance: ScheduledMaintenance = {
	id: "maint-1",
	name: "Planned Infrastructure Upgrade",
	status: "scheduled",
	impact: "minor",
	shortlink: "https://stspg.io/maint123",
	page_id: "test-page-id",
	scheduled_for: "2024-01-20T02:00:00Z",
	scheduled_until: "2024-01-20T06:00:00Z",
	created_at: "2024-01-10T00:00:00Z",
	updated_at: "2024-01-10T00:00:00Z",
	started_at: null,
	monitoring_at: null,
	resolved_at: null,
	incident_updates: [],
	components: [mockOperationalComponent],
};

export const mockInProgressMaintenance: ScheduledMaintenance = {
	id: "maint-2",
	name: "Database Migration",
	status: "in_progress",
	impact: "major",
	shortlink: "https://stspg.io/maint456",
	page_id: "test-page-id",
	scheduled_for: "2024-01-15T02:00:00Z",
	scheduled_until: "2024-01-15T06:00:00Z",
	created_at: "2024-01-08T00:00:00Z",
	updated_at: "2024-01-15T02:30:00Z",
	started_at: "2024-01-15T02:00:00Z",
	monitoring_at: null,
	resolved_at: null,
	incident_updates: [],
	components: [mockDegradedComponent],
};

export const mockCompletedMaintenance: ScheduledMaintenance = {
	id: "maint-3",
	name: "Security Patch Deployment",
	status: "completed",
	impact: "none",
	shortlink: "https://stspg.io/maint789",
	page_id: "test-page-id",
	scheduled_for: "2024-01-12T02:00:00Z",
	scheduled_until: "2024-01-12T04:00:00Z",
	created_at: "2024-01-05T00:00:00Z",
	updated_at: "2024-01-12T04:00:00Z",
	started_at: "2024-01-12T02:00:00Z",
	monitoring_at: "2024-01-12T03:30:00Z",
	resolved_at: "2024-01-12T04:00:00Z",
	incident_updates: [],
	components: [mockOperationalComponent],
};

// MaintenancesResponse fixtures
export const mockMaintenancesResponse: MaintenancesResponse = {
	page: mockPageInfo,
	scheduled_maintenances: [
		mockScheduledMaintenance,
		mockInProgressMaintenance,
		mockCompletedMaintenance,
	],
};

export const mockMaintenancesUpcomingOnly: MaintenancesResponse = {
	page: mockPageInfo,
	scheduled_maintenances: [mockScheduledMaintenance],
};

export const mockMaintenancesEmpty: MaintenancesResponse = {
	page: mockPageInfo,
	scheduled_maintenances: [],
};

// SummaryResponse fixture (combines all)
export const mockSummaryResponse: SummaryResponse = {
	page: mockPageInfo,
	status: mockStatusHealthy,
	components: [
		mockGroupComponent,
		mockOperationalComponent,
		mockDegradedComponent,
	],
	incidents: [mockActiveIncident, mockResolvedIncident],
	scheduled_maintenances: [
		mockScheduledMaintenance,
		mockInProgressMaintenance,
	],
};

export const mockSummaryResponseAllClear: SummaryResponse = {
	page: mockPageInfo,
	status: mockStatusHealthy,
	components: [mockGroupComponent, mockOperationalComponent],
	incidents: [],
	scheduled_maintenances: [],
};
