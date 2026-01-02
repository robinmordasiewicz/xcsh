#!/usr/bin/env npx tsx
/**
 * Description Loader Generator
 *
 * Reads descriptions from config/custom-domain-descriptions.yaml
 * and generates a TypeScript module that exports them for use by domain definitions.
 *
 * This runs as part of the build process to make YAML descriptions available at compile time.
 */

import { readFileSync, writeFileSync, existsSync } from "fs";
import { join, dirname } from "path";
import { fileURLToPath } from "url";
import { parse as parseYaml } from "yaml";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

interface SpecInfo {
	title: string | null;
	summary: string | null;
	description: string | null;
}

/**
 * Extract CLI info from upstream OpenAPI spec (build-time)
 * No local enrichment - use upstream verbatim as single source of truth
 * Returns title (short), summary (medium), and description (long) from spec
 */
function extractSpecInfo(): SpecInfo {
	const specPath = join(__dirname, "..", ".specs", "openapi.json");
	if (!existsSync(specPath)) {
		console.warn("  Warning: API spec not found, CLI info will use fallback");
		return { title: null, summary: null, description: null };
	}
	try {
		const content = readFileSync(specPath, "utf-8");
		const spec = JSON.parse(content);
		const title = spec.info?.title || null;
		const summary = spec.info?.summary || null;
		const description = spec.info?.description || null;
		if (title || summary || description) {
			console.log("  Extracted CLI info from .specs/openapi.json");
			if (title) console.log(`    title: "${title}"`);
			if (summary) console.log(`    summary: ${summary.length} chars`);
			if (description) console.log(`    description: ${description.length} chars`);
		}
		return { title, summary, description };
	} catch (error) {
		console.warn(`  Warning: Failed to parse API spec: ${error}`);
		return { title: null, summary: null, description: null };
	}
}

interface DescriptionTiers {
	short: string;
	medium: string;
	long: string;
}

interface CommandDescription extends DescriptionTiers {}

interface SubcommandDescription extends DescriptionTiers {
	commands?: Record<string, CommandDescription>;
}

interface DomainDescription extends DescriptionTiers {
	source_patterns_hash: string;
	subcommands?: Record<string, SubcommandDescription>;
	commands?: Record<string, CommandDescription>;
}

interface CliDescription extends DescriptionTiers {
	source_patterns_hash: string;
}

interface GeneratedDescriptions {
	version: string;
	generated_at: string;
	cli?: Record<string, CliDescription>;
	domains: Record<string, DomainDescription>;
}

/**
 * Escape string for TypeScript
 */
function escapeTs(str: string): string {
	return str
		.replace(/\\/g, "\\\\")
		.replace(/"/g, '\\"')
		.replace(/\n/g, "\\n")
		.replace(/\r/g, "\\r")
		.replace(/\t/g, "\\t");
}

/**
 * Generate TypeScript interface definitions
 */
function generateInterfaces(): string {
	return `/**
 * Generated Description Types
 * Auto-generated from config/custom-domain-descriptions.yaml
 */

export interface DescriptionTiers {
	short: string;
	medium: string;
	long: string;
}

export type CommandDescriptions = DescriptionTiers;

export interface SubcommandDescriptions extends DescriptionTiers {
	commands?: Record<string, CommandDescriptions>;
}

export interface DomainDescriptions extends DescriptionTiers {
	subcommands?: Record<string, SubcommandDescriptions>;
	commands?: Record<string, CommandDescriptions>;
}

export type CliDescriptions = DescriptionTiers;

export interface GeneratedDescriptionsData {
	version: string;
	generatedAt: string;
	cli?: Record<string, CliDescriptions>;
	domains: Record<string, DomainDescriptions>;
}
`;
}

/**
 * Generate description object for a tier
 */
function generateTierObject(desc: DescriptionTiers, indent: string): string {
	return `{
${indent}	short: "${escapeTs(desc.short)}",
${indent}	medium: "${escapeTs(desc.medium)}",
${indent}	long: "${escapeTs(desc.long)}",
${indent}}`;
}

/**
 * Generate commands object
 */
function generateCommandsObject(
	commands: Record<string, CommandDescription>,
	indent: string,
): string {
	const entries = Object.entries(commands).map(([name, cmd]) => {
		return `${indent}	"${name}": ${generateTierObject(cmd, indent + "\t")},`;
	});
	return `{
${entries.join("\n")}
${indent}}`;
}

/**
 * Generate subcommands object
 */
function generateSubcommandsObject(
	subcommands: Record<string, SubcommandDescription>,
	indent: string,
): string {
	const entries = Object.entries(subcommands).map(([name, sub]) => {
		let obj = `${indent}	"${name}": {
${indent}		short: "${escapeTs(sub.short)}",
${indent}		medium: "${escapeTs(sub.medium)}",
${indent}		long: "${escapeTs(sub.long)}",`;

		if (sub.commands && Object.keys(sub.commands).length > 0) {
			obj += `
${indent}		commands: ${generateCommandsObject(sub.commands, indent + "\t\t")},`;
		}

		obj += `
${indent}	},`;
		return obj;
	});
	return `{
${entries.join("\n")}
${indent}}`;
}

/**
 * Generate domain object
 */
function generateDomainObject(
	domain: DomainDescription,
	indent: string,
): string {
	let obj = `{
${indent}	short: "${escapeTs(domain.short)}",
${indent}	medium: "${escapeTs(domain.medium)}",
${indent}	long: "${escapeTs(domain.long)}",`;

	if (domain.subcommands && Object.keys(domain.subcommands).length > 0) {
		obj += `
${indent}	subcommands: ${generateSubcommandsObject(domain.subcommands, indent + "\t")},`;
	}

	if (domain.commands && Object.keys(domain.commands).length > 0) {
		obj += `
${indent}	commands: ${generateCommandsObject(domain.commands, indent + "\t")},`;
	}

	obj += `
${indent}}`;
	return obj;
}

/**
 * Generate CLI object
 */
function generateCliObject(cli: CliDescription, indent: string): string {
	return `{
${indent}	short: "${escapeTs(cli.short)}",
${indent}	medium: "${escapeTs(cli.medium)}",
${indent}	long: "${escapeTs(cli.long)}",
${indent}}`;
}

/**
 * Generate the full TypeScript module
 */
function generateModule(data: GeneratedDescriptions): string {
	const interfaces = generateInterfaces();

	// Extract CLI info from upstream spec (single source of truth)
	const specInfo = extractSpecInfo();

	// Generate exports for title (short), summary (medium), and description (long)
	const specInfoExport = `
/**
 * CLI Title from upstream OpenAPI spec (short description)
 * Extracted at build time from .specs/openapi.json info.title
 * This is the single source of truth - no local enrichment
 */
export const CLI_TITLE_FROM_SPEC: string | null = ${JSON.stringify(specInfo.title)};

/**
 * CLI Summary from upstream OpenAPI spec (medium description)
 * Extracted at build time from .specs/openapi.json info.summary
 * This is the single source of truth - no local enrichment
 */
export const CLI_SUMMARY_FROM_SPEC: string | null = ${JSON.stringify(specInfo.summary)};

/**
 * CLI Description from upstream OpenAPI spec (long description)
 * Extracted at build time from .specs/openapi.json info.description
 * This is the single source of truth - no local enrichment
 */
export const CLI_DESCRIPTION_FROM_SPEC: string | null = ${JSON.stringify(specInfo.description)};
`;

	// Generate CLI entries if present (for backwards compatibility)
	let cliSection = "";
	if (data.cli && Object.keys(data.cli).length > 0) {
		const cliEntries = Object.entries(data.cli).map(([name, cli]) => {
			return `		"${name}": ${generateCliObject(cli, "\t\t")},`;
		});
		cliSection = `	cli: {
${cliEntries.join("\n")}
	},`;
	}

	const domainEntries = Object.entries(data.domains).map(([name, domain]) => {
		return `		"${name}": ${generateDomainObject(domain, "\t\t")},`;
	});

	return `${interfaces}
${specInfoExport}

/**
 * Generated Descriptions Data
 * Auto-generated from config/custom-domain-descriptions.yaml
 * Generated at: ${data.generated_at}
 *
 * DO NOT EDIT MANUALLY - Regenerate with: npm run generate:descriptions
 */
export const generatedDescriptions: GeneratedDescriptionsData = {
	version: "${data.version}",
	generatedAt: "${data.generated_at}",
${cliSection}
	domains: {
${domainEntries.join("\n")}
	},
};

/**
 * Get CLI descriptions
 */
export function getCliDescriptions(cliName: string = "xcsh"): CliDescriptions | undefined {
	return generatedDescriptions.cli?.[cliName];
}

/**
 * Get descriptions for a domain
 */
export function getDomainDescriptions(domainName: string): DomainDescriptions | undefined {
	return generatedDescriptions.domains[domainName];
}

/**
 * Get descriptions for a subcommand within a domain
 */
export function getSubcommandDescriptions(
	domainName: string,
	subcommandName: string,
): SubcommandDescriptions | undefined {
	const domain = generatedDescriptions.domains[domainName];
	return domain?.subcommands?.[subcommandName];
}

/**
 * Get descriptions for a command within a domain or subcommand
 */
export function getCommandDescriptions(
	domainName: string,
	commandName: string,
	subcommandName?: string,
): CommandDescriptions | undefined {
	const domain = generatedDescriptions.domains[domainName];
	if (!domain) return undefined;

	if (subcommandName) {
		const subcommand = domain.subcommands?.[subcommandName];
		return subcommand?.commands?.[commandName];
	}

	return domain.commands?.[commandName];
}
`;
}

/**
 * Main function
 */
async function main(): Promise<void> {
	console.log("=== Description Loader Generator ===\n");

	const projectRoot = join(__dirname, "..");
	const yamlPath = join(projectRoot, "config", "custom-domain-descriptions.yaml");
	const outputPath = join(projectRoot, "src", "domains", "descriptions.generated.ts");

	// Check if YAML exists
	if (!existsSync(yamlPath)) {
		console.log("No descriptions YAML found. Creating empty module.");
		const emptyModule = `${generateInterfaces()}

/**
 * Empty descriptions - run 'npm run generate:descriptions' to populate
 */
export const generatedDescriptions: GeneratedDescriptionsData = {
	version: "1.0.0",
	generatedAt: "",
	cli: {},
	domains: {},
};

export function getCliDescriptions(_cliName: string = "xcsh"): CliDescriptions | undefined {
	return undefined;
}

export function getDomainDescriptions(_domainName: string): DomainDescriptions | undefined {
	return undefined;
}

export function getSubcommandDescriptions(
	_domainName: string,
	_subcommandName: string,
): SubcommandDescriptions | undefined {
	return undefined;
}

export function getCommandDescriptions(
	_domainName: string,
	_commandName: string,
	_subcommandName?: string,
): CommandDescriptions | undefined {
	return undefined;
}
`;
		writeFileSync(outputPath, emptyModule, "utf-8");
		console.log(`  Generated empty module: ${outputPath}`);
		return;
	}

	// Read and parse YAML
	console.log(`Reading: ${yamlPath}`);
	const yamlContent = readFileSync(yamlPath, "utf-8");
	const data = parseYaml(yamlContent) as GeneratedDescriptions;

	if (!data || !data.domains) {
		console.error("Invalid YAML structure");
		process.exit(1);
	}

	// Generate TypeScript module
	console.log("Generating TypeScript module...");
	const tsModule = generateModule(data);

	// Write output
	console.log(`Writing: ${outputPath}`);
	writeFileSync(outputPath, tsModule, "utf-8");

	// Count what was generated
	const cliCount = data.cli ? Object.keys(data.cli).length : 0;
	const domainCount = Object.keys(data.domains).length;
	let subcommandCount = 0;
	for (const domain of Object.values(data.domains)) {
		if (domain.subcommands) {
			subcommandCount += Object.keys(domain.subcommands).length;
		}
	}

	console.log(`\n=== Complete ===`);
	if (cliCount > 0) {
		console.log(`  CLI: ${cliCount}`);
	}
	console.log(`  Domains: ${domainCount}`);
	console.log(`  Subcommands: ${subcommandCount}`);
	console.log(`  Output: ${outputPath}`);
}

// Run main
main().catch((err) => {
	console.error("Fatal error:", err);
	process.exit(1);
});
