/**
 * Output Types
 * Unified type definitions for output formatting
 */

/**
 * Supported output format types
 */
export type OutputFormat =
	| "json"
	| "yaml"
	| "table"
	| "text"
	| "tsv"
	| "none"
	| "spec";

/**
 * Table styling configuration
 */
export interface TableStyle {
	/** Use Unicode box-drawing characters (╭╮) vs ASCII (+|-) */
	unicode: boolean;
	/** Apply F5 brand red color to borders */
	coloredBorders: boolean;
	/** Border color ANSI code (defaults to F5 red) */
	borderColor?: string;
	/** Header styling */
	headerStyle?: "bold" | "normal" | "dim";
}

/**
 * Column definition for tables
 */
export interface ColumnDefinition {
	/** Column header name */
	header: string;
	/** Data accessor key or function */
	accessor: string | ((row: Record<string, unknown>) => string);
	/** Fixed column width */
	width?: number;
	/** Minimum column width */
	minWidth?: number;
	/** Maximum column width */
	maxWidth?: number;
	/** Column alignment */
	align?: "left" | "center" | "right";
	/** Priority for column display (lower = higher priority) */
	priority?: number;
}

/**
 * Table configuration
 */
export interface TableConfig {
	/** Column definitions */
	columns: ColumnDefinition[];
	/** Table style options */
	style?: TableStyle;
	/** Maximum total table width */
	maxWidth?: number;
	/** Enable text wrapping in cells */
	wrapText?: boolean;
	/** Show row separators between data rows */
	rowSeparators?: boolean;
	/** Title to display in the table header */
	title?: string;
}

/**
 * Output context for format resolution
 */
export interface OutputContext {
	/** Explicit format from CLI flag (--output) */
	cliFormat?: OutputFormat;
	/** Format from environment variable (F5XC_OUTPUT_FORMAT) */
	envFormat?: OutputFormat;
	/** Format from config file */
	configFormat?: OutputFormat;
	/** Is this an interactive session (REPL) */
	isInteractive: boolean;
	/** Is output to a terminal (TTY) */
	isTTY: boolean;
	/** No color flag (--no-color) */
	noColor?: boolean;
}

/**
 * Flag specification for command documentation
 */
export interface FlagSpec {
	/** Flag name (long form, e.g., "--output") */
	name: string;
	/** Short alias (e.g., "-o") */
	alias?: string;
	/** Flag description */
	description: string;
	/** Expected value type */
	type?: "string" | "boolean" | "number";
	/** Default value */
	default?: string;
	/** Is this flag required */
	required?: boolean;
	/** Allowed values for string type */
	choices?: string[];
}

/**
 * Example specification for command documentation
 */
export interface ExampleSpec {
	/** Command example */
	command: string;
	/** Description of what it does */
	description: string;
}

/**
 * Command specification for AI assistants
 * JSON schema output format for --spec flag
 */
export interface CommandSpec {
	/** Command name (e.g., "cloudstatus status") */
	command: string;
	/** Command description */
	description: string;
	/** Usage pattern */
	usage: string;
	/** Available flags/options */
	flags: FlagSpec[];
	/** Example invocations */
	examples: ExampleSpec[];
	/** Supported output formats */
	outputFormats: string[];
	/** Related commands */
	related?: string[];
	/** Command category/domain */
	category?: string;
}

/**
 * Box drawing character sets
 */
export interface BoxCharacters {
	topLeft: string;
	topRight: string;
	bottomLeft: string;
	bottomRight: string;
	horizontal: string;
	vertical: string;
	leftT: string;
	rightT: string;
	topT: string;
	bottomT: string;
	cross: string;
}

/**
 * Default table style using Unicode characters with F5 red borders
 */
export const DEFAULT_TABLE_STYLE: TableStyle = {
	unicode: true,
	coloredBorders: true,
	headerStyle: "bold",
};

/**
 * Plain ASCII table style (no colors, ASCII characters)
 */
export const PLAIN_TABLE_STYLE: TableStyle = {
	unicode: false,
	coloredBorders: false,
	headerStyle: "normal",
};

/**
 * Check if a string is a valid output format
 */
export function isValidOutputFormat(format: string): format is OutputFormat {
	return ["json", "yaml", "table", "text", "tsv", "none", "spec"].includes(
		format.toLowerCase(),
	);
}
