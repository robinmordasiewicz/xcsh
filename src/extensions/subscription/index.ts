/**
 * Subscription Domain Extension
 *
 * xcsh-specific commands for F5 XC subscription management.
 * These commands provide CLI-specific functionality that doesn't
 * belong in the upstream API (formatted views, validation helpers).
 *
 * Commands:
 * - overview: Formatted subscription overview (xcsh-specific view)
 * - addons: List addon services with filtering
 * - quota: Tenant-level quota analysis
 * - validate: Pre-deployment validation helper
 * - activation-status: Check pending addon activations
 *
 * Note: Standard API actions (list, get, create, delete) will come
 * from upstream specs when the subscription domain is added there.
 */

import type { DomainExtension } from "../types.js";
import type {
	CommandDefinition,
	DomainCommandResult,
} from "../../domains/registry.js";
import { successResult, errorResult } from "../../domains/registry.js";
import { SubscriptionClient } from "../../subscription/client.js";
import type { APIClient } from "../../api/client.js";
import {
	isAddonActive,
	isAddonAvailable,
	isAddonDenied,
	getStateDescription,
	getAccessStatusDescription,
	getTierDescription,
	QuotaStatus,
} from "../../subscription/types.js";

// Lazy client cache (keyed by server URL to handle different sessions)
const clientCache = new Map<string, SubscriptionClient>();

/**
 * Get or create subscription client for API client
 */
function getClient(apiClient: APIClient): SubscriptionClient {
	const key = apiClient.getServerUrl();
	let client = clientCache.get(key);
	if (!client) {
		client = new SubscriptionClient(apiClient);
		clientCache.set(key, client);
	}
	return client;
}

/**
 * Overview command - Display subscription overview
 * xcsh-specific: formatted view not available via API
 */
const overviewCommand: CommandDefinition = {
	name: "overview",
	description:
		"Display a comprehensive overview of your tenant subscription including current tier level, activated addon services, and quota usage summary. Provides at-a-glance subscription health assessment.",
	descriptionShort: "Display subscription tier and summary",
	descriptionMedium:
		"Show subscription tier, active addons, and quota usage summary for tenant health assessment.",
	usage: "[--json]",
	aliases: ["show", "info"],

	async execute(args, session): Promise<DomainCommandResult> {
		const jsonOutput = args.includes("--json");
		const format = session.getOutputFormat();

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' command first.");
		}

		try {
			const client = getClient(apiClient);
			const info = await client.getSubscriptionInfo();

			if (jsonOutput || format === "json") {
				return successResult([JSON.stringify(info, null, 2)]);
			}

			// Format human-readable output
			const lines: string[] = [];
			lines.push("=== Subscription Overview ===");
			lines.push("");
			lines.push(`Tier: ${info.tier}`);

			if (info.plan) {
				lines.push(`Plan: ${info.plan.displayName}`);
				if (info.plan.description) {
					lines.push(`Description: ${info.plan.description}`);
				}
			}

			lines.push("");
			lines.push("--- Active Addons ---");
			if (info.activeAddons.length === 0) {
				lines.push("  No active addons");
			} else {
				for (const addon of info.activeAddons) {
					lines.push(
						`  - ${addon.displayName} (${getTierDescription(addon.tier)})`,
					);
				}
			}

			lines.push("");
			lines.push("--- Quota Summary ---");
			lines.push(`  Total Limits: ${info.quotaSummary.totalLimits}`);
			lines.push(`  At Risk (>80%): ${info.quotaSummary.limitsAtRisk}`);
			lines.push(`  Exceeded: ${info.quotaSummary.limitsExceeded}`);

			if (info.quotaSummary.limitsExceeded > 0) {
				lines.push("");
				lines.push("  Warning: Some quotas are exceeded!");
			}

			return successResult(lines);
		} catch (err) {
			const message = err instanceof Error ? err.message : String(err);
			return errorResult(`Failed to get subscription info: ${message}`);
		}
	},
};

/**
 * Addons command - List addon services
 * xcsh-specific: filtering and formatted display
 */
const addonsCommand: CommandDefinition = {
	name: "addons",
	description:
		"List all addon services with their activation state and access status. Filter by active, available, or denied status. Shows tier requirements and enables identification of upgradable features.",
	descriptionShort: "List addon services and status",
	descriptionMedium:
		"Display addon services with activation state. Filter by active, available, or denied status.",
	usage: "[--filter active|available|denied] [--all] [--json]",
	aliases: ["services"],

	async execute(args, session): Promise<DomainCommandResult> {
		const jsonOutput = args.includes("--json");
		const showAll = args.includes("--all");
		const format = session.getOutputFormat();

		// Parse filter
		let filter = "";
		const filterIdx = args.indexOf("--filter");
		if (filterIdx !== -1) {
			const filterArg = args[filterIdx + 1];
			if (filterArg) {
				filter = filterArg;
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' command first.");
		}

		try {
			const client = getClient(apiClient);
			let addons = await client.getAddonServices("system");

			// Apply filter
			if (filter) {
				addons = client.filterAddons(addons, filter);
			} else if (!showAll) {
				// Default: show only active and available
				addons = addons.filter(
					(a) => isAddonActive(a) || isAddonAvailable(a),
				);
			}

			if (jsonOutput || format === "json") {
				return successResult([JSON.stringify(addons, null, 2)]);
			}

			if (addons.length === 0) {
				return successResult([
					`No addon services found${filter ? ` (filter: ${filter})` : ""}`,
				]);
			}

			// Format human-readable output
			const lines: string[] = [];
			lines.push("=== Addon Services ===");
			lines.push("");

			// Group by status
			const active = addons.filter(isAddonActive);
			const available = addons.filter(isAddonAvailable);
			const denied = addons.filter(isAddonDenied);
			const other = addons.filter(
				(a) =>
					!isAddonActive(a) &&
					!isAddonAvailable(a) &&
					!isAddonDenied(a),
			);

			if (active.length > 0) {
				lines.push("--- Active ---");
				for (const addon of active) {
					lines.push(`  [OK] ${addon.displayName}`);
					lines.push(
						`       Tier: ${getTierDescription(addon.tier)}`,
					);
				}
				lines.push("");
			}

			if (available.length > 0) {
				lines.push("--- Available ---");
				for (const addon of available) {
					lines.push(`  [ ] ${addon.displayName}`);
					lines.push(
						`       Status: ${getStateDescription(addon.state)}`,
					);
				}
				lines.push("");
			}

			if (denied.length > 0 && (showAll || filter === "denied")) {
				lines.push("--- Access Denied ---");
				for (const addon of denied) {
					lines.push(`  [X] ${addon.displayName}`);
					lines.push(
						`       Reason: ${getAccessStatusDescription(addon.accessStatus)}`,
					);
				}
				lines.push("");
			}

			if (other.length > 0 && showAll) {
				lines.push("--- Other ---");
				for (const addon of other) {
					lines.push(`  [-] ${addon.displayName}`);
					lines.push(
						`       State: ${getStateDescription(addon.state)}`,
					);
					lines.push(
						`       Access: ${getAccessStatusDescription(addon.accessStatus)}`,
					);
				}
			}

			return successResult(lines);
		} catch (err) {
			const message = err instanceof Error ? err.message : String(err);
			return errorResult(`Failed to get addon services: ${message}`);
		}
	},
};

/**
 * Quota command - Display quota limits and usage
 * xcsh-specific: formatted analysis with warnings
 */
const quotaCommand: CommandDefinition = {
	name: "quota",
	description:
		"Display tenant-level quota limits and current usage for all resource types. Identify resources approaching limits with warning thresholds. Use --warnings to filter for at-risk quotas only.",
	descriptionShort: "Display quota limits and usage",
	descriptionMedium:
		"Show tenant-level quota limits with usage percentages. Filter for at-risk resources.",
	usage: "[--warnings] [--json]",
	aliases: ["quotas", "limits"],

	async execute(args, session): Promise<DomainCommandResult> {
		const jsonOutput = args.includes("--json");
		const showWarnings = args.includes("--warnings");
		const format = session.getOutputFormat();

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' command first.");
		}

		try {
			const client = getClient(apiClient);
			const quotaInfo = await client.getQuotaInfo();

			let objects = quotaInfo.objects;

			// Filter to only warnings/exceeded if requested
			if (showWarnings) {
				objects = objects.filter(
					(q) =>
						q.status === QuotaStatus.Warning ||
						q.status === QuotaStatus.Exceeded,
				);
			}

			if (jsonOutput || format === "json") {
				return successResult([
					JSON.stringify({ ...quotaInfo, objects }, null, 2),
				]);
			}

			if (objects.length === 0) {
				return successResult([
					showWarnings
						? "No quota warnings or exceeded limits"
						: "No quota information available",
				]);
			}

			// Format human-readable output
			const lines: string[] = [];
			lines.push("=== Quota Usage (Tenant-Level) ===");
			lines.push("");

			// Group by status
			const exceeded = objects.filter(
				(q) => q.status === QuotaStatus.Exceeded,
			);
			const warning = objects.filter(
				(q) => q.status === QuotaStatus.Warning,
			);
			const ok = objects.filter((q) => q.status === QuotaStatus.OK);

			if (exceeded.length > 0) {
				lines.push("--- Exceeded ---");
				for (const q of exceeded) {
					lines.push(`  [!!] ${q.displayName}`);
					lines.push(
						`       Usage: ${q.usage} / ${q.limit} (${Math.round(q.percentage)}%)`,
					);
				}
				lines.push("");
			}

			if (warning.length > 0) {
				lines.push("--- Warning (>80%) ---");
				for (const q of warning) {
					lines.push(`  [!] ${q.displayName}`);
					lines.push(
						`      Usage: ${q.usage} / ${q.limit} (${Math.round(q.percentage)}%)`,
					);
				}
				lines.push("");
			}

			if (ok.length > 0 && !showWarnings) {
				lines.push("--- OK ---");
				for (const q of ok) {
					lines.push(
						`  [OK] ${q.displayName}: ${q.usage} / ${q.limit} (${Math.round(q.percentage)}%)`,
					);
				}
			}

			return successResult(lines);
		} catch (err) {
			const message = err instanceof Error ? err.message : String(err);
			return errorResult(`Failed to get quota info: ${message}`);
		}
	},
};

/**
 * Validate command - Pre-deployment validation
 * xcsh-specific: CLI validation helper
 */
const validateCommand: CommandDefinition = {
	name: "validate",
	description:
		"Run pre-deployment validation checks against quota limits and feature availability. Verify resource creation will succeed before deployment. Check if specific features are enabled for your subscription tier.",
	descriptionShort: "Validate quotas and feature access",
	descriptionMedium:
		"Pre-deployment validation for quota capacity and feature availability checks.",
	usage: "[--resource-type <type> --count <n>] [--feature <name>] [--json]",

	async execute(args, session): Promise<DomainCommandResult> {
		const jsonOutput = args.includes("--json");
		const format = session.getOutputFormat();

		// Parse resource type
		let resourceType: string | undefined;
		const resourceIdx = args.indexOf("--resource-type");
		if (resourceIdx !== -1 && args[resourceIdx + 1]) {
			resourceType = args[resourceIdx + 1];
		}

		// Parse count
		let count: number | undefined;
		const countIdx = args.indexOf("--count");
		if (countIdx !== -1) {
			const countArg = args[countIdx + 1];
			if (countArg) {
				count = parseInt(countArg, 10);
			}
		}

		// Parse feature
		let feature: string | undefined;
		const featureIdx = args.indexOf("--feature");
		if (featureIdx !== -1 && args[featureIdx + 1]) {
			feature = args[featureIdx + 1];
		}

		if (!resourceType && !feature) {
			return errorResult(
				"Please specify --resource-type and --count, or --feature to validate",
			);
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' command first.");
		}

		try {
			const client = getClient(apiClient);
			const result = await client.validateResource({
				resourceType: resourceType ?? "",
				count: count ?? 0,
				feature: feature ?? "",
			});

			if (jsonOutput || format === "json") {
				return successResult([JSON.stringify(result, null, 2)]);
			}

			// Format human-readable output
			const lines: string[] = [];
			lines.push("=== Validation Result ===");
			lines.push("");
			lines.push(`Status: ${result.valid ? "[PASS]" : "[FAIL]"}`);
			lines.push("");

			if (result.checks.length > 0) {
				lines.push("--- Checks ---");
				for (const check of result.checks) {
					const icon =
						check.result === "PASS"
							? "[OK]"
							: check.result === "WARNING"
								? "[!]"
								: "[X]";
					lines.push(
						`  ${icon} [${check.type}] ${check.message ?? ""}`,
					);
				}
			}

			if (result.errors && result.errors.length > 0) {
				lines.push("");
				lines.push("--- Errors ---");
				for (const error of result.errors) {
					lines.push(`  [X] ${error}`);
				}
			}

			if (result.warnings && result.warnings.length > 0) {
				lines.push("");
				lines.push("--- Warnings ---");
				for (const warning of result.warnings) {
					lines.push(`  [!] ${warning}`);
				}
			}

			const commandResult: DomainCommandResult = {
				output: lines,
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
			};
			if (!result.valid) {
				commandResult.error = "Validation failed";
			}
			return commandResult;
		} catch (err) {
			const message = err instanceof Error ? err.message : String(err);
			return errorResult(`Validation failed: ${message}`);
		}
	},
};

/**
 * Activation-status command - Check pending activations
 * xcsh-specific: activation monitoring helper
 */
const activationStatusCommand: CommandDefinition = {
	name: "activation-status",
	description:
		"Monitor the status of pending addon service activation requests. Track which addons are awaiting approval or provisioning. View currently active addons alongside pending requests.",
	descriptionShort: "Check pending addon activations",
	descriptionMedium:
		"Monitor pending addon activation requests and track provisioning status.",
	usage: "[--json]",

	async execute(args, session): Promise<DomainCommandResult> {
		const jsonOutput = args.includes("--json");
		const format = session.getOutputFormat();

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' command first.");
		}

		try {
			const client = getClient(apiClient);
			const result = await client.getPendingActivations("system");

			if (jsonOutput || format === "json") {
				return successResult([JSON.stringify(result, null, 2)]);
			}

			// Format human-readable output
			const lines: string[] = [];
			lines.push("=== Activation Status ===");
			lines.push("");

			if (result.pendingActivations.length === 0) {
				lines.push("No pending activation requests.");
			} else {
				lines.push(`Pending Activations: ${result.totalPending}`);
				lines.push("");
				for (const pending of result.pendingActivations) {
					lines.push(`  [...] ${pending.addonService}`);
					if (pending.message) {
						lines.push(`        Status: ${pending.message}`);
					}
				}
			}

			lines.push("");
			lines.push(`Active Addons: ${result.activeAddons.length}`);
			if (result.activeAddons.length > 0) {
				for (const addon of result.activeAddons) {
					lines.push(`  [OK] ${addon}`);
				}
			}

			return successResult(lines);
		} catch (err) {
			const message = err instanceof Error ? err.message : String(err);
			return errorResult(`Failed to get activation status: ${message}`);
		}
	},
};

/**
 * Subscription domain extension
 *
 * Provides xcsh-specific subscription management commands.
 * Standard API actions (list, get, etc.) will come from upstream.
 */
export const subscriptionExtension: DomainExtension = {
	targetDomain: "subscription",
	description:
		"xcsh-specific subscription management commands (overview, quota analysis, validation)",
	standalone: true, // Works even before upstream adds subscription domain

	commands: new Map([
		["overview", overviewCommand],
		["addons", addonsCommand],
		["quota", quotaCommand],
		["validate", validateCommand],
		["activation-status", activationStatusCommand],
	]),

	subcommands: new Map(),
};
