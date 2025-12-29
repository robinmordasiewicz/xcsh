#!/usr/bin/env npx tsx
/**
 * Domain Generator Script
 * Generates src/types/domains_generated.ts from .specs/index.json
 *
 * Run: npx tsx scripts/generate-domains.ts
 */

import * as fs from "fs";
import * as path from "path";
import * as yaml from "yaml";

// Types matching the upstream spec structure
interface SpecIndexEntry {
	domain: string;
	title: string;
	description: string;
	description_short: string; // ~60 chars for completions, badges
	description_medium: string; // ~150 chars for tooltips, summaries
	file: string;
	path_count: number;
	schema_count: number;
	complexity: string;
	is_preview: boolean;
	requires_tier: string;
	domain_category: string;
	use_cases: string[];
	related_domains: string[];
	aliases?: string[]; // Added in upstream spec (issue #178)
	cli_metadata?: Record<string, unknown>;
}

interface SpecIndex {
	version: string;
	timestamp: string;
	specifications: SpecIndexEntry[];
}

interface DomainConfig {
	version?: string;
	aliases?: Record<string, string[]>;
	deprecated_domains?: Record<
		string,
		{
			maps_to: string;
			reason: string;
			deprecated_since: string;
		}
	>;
}

interface DomainInfo {
	name: string;
	displayName: string;
	description: string;
	descriptionShort: string;
	descriptionMedium: string;
	aliases: string[];
	complexity: string;
	isPreview: boolean;
	requiresTier: string;
	category: string;
	useCases: string[];
	relatedDomains: string[];
	cliMetadata?: Record<string, unknown>;
}

/**
 * Convert snake_case to Title Case
 */
function titleCase(s: string): string {
	return s
		.split("_")
		.map((part) =>
			part.length > 0 ? part[0].toUpperCase() + part.slice(1) : "",
		)
		.join(" ");
}

/**
 * Escape string for TypeScript output
 */
function escapeString(s: string): string {
	return s
		.replace(/\\/g, "\\\\")
		.replace(/"/g, '\\"')
		.replace(/\n/g, "\\n")
		.replace(/\r/g, "\\r")
		.replace(/\t/g, "\\t");
}

/**
 * Generate TypeScript code for a domain
 */
function generateDomainEntry(domain: DomainInfo): string {
	const aliasArray = domain.aliases
		.map((a) => `"${escapeString(a)}"`)
		.join(", ");
	const useCasesArray = domain.useCases
		.map((u) => `"${escapeString(u)}"`)
		.join(", ");
	const relatedArray = domain.relatedDomains
		.map((r) => `"${escapeString(r)}"`)
		.join(", ");

	let code = `\t["${domain.name}", {\n`;
	code += `\t\tname: "${domain.name}",\n`;
	code += `\t\tdisplayName: "${escapeString(domain.displayName)}",\n`;
	code += `\t\tdescription: "${escapeString(domain.description)}",\n`;
	code += `\t\tdescriptionShort: "${escapeString(domain.descriptionShort)}",\n`;
	code += `\t\tdescriptionMedium: "${escapeString(domain.descriptionMedium)}",\n`;
	code += `\t\taliases: [${aliasArray}],\n`;
	code += `\t\tcomplexity: "${domain.complexity}" as const,\n`;
	code += `\t\tisPreview: ${domain.isPreview},\n`;
	code += `\t\trequiresTier: "${domain.requiresTier}",\n`;
	code += `\t\tcategory: "${domain.category}",\n`;
	code += `\t\tuseCases: [${useCasesArray}],\n`;
	code += `\t\trelatedDomains: [${relatedArray}],\n`;

	if (domain.cliMetadata && Object.keys(domain.cliMetadata).length > 0) {
		code += `\t\tcliMetadata: ${JSON.stringify(domain.cliMetadata, null, 2).replace(/\n/g, "\n\t\t")},\n`;
	}

	code += `\t}]`;
	return code;
}

/**
 * Main generator function
 */
async function main(): Promise<void> {
	console.log("🏗️  Generating domains from upstream specs...");

	const specsDir = ".specs";
	const indexPath = path.join(specsDir, "index.json");
	const configPath = path.join(specsDir, "domain_config.yaml");
	const outputPath = path.join("src", "types", "domains_generated.ts");

	// Read spec index
	if (!fs.existsSync(indexPath)) {
		console.error(`❌ Spec index not found: ${indexPath}`);
		console.error(
			"   Run 'make download-specs' first to download API specifications.",
		);
		process.exit(1);
	}

	const indexData = fs.readFileSync(indexPath, "utf-8");
	const specIndex: SpecIndex = JSON.parse(indexData);
	console.log(
		`✓ Loaded spec index v${specIndex.version} with ${specIndex.specifications.length} domains`,
	);

	// Read domain config (optional fallback for aliases)
	let config: DomainConfig = { aliases: {}, deprecated_domains: {} };
	if (fs.existsSync(configPath)) {
		const configData = fs.readFileSync(configPath, "utf-8");
		config = yaml.parse(configData) || config;
		console.log(
			`✓ Loaded domain config with ${Object.keys(config.aliases || {}).length} alias mappings (fallback)`,
		);
	}

	// Check if upstream specs include aliases
	const hasUpstreamAliases = specIndex.specifications.some(
		(spec) => spec.aliases && spec.aliases.length > 0,
	);
	if (hasUpstreamAliases) {
		console.log("✓ Using aliases from upstream specs");
	} else if (Object.keys(config.aliases || {}).length > 0) {
		console.log(
			"⚠️  Upstream specs have no aliases, using local domain_config.yaml fallback",
		);
	} else {
		console.log("ℹ️  No aliases configured (upstream or local)");
	}

	// Build domain registry
	const domains: DomainInfo[] = [];

	for (const spec of specIndex.specifications) {
		// Skip empty domains
		if (spec.path_count === 0 && spec.schema_count === 0) {
			console.log(`⊘ Skipping empty domain: ${spec.domain}`);
			continue;
		}

		// Prefer upstream aliases, fallback to local config
		const aliases = spec.aliases?.length
			? spec.aliases
			: config.aliases?.[spec.domain] || [];

		const domainInfo: DomainInfo = {
			name: spec.domain,
			displayName: titleCase(spec.domain),
			description: spec.description,
			descriptionShort: spec.description_short,
			descriptionMedium: spec.description_medium,
			aliases,
			complexity: spec.complexity || "moderate",
			isPreview: spec.is_preview || false,
			requiresTier: spec.requires_tier || "Standard",
			category: spec.domain_category || "Other",
			useCases: spec.use_cases || [],
			relatedDomains: spec.related_domains || [],
			cliMetadata: spec.cli_metadata,
		};

		domains.push(domainInfo);
	}

	// Sort domains alphabetically for consistent output
	domains.sort((a, b) => a.name.localeCompare(b.name));

	console.log(`✓ Generated registry with ${domains.length} active domains`);

	// Generate TypeScript file
	const domainEntries = domains.map(generateDomainEntry).join(",\n");

	const outputContent = `/**
 * AUTO-GENERATED FILE - DO NOT EDIT
 * Generated from .specs/index.json v${specIndex.version}
 * Run: npx tsx scripts/generate-domains.ts
 */

import type { DomainInfo } from "./domains.js";

/**
 * Spec version used for generation
 */
export const SPEC_VERSION = "${specIndex.version}";

/**
 * Generated domain data from upstream API specifications
 */
export const generatedDomains: Map<string, DomainInfo> = new Map([
${domainEntries}
]);

/**
 * Total domain count
 */
export const DOMAIN_COUNT = ${domains.length};
`;

	// Ensure output directory exists
	const outputDir = path.dirname(outputPath);
	if (!fs.existsSync(outputDir)) {
		fs.mkdirSync(outputDir, { recursive: true });
	}

	// Write output file
	fs.writeFileSync(outputPath, outputContent, "utf-8");
	console.log(`✓ Generated: ${outputPath}`);
	console.log(`✅ Domain generation complete! (${domains.length} domains)`);
}

main().catch((err) => {
	console.error("❌ Generation failed:", err);
	process.exit(1);
});
