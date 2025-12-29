/**
 * Cloudstatus Domain - F5 Distributed Cloud service status monitoring
 *
 * Monitor service health, incidents, and maintenance windows.
 * No authentication required - the status API is publicly accessible.
 */

import type {
	DomainDefinition,
	CommandDefinition,
	DomainCommandResult,
} from "../registry.js";
import { successResult, errorResult } from "../registry.js";
import {
	CloudstatusClient,
	isIncidentActive,
	isMaintenanceActive,
	isMaintenanceUpcoming,
	isMaintenanceCompleted,
	isComponentDegraded,
	isComponentOperational,
	statusIndicatorToExitCode,
} from "../../cloudstatus/index.js";
import type { SummaryResponse } from "../../cloudstatus/index.js";

// Lazy-initialized cloudstatus client
let cloudstatusClient: CloudstatusClient | null = null;

function getClient(): CloudstatusClient {
	if (!cloudstatusClient) {
		cloudstatusClient = new CloudstatusClient();
	}
	return cloudstatusClient;
}

/**
 * Status command - Get overall status indicator
 */
const statusCommand: CommandDefinition = {
	name: "status",
	description: "Get overall F5 Cloud status indicator",
	usage: "[--quiet]",
	aliases: ["st"],

	async execute(args, session): Promise<DomainCommandResult> {
		const quiet = args.includes("--quiet") || args.includes("-q");
		const format = session.getOutputFormat();

		try {
			const client = getClient();
			const response = await client.getStatus();

			// In quiet mode, return exit code info
			if (quiet) {
				const exitCode = statusIndicatorToExitCode(
					response.status.indicator,
				);
				return successResult([
					`Exit code: ${exitCode} (${response.status.indicator})`,
				]);
			}

			// Format based on output format
			if (format === "json") {
				return successResult([JSON.stringify(response, null, 2)]);
			}

			if (format === "yaml") {
				// Simple YAML formatting
				const yaml = [
					"page:",
					`  id: ${response.page.id}`,
					`  name: ${response.page.name}`,
					`  url: ${response.page.url}`,
					`  updated_at: ${response.page.updated_at}`,
					"status:",
					`  indicator: ${response.status.indicator}`,
					`  description: ${response.status.description}`,
				].join("\n");
				return successResult([yaml]);
			}

			// Table format (default)
			const lines: string[] = [
				"┌────────────────┬────────────────────────────────────┐",
				"│ STATUS         │ DESCRIPTION                        │",
				"├────────────────┼────────────────────────────────────┤",
				`│ ${padEnd(response.status.indicator, 14)} │ ${padEnd(response.status.description, 34)} │`,
				"└────────────────┴────────────────────────────────────┘",
			];

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to get status: ${message}`);
		}
	},
};

/**
 * Summary command - Get complete status summary
 */
const summaryCommand: CommandDefinition = {
	name: "summary",
	description:
		"Get complete status summary including components and incidents",
	usage: "[--brief]",
	aliases: ["sum"],

	async execute(args, session): Promise<DomainCommandResult> {
		const brief = args.includes("--brief") || args.includes("-b");
		const format = session.getOutputFormat();

		try {
			const client = getClient();
			const response = await client.getSummary();

			// Format based on output format
			if (format === "json") {
				return successResult([JSON.stringify(response, null, 2)]);
			}

			if (format === "yaml") {
				// Return structured YAML
				return successResult([formatSummaryYaml(response)]);
			}

			// Brief or full summary
			if (brief) {
				return successResult(formatBriefSummary(response));
			}

			return successResult(formatFullSummary(response));
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to get summary: ${message}`);
		}
	},
};

/**
 * Components command - List all components
 */
const componentsCommand: CommandDefinition = {
	name: "components",
	description: "List all components and their status",
	usage: "[--degraded-only]",
	aliases: ["comp"],

	async execute(args, session): Promise<DomainCommandResult> {
		const degradedOnly =
			args.includes("--degraded-only") || args.includes("-d");
		const format = session.getOutputFormat();

		try {
			const client = getClient();
			const response = await client.getComponents();

			let components = response.components.filter((c) => !c.group);
			if (degradedOnly) {
				components = components.filter((c) => isComponentDegraded(c));
			}

			if (format === "json") {
				return successResult([JSON.stringify(components, null, 2)]);
			}

			if (components.length === 0) {
				return successResult(["No components found."]);
			}

			// Table format
			const lines: string[] = [
				"┌─────────────────────────────────────────┬────────────────────┐",
				"│ NAME                                    │ STATUS             │",
				"├─────────────────────────────────────────┼────────────────────┤",
			];

			for (const comp of components) {
				lines.push(
					`│ ${padEnd(comp.name, 39)} │ ${padEnd(comp.status, 18)} │`,
				);
			}

			lines.push(
				"└─────────────────────────────────────────┴────────────────────┘",
			);

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to get components: ${message}`);
		}
	},
};

/**
 * Incidents command - List incidents
 */
const incidentsCommand: CommandDefinition = {
	name: "incidents",
	description: "List active and recent incidents",
	usage: "[--active-only]",
	aliases: ["inc"],

	async execute(args, session): Promise<DomainCommandResult> {
		const activeOnly =
			args.includes("--active-only") || args.includes("-a");
		const format = session.getOutputFormat();

		try {
			const client = getClient();
			const response = activeOnly
				? await client.getUnresolvedIncidents()
				: await client.getIncidents();

			if (format === "json") {
				return successResult([
					JSON.stringify(response.incidents, null, 2),
				]);
			}

			if (response.incidents.length === 0) {
				return successResult(["No incidents found."]);
			}

			const lines: string[] = [];
			for (const inc of response.incidents) {
				const status = isIncidentActive(inc)
					? "[ACTIVE]"
					: "[RESOLVED]";
				lines.push(`${status} ${inc.name} (Impact: ${inc.impact})`);
				lines.push(`  Started: ${formatDate(inc.started_at)}`);
				const latestUpdate = inc.incident_updates[0];
				if (latestUpdate) {
					lines.push(
						`  Latest: ${latestUpdate.body.slice(0, 80)}...`,
					);
				}
				lines.push("");
			}

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to get incidents: ${message}`);
		}
	},
};

/**
 * Maintenance command - List scheduled maintenances
 */
const maintenanceCommand: CommandDefinition = {
	name: "maintenance",
	description: "List scheduled and active maintenance windows",
	usage: "[--upcoming]",
	aliases: ["maint"],

	async execute(args, session): Promise<DomainCommandResult> {
		const upcomingOnly = args.includes("--upcoming") || args.includes("-u");
		const format = session.getOutputFormat();

		try {
			const client = getClient();
			const response = upcomingOnly
				? await client.getUpcomingMaintenances()
				: await client.getMaintenances();

			if (format === "json") {
				return successResult([
					JSON.stringify(response.scheduled_maintenances, null, 2),
				]);
			}

			if (response.scheduled_maintenances.length === 0) {
				return successResult(["No scheduled maintenance."]);
			}

			const lines: string[] = [];
			for (const maint of response.scheduled_maintenances) {
				if (!isMaintenanceCompleted(maint)) {
					const status = isMaintenanceActive(maint)
						? "[IN PROGRESS]"
						: "[SCHEDULED]";
					lines.push(`${status} ${maint.name}`);
					lines.push(
						`  Scheduled: ${formatDate(maint.scheduled_for)} to ${formatDate(maint.scheduled_until)}`,
					);
					lines.push("");
				}
			}

			if (lines.length === 0) {
				return successResult(["No upcoming maintenance."]);
			}

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to get maintenance: ${message}`);
		}
	},
};

// Helper functions

function padEnd(str: string, length: number): string {
	if (str.length >= length) {
		return str.slice(0, length);
	}
	return str + " ".repeat(length - str.length);
}

function formatDate(dateStr: string): string {
	try {
		const date = new Date(dateStr);
		return date.toLocaleString();
	} catch {
		return dateStr;
	}
}

function formatBriefSummary(summary: SummaryResponse): string[] {
	// Count statistics
	let operationalCount = 0;
	let degradedCount = 0;
	for (const comp of summary.components) {
		if (!comp.group) {
			if (isComponentOperational(comp)) {
				operationalCount++;
			} else {
				degradedCount++;
			}
		}
	}

	let activeIncidents = 0;
	for (const inc of summary.incidents) {
		if (isIncidentActive(inc)) {
			activeIncidents++;
		}
	}

	let upcomingMaint = 0;
	let activeMaint = 0;
	for (const maint of summary.scheduled_maintenances) {
		if (isMaintenanceUpcoming(maint)) {
			upcomingMaint++;
		} else if (isMaintenanceActive(maint)) {
			activeMaint++;
		}
	}

	return [
		`Status: ${summary.status.indicator} (${summary.status.description})`,
		`Components: ${operationalCount} operational, ${degradedCount} degraded`,
		`Incidents: ${activeIncidents} active`,
		`Maintenance: ${upcomingMaint} upcoming, ${activeMaint} active`,
	];
}

function formatFullSummary(summary: SummaryResponse): string[] {
	const lines: string[] = [];

	// Overall status
	lines.push("=== OVERALL STATUS ===");
	lines.push(`Indicator: ${summary.status.indicator}`);
	lines.push(`Description: ${summary.status.description}`);
	lines.push(`Last Updated: ${formatDate(summary.page.updated_at)}`);
	lines.push("");

	// Components summary
	lines.push("=== COMPONENTS ===");

	// Show degraded components first
	const degraded = summary.components.filter(
		(c) => !c.group && isComponentDegraded(c),
	);
	if (degraded.length > 0) {
		for (const comp of degraded) {
			lines.push(`  ${comp.name}: ${comp.status}`);
		}
	}

	// Count operational
	const operationalCount = summary.components.filter(
		(c) => !c.group && isComponentOperational(c),
	).length;
	lines.push("");
	lines.push(`(${operationalCount} components operational, not shown)`);
	lines.push("");

	// Active incidents
	lines.push("=== ACTIVE INCIDENTS ===");
	const activeIncidents = summary.incidents.filter((i) =>
		isIncidentActive(i),
	);
	if (activeIncidents.length > 0) {
		for (const inc of activeIncidents) {
			lines.push(`[${inc.status}] ${inc.name} (Impact: ${inc.impact})`);
			lines.push(`  Started: ${formatDate(inc.started_at)}`);
			const latestUpdate = inc.incident_updates[0];
			if (latestUpdate) {
				lines.push(`  Latest: ${latestUpdate.body.slice(0, 80)}...`);
			}
			lines.push("");
		}
	} else {
		lines.push("No active incidents");
		lines.push("");
	}

	// Scheduled maintenance
	lines.push("=== SCHEDULED MAINTENANCE ===");
	const activeMaint = summary.scheduled_maintenances.filter(
		(m) => !isMaintenanceCompleted(m),
	);
	if (activeMaint.length > 0) {
		for (const maint of activeMaint) {
			lines.push(`[${maint.status}] ${maint.name}`);
			lines.push(
				`  Scheduled: ${formatDate(maint.scheduled_for)} to ${formatDate(maint.scheduled_until)}`,
			);
			lines.push("");
		}
	} else {
		lines.push("No scheduled maintenance");
	}

	return lines;
}

function formatSummaryYaml(summary: SummaryResponse): string {
	const lines: string[] = [];
	lines.push("status:");
	lines.push(`  indicator: ${summary.status.indicator}`);
	lines.push(`  description: ${summary.status.description}`);
	lines.push("page:");
	lines.push(`  id: ${summary.page.id}`);
	lines.push(`  name: ${summary.page.name}`);
	lines.push(`  updated_at: ${summary.page.updated_at}`);
	lines.push(`components_count: ${summary.components.length}`);
	lines.push(`incidents_count: ${summary.incidents.length}`);
	lines.push(`maintenances_count: ${summary.scheduled_maintenances.length}`);
	return lines.join("\n");
}

/**
 * Cloudstatus domain definition
 */
export const cloudstatusDomain: DomainDefinition = {
	name: "cloudstatus",
	description:
		"Monitor F5 Distributed Cloud service status and incidents. Check overall status indicators, view component health, track active incidents and their updates, and monitor scheduled maintenance windows.",
	descriptionShort: "F5 XC service status and incidents",
	descriptionMedium:
		"Monitor F5 Distributed Cloud service health, active incidents, component status, and scheduled maintenance windows.",
	commands: new Map([
		["status", statusCommand],
		["summary", summaryCommand],
		["components", componentsCommand],
		["incidents", incidentsCommand],
		["maintenance", maintenanceCommand],
	]),
	subcommands: new Map(),
};

// Aliases for the domain
export const cloudstatusAliases = ["cs", "status"];
