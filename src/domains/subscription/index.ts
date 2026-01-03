/**
 * Subscription Domain - F5 XC subscription, billing, quota, and usage management
 *
 * Provides comprehensive access to subscription tier information, addon services,
 * quota limits and usage, billing details, and reporting capabilities.
 */

import type {
	DomainDefinition,
	SubcommandGroup,
	CommandDefinition,
	DomainCommandResult,
} from "../registry.js";
import { successResult, errorResult } from "../registry.js";
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
import { getSubscriptionClient } from "./client.js";

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
	const { options, remainingArgs } = parseDomainOutputFlags(
		args,
		session.getOutputFormat(),
	);

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

/**
 * Helper: pad string to fixed width
 */
function padEnd(str: string, length: number): string {
	if (str.length >= length) {
		return str.slice(0, length);
	}
	return str + " ".repeat(length - str.length);
}

/**
 * Helper: format currency
 */
function formatCurrency(amount: number | undefined, currency = "USD"): string {
	if (amount === undefined) return "N/A";
	return new Intl.NumberFormat("en-US", {
		style: "currency",
		currency,
	}).format(amount);
}

/**
 * Helper: format percentage
 */
function formatPercentage(value: number | undefined): string {
	if (value === undefined) return "N/A";
	return `${value.toFixed(1)}%`;
}

/**
 * Helper: format date
 */
function formatDate(dateStr: string | undefined): string {
	if (!dateStr) return "N/A";
	try {
		return new Date(dateStr).toLocaleDateString();
	} catch {
		return dateStr;
	}
}

/**
 * Helper: create progress bar
 */
function progressBar(percentage: number, width = 10): string {
	const filled = Math.round((percentage / 100) * width);
	const empty = width - filled;
	return "\u2588".repeat(filled) + "\u2591".repeat(empty);
}

// ============================================================================
// DEFAULT SHOW COMMAND
// ============================================================================

const showCommand: CommandDefinition = {
	name: "show",
	description:
		"Display comprehensive subscription overview including plan tier, addon summary, quota utilization, and current billing period usage. Provides a quick snapshot of your F5 XC subscription status.",
	descriptionShort: "Display subscription overview",
	descriptionMedium:
		"Show subscription tier, active addons, quota usage summary, and current billing status.",
	usage: "",
	aliases: ["overview", "info"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec } = parseOutputArgs(args, session);

		if (spec) {
			const cmdSpec = getCommandSpec("subscription show");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);
			const overview = await client.getOverview();

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatDomainOutput(overview, { format, noColor }),
				);
			}

			// Table format - visual display
			const lines: string[] = [];
			lines.push(
				"\u256D\u2500 Subscription Overview \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256E",
			);

			// Plan info
			const planName =
				overview.plan?.display_name ?? overview.plan?.name ?? "Unknown";
			const planTier =
				overview.plan?.tier ?? overview.plan?.plan_type ?? "N/A";
			lines.push(`\u2502 Plan:     ${padEnd(planName, 44)} \u2502`);
			lines.push(`\u2502 Tier:     ${padEnd(planTier, 44)} \u2502`);

			lines.push(
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			);

			// Addon summary
			if (overview.addon_summary) {
				const addonInfo = `${overview.addon_summary.active} active, ${overview.addon_summary.available} available`;
				lines.push(
					`\u2502 Addon Services: ${padEnd(addonInfo, 38)} \u2502`,
				);
			}

			// Quota summary
			if (overview.quota_summary) {
				const quotaInfo = `${overview.quota_summary.used} quotas tracked`;
				const criticalCount =
					overview.quota_summary.critical_quotas?.length ?? 0;
				const criticalInfo =
					criticalCount > 0 ? ` (${criticalCount} >80%)` : "";
				lines.push(
					`\u2502 Quota Usage:    ${padEnd(quotaInfo + criticalInfo, 38)} \u2502`,
				);
			}

			// Current usage
			if (overview.current_usage) {
				const costInfo = formatCurrency(
					overview.current_usage.total_cost,
					overview.current_usage.currency,
				);
				const projectedInfo = overview.current_usage.projected_cost
					? ` (projected: ${formatCurrency(overview.current_usage.projected_cost, overview.current_usage.currency)})`
					: "";
				lines.push(
					`\u2502 Current Cost:   ${padEnd(costInfo + projectedInfo, 38)} \u2502`,
				);
			}

			// Billing status
			if (overview.billing_status) {
				const pmStatus =
					overview.billing_status.payment_method_status ?? "Unknown";
				lines.push(
					`\u2502 Payment Status: ${padEnd(pmStatus, 38)} \u2502`,
				);
			}

			lines.push(
				"\u2570\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256F",
			);

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(
				`Failed to get subscription overview: ${message}`,
			);
		}
	},
};

// ============================================================================
// PLAN COMMANDS
// ============================================================================

const planShowCommand: CommandDefinition = {
	name: "show",
	description:
		"Display detailed information about your current subscription plan including tier, features, and limits.",
	descriptionShort: "Show current plan details",
	descriptionMedium:
		"Display detailed subscription plan information including tier level, features, and resource limits.",
	usage: "",

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec } = parseOutputArgs(args, session);

		if (spec) {
			const cmdSpec = getCommandSpec("subscription plan show");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);
			const plan = await client.getCurrentPlan();

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatDomainOutput(plan, { format, noColor }),
				);
			}

			const lines: string[] = [
				"=== Current Subscription Plan ===",
				"",
				`Name:        ${plan.display_name ?? plan.name}`,
				`Tier:        ${plan.tier ?? plan.plan_type ?? "N/A"}`,
				`Description: ${plan.description ?? "N/A"}`,
				`Status:      ${plan.status ?? "Active"}`,
			];

			if (plan.features && plan.features.length > 0) {
				lines.push("", "Features:");
				for (const feature of plan.features) {
					lines.push(`  - ${feature}`);
				}
			}

			if (plan.limits && Object.keys(plan.limits).length > 0) {
				lines.push("", "Limits:");
				for (const [key, value] of Object.entries(plan.limits)) {
					lines.push(`  ${key}: ${value}`);
				}
			}

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to get plan details: ${message}`);
		}
	},
};

const planListCommand: CommandDefinition = {
	name: "list",
	description:
		"List all available subscription plans with their features and pricing tiers.",
	descriptionShort: "List available plans",
	descriptionMedium:
		"Display all subscription plan options with tier information and feature comparisons.",
	usage: "",

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec } = parseOutputArgs(args, session);

		if (spec) {
			const cmdSpec = getCommandSpec("subscription plan list");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);
			const plans = await client.listPlans();

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatListOutput(plans, { format, noColor }),
				);
			}

			if (plans.length === 0) {
				return successResult(["No plans available."]);
			}

			const lines: string[] = [
				"\u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510",
				"\u2502 NAME                 \u2502 TIER          \u2502 DESCRIPTION        \u2502",
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			];

			for (const plan of plans) {
				const name = padEnd(plan.display_name ?? plan.name, 20);
				const tier = padEnd(plan.tier ?? plan.plan_type ?? "N/A", 13);
				const desc = padEnd((plan.description ?? "").slice(0, 18), 18);
				lines.push(
					`\u2502 ${name} \u2502 ${tier} \u2502 ${desc} \u2502`,
				);
			}

			lines.push(
				"\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518",
			);

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to list plans: ${message}`);
		}
	},
};

const planSubcommands: SubcommandGroup = {
	name: "plan",
	description:
		"View and manage subscription plan information. Show current plan details, list available plans, and initiate plan transitions.",
	descriptionShort: "Manage subscription plans",
	descriptionMedium:
		"View current plan, list available options, and manage plan transitions.",
	commands: new Map([
		["show", planShowCommand],
		["list", planListCommand],
	]),
	defaultCommand: planShowCommand,
};

// ============================================================================
// ADDON COMMANDS
// ============================================================================

const addonListCommand: CommandDefinition = {
	name: "list",
	description:
		"List all available addon services with their status, category, and access level. Use --subscribed to filter for active subscriptions only.",
	descriptionShort: "List addon services",
	descriptionMedium:
		"Display all addon services with status and access information. Use --subscribed for active only.",
	usage: "[--subscribed]",
	aliases: ["ls"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec, filteredArgs } = parseOutputArgs(
			args,
			session,
		);
		const subscribedOnly =
			filteredArgs.includes("--subscribed") ||
			filteredArgs.includes("-s");

		if (spec) {
			const cmdSpec = getCommandSpec("subscription addon list");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);
			let addons = await client.listAddonServices();

			if (subscribedOnly) {
				addons = addons.filter(
					(a) =>
						a.access_type === "SUBSCRIBED" ||
						a.access_type === "ENABLED",
				);
			}

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatListOutput(addons, { format, noColor }),
				);
			}

			if (addons.length === 0) {
				return successResult([
					subscribedOnly
						? "No subscribed addons."
						: "No addon services available.",
				]);
			}

			const lines: string[] = [
				"\u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510",
				"\u2502 NAME                         \u2502 STATUS       \u2502 ACCESS      \u2502",
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			];

			for (const addon of addons) {
				const name = padEnd(addon.display_name ?? addon.name, 28);
				const status = padEnd(addon.status ?? "N/A", 12);
				const access = padEnd(addon.access_type ?? "N/A", 11);
				lines.push(
					`\u2502 ${name} \u2502 ${status} \u2502 ${access} \u2502`,
				);
			}

			lines.push(
				"\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518",
			);

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to list addon services: ${message}`);
		}
	},
};

const addonShowCommand: CommandDefinition = {
	name: "show",
	description:
		"Display detailed information about a specific addon service including features, pricing, and requirements.",
	descriptionShort: "Show addon details",
	descriptionMedium:
		"Display detailed addon service information including features, pricing, and dependencies.",
	usage: "<addon-name>",

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec, filteredArgs } = parseOutputArgs(
			args,
			session,
		);

		if (spec) {
			const cmdSpec = getCommandSpec("subscription addon show");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const addonName = filteredArgs[0];
		if (!addonName) {
			return errorResult("Usage: subscription addon show <addon-name>");
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);
			const addon = await client.getAddonService(addonName);

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatDomainOutput(addon, { format, noColor }),
				);
			}

			const lines: string[] = [
				`=== Addon Service: ${addon.display_name ?? addon.name} ===`,
				"",
				`Name:        ${addon.name}`,
				`Category:    ${addon.category ?? "N/A"}`,
				`Status:      ${addon.status ?? "N/A"}`,
				`Access:      ${addon.access_type ?? "N/A"}`,
				`Description: ${addon.description ?? "N/A"}`,
			];

			if (addon.pricing) {
				lines.push("", "Pricing:");
				lines.push(`  Model:   ${addon.pricing.model ?? "N/A"}`);
				if (addon.pricing.base_price !== undefined) {
					lines.push(
						`  Base:    ${formatCurrency(addon.pricing.base_price, addon.pricing.currency)}`,
					);
				}
				if (addon.pricing.billing_period) {
					lines.push(`  Period:  ${addon.pricing.billing_period}`);
				}
			}

			if (addon.features && addon.features.length > 0) {
				lines.push("", "Features:");
				for (const feature of addon.features) {
					lines.push(`  - ${feature}`);
				}
			}

			if (addon.requires && addon.requires.length > 0) {
				lines.push("", "Requirements:");
				for (const req of addon.requires) {
					lines.push(`  - ${req}`);
				}
			}

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to get addon details: ${message}`);
		}
	},
};

const addonStatusCommand: CommandDefinition = {
	name: "status",
	description:
		"Display activation status for all addon services showing which are active, pending, or inactive.",
	descriptionShort: "Show addon activation status",
	descriptionMedium:
		"Display activation status for all addon services with detailed state information.",
	usage: "",

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec } = parseOutputArgs(args, session);

		if (spec) {
			const cmdSpec = getCommandSpec("subscription addon status");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);
			const statuses = await client.getAllAddonActivationStatus();

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatListOutput(statuses, { format, noColor }),
				);
			}

			if (statuses.length === 0) {
				return successResult(["No addon activation data available."]);
			}

			const lines: string[] = [
				"\u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510",
				"\u2502 ADDON                        \u2502 STATUS     \u2502 ACTIVATED          \u2502",
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			];

			for (const status of statuses) {
				const name = padEnd(status.addon_name, 28);
				const state = padEnd(status.status, 10);
				const activated = padEnd(formatDate(status.activated_at), 18);
				lines.push(
					`\u2502 ${name} \u2502 ${state} \u2502 ${activated} \u2502`,
				);
			}

			lines.push(
				"\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518",
			);

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to get addon status: ${message}`);
		}
	},
};

const addonSubcommands: SubcommandGroup = {
	name: "addon",
	description:
		"Manage addon services for your subscription. List available addons, view details, check activation status, and manage subscriptions.",
	descriptionShort: "Manage addon services",
	descriptionMedium:
		"List, view, and manage addon service subscriptions and activation status.",
	commands: new Map([
		["list", addonListCommand],
		["show", addonShowCommand],
		["status", addonStatusCommand],
	]),
	defaultCommand: addonListCommand,
};

// ============================================================================
// QUOTA COMMANDS
// ============================================================================

const quotaLimitsCommand: CommandDefinition = {
	name: "limits",
	description:
		"Display all tenant-level quota limits including resource caps for load balancers, origins, sites, and other configurable resources.",
	descriptionShort: "Show quota limits",
	descriptionMedium:
		"Display all tenant-level resource quota limits and their current values.",
	usage: "",

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec } = parseOutputArgs(args, session);

		if (spec) {
			const cmdSpec = getCommandSpec("subscription quota limits");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);
			const response = await client.getQuotaLimits();

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatDomainOutput(response, { format, noColor }),
				);
			}

			const limits = response.limits ?? [];
			if (limits.length === 0) {
				return successResult(["No quota limits defined."]);
			}

			const lines: string[] = [
				`Plan: ${response.plan_type ?? "N/A"}`,
				"",
				"\u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510",
				"\u2502 RESOURCE                           \u2502 LIMIT    \u2502 SCOPE    \u2502",
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			];

			for (const limit of limits) {
				const name = padEnd(limit.display_name ?? limit.name, 34);
				const value = padEnd(
					`${limit.limit}${limit.unit ? " " + limit.unit : ""}`,
					8,
				);
				const scope = padEnd(limit.scope ?? "TENANT", 8);
				lines.push(
					`\u2502 ${name} \u2502 ${value} \u2502 ${scope} \u2502`,
				);
			}

			lines.push(
				"\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518",
			);

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to get quota limits: ${message}`);
		}
	},
};

const quotaUsageCommand: CommandDefinition = {
	name: "usage",
	description:
		"Display current quota usage against limits with utilization percentages. Use --critical to show only quotas above 80% utilization.",
	descriptionShort: "Show quota usage",
	descriptionMedium:
		"Display current resource usage against quota limits with utilization percentage.",
	usage: "[--critical]",

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec, filteredArgs } = parseOutputArgs(
			args,
			session,
		);
		const criticalOnly =
			filteredArgs.includes("--critical") || filteredArgs.includes("-c");

		if (spec) {
			const cmdSpec = getCommandSpec("subscription quota usage");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);
			const response = await client.getQuotaUsage();

			let usage = response.usage ?? [];
			if (criticalOnly) {
				usage = usage.filter(
					(q) => q.percentage !== undefined && q.percentage > 80,
				);
			}

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatListOutput(usage, { format, noColor }),
				);
			}

			if (usage.length === 0) {
				return successResult([
					criticalOnly
						? "No quotas above 80% utilization."
						: "No quota usage data available.",
				]);
			}

			const lines: string[] = [
				`As of: ${response.as_of ?? "N/A"}`,
				"",
				"\u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510",
				"\u2502 RESOURCE                 \u2502 USAGE        \u2502 UTILIZATION  \u2502",
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			];

			for (const q of usage) {
				const name = padEnd(q.display_name ?? q.name, 24);
				const usageStr = padEnd(`${q.current}/${q.limit}`, 12);
				const pct = q.percentage ?? (q.current / q.limit) * 100;
				const bar = progressBar(pct);
				const pctStr = formatPercentage(pct);
				lines.push(
					`\u2502 ${name} \u2502 ${usageStr} \u2502 ${bar} ${padEnd(pctStr, 1)} \u2502`,
				);
			}

			lines.push(
				"\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518",
			);

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to get quota usage: ${message}`);
		}
	},
};

const quotaSubcommands: SubcommandGroup = {
	name: "quota",
	description:
		"View tenant-level quota limits and current usage. Monitor resource utilization to avoid quota exhaustion.",
	descriptionShort: "View quota limits and usage",
	descriptionMedium:
		"Display tenant-level quota limits and current usage with utilization metrics.",
	commands: new Map([
		["limits", quotaLimitsCommand],
		["usage", quotaUsageCommand],
	]),
	defaultCommand: quotaUsageCommand,
};

// ============================================================================
// USAGE COMMANDS
// ============================================================================

const usageCurrentCommand: CommandDefinition = {
	name: "current",
	description:
		"Display current billing period usage including itemized costs and projected totals.",
	descriptionShort: "Show current period usage",
	descriptionMedium:
		"Display current billing period usage with cost breakdown and projections.",
	usage: "",

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec } = parseOutputArgs(args, session);

		if (spec) {
			const cmdSpec = getCommandSpec("subscription usage current");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);
			const usage = await client.getCurrentUsage();

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatDomainOutput(usage, { format, noColor }),
				);
			}

			const lines: string[] = [
				"=== Current Billing Period Usage ===",
				"",
				`Period: ${formatDate(usage.billing_period_start)} - ${formatDate(usage.billing_period_end)}`,
				`Total Cost: ${formatCurrency(usage.total_cost, usage.currency)}`,
			];

			if (usage.projected_cost !== undefined) {
				lines.push(
					`Projected: ${formatCurrency(usage.projected_cost, usage.currency)}`,
				);
			}

			if (usage.usage_items && usage.usage_items.length > 0) {
				lines.push("", "Usage Items:");
				for (const item of usage.usage_items) {
					const cost = formatCurrency(
						item.total_cost,
						usage.currency,
					);
					lines.push(
						`  ${item.display_name ?? item.name}: ${item.quantity} ${item.unit ?? ""} = ${cost}`,
					);
				}
			}

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to get current usage: ${message}`);
		}
	},
};

const usageMonthlyCommand: CommandDefinition = {
	name: "monthly",
	description:
		"Display monthly usage summaries with cost breakdowns for historical billing periods.",
	descriptionShort: "Show monthly usage history",
	descriptionMedium:
		"Display monthly usage summaries with historical cost data and trends.",
	usage: "[--limit <n>]",

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec, filteredArgs } = parseOutputArgs(
			args,
			session,
		);

		// Parse --limit flag
		let limit = 6;
		const limitIdx = filteredArgs.indexOf("--limit");
		const limitValue =
			limitIdx !== -1 ? filteredArgs[limitIdx + 1] : undefined;
		if (limitValue) {
			limit = parseInt(limitValue, 10) || 6;
		}

		if (spec) {
			const cmdSpec = getCommandSpec("subscription usage monthly");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);
			let months = await client.getMonthlyUsage();
			months = months.slice(0, limit);

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatListOutput(months, { format, noColor }),
				);
			}

			if (months.length === 0) {
				return successResult(["No monthly usage data available."]);
			}

			const lines: string[] = [
				"\u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510",
				"\u2502 MONTH        \u2502 TOTAL COST    \u2502 STATUS       \u2502",
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			];

			for (const month of months) {
				const monthStr = padEnd(`${month.month}/${month.year}`, 12);
				const cost = padEnd(
					formatCurrency(month.total_cost, month.currency),
					13,
				);
				const status = padEnd(month.invoice_status ?? "N/A", 12);
				lines.push(
					`\u2502 ${monthStr} \u2502 ${cost} \u2502 ${status} \u2502`,
				);
			}

			lines.push(
				"\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518",
			);

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to get monthly usage: ${message}`);
		}
	},
};

const usageSubcommands: SubcommandGroup = {
	name: "usage",
	description:
		"View usage metrics and cost data for current and historical billing periods.",
	descriptionShort: "View usage metrics",
	descriptionMedium:
		"Display usage metrics, cost breakdowns, and historical billing data.",
	commands: new Map([
		["current", usageCurrentCommand],
		["monthly", usageMonthlyCommand],
	]),
	defaultCommand: usageCurrentCommand,
};

// ============================================================================
// BILLING COMMANDS
// ============================================================================

const billingPaymentListCommand: CommandDefinition = {
	name: "list",
	description:
		"List all configured payment methods showing type, status, and primary designation.",
	descriptionShort: "List payment methods",
	descriptionMedium:
		"Display all payment methods with type, status, and primary/secondary designation.",
	usage: "",

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec } = parseOutputArgs(args, session);

		if (spec) {
			const cmdSpec = getCommandSpec("subscription billing payment list");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);
			const methods = await client.listPaymentMethods();

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatListOutput(methods, { format, noColor }),
				);
			}

			if (methods.length === 0) {
				return successResult(["No payment methods configured."]);
			}

			const lines: string[] = [
				"\u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510",
				"\u2502 NAME               \u2502 TYPE         \u2502 STATUS   \u2502 PRIMARY  \u2502",
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			];

			for (const method of methods) {
				const name = padEnd(method.name, 18);
				const type = padEnd(method.type, 12);
				const status = padEnd(method.status ?? "N/A", 8);
				const primary = padEnd(method.is_primary ? "Yes" : "No", 8);
				lines.push(
					`\u2502 ${name} \u2502 ${type} \u2502 ${status} \u2502 ${primary} \u2502`,
				);
			}

			lines.push(
				"\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518",
			);

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to list payment methods: ${message}`);
		}
	},
};

const billingInvoiceListCommand: CommandDefinition = {
	name: "list",
	description:
		"List all invoices with their status, amount, and billing period information.",
	descriptionShort: "List invoices",
	descriptionMedium:
		"Display all invoices with status, amount, and billing period details.",
	usage: "[--limit <n>]",

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec, filteredArgs } = parseOutputArgs(
			args,
			session,
		);

		// Parse --limit flag
		let limit = 10;
		const limitIdx = filteredArgs.indexOf("--limit");
		const limitValue =
			limitIdx !== -1 ? filteredArgs[limitIdx + 1] : undefined;
		if (limitValue) {
			limit = parseInt(limitValue, 10) || 10;
		}

		if (spec) {
			const cmdSpec = getCommandSpec("subscription billing invoice list");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);
			let invoices = await client.listInvoices();
			invoices = invoices.slice(0, limit);

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatListOutput(invoices, { format, noColor }),
				);
			}

			if (invoices.length === 0) {
				return successResult(["No invoices found."]);
			}

			const lines: string[] = [
				"\u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u252C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510",
				"\u2502 INVOICE        \u2502 DATE         \u2502 AMOUNT        \u2502 STATUS   \u2502",
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u253C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			];

			for (const invoice of invoices) {
				const num = padEnd(
					invoice.invoice_number ?? invoice.invoice_id,
					14,
				);
				const date = padEnd(formatDate(invoice.issue_date), 12);
				const amount = padEnd(
					formatCurrency(invoice.total_amount, invoice.currency),
					13,
				);
				const status = padEnd(invoice.status ?? "N/A", 8);
				lines.push(
					`\u2502 ${num} \u2502 ${date} \u2502 ${amount} \u2502 ${status} \u2502`,
				);
			}

			lines.push(
				"\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2534\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518",
			);

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to list invoices: ${message}`);
		}
	},
};

// Payment subcommand group
const billingPaymentSubcommands: SubcommandGroup = {
	name: "payment",
	description: "Manage payment methods for your subscription billing.",
	descriptionShort: "Manage payment methods",
	descriptionMedium:
		"List, view, and manage payment methods for subscription billing.",
	commands: new Map([["list", billingPaymentListCommand]]),
	defaultCommand: billingPaymentListCommand,
};

// Invoice subcommand group
const billingInvoiceSubcommands: SubcommandGroup = {
	name: "invoice",
	description: "View and download invoices for your subscription.",
	descriptionShort: "Manage invoices",
	descriptionMedium:
		"List invoices and download invoice PDFs for billing records.",
	commands: new Map([["list", billingInvoiceListCommand]]),
	defaultCommand: billingInvoiceListCommand,
};

// Billing commands that wrap subcommand groups
const billingPaymentCommand: CommandDefinition = {
	name: "payment",
	description: billingPaymentSubcommands.description,
	descriptionShort: billingPaymentSubcommands.descriptionShort,
	descriptionMedium: billingPaymentSubcommands.descriptionMedium,
	usage: "<list>",

	async execute(args, session): Promise<DomainCommandResult> {
		// Default to list
		return billingPaymentListCommand.execute(args, session);
	},
};

const billingInvoiceCommand: CommandDefinition = {
	name: "invoice",
	description: billingInvoiceSubcommands.description,
	descriptionShort: billingInvoiceSubcommands.descriptionShort,
	descriptionMedium: billingInvoiceSubcommands.descriptionMedium,
	usage: "<list>",

	async execute(args, session): Promise<DomainCommandResult> {
		// Default to list
		return billingInvoiceListCommand.execute(args, session);
	},
};

const billingSubcommands: SubcommandGroup = {
	name: "billing",
	description:
		"Manage billing information including payment methods and invoices.",
	descriptionShort: "Manage billing",
	descriptionMedium:
		"View and manage payment methods, invoices, and billing details.",
	commands: new Map([
		["payment", billingPaymentCommand],
		["invoice", billingInvoiceCommand],
	]),
};

// ============================================================================
// REPORT COMMANDS
// ============================================================================

const reportSummaryCommand: CommandDefinition = {
	name: "summary",
	description:
		"Generate comprehensive subscription report combining plan details, addon status, quota utilization, usage metrics, and billing summary.",
	descriptionShort: "Generate subscription report",
	descriptionMedium:
		"Create comprehensive report with plan, addons, quotas, usage, and billing data.",
	usage: "",
	aliases: ["full"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec } = parseOutputArgs(args, session);

		if (spec) {
			const cmdSpec = getCommandSpec("subscription report summary");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult("Not authenticated. Use 'login' first.");
		}

		try {
			const client = getSubscriptionClient(apiClient);

			// Fetch all data in parallel
			const [
				plan,
				addons,
				subscriptions,
				quotaUsage,
				currentUsage,
				paymentMethods,
				invoices,
			] = await Promise.allSettled([
				client.getCurrentPlan(),
				client.listAddonServices(),
				client.listAddonSubscriptions(),
				client.getQuotaUsage(),
				client.getCurrentUsage(),
				client.listPaymentMethods(),
				client.listInvoices(),
			]);

			const report: Record<string, unknown> = {
				generated_at: new Date().toISOString(),
			};

			// Add plan
			if (plan.status === "fulfilled") {
				report.plan = plan.value;
			}

			// Add addons
			if (
				addons.status === "fulfilled" ||
				subscriptions.status === "fulfilled"
			) {
				report.addons = {
					available:
						addons.status === "fulfilled" ? addons.value : [],
					subscribed:
						subscriptions.status === "fulfilled"
							? subscriptions.value
							: [],
				};
			}

			// Add quotas
			if (quotaUsage.status === "fulfilled") {
				const usage = quotaUsage.value.usage ?? [];
				const totalUsed = usage.reduce(
					(sum, q) => sum + (q.percentage ?? 0),
					0,
				);
				report.quotas = {
					usage,
					average_utilization:
						usage.length > 0 ? totalUsed / usage.length : 0,
				};
			}

			// Add usage
			if (currentUsage.status === "fulfilled") {
				report.usage = currentUsage.value;
			}

			// Add billing
			if (
				paymentMethods.status === "fulfilled" ||
				invoices.status === "fulfilled"
			) {
				report.billing = {
					payment_methods:
						paymentMethods.status === "fulfilled"
							? paymentMethods.value
							: [],
					recent_invoices:
						invoices.status === "fulfilled"
							? invoices.value.slice(0, 3)
							: [],
				};
			}

			if (format === "none") {
				return successResult([]);
			}

			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatDomainOutput(report, { format, noColor }),
				);
			}

			// Table format - comprehensive visual report
			const lines: string[] = [];
			lines.push(
				"\u256D\u2500 Subscription Summary Report \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256E",
			);
			lines.push(
				`\u2502 Generated: ${new Date().toLocaleString()}${" ".repeat(28)}\u2502`,
			);
			lines.push(
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			);

			// Plan section
			lines.push("\u2502 PLAN DETAILS" + " ".repeat(43) + "\u2502");
			if (plan.status === "fulfilled") {
				const p = plan.value;
				lines.push(
					`\u2502   Tier: ${padEnd(p.tier ?? p.plan_type ?? "N/A", 45)}\u2502`,
				);
				lines.push(
					`\u2502   Name: ${padEnd(p.display_name ?? p.name, 45)}\u2502`,
				);
			}

			lines.push(
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			);

			// Addons section
			lines.push("\u2502 ADDON SERVICES" + " ".repeat(41) + "\u2502");
			if (addons.status === "fulfilled") {
				const addonList = addons.value;
				const active = addonList.filter(
					(a) =>
						a.access_type === "SUBSCRIBED" ||
						a.access_type === "ENABLED",
				);
				lines.push(
					`\u2502   Active: ${active.length} of ${addonList.length} available${" ".repeat(30)}\u2502`,
				);
			}

			lines.push(
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			);

			// Quota section
			lines.push("\u2502 QUOTA UTILIZATION" + " ".repeat(38) + "\u2502");
			if (quotaUsage.status === "fulfilled") {
				const usage = quotaUsage.value.usage ?? [];
				const critical = usage.filter((q) => (q.percentage ?? 0) > 80);
				const avgUtil =
					usage.length > 0
						? usage.reduce(
								(sum, q) => sum + (q.percentage ?? 0),
								0,
							) / usage.length
						: 0;
				lines.push(
					`\u2502   Average: ${formatPercentage(avgUtil)}${" ".repeat(40)}\u2502`,
				);
				lines.push(
					`\u2502   Critical (>80%): ${critical.length} quotas${" ".repeat(29)}\u2502`,
				);
			}

			lines.push(
				"\u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524",
			);

			// Usage section
			lines.push("\u2502 CURRENT USAGE" + " ".repeat(42) + "\u2502");
			if (currentUsage.status === "fulfilled") {
				const u = currentUsage.value;
				lines.push(
					`\u2502   Total Cost: ${formatCurrency(u.total_cost, u.currency)}${" ".repeat(35)}\u2502`,
				);
				if (u.projected_cost) {
					lines.push(
						`\u2502   Projected: ${formatCurrency(u.projected_cost, u.currency)}${" ".repeat(36)}\u2502`,
					);
				}
			}

			lines.push(
				"\u2570\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256F",
			);

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to generate report: ${message}`);
		}
	},
};

const reportSubcommands: SubcommandGroup = {
	name: "report",
	description:
		"Generate comprehensive subscription reports for analysis and planning.",
	descriptionShort: "Generate reports",
	descriptionMedium:
		"Create detailed subscription reports combining multiple data sources.",
	commands: new Map([["summary", reportSummaryCommand]]),
	defaultCommand: reportSummaryCommand,
};

// ============================================================================
// DOMAIN DEFINITION
// ============================================================================

/**
 * Subscription domain definition
 */
export const subscriptionDomain: DomainDefinition = {
	name: "subscription",
	description:
		"Manage F5 Distributed Cloud subscription, billing, quotas, and usage. View plan details, addon services, resource limits, usage metrics, payment methods, invoices, and generate comprehensive reports.",
	descriptionShort: "Subscription and billing management",
	descriptionMedium:
		"Manage subscription tier, addon services, quota limits, usage metrics, and billing information.",
	defaultCommand: showCommand,
	commands: new Map([["show", showCommand]]),
	subcommands: new Map([
		["plan", planSubcommands],
		["addon", addonSubcommands],
		["quota", quotaSubcommands],
		["usage", usageSubcommands],
		["billing", billingSubcommands],
		["report", reportSubcommands],
	]),
};

// Domain aliases
export const subscriptionAliases = ["sub", "billing", "quota"];
