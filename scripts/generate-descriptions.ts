#!/usr/bin/env npx tsx
/**
 * Description Generator for Custom Domains
 *
 * Generates 3-tier descriptions (short ~60, medium ~150, long ~500 chars)
 * for custom domains and subcommands using Claude CLI with self-refine loop.
 *
 * Based on upstream pattern from f5xc-api-enriched:
 * - Claude CLI with --json-schema for structured JSON output
 * - Validation: character limits, banned terms, action verb starters
 * - Self-refine loop with specific feedback on violations
 * - YAML storage with source hash for change detection
 */

import { spawnSync } from "child_process";
import { readFileSync, writeFileSync, existsSync } from "fs";
import { createHash } from "crypto";
import { join, dirname } from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Configuration
const MAX_SHORT = 60;
const MAX_MEDIUM = 150;
const MAX_LONG = 500;
const MAX_RETRIES = 3;

// Banned terms (case-insensitive)
const BANNED_TERMS = [
	"f5",
	"f5 xc",
	"f5xc",
	"distributed cloud",
	"xc ",
	" xc",
	" api",
	"api ",
	"specifications",
	" spec",
	"spec ",
	"endpoint",
	"comprehensive",
	"complete",
	"full ",
	" full",
	"various",
	"extensive",
	"robust",
	"powerful",
	"seamless",
	"cutting-edge",
	"state-of-the-art",
	"world-class",
	"best-in-class",
];

// Bad starters (case-insensitive)
const BAD_STARTERS = [
	"this",
	"the ",
	"a ",
	"an ",
	"provides",
	"enables",
	"allows",
	"offers",
	"it ",
	"there ",
];

// Good action verb starters
const ACTION_VERBS = [
	"configure",
	"create",
	"manage",
	"define",
	"deploy",
	"monitor",
	"display",
	"generate",
	"list",
	"show",
	"set",
	"get",
	"add",
	"remove",
	"delete",
	"update",
	"check",
	"validate",
	"verify",
	"run",
	"execute",
	"build",
	"install",
	"switch",
	"activate",
	"view",
	"track",
	"retrieve",
];

// Interfaces
interface DescriptionTiers {
	short: string;
	medium: string;
	long: string;
}

interface DomainContext {
	name: string;
	type: "domain" | "subcommand" | "command" | "cli_root";
	parentName?: string;
	existingDescription?: string;
	commands?: string[];
	features?: string[];
}

interface CliDescription {
	source_patterns_hash: string;
	short: string;
	medium: string;
	long: string;
}

interface GeneratedDescriptions {
	version: string;
	generated_at: string;
	cli?: Record<string, CliDescription>;
	domains: Record<string, DomainDescription>;
}

interface DomainDescription {
	source_patterns_hash: string;
	short: string;
	medium: string;
	long: string;
	subcommands?: Record<string, SubcommandDescription>;
	commands?: Record<string, CommandDescription>;
}

interface SubcommandDescription {
	short: string;
	medium: string;
	long: string;
	commands?: Record<string, CommandDescription>;
}

interface CommandDescription {
	short: string;
	medium: string;
	long: string;
}

// JSON Schema for Claude CLI
const DESCRIPTION_SCHEMA = {
	type: "object",
	properties: {
		short: {
			type: "string",
			description: `Short description, max ${MAX_SHORT} characters. For completions and badges.`,
		},
		medium: {
			type: "string",
			description: `Medium description, max ${MAX_MEDIUM} characters. For tooltips and summaries.`,
		},
		long: {
			type: "string",
			description: `Long description, max ${MAX_LONG} characters. For detailed help output.`,
		},
	},
	required: ["short", "medium", "long"],
	additionalProperties: false,
};

/**
 * Calculate hash of source content for change detection
 */
function calculateHash(content: string): string {
	return createHash("sha256").update(content).digest("hex").slice(0, 16);
}

/**
 * Check character limits
 */
function checkCharacterLimits(desc: DescriptionTiers): string[] {
	const violations: string[] = [];

	if (desc.short.length > MAX_SHORT) {
		violations.push(
			`Short description exceeds ${MAX_SHORT} chars (got ${desc.short.length}): "${desc.short.slice(0, 50)}..."`,
		);
	}

	if (desc.medium.length > MAX_MEDIUM) {
		violations.push(
			`Medium description exceeds ${MAX_MEDIUM} chars (got ${desc.medium.length}): "${desc.medium.slice(0, 80)}..."`,
		);
	}

	if (desc.long.length > MAX_LONG) {
		violations.push(
			`Long description exceeds ${MAX_LONG} chars (got ${desc.long.length}): "${desc.long.slice(0, 100)}..."`,
		);
	}

	return violations;
}

/**
 * Check for banned terms
 */
function checkBannedTerms(desc: DescriptionTiers): string[] {
	const violations: string[] = [];
	const allText = `${desc.short} ${desc.medium} ${desc.long}`.toLowerCase();

	for (const term of BANNED_TERMS) {
		if (allText.includes(term.toLowerCase())) {
			violations.push(`Contains banned term: "${term}"`);
		}
	}

	return violations;
}

/**
 * Check for good starters (action verbs)
 */
function checkStyleCompliance(desc: DescriptionTiers): string[] {
	const violations: string[] = [];

	// Check each tier starts appropriately
	for (const [tier, text] of Object.entries(desc)) {
		const lowerText = text.toLowerCase().trim();

		// Check for bad starters
		for (const badStart of BAD_STARTERS) {
			if (lowerText.startsWith(badStart.toLowerCase())) {
				violations.push(
					`${tier} starts with weak word "${badStart.trim()}". Use action verb instead.`,
				);
			}
		}

		// Check if it starts with an action verb (recommended but not required)
		const startsWithVerb = ACTION_VERBS.some((verb) =>
			lowerText.startsWith(verb.toLowerCase()),
		);

		if (!startsWithVerb && tier !== "long") {
			// Long descriptions can be more flexible
			// This is a soft warning, not a hard violation
		}
	}

	return violations;
}

/**
 * Check for cross-tier repetition
 */
function checkCrossTierRepetition(desc: DescriptionTiers): string[] {
	const violations: string[] = [];

	// Extract significant words (>3 chars, not common words)
	const commonWords = new Set([
		"the",
		"and",
		"for",
		"with",
		"from",
		"that",
		"this",
		"your",
		"will",
		"can",
		"are",
		"was",
		"has",
		"have",
		"been",
	]);

	function getSignificantWords(text: string): Set<string> {
		return new Set(
			text
				.toLowerCase()
				.replace(/[^a-z\s]/g, "")
				.split(/\s+/)
				.filter((w) => w.length > 3 && !commonWords.has(w)),
		);
	}

	const shortWords = getSignificantWords(desc.short);
	const mediumWords = getSignificantWords(desc.medium);
	const longWords = getSignificantWords(desc.long);

	// Check overlap between short and medium
	const shortMediumOverlap = [...shortWords].filter((w) =>
		mediumWords.has(w),
	);
	if (shortMediumOverlap.length >= 4) {
		violations.push(
			`Too much overlap between short and medium (${shortMediumOverlap.length} words): ${shortMediumOverlap.slice(0, 4).join(", ")}`,
		);
	}

	// Check overlap between medium and long
	const mediumLongOverlap = [...mediumWords].filter((w) => longWords.has(w));
	if (mediumLongOverlap.length >= 6) {
		// Allow more overlap for longer descriptions
		violations.push(
			`Too much overlap between medium and long (${mediumLongOverlap.length} words)`,
		);
	}

	return violations;
}

/**
 * Check for domain name self-references
 */
function checkDryCompliance(
	domainName: string,
	desc: DescriptionTiers,
): string[] {
	const violations: string[] = [];
	const allText = `${desc.short} ${desc.medium} ${desc.long}`.toLowerCase();

	// Check if domain name is repeated excessively
	const nameRegex = new RegExp(domainName.toLowerCase(), "g");
	const matches = allText.match(nameRegex);
	if (matches && matches.length > 2) {
		violations.push(
			`Domain name "${domainName}" appears ${matches.length} times (max 2)`,
		);
	}

	return violations;
}

/**
 * Run all validations
 */
function validateDescriptions(
	domainName: string,
	desc: DescriptionTiers,
): string[] {
	return [
		...checkCharacterLimits(desc),
		...checkBannedTerms(desc),
		...checkStyleCompliance(desc),
		...checkCrossTierRepetition(desc),
		...checkDryCompliance(domainName, desc),
	];
}

/**
 * Build initial prompt for Claude
 */
function buildInitialPrompt(context: DomainContext): string {
	const typeLabel =
		context.type === "domain"
			? "CLI domain"
			: context.type === "subcommand"
				? "subcommand group"
				: context.type === "cli_root"
					? "command-line interface tool"
					: "command";

	let prompt = `Generate 3-tier descriptions for a ${typeLabel} named "${context.name}"`;

	if (context.parentName) {
		prompt += ` (part of "${context.parentName}")`;
	}

	prompt += ".\n\n";

	if (context.existingDescription) {
		prompt += `Current description: "${context.existingDescription}"\n\n`;
	}

	if (context.features && context.features.length > 0) {
		prompt += `Key features:\n`;
		for (const feature of context.features) {
			prompt += `- ${feature}\n`;
		}
		prompt += "\n";
	}

	if (context.commands && context.commands.length > 0) {
		prompt += `Contains commands: ${context.commands.join(", ")}\n\n`;
	}

	prompt += `Requirements:
1. SHORT (max ${MAX_SHORT} chars): Concise label for shell completions
2. MEDIUM (max ${MAX_MEDIUM} chars): Brief summary for tooltips
3. LONG (max ${MAX_LONG} chars): Detailed help text

Style guidelines:
- Start with action verbs (Configure, Create, Manage, Display, Generate, etc.)
- Avoid: "This", "The", "A", "An", "Provides", "Enables", "Allows"
- Avoid marketing terms: "comprehensive", "seamless", "powerful", "robust"
- Avoid product-specific terms: "F5", "XC", "API"
- Each tier should provide progressively more detail
- Minimize word repetition across tiers
- Be specific about what the ${typeLabel} does`;

	return prompt;
}

/**
 * Build refinement prompt with violations
 */
function buildRefinementPrompt(
	context: DomainContext,
	prevDesc: DescriptionTiers,
	violations: string[],
): string {
	let prompt = buildInitialPrompt(context);

	prompt += "\n\n--- PREVIOUS ATTEMPT ---\n";
	prompt += `Short: "${prevDesc.short}"\n`;
	prompt += `Medium: "${prevDesc.medium}"\n`;
	prompt += `Long: "${prevDesc.long}"\n`;

	prompt += "\n--- VIOLATIONS TO FIX ---\n";
	for (const v of violations) {
		prompt += `- ${v}\n`;
	}

	prompt += "\nPlease fix these violations while maintaining quality descriptions.";

	return prompt;
}

/**
 * Call Claude CLI with structured output
 */
function callClaude(prompt: string): DescriptionTiers | null {
	const cmd = [
		"claude",
		"-p",
		prompt,
		"--output-format",
		"json",
		"--json-schema",
		JSON.stringify(DESCRIPTION_SCHEMA),
		"--tools",
		"", // Disable all tools
		"--no-session-persistence",
		"--strict-mcp-config", // Ignore MCP configurations
		"--disable-slash-commands", // Disable skills
		"--append-system-prompt",
		"You are generating CLI command descriptions. Respond ONLY with JSON matching the schema. No explanations.",
	];

	try {
		const result = spawnSync(cmd[0], cmd.slice(1), {
			encoding: "utf-8",
			timeout: 120000,
			maxBuffer: 1024 * 1024,
		});

		if (result.error) {
			console.error("Claude CLI error:", result.error.message);
			return null;
		}

		if (result.status !== 0) {
			console.error("Claude CLI exited with code:", result.status);
			console.error("stderr:", result.stderr);
			return null;
		}

		// Parse JSON response
		const output = result.stdout.trim();
		try {
			const parsed = JSON.parse(output);

			// Handle nested result structure from Claude CLI
			// When using --json-schema, output is in "structured_output" field
			let descriptions: DescriptionTiers | null = null;

			if (parsed.structured_output) {
				// Primary path: --json-schema puts result in structured_output
				descriptions = parsed.structured_output as DescriptionTiers;
			} else if (parsed.result && typeof parsed.result === "object") {
				// Fallback: check result field if it's an object
				descriptions = parsed.result as DescriptionTiers;
			} else if (parsed.short && parsed.medium && parsed.long) {
				// Direct response without wrapper
				descriptions = parsed as DescriptionTiers;
			}

			// Validate the parsed response has all required fields
			if (
				descriptions &&
				typeof descriptions.short === "string" &&
				typeof descriptions.medium === "string" &&
				typeof descriptions.long === "string"
			) {
				return descriptions;
			}

			console.error(
				"Invalid response structure:",
				JSON.stringify(parsed).slice(0, 300),
			);
			return null;
		} catch {
			console.error("Failed to parse Claude response:", output.slice(0, 300));
			return null;
		}
	} catch (err) {
		console.error("Error calling Claude:", err);
		return null;
	}
}

/**
 * Generate descriptions with retry loop
 */
async function generateWithRetry(
	context: DomainContext,
): Promise<DescriptionTiers | null> {
	let lastResponse: DescriptionTiers | null = null;
	let lastViolations: string[] = [];

	for (let attempt = 0; attempt < MAX_RETRIES; attempt++) {
		const prompt =
			lastResponse && lastViolations.length > 0
				? buildRefinementPrompt(context, lastResponse, lastViolations)
				: buildInitialPrompt(context);

		console.log(
			`  Attempt ${attempt + 1}/${MAX_RETRIES} for ${context.type} "${context.name}"...`,
		);

		const response = callClaude(prompt);
		if (!response) {
			console.error(`  Failed to get valid response from Claude`);
			continue;
		}

		// Validate structure before checking violations
		if (!response.short || !response.medium || !response.long) {
			console.error(`  Incomplete response structure`);
			continue;
		}

		const violations = validateDescriptions(context.name, response);

		if (violations.length === 0) {
			console.log(`  Success!`);
			return response;
		}

		console.log(`  Found ${violations.length} violations:`);
		for (const v of violations.slice(0, 3)) {
			console.log(`    - ${v}`);
		}

		lastResponse = response;
		lastViolations = violations;
	}

	// Return last attempt even if it has violations
	if (lastResponse) {
		console.log(
			`  Using last response despite ${lastViolations.length} violations`,
		);
		return lastResponse;
	}

	return null;
}

/**
 * Load existing descriptions from YAML
 */
function loadExistingDescriptions(
	yamlPath: string,
): GeneratedDescriptions | null {
	if (!existsSync(yamlPath)) {
		return null;
	}

	try {
		const content = readFileSync(yamlPath, "utf-8");
		// Simple YAML parsing for our structure
		return parseYaml(content);
	} catch {
		return null;
	}
}

/**
 * Simple YAML parser for our specific structure
 */
function parseYaml(content: string): GeneratedDescriptions | null {
	// For now, return null to force regeneration
	// A proper YAML parser would be needed for production
	return null;
}

/**
 * Format descriptions to YAML
 */
function formatToYaml(data: GeneratedDescriptions): string {
	let yaml = `# Custom Domain Descriptions
# CLI-specific domains only (login, cloudstatus, completion)
# These domains are NOT in the upstream OpenAPI spec
#
# CLI root description comes from .specs/openapi.json info.description
# (extracted at build time by generate-description-loader.ts - no local enrichment)
#
# Auto-generated by scripts/generate-descriptions.ts
# Do not edit manually - regenerate with: npm run generate:descriptions

version: "${data.version}"
generated_at: "${data.generated_at}"

# NOTE: No cli: section - CLI description extracted from upstream spec at build time

`;

	// Output CLI section first if present
	if (data.cli && Object.keys(data.cli).length > 0) {
		yaml += `cli:\n`;
		for (const [cliName, cli] of Object.entries(data.cli)) {
			yaml += `  ${cliName}:\n`;
			yaml += `    source_patterns_hash: "${cli.source_patterns_hash}"\n`;
			yaml += `    short: "${escapeYaml(cli.short)}"\n`;
			yaml += `    medium: >-\n`;
			yaml += `      ${escapeYamlMultiline(cli.medium)}\n`;
			yaml += `    long: >-\n`;
			yaml += `      ${escapeYamlMultiline(cli.long)}\n`;
		}
		yaml += "\n";
	}

	yaml += `domains:\n`;

	for (const [domainName, domain] of Object.entries(data.domains)) {
		yaml += `  ${domainName}:\n`;
		yaml += `    source_patterns_hash: "${domain.source_patterns_hash}"\n`;
		yaml += `    short: "${escapeYaml(domain.short)}"\n`;
		yaml += `    medium: "${escapeYaml(domain.medium)}"\n`;
		yaml += `    long: "${escapeYaml(domain.long)}"\n`;

		if (domain.subcommands) {
			yaml += `    subcommands:\n`;
			for (const [subName, sub] of Object.entries(domain.subcommands)) {
				yaml += `      ${subName}:\n`;
				yaml += `        short: "${escapeYaml(sub.short)}"\n`;
				yaml += `        medium: "${escapeYaml(sub.medium)}"\n`;
				yaml += `        long: "${escapeYaml(sub.long)}"\n`;

				if (sub.commands) {
					yaml += `        commands:\n`;
					for (const [cmdName, cmd] of Object.entries(sub.commands)) {
						yaml += `          ${cmdName}:\n`;
						yaml += `            short: "${escapeYaml(cmd.short)}"\n`;
						yaml += `            medium: "${escapeYaml(cmd.medium)}"\n`;
						yaml += `            long: "${escapeYaml(cmd.long)}"\n`;
					}
				}
			}
		}

		if (domain.commands) {
			yaml += `    commands:\n`;
			for (const [cmdName, cmd] of Object.entries(domain.commands)) {
				yaml += `      ${cmdName}:\n`;
				yaml += `        short: "${escapeYaml(cmd.short)}"\n`;
				yaml += `        medium: "${escapeYaml(cmd.medium)}"\n`;
				yaml += `        long: "${escapeYaml(cmd.long)}"\n`;
			}
		}

		yaml += "\n";
	}

	return yaml;
}

/**
 * Escape YAML string for inline use
 */
function escapeYaml(str: string): string {
	return str.replace(/"/g, '\\"').replace(/\n/g, "\\n");
}

/**
 * Escape YAML string for multiline (>-) block scalar use
 * Wraps long lines for readability while preserving content
 */
function escapeYamlMultiline(str: string): string {
	// For >- block scalar, we just need to handle line breaks
	// The >- indicator will fold lines, so we can use line breaks for wrapping
	const words = str.split(/\s+/);
	let result = "";
	let lineLength = 0;
	const maxLineLength = 76; // Leave room for indentation

	for (const word of words) {
		if (lineLength + word.length + 1 > maxLineLength && lineLength > 0) {
			result += "\n      " + word;
			lineLength = word.length;
		} else {
			result += (lineLength > 0 ? " " : "") + word;
			lineLength += word.length + 1;
		}
	}

	return result;
}

/**
 * Load API spec info for CLI description source
 */
function loadApiSpecInfo(): { title: string; description: string; version: string } | null {
	const projectRoot = join(__dirname, "..");
	const specPath = join(projectRoot, ".specs", "openapi.json");

	try {
		if (!existsSync(specPath)) {
			console.warn(`API spec not found at ${specPath}`);
			return null;
		}
		const specContent = readFileSync(specPath, "utf-8");
		const spec = JSON.parse(specContent);
		return spec.info || null;
	} catch (error) {
		console.warn(`Failed to load API spec: ${error}`);
		return null;
	}
}

/**
 * Get CLI root context for generating CLI-level descriptions
 * Sources description from upstream API spec's info.description field
 */
function getCliRootContext(): DomainContext {
	// Read description from upstream API spec (single source of truth)
	const apiInfo = loadApiSpecInfo();

	if (apiInfo?.description) {
		console.log("  Using API spec info.description as source");
		return {
			name: "xcsh",
			type: "cli_root",
			existingDescription: apiInfo.description,
			features: [], // Features derived from API description
		};
	}

	// Fallback only if API spec is unavailable
	console.warn("  API spec unavailable, using fallback context");
	return {
		name: "xcsh",
		type: "cli_root",
		existingDescription:
			"Interactive CLI for F5 Distributed Cloud services.",
		features: [
			"Interactive REPL with intelligent tab completion",
			"Multi-tenant profile management for environment switching",
			"100+ API domain operations across cloud services",
			"JSON, YAML, and table output format options",
			"Environment-based configuration with profile persistence",
			"Shell completion generation for bash, zsh, and fish",
		],
	};
}

/**
 * Get custom domain definitions
 */
function getCustomDomainContexts(): DomainContext[] {
	// Define custom domains that need descriptions
	return [
		{
			name: "login",
			type: "domain",
			existingDescription:
				"Manage authentication, profiles, and connection context for F5 Distributed Cloud CLI sessions.",
			commands: ["show", "profile", "context", "banner"],
		},
		{
			name: "profile",
			type: "subcommand",
			parentName: "login",
			existingDescription:
				"Manage saved connection profiles for tenant authentication.",
			commands: ["list", "show", "create", "delete", "active", "use"],
		},
		{
			name: "context",
			type: "subcommand",
			parentName: "login",
			existingDescription:
				"Manage default namespace context for scoping operations.",
			commands: ["show", "set", "list"],
		},
		{
			name: "cloudstatus",
			type: "domain",
			existingDescription:
				"Monitor F5 Distributed Cloud service status and incidents.",
			commands: ["status", "summary", "components", "incidents", "maintenance"],
		},
		{
			name: "completion",
			type: "domain",
			existingDescription:
				"Generate shell completion scripts for bash, zsh, and fish shells.",
			commands: ["bash", "zsh", "fish"],
		},
	];
}

/**
 * Main function
 */
async function main(): Promise<void> {
	console.log("Description Generator for CLI and Custom Domains");
	console.log("================================================\n");

	const projectRoot = join(__dirname, "..");
	const outputPath = join(projectRoot, "config", "custom-domain-descriptions.yaml");

	// Ensure config directory exists
	const configDir = dirname(outputPath);
	if (!existsSync(configDir)) {
		const { mkdirSync } = await import("fs");
		mkdirSync(configDir, { recursive: true });
	}

	const generatedData: GeneratedDescriptions = {
		version: "1.0.0",
		generated_at: new Date().toISOString(),
		// NOTE: No cli section - CLI description extracted from upstream spec
		// at build time by generate-description-loader.ts
		domains: {},
	};

	// NOTE: CLI root description is extracted from .specs/openapi.json
	// by the generate-description-loader.ts script (single source of truth)
	// No local Claude-based enrichment for CLI descriptions

	// Get domain contexts (CLI-specific domains only: login, cloudstatus, completion)
	const contexts = getCustomDomainContexts();

	console.log(`\nProcessing ${contexts.length} custom domains/subcommands...\n`);

	// Process each domain context
	for (const context of contexts) {
		console.log(`\nProcessing: ${context.type} "${context.name}"`);

		const descriptions = await generateWithRetry(context);

		if (!descriptions) {
			console.error(`  FAILED to generate descriptions for ${context.name}`);
			continue;
		}

		// Calculate hash from context
		const hash = calculateHash(JSON.stringify(context));

		if (context.type === "domain") {
			generatedData.domains[context.name] = {
				source_patterns_hash: hash,
				short: descriptions.short,
				medium: descriptions.medium,
				long: descriptions.long,
			};
		} else if (context.type === "subcommand" && context.parentName) {
			if (!generatedData.domains[context.parentName]) {
				generatedData.domains[context.parentName] = {
					source_patterns_hash: hash,
					short: "",
					medium: "",
					long: "",
					subcommands: {},
				};
			}
			if (!generatedData.domains[context.parentName].subcommands) {
				generatedData.domains[context.parentName].subcommands = {};
			}
			generatedData.domains[context.parentName].subcommands![context.name] = {
				short: descriptions.short,
				medium: descriptions.medium,
				long: descriptions.long,
			};
		}

		console.log(`  Short: "${descriptions.short.slice(0, 50)}..."`);
	}

	// Write YAML output
	console.log(`\nWriting output to: ${outputPath}`);
	const yamlContent = formatToYaml(generatedData);
	writeFileSync(outputPath, yamlContent, "utf-8");

	console.log("\nDone!");
}

// Run main
main().catch((err) => {
	console.error("Fatal error:", err);
	process.exit(1);
});
