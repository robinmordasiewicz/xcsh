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
	return formatWhoamiBox(info);
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
	output.namespace = info.namespace;
	output.serverUrl = info.serverUrl;
	output.isAuthenticated = info.isAuthenticated;

	return [JSON.stringify(output, null, 2)];
}

/**
 * Get auth status display string based on validation state
 */
function getAuthStatusDisplay(info: WhoamiInfo): string {
	if (!info.isAuthenticated) {
		return "Not authenticated";
	}
	if (info.isValidated) {
		return "\u2713 Authenticated";
	}
	if (info.validationError) {
		return `\u2717 ${info.validationError}`;
	}
	return "\u26A0 Token not verified";
}

/**
 * Format whoami with box decoration
 * Uses F5 brand red for the box frame
 */
function formatWhoamiBox(info: WhoamiInfo): string[] {
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
	contentLines.push({ label: "Namespace", value: info.namespace });
	contentLines.push({ label: "Server", value: info.serverUrl });
	contentLines.push({
		label: "Auth",
		value: getAuthStatusDisplay(info),
	});

	// Find max label width for alignment
	const maxLabelWidth = Math.max(...contentLines.map((c) => c.label.length));

	// Calculate formatted content lines
	const formattedContent = contentLines.map((c) => {
		const paddedLabel = c.label.padEnd(maxLabelWidth);
		return `${paddedLabel}:  ${c.value}`;
	});

	// Calculate dynamic width based on longest content
	const headerTitle = "Connection Info";
	const allTextLines = [...formattedContent, headerTitle];
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

	return parts.join(" | ");
}
