/**
 * Environment Variable Registry for xcsh.
 * Centralized registry matching the Go version's professional structure.
 * Used for dynamic help generation with column alignment and flag notation.
 */

import { ENV_PREFIX, CONFIG_FILE_NAME, DOCS_URL } from "../branding/index.js";

export interface EnvVar {
	name: string;
	description: string;
	relatedFlag: string; // Empty string if no related flag
	required?: boolean;
}

/**
 * Registry of all environment variables used by the CLI.
 * This is the single source of truth for environment variable documentation.
 */
export const EnvVarRegistry: EnvVar[] = [
	{
		name: `${ENV_PREFIX}_API_URL`,
		description: "API endpoint URL",
		relatedFlag: "",
		required: true,
	},
	{
		name: `${ENV_PREFIX}_API_TOKEN`,
		description: "API authentication token",
		relatedFlag: "",
		required: true,
	},
	{
		name: `${ENV_PREFIX}_NAMESPACE`,
		description: "Default namespace",
		relatedFlag: "-ns",
	},
	{
		name: `${ENV_PREFIX}_OUTPUT_FORMAT`,
		description: "Output format (json, yaml, table)",
		relatedFlag: "-o",
	},
	{
		name: `${ENV_PREFIX}_LOGO`,
		description: "Logo display mode (auto, image, ascii, both, none)",
		relatedFlag: "--logo",
	},
	{
		name: "NO_COLOR",
		description: "Disable color output",
		relatedFlag: "--no-color",
	},
];

/**
 * Format environment variables section with dynamic column alignment.
 * Matches the Go version's professional formatting with [flag] notation.
 */
export function formatEnvVarsSection(): string[] {
	const maxLen = Math.max(...EnvVarRegistry.map((e) => e.name.length));
	const lines: string[] = ["ENVIRONMENT VARIABLES"];

	for (const env of EnvVarRegistry) {
		const padding = " ".repeat(maxLen - env.name.length + 3);
		const flagNote = env.relatedFlag ? ` [${env.relatedFlag}]` : "";
		lines.push(`  ${env.name}${padding}${env.description}${flagNote}`);
	}

	return lines;
}

/**
 * Format configuration section with precedence order.
 * Matches the Go version's layout.
 */
export function formatConfigSection(): string[] {
	return [
		"CONFIGURATION",
		`  Config file:  ~/${CONFIG_FILE_NAME}`,
		"  Priority:     CLI flags > environment variables > config file > defaults",
		"",
		"DOCUMENTATION",
		`  ${DOCS_URL}`,
	];
}
