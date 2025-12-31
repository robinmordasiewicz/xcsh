/**
 * Connection Summary Table
 * Displays profile connection status in a beautiful F5-branded table
 */

import { colors } from "../../../branding/index.js";
import { shouldUseColors } from "../../../output/resolver.js";

/**
 * Unicode box drawing characters (rounded corners)
 */
const UNICODE_BOX = {
	topLeft: "\u256D", // ╭
	topRight: "\u256E", // ╮
	bottomLeft: "\u2570", // ╰
	bottomRight: "\u256F", // ╯
	horizontal: "\u2500", // ─
	vertical: "\u2502", // │
	leftT: "\u251C", // ├
	rightT: "\u2524", // ┤
};

/**
 * ASCII box drawing characters (fallback)
 */
const ASCII_BOX = {
	topLeft: "+",
	topRight: "+",
	bottomLeft: "+",
	bottomRight: "+",
	horizontal: "-",
	vertical: "|",
	leftT: "+",
	rightT: "+",
};

/**
 * Connection info for display
 */
export interface ConnectionInfo {
	profileName: string;
	tenant: string;
	apiUrl: string;
	hasToken: boolean;
	namespace: string;
	isConnected: boolean;
	isValidated?: boolean;
	validationError?: string;
}

/**
 * Extract tenant from API URL
 */
export function extractTenantFromUrl(url: string): string {
	try {
		const parsed = new URL(url);
		const hostname = parsed.hostname;
		const parts = hostname.split(".");
		if (parts.length > 0 && parts[0]) {
			return parts[0];
		}
		return hostname;
	} catch {
		return "unknown";
	}
}

/**
 * Get auth status display value based on validation state
 */
function getAuthStatusValue(
	info: ConnectionInfo,
	colorStatus: (text: string, isGood: boolean) => string,
): string {
	if (!info.hasToken) {
		return colorStatus("\u2717 No token", false);
	}
	if (info.isValidated) {
		return colorStatus("\u2713 Authenticated", true);
	}
	if (info.validationError) {
		return colorStatus(`\u2717 ${info.validationError}`, false);
	}
	return colorStatus("\u26A0 Token not verified", false);
}

/**
 * Format connection table
 */
export function formatConnectionTable(
	info: ConnectionInfo,
	noColor: boolean = false,
): string[] {
	const useColors = shouldUseColors(undefined, noColor);
	const box = useColors ? UNICODE_BOX : ASCII_BOX;
	const borderColor = colors.red;
	const successColor = colors.green;
	const errorColor = colors.red;

	// Helper to colorize borders
	const colorBorder = (text: string) =>
		useColors ? `${borderColor}${text}${colors.reset}` : text;

	// Helper to colorize status
	const colorStatus = (text: string, isGood: boolean) =>
		useColors
			? `${isGood ? successColor : errorColor}${text}${colors.reset}`
			: text;

	// Build rows
	const rows = [
		{ label: "Profile", value: info.profileName },
		{ label: "Tenant", value: info.tenant },
		{ label: "API URL", value: info.apiUrl },
		{
			label: "Auth",
			value: getAuthStatusValue(info, colorStatus),
		},
		{ label: "Namespace", value: info.namespace || "default" },
		{
			label: "Status",
			value: info.isConnected
				? colorStatus("\u25CF Connected", true)
				: colorStatus("\u25CB Not connected", false),
		},
	];

	// Calculate widths
	const labelWidth = Math.max(...rows.map((r) => r.label.length));
	const valueWidth = Math.max(
		...rows.map((r) => stripAnsi(r.value).length),
		30,
	);
	const innerWidth = labelWidth + valueWidth + 5; // " | " separator + padding

	// Title
	const title = " Connection Summary ";
	const remainingWidth = innerWidth - title.length;
	const leftDashes = 1;
	const rightDashes = Math.max(0, remainingWidth - leftDashes);

	const lines: string[] = [];

	// Top border with title
	lines.push(
		colorBorder(box.topLeft + box.horizontal.repeat(leftDashes)) +
			title +
			colorBorder(box.horizontal.repeat(rightDashes) + box.topRight),
	);

	// Content rows
	for (const row of rows) {
		const paddedLabel = row.label.padEnd(labelWidth);
		const valueLen = stripAnsi(row.value).length;
		const paddedValue =
			row.value + " ".repeat(Math.max(0, valueWidth - valueLen));
		const content = ` ${paddedLabel} ${colorBorder(box.vertical)} ${paddedValue} `;
		lines.push(
			colorBorder(box.vertical) + content + colorBorder(box.vertical),
		);
	}

	// Bottom border
	lines.push(
		colorBorder(
			box.bottomLeft +
				box.horizontal.repeat(innerWidth) +
				box.bottomRight,
		),
	);

	return lines;
}

/**
 * Strip ANSI escape codes for length calculation
 */
function stripAnsi(str: string): string {
	// eslint-disable-next-line no-control-regex
	return str.replace(/\x1b\[[0-9;]*m/g, "");
}

/**
 * Build connection info from session
 */
export function buildConnectionInfo(
	profileName: string,
	apiUrl: string,
	hasToken: boolean,
	namespace: string,
	isConnected: boolean,
	isValidated?: boolean,
	validationError?: string,
): ConnectionInfo {
	const info: ConnectionInfo = {
		profileName,
		tenant: extractTenantFromUrl(apiUrl),
		apiUrl,
		hasToken,
		namespace,
		isConnected,
	};

	// Add optional validation fields only if they have values
	if (isValidated !== undefined) {
		info.isValidated = isValidated;
	}
	if (validationError) {
		info.validationError = validationError;
	}

	return info;
}
