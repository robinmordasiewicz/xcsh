/**
 * Command Specification Generator
 * Generates AI-friendly JSON schema output for --spec flag
 */

import type { CommandSpec, ExampleSpec, FlagSpec } from "./types.js";
import { CLI_VERSION } from "../branding/index.js";
import { customDomains } from "../domains/registry.js";
import { domainRegistry, validActions } from "../types/domains.js";

/**
 * Full CLI specification for documentation generation
 * Matches the format expected by scripts/generate-docs.py
 */
export interface CLISpec {
	version: string;
	global_flags: GlobalFlagSpec[];
	commands: CLICommandSpec[];
}

/**
 * Global flag specification for CLI documentation
 */
export interface GlobalFlagSpec {
	name: string;
	type: string;
	description: string;
	shorthand: string;
	default: string;
}

/**
 * Command specification for CLI documentation
 * Matches Python dataclass structure in generate-docs.py
 */
export interface CLICommandSpec {
	path: string[];
	use: string;
	short: string;
	long: string;
	example: string;
	aliases: string[];
	flags: GlobalFlagSpec[];
	subcommands: CLICommandSpec[];
}

/**
 * Build a command specification object
 */
export function buildCommandSpec(options: {
	command: string;
	description: string;
	usage?: string;
	flags?: FlagSpec[];
	examples?: ExampleSpec[];
	outputFormats?: string[];
	related?: string[];
	category?: string;
}): CommandSpec {
	const spec: CommandSpec = {
		command: options.command,
		description: options.description,
		usage: options.usage ?? `xcsh ${options.command} [options]`,
		flags: options.flags ?? [],
		examples: options.examples ?? [],
		outputFormats: options.outputFormats ?? ["table", "json", "yaml"],
	};

	// Only add optional properties if they have values (exactOptionalPropertyTypes)
	if (options.related !== undefined) {
		spec.related = options.related;
	}
	if (options.category !== undefined) {
		spec.category = options.category;
	}

	return spec;
}

/**
 * Format command spec as JSON string
 */
export function formatSpec(spec: CommandSpec): string {
	return JSON.stringify(spec, null, 2);
}

/**
 * Global flags that apply to all commands
 */
export const GLOBAL_FLAGS: FlagSpec[] = [
	{
		name: "--output",
		alias: "-o",
		description: "Output format",
		type: "string",
		default: "table",
		choices: ["table", "json", "yaml", "tsv"],
	},
	{
		name: "--namespace",
		alias: "-ns",
		description: "Namespace to use for the operation",
		type: "string",
	},
	{
		name: "--no-color",
		description: "Disable colored output",
		type: "boolean",
		default: "false",
	},
	{
		name: "--spec",
		description: "Output command specification as JSON (for AI assistants)",
		type: "boolean",
	},
	{
		name: "--help",
		alias: "-h",
		description: "Show help information",
		type: "boolean",
	},
	{
		name: "--version",
		alias: "-v",
		description: "Show version number",
		type: "boolean",
	},
];

/**
 * Build spec for cloudstatus domain commands
 */
export function buildCloudstatusSpecs(): Record<string, CommandSpec> {
	return {
		status: buildCommandSpec({
			command: "cloudstatus status",
			description:
				"Get the overall health indicator for F5 Distributed Cloud services. Returns status level (operational, degraded, major outage) with description.",
			usage: "xcsh cloudstatus status [--quiet]",
			flags: [
				{
					name: "--quiet",
					alias: "-q",
					description:
						"Return exit code only (0=operational, 1=degraded, 2=outage)",
					type: "boolean",
				},
			],
			examples: [
				{
					command: "xcsh cloudstatus status",
					description: "Check current F5 XC service status",
				},
				{
					command:
						"xcsh cloudstatus status --quiet && echo 'All systems operational'",
					description: "Use in scripts for health checks",
				},
				{
					command: "xcsh cloudstatus status --output json",
					description: "Get status as JSON for automation",
				},
			],
			category: "cloudstatus",
			related: [
				"cloudstatus summary",
				"cloudstatus components",
				"cloudstatus incidents",
			],
		}),

		summary: buildCommandSpec({
			command: "cloudstatus summary",
			description:
				"Get complete status summary including overall health, component status, and active incidents.",
			usage: "xcsh cloudstatus summary",
			examples: [
				{
					command: "xcsh cloudstatus summary",
					description: "View full infrastructure health overview",
				},
				{
					command: "xcsh cloudstatus summary --output json",
					description: "Get complete summary as JSON",
				},
			],
			category: "cloudstatus",
			related: ["cloudstatus status", "cloudstatus components"],
		}),

		components: buildCommandSpec({
			command: "cloudstatus components",
			description:
				"List all infrastructure components and their current operational status.",
			usage: "xcsh cloudstatus components",
			examples: [
				{
					command: "xcsh cloudstatus components",
					description: "List all components with status",
				},
				{
					command: "xcsh cloudstatus components --output json",
					description:
						"Get components as JSON for monitoring integration",
				},
			],
			category: "cloudstatus",
			related: ["cloudstatus status", "cloudstatus summary"],
		}),

		incidents: buildCommandSpec({
			command: "cloudstatus incidents",
			description:
				"List active and recent incidents affecting F5 Distributed Cloud services.",
			usage: "xcsh cloudstatus incidents",
			examples: [
				{
					command: "xcsh cloudstatus incidents",
					description: "View active incidents",
				},
				{
					command: "xcsh cloudstatus incidents --output json",
					description: "Get incidents as JSON for alerting systems",
				},
			],
			category: "cloudstatus",
			related: ["cloudstatus status", "cloudstatus maintenance"],
		}),

		maintenance: buildCommandSpec({
			command: "cloudstatus maintenance",
			description:
				"List scheduled maintenance windows for F5 Distributed Cloud services.",
			usage: "xcsh cloudstatus maintenance",
			examples: [
				{
					command: "xcsh cloudstatus maintenance",
					description: "View upcoming maintenance windows",
				},
				{
					command: "xcsh cloudstatus maintenance --output json",
					description: "Get maintenance schedule as JSON",
				},
			],
			category: "cloudstatus",
			related: ["cloudstatus status", "cloudstatus incidents"],
		}),
	};
}

/**
 * Build spec for login domain commands
 */
export function buildLoginSpecs(): Record<string, CommandSpec> {
	return {
		banner: buildCommandSpec({
			command: "login banner",
			description:
				"Display xcsh banner with logo and connection information.",
			usage: "xcsh login banner",
			examples: [
				{
					command: "xcsh login banner",
					description: "Show the xcsh welcome banner",
				},
			],
			category: "login",
			related: ["login profile show"],
		}),

		"profile list": buildCommandSpec({
			command: "login profile list",
			description: "List all saved connection profiles.",
			usage: "xcsh login profile list",
			examples: [
				{
					command: "xcsh login profile list",
					description: "List saved profiles",
				},
				{
					command: "xcsh login profile list --output json",
					description: "Get profiles as JSON",
				},
			],
			category: "login",
			related: [
				"login profile show",
				"login profile create",
				"login profile use",
			],
		}),

		"profile show": buildCommandSpec({
			command: "login profile show",
			description:
				"Show current connection profile and authentication status.",
			usage: "xcsh login profile show [name]",
			flags: [
				{
					name: "name",
					description:
						"Profile name to show (optional, defaults to active)",
					type: "string",
				},
			],
			examples: [
				{
					command: "xcsh login profile show",
					description: "Show active profile",
				},
				{
					command: "xcsh login profile show production",
					description: "Show specific profile",
				},
			],
			category: "login",
			related: ["login profile list", "login profile use"],
		}),

		"profile create": buildCommandSpec({
			command: "login profile create",
			description:
				"Create a new connection profile with URL and credentials.",
			usage: "xcsh login profile create <name>",
			flags: [
				{
					name: "name",
					description: "Profile name",
					type: "string",
					required: true,
				},
			],
			examples: [
				{
					command: "xcsh login profile create production",
					description: "Create a new profile named 'production'",
				},
			],
			category: "login",
			related: ["login profile list", "login profile use"],
		}),

		"profile use": buildCommandSpec({
			command: "login profile use",
			description: "Switch to a different connection profile.",
			usage: "xcsh login profile use <name>",
			flags: [
				{
					name: "name",
					description: "Profile name to activate",
					type: "string",
					required: true,
				},
			],
			examples: [
				{
					command: "xcsh login profile use staging",
					description: "Switch to staging profile",
				},
			],
			category: "login",
			related: ["login profile list", "login profile show"],
		}),

		"context show": buildCommandSpec({
			command: "login context show",
			description: "Show the current default namespace context.",
			usage: "xcsh login context show",
			examples: [
				{
					command: "xcsh login context show",
					description: "Display current namespace",
				},
			],
			category: "login",
			related: ["login context set", "login context list"],
		}),

		"context set": buildCommandSpec({
			command: "login context set",
			description: "Set the default namespace for subsequent operations.",
			usage: "xcsh login context set <namespace>",
			flags: [
				{
					name: "namespace",
					description: "Namespace to set as default",
					type: "string",
					required: true,
				},
			],
			examples: [
				{
					command: "xcsh login context set production",
					description: "Set production as default namespace",
				},
			],
			category: "login",
			related: ["login context show", "login context list"],
		}),
	};
}

/**
 * Get command spec by command path
 */
export function getCommandSpec(commandPath: string): CommandSpec | undefined {
	const cloudstatusSpecs = buildCloudstatusSpecs();
	const loginSpecs = buildLoginSpecs();

	// Normalize command path
	const normalized = commandPath.toLowerCase().trim();

	// Check cloudstatus commands
	if (normalized.startsWith("cloudstatus ")) {
		const subcommand = normalized.replace("cloudstatus ", "");
		return cloudstatusSpecs[subcommand];
	}

	// Check login commands
	if (normalized.startsWith("login ")) {
		const subcommand = normalized.replace("login ", "");
		return loginSpecs[subcommand];
	}

	return undefined;
}

/**
 * List all available command specs
 */
export function listAllCommandSpecs(): CommandSpec[] {
	return [
		...Object.values(buildCloudstatusSpecs()),
		...Object.values(buildLoginSpecs()),
	];
}

/**
 * Convert FlagSpec to GlobalFlagSpec format
 */
function toGlobalFlagSpec(flag: FlagSpec): GlobalFlagSpec {
	return {
		name: flag.name,
		type: flag.type ?? "string",
		description: flag.description,
		shorthand: flag.alias ?? "",
		default: flag.default ?? "",
	};
}

/**
 * Build spec for a custom domain from the registry
 */
function buildCustomDomainSpec(domainName: string): CLICommandSpec | null {
	const domain = customDomains.get(domainName);
	if (!domain) return null;

	const subcommands: CLICommandSpec[] = [];

	// Add direct commands
	for (const [cmdName, cmd] of domain.commands) {
		subcommands.push({
			path: [domainName, cmdName],
			use: cmd.usage ?? cmdName,
			short: cmd.descriptionShort,
			long: cmd.description,
			example: "",
			aliases: cmd.aliases ?? [],
			flags: [],
			subcommands: [],
		});
	}

	// Add subcommand groups
	for (const [groupName, group] of domain.subcommands) {
		const groupSubcommands: CLICommandSpec[] = [];

		for (const [cmdName, cmd] of group.commands) {
			groupSubcommands.push({
				path: [domainName, groupName, cmdName],
				use: cmd.usage ?? cmdName,
				short: cmd.descriptionShort,
				long: cmd.description,
				example: "",
				aliases: cmd.aliases ?? [],
				flags: [],
				subcommands: [],
			});
		}

		subcommands.push({
			path: [domainName, groupName],
			use: groupName,
			short: group.descriptionShort,
			long: group.description,
			example: "",
			aliases: [],
			flags: [],
			subcommands: groupSubcommands,
		});
	}

	return {
		path: [domainName],
		use: domainName,
		short: domain.descriptionShort,
		long: domain.description,
		example: "",
		aliases: [],
		flags: [],
		subcommands,
	};
}

/**
 * Build spec for an API-generated domain
 */
function buildApiDomainSpec(domainName: string): CLICommandSpec | null {
	const info = domainRegistry.get(domainName);
	if (!info) return null;

	// API domains have standard CRUD actions
	const actions = Array.from(validActions);
	const subcommands: CLICommandSpec[] = actions.map((action) => ({
		path: [domainName, action],
		use: action,
		short: `${action.charAt(0).toUpperCase() + action.slice(1)} ${info.displayName} resources`,
		long: `${action.charAt(0).toUpperCase() + action.slice(1)} ${info.displayName} resources in F5 Distributed Cloud`,
		example: `xcsh ${domainName} ${action}`,
		aliases: [],
		flags: [],
		subcommands: [],
	}));

	return {
		path: [domainName],
		use: domainName,
		short: info.descriptionShort,
		long: info.description,
		example: "",
		aliases: info.aliases,
		flags: [],
		subcommands,
	};
}

/**
 * Build the complete CLI specification for documentation generation
 * This is used by scripts/generate-docs.py when running `xcsh --spec`
 */
export function buildFullCLISpec(): CLISpec {
	const commands: CLICommandSpec[] = [];

	// Add custom domains first
	for (const domain of customDomains.all()) {
		const spec = buildCustomDomainSpec(domain.name);
		if (spec) {
			commands.push(spec);
		}
	}

	// Add API-generated domains (skip if already added as custom domain)
	const customDomainNames = new Set(customDomains.list());
	for (const [domainName] of domainRegistry) {
		if (!customDomainNames.has(domainName)) {
			const spec = buildApiDomainSpec(domainName);
			if (spec) {
				commands.push(spec);
			}
		}
	}

	// Sort commands by name for consistent output
	commands.sort((a, b) => (a.path[0] ?? "").localeCompare(b.path[0] ?? ""));

	return {
		version: CLI_VERSION,
		global_flags: GLOBAL_FLAGS.map(toGlobalFlagSpec),
		commands,
	};
}

/**
 * Format the full CLI spec as JSON string
 */
export function formatFullCLISpec(): string {
	return JSON.stringify(buildFullCLISpec(), null, 2);
}
