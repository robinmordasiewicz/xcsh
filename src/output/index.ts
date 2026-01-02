/**
 * Output Module
 * Unified output formatting with beautiful tables and color auto-detection
 */

// Core formatter functions
export {
	formatOutput,
	formatJSON,
	formatYAML,
	formatTable,
	formatTSV,
	parseOutputFormat,
	formatAPIError,
} from "./formatter.js";
export type { FormatterConfig } from "./formatter.js";

// Domain-specific formatting (for unified output handling)
export {
	formatDomainOutput,
	parseDomainOutputFlags,
	type DomainFormatOptions,
} from "./domain-formatter.js";

// Types
export type {
	OutputFormat,
	TableStyle,
	TableConfig,
	ColumnDefinition,
	OutputContext,
	CommandSpec,
	FlagSpec,
	ExampleSpec,
	BoxCharacters,
} from "./types.js";
export {
	isValidOutputFormat,
	DEFAULT_TABLE_STYLE,
	PLAIN_TABLE_STYLE,
	ALL_OUTPUT_FORMATS,
	OUTPUT_FORMAT_HELP,
} from "./types.js";

// Format resolution
export {
	resolveOutputFormat,
	getOutputFormatFromEnv,
	parseOutputFlag,
	parseSpecFlag,
	parseOutputFlags,
	shouldUseColors,
	buildOutputContext,
	OUTPUT_FORMAT_ENV_VAR,
} from "./resolver.js";

// Beautiful table formatting
export {
	formatBeautifulTable,
	formatResourceTable,
	formatKeyValueBox,
	wrapText,
	DEFAULT_RESOURCE_COLUMNS,
} from "./table.js";

// Command spec generation (for --spec flag)
export {
	buildCommandSpec,
	formatSpec,
	getCommandSpec,
	listAllCommandSpecs,
	GLOBAL_FLAGS,
	buildCloudstatusSpecs,
	buildLoginSpecs,
} from "./spec.js";
