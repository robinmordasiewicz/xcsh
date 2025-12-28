/**
 * Whoami Formatter
 * Format whoami info for display in various formats
 */

import type { WhoamiInfo, WhoamiOptions } from "./types.js";
import { colors } from "../../../branding/index.js";

/**
 * Format whoami info as text for display
 * Returns array of lines to display
 * Omits fields that are undefined (no "Unknown" placeholders)
 */
export function formatWhoami(
	info: WhoamiInfo,
	options: WhoamiOptions = {},
): string[] {
	// Not authenticated - show login prompt
	if (!info.isAuthenticated) {
		return ["Enter '/login' to authenticate"];
	}

	// JSON output
	if (options.json) {
		return formatWhoamiJson(info);
	}

	// Text output with box
	return formatWhoamiBox(info, options);
}

/**
 * Format whoami as JSON
 */
function formatWhoamiJson(info: WhoamiInfo): string[] {
	// Build clean object without undefined values
	const output: Record<string, unknown> = {};

	if (info.tenant) output.tenant = info.tenant;
	if (info.username) output.username = info.username;
	if (info.email) output.email = info.email;
	if (info.tier) output.tier = info.tier;
	output.namespace = info.namespace;
	output.serverUrl = info.serverUrl;
	output.isAuthenticated = info.isAuthenticated;

	if (info.quotas) {
		output.quotas = {
			totalLimits: info.quotas.totalLimits,
			limitsAtRisk: info.quotas.limitsAtRisk,
			limitsExceeded: info.quotas.limitsExceeded,
			objects: info.quotas.objects?.map((q) => ({
				name: q.name,
				displayName: q.displayName,
				usage: q.usage,
				limit: q.limit,
				percentage: Math.round(q.percentage),
			})),
		};
	}

	if (info.addons && info.addons.length > 0) {
		output.addons = info.addons.map((a) => ({
			name: a.name,
			displayName: a.displayName,
			state: a.state,
		}));
	}

	return [JSON.stringify(output, null, 2)];
}

/**
 * Format whoami with box decoration
 * Uses F5 brand red for the box frame
 */
function formatWhoamiBox(info: WhoamiInfo, options: WhoamiOptions): string[] {
	const lines: string[] = [];
	const red = colors.red;
	const reset = colors.reset;

	// Box drawing characters
	const BOX = {
		topLeft: "\u256D",
		topRight: "\u256E",
		bottomLeft: "\u2570",
		bottomRight: "\u256F",
		horizontal: "\u2500",
		vertical: "\u2502",
		leftT: "\u251C",
		rightT: "\u2524",
	};

	// Build all content lines first to calculate width
	const contentLines: Array<{ label: string; value: string }> = [];

	if (info.tenant) {
		contentLines.push({ label: "Tenant", value: info.tenant });
	}
	if (info.email) {
		contentLines.push({ label: "User", value: info.email });
	} else if (info.username) {
		contentLines.push({ label: "User", value: info.username });
	}
	if (info.tier) {
		contentLines.push({ label: "Tier", value: info.tier });
	}
	contentLines.push({ label: "Namespace", value: info.namespace });
	contentLines.push({ label: "Server", value: info.serverUrl });
	contentLines.push({
		label: "Auth",
		value: info.isAuthenticated
			? "\u2713 Authenticated"
			: "Not authenticated",
	});

	// Build quota lines if included
	const quotaLines: string[] = [];
	if (
		info.quotas &&
		info.quotas.objects &&
		(options.includeQuotas || options.verbose)
	) {
		for (const quota of info.quotas.objects.slice(0, 10)) {
			const pct = Math.round(quota.percentage);
			quotaLines.push(
				`${quota.displayName}:  ${quota.usage}/${quota.limit} (${pct}%)`,
			);
		}
		if (info.quotas.objects.length > 10) {
			const remaining = info.quotas.objects.length - 10;
			quotaLines.push(`... and ${remaining} more`);
		}
	}

	// Build addon lines if included
	const addonLines: string[] = [];
	if (
		info.addons &&
		info.addons.length > 0 &&
		(options.includeAddons || options.verbose)
	) {
		for (const addon of info.addons) {
			addonLines.push(`\u2713 ${addon.displayName}`);
		}
	}

	// Find max label width for alignment
	const maxLabelWidth = Math.max(...contentLines.map((c) => c.label.length));

	// Calculate formatted content lines
	const formattedContent = contentLines.map((c) => {
		const paddedLabel = c.label.padEnd(maxLabelWidth);
		return `${paddedLabel}:  ${c.value}`;
	});

	// Calculate dynamic width based on longest content
	const headerTitles = ["Connection Info", "Quota Usage", "Active Addons"];
	const allTextLines = [
		...formattedContent,
		...quotaLines,
		...addonLines,
		...headerTitles,
	];
	const maxContentWidth = Math.max(...allTextLines.map((l) => l.length));
	const innerWidth = maxContentWidth + 2; // Add padding

	// Top border with title
	const title = " Connection Info ";
	const remainingWidth = innerWidth - title.length;
	const leftDashes = 1;
	const rightDashes = Math.max(0, remainingWidth - leftDashes);
	lines.push(
		`${red}${BOX.topLeft}${BOX.horizontal.repeat(leftDashes)}${reset}${title}${red}${BOX.horizontal.repeat(rightDashes)}${BOX.topRight}${reset}`,
	);

	// Content lines
	for (const text of formattedContent) {
		const padding = innerWidth - text.length;
		lines.push(
			`${red}${BOX.vertical}${reset} ${text}${" ".repeat(Math.max(0, padding - 1))}${red}${BOX.vertical}${reset}`,
		);
	}

	// Quota section if included
	if (quotaLines.length > 0) {
		const quotaTitle = " Quota Usage ";
		const quotaRemaining = innerWidth - quotaTitle.length;
		const quotaLeft = 1;
		const quotaRight = Math.max(0, quotaRemaining - quotaLeft);
		lines.push(
			`${red}${BOX.leftT}${BOX.horizontal.repeat(quotaLeft)}${reset}${quotaTitle}${red}${BOX.horizontal.repeat(quotaRight)}${BOX.rightT}${reset}`,
		);

		for (const text of quotaLines) {
			const padding = innerWidth - text.length;
			lines.push(
				`${red}${BOX.vertical}${reset} ${text}${" ".repeat(Math.max(0, padding - 1))}${red}${BOX.vertical}${reset}`,
			);
		}
	}

	// Addons section if included
	if (addonLines.length > 0) {
		const addonTitle = " Active Addons ";
		const addonRemaining = innerWidth - addonTitle.length;
		const addonLeft = 1;
		const addonRight = Math.max(0, addonRemaining - addonLeft);
		lines.push(
			`${red}${BOX.leftT}${BOX.horizontal.repeat(addonLeft)}${reset}${addonTitle}${red}${BOX.horizontal.repeat(addonRight)}${BOX.rightT}${reset}`,
		);

		for (const text of addonLines) {
			const padding = innerWidth - text.length;
			lines.push(
				`${red}${BOX.vertical}${reset} ${text}${" ".repeat(Math.max(0, padding - 1))}${red}${BOX.vertical}${reset}`,
			);
		}
	}

	// Bottom border
	lines.push(
		`${red}${BOX.bottomLeft}${BOX.horizontal.repeat(innerWidth)}${BOX.bottomRight}${reset}`,
	);

	return lines;
}

/**
 * Format whoami for compact banner display
 * Single line or minimal output
 */
export function formatWhoamiCompact(info: WhoamiInfo): string {
	if (!info.isAuthenticated) {
		return "Not authenticated - Enter '/login' to authenticate";
	}

	const parts: string[] = [];

	if (info.tenant) {
		parts.push(info.tenant);
	}

	if (info.email) {
		parts.push(info.email);
	} else if (info.username) {
		parts.push(info.username);
	}

	if (info.tier) {
		parts.push(`[${info.tier}]`);
	}

	return parts.join(" | ");
}
