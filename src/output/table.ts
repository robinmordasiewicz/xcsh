/**
 * Beautiful Table Formatter
 * Unicode tables with F5 brand red borders
 * Auto-disables colors when piped/non-TTY
 */

import { colors } from "../branding/index.js";
import type {
	BoxCharacters,
	ColumnDefinition,
	TableConfig,
	TableStyle,
} from "./types.js";
import { DEFAULT_TABLE_STYLE, PLAIN_TABLE_STYLE } from "./types.js";
import { shouldUseColors } from "./resolver.js";

/**
 * Unicode box drawing characters (rounded corners)
 */
const UNICODE_BOX: BoxCharacters = {
	topLeft: "\u256D", // ╭
	topRight: "\u256E", // ╮
	bottomLeft: "\u2570", // ╰
	bottomRight: "\u256F", // ╯
	horizontal: "\u2500", // ─
	vertical: "\u2502", // │
	leftT: "\u251C", // ├
	rightT: "\u2524", // ┤
	topT: "\u252C", // ┬
	bottomT: "\u2534", // ┴
	cross: "\u253C", // ┼
};

/**
 * ASCII box drawing characters
 */
const ASCII_BOX: BoxCharacters = {
	topLeft: "+",
	topRight: "+",
	bottomLeft: "+",
	bottomRight: "+",
	horizontal: "-",
	vertical: "|",
	leftT: "+",
	rightT: "+",
	topT: "+",
	bottomT: "+",
	cross: "+",
};

/**
 * Get box characters based on style
 */
function getBoxCharacters(style: TableStyle): BoxCharacters {
	return style.unicode ? UNICODE_BOX : ASCII_BOX;
}

/**
 * Apply color to text if colors are enabled
 */
function applyColor(text: string, color: string, useColors: boolean): string {
	if (!useColors) {
		return text;
	}
	return `${color}${text}${colors.reset}`;
}

/**
 * Wrap text to fit within max width, respecting word boundaries
 */
export function wrapText(text: string, maxWidth: number): string[] {
	if (text.length <= maxWidth) {
		return [text];
	}

	const lines: string[] = [];
	let remaining = text;

	while (remaining.length > 0) {
		if (remaining.length <= maxWidth) {
			lines.push(remaining);
			break;
		}

		// Find break point (prefer space)
		let breakPoint = maxWidth;
		for (let i = maxWidth - 1; i > 0; i--) {
			if (remaining[i] === " ") {
				breakPoint = i;
				break;
			}
		}

		lines.push(remaining.slice(0, breakPoint));
		remaining = remaining.slice(breakPoint).trimStart();
	}

	return lines;
}

/**
 * Get value from object using accessor
 */
function getValue(
	row: Record<string, unknown>,
	accessor: string | ((row: Record<string, unknown>) => string),
): string {
	if (typeof accessor === "function") {
		return accessor(row);
	}

	const value = row[accessor];
	if (value === null || value === undefined) {
		return "";
	}
	if (typeof value === "object") {
		// Handle labels object specially
		if (accessor === "labels") {
			return formatLabelsValue(value as Record<string, unknown>);
		}
		return JSON.stringify(value);
	}
	return String(value);
}

/**
 * Format labels as map[key:value key:value]
 */
function formatLabelsValue(labels: Record<string, unknown>): string {
	const entries = Object.entries(labels)
		.sort(([a], [b]) => a.localeCompare(b))
		.map(([k, v]) => `${k}:${v}`);

	if (entries.length === 0) {
		return "";
	}

	return `map[${entries.join(" ")}]`;
}

/**
 * Calculate column widths based on content
 */
function calculateColumnWidths(
	columns: ColumnDefinition[],
	rows: Record<string, unknown>[],
	maxTableWidth?: number,
): number[] {
	const widths = columns.map((col) => {
		// Start with header width
		let width = col.header.length;

		// Check all row values
		for (const row of rows) {
			const value = getValue(row, col.accessor);
			width = Math.max(width, value.length);
		}

		// Apply min/max constraints
		if (col.minWidth) {
			width = Math.max(width, col.minWidth);
		}
		if (col.maxWidth) {
			width = Math.min(width, col.maxWidth);
		}
		if (col.width) {
			width = col.width;
		}

		return width;
	});

	// If max table width is set, proportionally reduce columns
	if (maxTableWidth) {
		const totalWidth =
			widths.reduce((a, b) => a + b, 0) + columns.length * 3 + 1;
		if (totalWidth > maxTableWidth) {
			const ratio = (maxTableWidth - columns.length * 3 - 1) / totalWidth;
			for (let i = 0; i < widths.length; i++) {
				widths[i] = Math.max(5, Math.floor(widths[i]! * ratio));
			}
		}
	}

	return widths;
}

/**
 * Format data as a beautiful table with F5 branding
 */
export function formatBeautifulTable(
	data: Record<string, unknown>[],
	config: TableConfig,
	noColor: boolean = false,
): string {
	if (data.length === 0) {
		return "";
	}

	// Determine style and colors
	const useColors = shouldUseColors(undefined, noColor);
	const style = useColors
		? (config.style ?? DEFAULT_TABLE_STYLE)
		: PLAIN_TABLE_STYLE;
	const box = getBoxCharacters(style);
	const borderColor = style.borderColor ?? colors.red;

	// Calculate column widths
	const widths = calculateColumnWidths(config.columns, data, config.maxWidth);

	// Build table lines
	const lines: string[] = [];

	// Helper to colorize borders
	const colorBorder = (text: string) =>
		applyColor(text, borderColor, useColors && style.coloredBorders);

	// Build horizontal line
	const buildHorizontalLine = (
		left: string,
		mid: string,
		right: string,
		fill: string,
	): string => {
		const segments = widths.map((w) => fill.repeat(w + 2));
		return colorBorder(left + segments.join(mid) + right);
	};

	// Top border with optional title
	if (config.title) {
		const title = ` ${config.title} `;
		const totalWidth =
			widths.reduce((a, b) => a + b, 0) + widths.length * 3 - 1;
		const remainingWidth = totalWidth - title.length;
		const leftDashes = 1;
		const rightDashes = Math.max(0, remainingWidth - leftDashes);

		lines.push(
			colorBorder(box.topLeft + box.horizontal.repeat(leftDashes)) +
				title +
				colorBorder(box.horizontal.repeat(rightDashes) + box.topRight),
		);
	} else {
		lines.push(
			buildHorizontalLine(
				box.topLeft,
				box.topT,
				box.topRight,
				box.horizontal,
			),
		);
	}

	// Header row
	const headerCells = config.columns.map((col, i) => {
		const padding = widths[i]! - col.header.length;
		const leftPad = Math.floor(padding / 2);
		const rightPad = padding - leftPad;
		const content = " ".repeat(leftPad) + col.header + " ".repeat(rightPad);
		return ` ${content} `;
	});
	lines.push(
		colorBorder(box.vertical) +
			headerCells.join(colorBorder(box.vertical)) +
			colorBorder(box.vertical),
	);

	// Header separator
	lines.push(
		buildHorizontalLine(box.leftT, box.cross, box.rightT, box.horizontal),
	);

	// Data rows
	for (let rowIndex = 0; rowIndex < data.length; rowIndex++) {
		const row = data[rowIndex]!;

		// Get cell values and wrap text
		const cellValues = config.columns.map((col, i) => {
			const value = getValue(row, col.accessor) || "<None>";
			return config.wrapText !== false
				? wrapText(value, widths[i]!)
				: [value.slice(0, widths[i]!)];
		});

		// Find max lines needed for this row
		const maxLines = Math.max(...cellValues.map((c) => c.length));

		// Output each line of the row
		for (let lineIndex = 0; lineIndex < maxLines; lineIndex++) {
			const cells = cellValues.map((cellLines, i) => {
				const text = cellLines[lineIndex] ?? "";
				const padding = widths[i]! - text.length;
				const align = config.columns[i]?.align ?? "left";

				let content: string;
				if (align === "center") {
					const leftPad = Math.floor(padding / 2);
					const rightPad = padding - leftPad;
					content = " ".repeat(leftPad) + text + " ".repeat(rightPad);
				} else if (align === "right") {
					content = " ".repeat(padding) + text;
				} else {
					content = text + " ".repeat(padding);
				}

				return ` ${content} `;
			});

			lines.push(
				colorBorder(box.vertical) +
					cells.join(colorBorder(box.vertical)) +
					colorBorder(box.vertical),
			);
		}

		// Row separator (optional, or after last row as bottom border)
		if (rowIndex < data.length - 1 && config.rowSeparators) {
			lines.push(
				buildHorizontalLine(
					box.leftT,
					box.cross,
					box.rightT,
					box.horizontal,
				),
			);
		}
	}

	// Bottom border
	lines.push(
		buildHorizontalLine(
			box.bottomLeft,
			box.bottomT,
			box.bottomRight,
			box.horizontal,
		),
	);

	return lines.join("\n");
}

/**
 * Default columns for F5 XC resources
 */
export const DEFAULT_RESOURCE_COLUMNS: ColumnDefinition[] = [
	{ header: "NAMESPACE", accessor: "namespace", minWidth: 9 },
	{ header: "NAME", accessor: "name", minWidth: 10, maxWidth: 40 },
	{ header: "LABELS", accessor: "labels", minWidth: 10, maxWidth: 35 },
];

/**
 * Format F5 XC resource list as beautiful table
 * Convenience function for common use case
 */
export function formatResourceTable(
	data: unknown,
	noColor: boolean = false,
): string {
	const items = extractItems(data);
	if (items.length === 0) {
		return "";
	}

	return formatBeautifulTable(
		items,
		{
			columns: DEFAULT_RESOURCE_COLUMNS,
			wrapText: true,
		},
		noColor,
	);
}

/**
 * Extract items from various data formats
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

/**
 * Format a simple key-value box (like whoami)
 * Uses F5 brand red borders with title
 */
export function formatKeyValueBox(
	data: Array<{ label: string; value: string }>,
	title: string,
	noColor: boolean = false,
): string {
	if (data.length === 0) {
		return "";
	}

	const useColors = shouldUseColors(undefined, noColor);
	const box = useColors ? UNICODE_BOX : ASCII_BOX;
	const borderColor = colors.red;

	const colorBorder = (text: string) =>
		applyColor(text, borderColor, useColors);

	// Find max label width for alignment
	const maxLabelWidth = Math.max(...data.map((d) => d.label.length));

	// Format content lines
	const contentLines = data.map((d) => {
		const paddedLabel = d.label.padEnd(maxLabelWidth);
		return `${paddedLabel}:  ${d.value}`;
	});

	// Calculate box width
	const titleText = ` ${title} `;
	const maxContentWidth = Math.max(
		...contentLines.map((l) => l.length),
		titleText.length,
	);
	const innerWidth = maxContentWidth + 2;

	const lines: string[] = [];

	// Top border with title
	const remainingWidth = innerWidth - titleText.length;
	const leftDashes = 1;
	const rightDashes = Math.max(0, remainingWidth - leftDashes);
	lines.push(
		colorBorder(box.topLeft + box.horizontal.repeat(leftDashes)) +
			titleText +
			colorBorder(box.horizontal.repeat(rightDashes) + box.topRight),
	);

	// Content lines
	for (const text of contentLines) {
		const padding = innerWidth - text.length;
		lines.push(
			colorBorder(box.vertical) +
				` ${text}${" ".repeat(Math.max(0, padding - 1))}` +
				colorBorder(box.vertical),
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

	return lines.join("\n");
}
