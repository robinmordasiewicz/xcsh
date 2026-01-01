/**
 * Domain Formatter
 * Unified output formatting for custom domains (login, cloudstatus, etc.)
 * Provides DRY formatting that all domains can use consistently
 */

import type { OutputFormat } from "./types.js";
import { formatOutput, formatJSON, formatYAML } from "./formatter.js";
import { formatKeyValueBox, formatResourceTable } from "./table.js";
import { parseOutputFlag, shouldUseColors } from "./resolver.js";

/**
 * Format options for domain commands
 */
export interface DomainFormatOptions {
	/** Output format (json, yaml, table, tsv, none) */
	format: OutputFormat;
	/** Disable colors (for piped output) */
	noColor: boolean;
	/** Title for key-value box displays */
	title?: string;
}

/**
 * Key-value data structure for displays like whoami
 */
export interface KeyValueData {
	[key: string]: string | number | boolean | null | undefined;
}

/**
 * Parse output flags from domain command arguments
 * Returns format options with session default fallback
 *
 * @param args - Command arguments array
 * @param sessionFormat - Default format from session (if any)
 * @returns Format options and remaining arguments
 */
export function parseDomainOutputFlags(
	args: string[],
	sessionFormat: OutputFormat = "table",
): { options: DomainFormatOptions; remainingArgs: string[] } {
	const { format: parsedFormat, remainingArgs } = parseOutputFlag(args);

	// Also check for --no-color flag
	let noColor = false;
	const finalArgs: string[] = [];
	for (const arg of remainingArgs) {
		if (arg === "--no-color") {
			noColor = true;
		} else {
			finalArgs.push(arg);
		}
	}

	// Auto-detect color support if not explicitly disabled
	const effectiveNoColor = noColor || !shouldUseColors();

	return {
		options: {
			format: parsedFormat ?? sessionFormat,
			noColor: effectiveNoColor,
		},
		remainingArgs: finalArgs,
	};
}

/**
 * Format domain output data
 * Routes to appropriate formatter based on format type
 * Returns string[] for DomainCommandResult compatibility
 *
 * @param data - Data to format (object, array, or primitive)
 * @param options - Format options
 * @returns Formatted output as string array
 */
export function formatDomainOutput(
	data: unknown,
	options: DomainFormatOptions,
): string[] {
	const { format, noColor } = options;

	// Handle none format - return empty array
	if (format === "none") {
		return [];
	}

	// Use existing formatOutput for most formats
	const formatted = formatOutput(data, format, noColor);

	// Split into lines for DomainCommandResult
	if (formatted === "") {
		return [];
	}

	return formatted.split("\n");
}

/**
 * Format key-value data (like whoami, context show, profile show)
 * Uses beautiful box format for table, structured format for json/yaml
 *
 * @param data - Key-value data object
 * @param options - Format options (must include title for table format)
 * @returns Formatted output as string array
 */
export function formatKeyValueOutput(
	data: KeyValueData,
	options: DomainFormatOptions,
): string[] {
	const { format, noColor, title } = options;

	// Handle none format
	if (format === "none") {
		return [];
	}

	// For JSON/YAML - use standard formatters with the raw data
	if (format === "json") {
		return formatJSON(data).split("\n");
	}

	if (format === "yaml") {
		return formatYAML(data).split("\n");
	}

	// For TSV - format as tab-separated key-value pairs
	if (format === "tsv") {
		const lines: string[] = [];
		for (const [key, value] of Object.entries(data)) {
			if (value !== null && value !== undefined) {
				lines.push(`${key}\t${String(value)}`);
			}
		}
		return lines;
	}

	// For table/text - use the beautiful key-value box
	const boxData: Array<{ label: string; value: string }> = [];
	for (const [key, value] of Object.entries(data)) {
		if (value !== null && value !== undefined) {
			// Convert camelCase to Title Case for display
			const label = formatLabel(key);
			boxData.push({ label, value: String(value) });
		}
	}

	if (boxData.length === 0) {
		return [];
	}

	const boxOutput = formatKeyValueBox(boxData, title ?? "Info", noColor);
	return boxOutput.split("\n");
}

/**
 * Format array/list data for domain commands
 * Uses table format for table/text, structured for json/yaml
 *
 * @param data - Array of objects to format
 * @param options - Format options
 * @returns Formatted output as string array
 */
export function formatListOutput(
	data: unknown[],
	options: DomainFormatOptions,
): string[] {
	const { format, noColor } = options;

	// Handle none format
	if (format === "none") {
		return [];
	}

	// For JSON/YAML - use standard formatters
	if (format === "json") {
		return formatJSON(data).split("\n");
	}

	if (format === "yaml") {
		return formatYAML(data).split("\n");
	}

	// For TSV - format as tab-separated values
	if (format === "tsv") {
		if (data.length === 0) {
			return [];
		}

		// Get all keys from first item
		const firstItem = data[0];
		if (typeof firstItem !== "object" || firstItem === null) {
			return data.map((item) => String(item));
		}

		const keys = Object.keys(firstItem as Record<string, unknown>);
		const lines: string[] = [];

		for (const item of data) {
			if (typeof item === "object" && item !== null) {
				const record = item as Record<string, unknown>;
				const values = keys.map((k) => {
					const val = record[k];
					if (val === null || val === undefined) return "";
					if (typeof val === "object") return JSON.stringify(val);
					return String(val);
				});
				lines.push(values.join("\t"));
			}
		}

		return lines;
	}

	// For table/text - use resource table formatter
	const tableOutput = formatResourceTable(data, noColor);
	if (tableOutput === "") {
		return [];
	}

	return tableOutput.split("\n");
}

/**
 * Convert camelCase or snake_case key to Title Case label
 */
function formatLabel(key: string): string {
	// Handle snake_case
	if (key.includes("_")) {
		return key
			.split("_")
			.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
			.join(" ");
	}

	// Handle camelCase
	const withSpaces = key.replace(/([A-Z])/g, " $1");
	return withSpaces.charAt(0).toUpperCase() + withSpaces.slice(1);
}

/**
 * Re-export types for convenience
 */
export type { OutputFormat } from "./types.js";
