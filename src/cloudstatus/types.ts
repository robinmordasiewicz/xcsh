/**
 * Cloudstatus Types - F5 Cloud Status API type definitions
 *
 * Based on Atlassian Statuspage v2 API format
 */

// Base URL for the F5 Cloud Status API
export const BASE_URL = "https://www.f5cloudstatus.com/api/v2";

// Status indicator values
export const StatusIndicator = {
	None: "none",
	Minor: "minor",
	Major: "major",
	Critical: "critical",
	Maintenance: "maintenance",
} as const;

export type StatusIndicatorType =
	(typeof StatusIndicator)[keyof typeof StatusIndicator];

// Component status values
export const ComponentStatus = {
	Operational: "operational",
	DegradedPerformance: "degraded_performance",
	PartialOutage: "partial_outage",
	MajorOutage: "major_outage",
	UnderMaintenance: "under_maintenance",
} as const;

export type ComponentStatusType =
	(typeof ComponentStatus)[keyof typeof ComponentStatus];

// Incident status values
export const IncidentStatus = {
	Investigating: "investigating",
	Identified: "identified",
	Monitoring: "monitoring",
	Resolved: "resolved",
	Postmortem: "postmortem",
} as const;

export type IncidentStatusType =
	(typeof IncidentStatus)[keyof typeof IncidentStatus];

// Impact values
export const Impact = {
	None: "none",
	Minor: "minor",
	Major: "major",
	Critical: "critical",
} as const;

export type ImpactType = (typeof Impact)[keyof typeof Impact];

// Maintenance status values
export const MaintenanceStatus = {
	Scheduled: "scheduled",
	InProgress: "in_progress",
	Verifying: "verifying",
	Completed: "completed",
} as const;

export type MaintenanceStatusType =
	(typeof MaintenanceStatus)[keyof typeof MaintenanceStatus];

// Exit codes for CI/CD integration
export const ExitCode = {
	Healthy: 0,
	Minor: 1,
	Major: 2,
	Critical: 3,
	Maintenance: 4,
	APIError: 10,
	ParseError: 11,
} as const;

export type ExitCodeType = (typeof ExitCode)[keyof typeof ExitCode];

/**
 * Helper string for help text generation
 * Automatically derived from ExitCode constant (status codes only, excludes error codes)
 */
export const EXIT_CODE_HELP = Object.entries(ExitCode)
	.filter(([, v]) => typeof v === "number" && v < 10) // Exclude error codes
	.map(([k, v]) => `${v}=${k.toLowerCase()}`)
	.join(", ");

// Page info metadata
export interface PageInfo {
	id: string;
	name: string;
	url: string;
	time_zone: string;
	updated_at: string; // ISO 8601 date string
}

// Status indicator
export interface Status {
	indicator: StatusIndicatorType;
	description: string;
}

// Component on the status page
export interface Component {
	id: string;
	name: string;
	status: ComponentStatusType;
	description: string;
	group_id: string | null;
	group: boolean;
	components?: string[];
	position: number;
	showcase: boolean;
	only_show_if_degraded: boolean;
	page_id: string;
	start_date: string | null;
	created_at: string;
	updated_at: string;
}

// Affected component in an incident update
export interface AffectedComponent {
	code: string;
	name: string;
	old_status: string;
	new_status: string;
}

// Incident update
export interface IncidentUpdate {
	id: string;
	status: IncidentStatusType;
	body: string;
	incident_id: string;
	created_at: string;
	updated_at: string;
	display_at: string;
	affected_components: AffectedComponent[];
	deliver_notifications: boolean;
}

// Incident on the status page
export interface Incident {
	id: string;
	name: string;
	status: IncidentStatusType;
	impact: ImpactType;
	shortlink: string;
	page_id: string;
	created_at: string;
	updated_at: string;
	started_at: string;
	monitoring_at: string | null;
	resolved_at: string | null;
	incident_updates: IncidentUpdate[];
	components: Component[];
}

// Scheduled maintenance window
export interface ScheduledMaintenance {
	id: string;
	name: string;
	status: MaintenanceStatusType;
	impact: ImpactType;
	shortlink: string;
	page_id: string;
	scheduled_for: string;
	scheduled_until: string;
	created_at: string;
	updated_at: string;
	started_at: string | null;
	monitoring_at: string | null;
	resolved_at: string | null;
	incident_updates: IncidentUpdate[];
	components: Component[];
}

// API response types
export interface StatusResponse {
	page: PageInfo;
	status: Status;
}

export interface SummaryResponse {
	page: PageInfo;
	components: Component[];
	incidents: Incident[];
	scheduled_maintenances: ScheduledMaintenance[];
	status: Status;
}

export interface ComponentsResponse {
	page: PageInfo;
	components: Component[];
}

export interface ComponentResponse {
	page: PageInfo;
	component: Component;
}

export interface IncidentsResponse {
	page: PageInfo;
	incidents: Incident[];
}

export interface MaintenancesResponse {
	page: PageInfo;
	scheduled_maintenances: ScheduledMaintenance[];
}

// Component group for organization
export interface ComponentGroup {
	id: string;
	name: string;
	description: string;
	components: Component[];
	componentCount: number;
}

// Region for geographic grouping
export interface Region {
	id: string;
	name: string;
	displayName: string;
}

// Regional status aggregation
export interface RegionalStatus {
	region: Region;
	overallStatus: StatusIndicatorType;
	operationalCount: number;
	degradedCount: number;
	totalCount: number;
	components: Component[];
}

// Predefined regions for F5 XC PoPs
export const PredefinedRegions: Region[] = [
	{
		id: "north-america",
		name: "north-america",
		displayName: "North America",
	},
	{
		id: "south-america",
		name: "south-america",
		displayName: "South America",
	},
	{ id: "europe", name: "europe", displayName: "Europe" },
	{ id: "asia", name: "asia", displayName: "Asia" },
	{ id: "oceania", name: "oceania", displayName: "Oceania" },
	{ id: "middle-east", name: "middle-east", displayName: "Middle East" },
];

// Helper functions

export function isComponentOperational(component: Component): boolean {
	return component.status === ComponentStatus.Operational;
}

export function isComponentDegraded(component: Component): boolean {
	return component.status !== ComponentStatus.Operational;
}

export function isIncidentResolved(incident: Incident): boolean {
	return (
		incident.status === IncidentStatus.Resolved ||
		incident.status === IncidentStatus.Postmortem
	);
}

export function isIncidentActive(incident: Incident): boolean {
	return !isIncidentResolved(incident);
}

export function isMaintenanceUpcoming(maint: ScheduledMaintenance): boolean {
	return maint.status === MaintenanceStatus.Scheduled;
}

export function isMaintenanceActive(maint: ScheduledMaintenance): boolean {
	return (
		maint.status === MaintenanceStatus.InProgress ||
		maint.status === MaintenanceStatus.Verifying
	);
}

export function isMaintenanceCompleted(maint: ScheduledMaintenance): boolean {
	return maint.status === MaintenanceStatus.Completed;
}

export function statusIndicatorToExitCode(
	indicator: StatusIndicatorType,
): number {
	switch (indicator) {
		case StatusIndicator.None:
			return ExitCode.Healthy;
		case StatusIndicator.Minor:
			return ExitCode.Minor;
		case StatusIndicator.Major:
			return ExitCode.Major;
		case StatusIndicator.Critical:
			return ExitCode.Critical;
		case StatusIndicator.Maintenance:
			return ExitCode.Maintenance;
		default:
			return ExitCode.Healthy;
	}
}

export function statusIndicatorDescription(indicator: string): string {
	switch (indicator) {
		case StatusIndicator.None:
			return "All Systems Operational";
		case StatusIndicator.Minor:
			return "Minor System Issue";
		case StatusIndicator.Major:
			return "Major System Issue";
		case StatusIndicator.Critical:
			return "Critical System Outage";
		case StatusIndicator.Maintenance:
			return "System Under Maintenance";
		default:
			return "Unknown Status";
	}
}

export function componentStatusDescription(status: string): string {
	switch (status) {
		case ComponentStatus.Operational:
			return "Operational";
		case ComponentStatus.DegradedPerformance:
			return "Degraded Performance";
		case ComponentStatus.PartialOutage:
			return "Partial Outage";
		case ComponentStatus.MajorOutage:
			return "Major Outage";
		case ComponentStatus.UnderMaintenance:
			return "Under Maintenance";
		default:
			return "Unknown";
	}
}
