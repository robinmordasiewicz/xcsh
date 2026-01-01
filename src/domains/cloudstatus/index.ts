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
import {
	type OutputFormat,
	getCommandSpec,
	formatSpec,
} from "../../output/index.js";
import {
	parseDomainOutputFlags,
	formatDomainOutput,
	formatListOutput,
} from "../../output/domain-formatter.js";
import type { REPLSession } from "../../repl/session.js";

/**
 * Parse output format, spec flag, and command-specific flags from args
 */
function parseOutputArgs(
	args: string[],
	session: REPLSession,
): {
	format: OutputFormat;
	spec: boolean;
	noColor: boolean;
	filteredArgs: string[];
} {
	// Use unified domain formatter for output flags
	const { options, remainingArgs } = parseDomainOutputFlags(
		args,
		session.getOutputFormat(),
	);

	// Check for --spec flag separately
	let spec = false;
	const filteredArgs: string[] = [];
	for (const arg of remainingArgs) {
		if (arg === "--spec") {
			spec = true;
		} else {
			filteredArgs.push(arg);
		}
	}

	return {
		format: options.format,
		noColor: options.noColor,
		spec,
		filteredArgs,
	};
}

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
	description:
		"Retrieve the current overall health indicator for F5 Distributed Cloud services. Returns status level (operational, degraded, major outage) with description. Use --quiet for script-friendly exit code output.",
	descriptionShort: "Get overall cloud status indicator",
	descriptionMedium:
		"Check overall service health status. Use --quiet for exit code suitable for scripts.",
	usage: "[--quiet]",
	aliases: ["st"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec, filteredArgs } = parseOutputArgs(
			args,
			session,
		);
		const quiet =
			filteredArgs.includes("--quiet") || filteredArgs.includes("-q");

		// Handle --spec flag: return command specification for AI assistants
		if (spec) {
			const cmdSpec = getCommandSpec("cloudstatus status");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

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

			// Handle none format - return empty
			if (format === "none") {
				return successResult([]);
			}

			// Use unified formatter for json/yaml/tsv
			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatDomainOutput(response, { format, noColor }),
				);
			}

			// Table format (default) - use custom formatting for visual display
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
		"Display comprehensive status overview combining overall health, component statuses, active incidents, and scheduled maintenance in a single report. Use --brief for condensed statistics output.",
	descriptionShort: "Get complete status summary",
	descriptionMedium:
		"Show combined overview of health, components, incidents, and maintenance windows.",
	usage: "[--brief]",
	aliases: ["sum"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec, filteredArgs } = parseOutputArgs(
			args,
			session,
		);
		const brief =
			filteredArgs.includes("--brief") || filteredArgs.includes("-b");

		// Handle --spec flag: return command specification for AI assistants
		if (spec) {
			const cmdSpec = getCommandSpec("cloudstatus summary");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		try {
			const client = getClient();
			const response = await client.getSummary();

			// Handle none format - return empty
			if (format === "none") {
				return successResult([]);
			}

			// Use unified formatter for json/yaml/tsv
			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatDomainOutput(response, { format, noColor }),
				);
			}

			// Brief or full summary (table format)
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
	description:
		"List all service components with their current operational status. Shows each component's health level. Use --degraded-only to filter for components experiencing issues.",
	descriptionShort: "List all components and status",
	descriptionMedium:
		"Display service component health. Use --degraded-only to show only affected components.",
	usage: "[--degraded-only]",
	aliases: ["comp"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec, filteredArgs } = parseOutputArgs(
			args,
			session,
		);
		const degradedOnly =
			filteredArgs.includes("--degraded-only") ||
			filteredArgs.includes("-d");

		// Handle --spec flag: return command specification for AI assistants
		if (spec) {
			const cmdSpec = getCommandSpec("cloudstatus components");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		try {
			const client = getClient();
			const response = await client.getComponents();

			let components = response.components.filter((c) => !c.group);
			if (degradedOnly) {
				components = components.filter((c) => isComponentDegraded(c));
			}

			// Handle none format - return empty
			if (format === "none") {
				return successResult([]);
			}

			// Use unified formatter for json/yaml/tsv
			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatListOutput(components, { format, noColor }),
				);
			}

			if (components.length === 0) {
				return successResult(["No components found."]);
			}

			// Table format - custom visual display
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
	description:
		"Track service incidents with their impact levels, status, and latest updates. Shows both active and recently resolved incidents. Use --active-only to filter for ongoing issues requiring attention.",
	descriptionShort: "List active and recent incidents",
	descriptionMedium:
		"Display incidents with impact levels and updates. Use --active-only for ongoing issues.",
	usage: "[--active-only]",
	aliases: ["inc"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec, filteredArgs } = parseOutputArgs(
			args,
			session,
		);
		const activeOnly =
			filteredArgs.includes("--active-only") ||
			filteredArgs.includes("-a");

		// Handle --spec flag: return command specification for AI assistants
		if (spec) {
			const cmdSpec = getCommandSpec("cloudstatus incidents");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		try {
			const client = getClient();
			const response = activeOnly
				? await client.getUnresolvedIncidents()
				: await client.getIncidents();

			// Handle none format - return empty
			if (format === "none") {
				return successResult([]);
			}

			// Use unified formatter for json/yaml/tsv
			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatListOutput(response.incidents, { format, noColor }),
				);
			}

			if (response.incidents.length === 0) {
				return successResult(["No incidents found."]);
			}

			// Table format - custom visual display
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
	description:
		"View scheduled and in-progress maintenance windows with their timing and affected services. Plan around downtime windows. Use --upcoming to filter for future maintenance only.",
	descriptionShort: "List scheduled maintenance windows",
	descriptionMedium:
		"Show maintenance schedules and timing. Use --upcoming for future windows only.",
	usage: "[--upcoming]",
	aliases: ["maint"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec, filteredArgs } = parseOutputArgs(
			args,
			session,
		);
		const upcomingOnly =
			filteredArgs.includes("--upcoming") || filteredArgs.includes("-u");

		// Handle --spec flag: return command specification for AI assistants
		if (spec) {
			const cmdSpec = getCommandSpec("cloudstatus maintenance");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		try {
			const client = getClient();
			const response = upcomingOnly
				? await client.getUpcomingMaintenances()
				: await client.getMaintenances();

			// Handle none format - return empty
			if (format === "none") {
				return successResult([]);
			}

			// Use unified formatter for json/yaml/tsv
			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatListOutput(response.scheduled_maintenances, {
						format,
						noColor,
					}),
				);
			}

			if (response.scheduled_maintenances.length === 0) {
				return successResult(["No scheduled maintenance."]);
			}

			// Table format - custom visual display
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
