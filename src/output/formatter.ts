/**
 * Output Formatter
 * Formats API responses as JSON, YAML, or table
 * Unified formatting with beautiful tables and color auto-detection
 */

import YAML from "yaml";
import { formatResourceTable } from "./table.js";
import { shouldUseColors } from "./resolver.js";

// Re-export types from types.ts for backward compatibility
export type { OutputFormat } from "./types.js";
export { isValidOutputFormat } from "./types.js";

// Import OutputFormat for local use
import type { OutputFormat } from "./types.js";

/**
 * Formatter configuration
 */
export interface FormatterConfig {
	/** Output format */
	format: import("./types.js").OutputFormat;
	/** Column widths for table format (optional) */
	columnWidths?: number[];
	/** Priority columns to show first */
	priorityColumns?: string[];
	/** Disable colors (for piped output) */
	noColor?: boolean;
}

/**
 * Format data according to specified format
 */
export function formatOutput(
	data: unknown,
	format: import("./types.js").OutputFormat = "table",
	noColor: boolean = false,
): string {
	if (format === "none") {
		return "";
	}

	// Auto-detect color support if not explicitly disabled
	const useNoColor = noColor || !shouldUseColors();

	switch (format) {
		case "json":
			return formatJSON(data);
		case "yaml":
			return formatYAML(data);
		case "table":
		case "text":
			return formatTable(data, useNoColor);
		case "tsv":
			return formatTSV(data);
		case "spec":
			// Spec format is handled at command level
			return formatJSON(data);
		default:
			return formatTable(data, useNoColor);
	}
}

/**
 * Format as pretty-printed JSON
 */
export function formatJSON(data: unknown): string {
	return JSON.stringify(data, null, 2);
}

/**
 * Format as YAML
 */
export function formatYAML(data: unknown): string {
	return YAML.stringify(data, { indent: 2 });
}

/**
 * Extract items from list response
 */
function extractItems(data: unknown): Record<string, unknown>[] {
	// Handle map with "items" key
	if (data && typeof data === "object" && "items" in data) {
		const items = (data as { items: unknown[] }).items;
		if (Array.isArray(items)) {
			return items.filter(
				(item): item is Record<string, unknown> =>
					item !== null && typeof item === "object",
			);
		}
	}

	// Handle array directly
	if (Array.isArray(data)) {
		return data.filter(
			(item): item is Record<string, unknown> =>
				item !== null && typeof item === "object",
		);
	}

	// Single item, wrap it
	if (data && typeof data === "object") {
		return [data as Record<string, unknown>];
	}

	return [];
}

// NOTE: Helper functions getStringField, formatLabels, and wrapText
// have been moved to table.ts for direct use in beautiful table formatting.

/**
 * Format as beautiful Unicode table with F5 red borders
 * Falls back to plain ASCII when colors are disabled
 */
export function formatTable(data: unknown, noColor: boolean = false): string {
	return formatResourceTable(data, noColor);
}

/**
 * Format as tab-separated values
 */
export function formatTSV(data: unknown): string {
	const items = extractItems(data);
	if (items.length === 0) {
		return "";
	}

	// Get all keys
	const allKeys = new Set<string>();
	for (const item of items) {
		Object.keys(item).forEach((k) => allKeys.add(k));
	}

	// Priority order
	const priority = ["name", "namespace", "status", "created", "modified"];
	const headers = [
		...priority.filter((p) => allKeys.has(p)),
		...[...allKeys].filter((k) => !priority.includes(k)).sort(),
	];

	// Build rows
	const lines: string[] = [];
	for (const item of items) {
		const values = headers.map((h) => {
			const val = item[h];
			if (val === null || val === undefined) return "";
			if (typeof val === "object") return JSON.stringify(val);
			return String(val);
		});
		lines.push(values.join("\t"));
	}

	return lines.join("\n");
}

/**
 * Parse output format from string
 */
export function parseOutputFormat(format: string): OutputFormat {
	switch (format.toLowerCase()) {
		case "json":
			return "json";
		case "yaml":
			return "yaml";
		case "table":
		case "text":
		case "":
			return "table";
		case "tsv":
			return "tsv";
		case "none":
			return "none";
		default:
			return "table";
	}
}

/**
 * Format API error with helpful context
 */
export function formatAPIError(
	statusCode: number,
	body: unknown,
	operation: string,
): string {
	const lines: string[] = [];
	lines.push(`ERROR: ${operation} failed (HTTP ${statusCode})`);

	// Try to extract error details
	if (body && typeof body === "object") {
		const errResp = body as Record<string, unknown>;
		if (errResp.message) {
			lines.push(`  Message: ${errResp.message}`);
		}
		if (errResp.code) {
			lines.push(`  Code: ${errResp.code}`);
		}
		if (errResp.details) {
			lines.push(`  Details: ${errResp.details}`);
		}
	}

	// Add hints based on status code
	switch (statusCode) {
		case 401:
			lines.push(
				"\nHint: Authentication failed. Check your credentials with 'login profile show'",
			);
			break;
		case 403:
			lines.push(
				"\nHint: Permission denied. You may not have access to this resource.",
			);
			break;
		case 404:
			lines.push(
				"\nHint: Resource not found. Verify the name and namespace are correct.",
			);
			break;
		case 409:
			lines.push(
				"\nHint: Conflict - resource may already exist or be in a conflicting state.",
			);
			break;
		case 429:
			lines.push("\nHint: Rate limited. Please wait and try again.");
			break;
		case 500:
		case 502:
		case 503:
			lines.push(
				"\nHint: Server error. Please try again later or contact support.",
			);
			break;
	}

	return lines.join("\n");
}
