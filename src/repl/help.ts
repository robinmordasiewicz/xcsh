/**
 * Multi-context help system for xcsh.
 * Provides context-aware help at root, domain, and action levels.
 *
 * Uses the unified description resolver to source descriptions from:
 * - Upstream API specs (single source of truth for API domains)
 * - Custom descriptions for CLI-only domains (login, cloudstatus, completion)
 */

import {
	CLI_NAME,
	CLI_VERSION,
	CLI_FULL_NAME,
	CLI_DESCRIPTION_LONG,
	colorBoldWhite,
	colorDim,
	colorYellow,
	colorRed,
} from "../branding/index.js";
import {
	formatEnvVarsSection,
	formatConfigSection,
} from "../config/envvars.js";
import { OUTPUT_FORMAT_HELP } from "../output/types.js";
import {
	type DomainInfo,
	domainRegistry,
	validActions,
	getDomainInfo,
} from "../types/domains.js";
import type { DomainDefinition, SubcommandGroup } from "../domains/registry.js";
import {
	descriptionResolver,
	getOperationDescription,
} from "../descriptions/index.js";

/**
 * Wrap text to specified width with indentation.
 * Preserves words and wraps at word boundaries.
 */
function wrapText(text: string, width: number, indent: number): string[] {
	const prefix = " ".repeat(indent);
	const words = text.split(/\s+/);
	const lines: string[] = [];
	let currentLine = prefix;

	for (const word of words) {
		if (
			currentLine.length + word.length + 1 > width &&
			currentLine !== prefix
		) {
			lines.push(currentLine);
			currentLine = prefix + word;
		} else {
			currentLine += (currentLine === prefix ? "" : " ") + word;
		}
	}
	if (currentLine.trim()) {
		lines.push(currentLine);
	}
	return lines;
}

/**
 * Format root-level help with full details.
 * Includes: description, usage, examples, global flags, environment variables.
 * Matches the Go version's professional structure.
 */
export function formatRootHelp(): string[] {
	return [
		"",
		colorBoldWhite(`${CLI_NAME} - ${CLI_FULL_NAME} ${CLI_VERSION}`),
		"",
		"DESCRIPTION",
		...wrapText(CLI_DESCRIPTION_LONG, 80, 2),
		"",
		"USAGE",
		`  ${CLI_NAME}                              Enter interactive REPL mode`,
		`  ${CLI_NAME} <domain> <action>            Execute command non-interactively`,
		`  ${CLI_NAME} help [topic]                 Show help for a topic`,
		"",
		"EXAMPLES",
		`  ${CLI_NAME} tenant_and_identity list namespace   List all namespaces`,
		`  ${CLI_NAME} virtual get http_loadbalancer        Get a specific load balancer`,
		`  ${CLI_NAME} dns list                             List DNS zones`,
		`  ${CLI_NAME} waf list -ns prod                    List WAF policies in prod`,
		"",
		...formatDomainsSection(),
		"",
		...formatGlobalFlags(),
		"",
		...formatEnvVarsSection(),
		"",
		...formatConfigSection(),
		"",
		"NAVIGATION (Interactive Mode)",
		"  <domain>              Navigate into a domain (e.g., 'dns', 'lb')",
		"  /domain               Navigate directly to domain from anywhere",
		"  ..                    Go up one level",
		"  /                     Return to root",
		"  context               Show current navigation context",
		"",
		"BUILTINS",
		"  help                  Show this help",
		"  domains               List all available domains",
		"  clear                 Clear the screen",
		"  history               Show command history",
		"  quit, exit            Exit the shell",
		"",
	];
}

/**
 * Format global flags section.
 */
export function formatGlobalFlags(): string[] {
	return [
		"GLOBAL FLAGS",
		"  -v, --version         Show version number",
		"  -h, --help            Show this help",
		"  --no-color            Disable color output",
		`  -o, --output <fmt>    Output format (${OUTPUT_FORMAT_HELP})`,
		"  -ns, --namespace <ns> Target namespace",
		"  --spec                Output command specification as JSON (for AI)",
	];
}

/**
 * Format environment variables section.
 * Re-exports the centralized function from config/envvars for backwards compatibility.
 */
export function formatEnvironmentVariables(): string[] {
	return formatEnvVarsSection();
}

/**
 * Format domain-level help.
 * Shows domain-specific information WITHOUT global flags or environment variables.
 * Uses the description resolver to source descriptions from upstream API specs.
 */
export function formatDomainHelp(domain: DomainInfo): string[] {
	// Build header with optional icon
	const icon = domain.icon ? `${domain.icon} ` : "";
	const output: string[] = [
		"",
		colorBoldWhite(`${icon}${domain.displayName}`),
		"",
	];

	// Get description from resolver (uses upstream specs for API domains)
	const descriptions = descriptionResolver.getDomainDescription(domain.name);
	output.push(`  ${descriptions.long}`);
	output.push("");

	// Category, UI category, and complexity if available
	const meta: string[] = [];
	if (domain.uiCategory) meta.push(`Category: ${domain.uiCategory}`);
	else if (domain.category) meta.push(`Category: ${domain.category}`);
	if (domain.complexity) meta.push(`Complexity: ${domain.complexity}`);
	if (domain.requiresTier && domain.requiresTier !== "Standard") {
		meta.push(`Tier: ${domain.requiresTier}`);
	}
	if (meta.length > 0) {
		output.push(colorDim(`  ${meta.join("  |  ")}`));
		output.push("");
	}

	// Usage
	output.push("USAGE");
	output.push(`  ${CLI_NAME} ${domain.name} <action> [options]`);
	output.push("");

	// Actions - use operation descriptions from OpenAPI when available
	output.push("ACTIONS");

	for (const action of validActions) {
		// Try to get description from upstream OpenAPI spec first
		const opInfo = getOperationDescription(domain.name, action);
		let desc: string;

		if (opInfo) {
			// Use upstream operation summary/description
			desc = opInfo.summary || opInfo.description || action;
		} else {
			// Fallback to generic action description
			const actionDesc = descriptionResolver.getActionDescription(action);
			desc = actionDesc.short;
		}

		output.push(`  ${action.padEnd(16)} ${desc}`);
	}
	output.push("");

	// Primary resources with tier and dependencies (from upstream enrichment)
	if (domain.primaryResources && domain.primaryResources.length > 0) {
		output.push("RESOURCES");
		for (const resource of domain.primaryResources) {
			const tierBadge =
				resource.tier !== "Standard" ? ` [${resource.tier}]` : "";
			const icon = resource.icon ? `${resource.icon} ` : "";
			output.push(`  ${icon}${resource.name}${colorDim(tierBadge)}`);
			output.push(`      ${resource.descriptionShort}`);

			// Show dependencies if any
			if (resource.dependencies?.required?.length) {
				output.push(
					colorDim(
						`      Requires: ${resource.dependencies.required.join(", ")}`,
					),
				);
			}
			if (resource.dependencies?.optional?.length) {
				output.push(
					colorDim(
						`      Optional: ${resource.dependencies.optional.join(", ")}`,
					),
				);
			}
		}
		output.push("");
	}

	// Examples
	output.push("EXAMPLES");
	output.push(`  ${CLI_NAME} ${domain.name} list`);
	output.push(`  ${CLI_NAME} ${domain.name} get my-resource`);
	output.push(
		`  ${CLI_NAME} ${domain.name} create my-resource -f config.yaml`,
	);
	output.push(`  ${CLI_NAME} ${domain.name} delete my-resource`);
	output.push("");

	// Use cases if available
	if (domain.useCases && domain.useCases.length > 0) {
		output.push("USE CASES");
		for (const useCase of domain.useCases.slice(0, 5)) {
			output.push(`  - ${useCase}`);
		}
		output.push("");
	}

	// Related domains if available
	if (domain.relatedDomains && domain.relatedDomains.length > 0) {
		output.push("RELATED DOMAINS");
		output.push(`  ${domain.relatedDomains.join(", ")}`);
		output.push("");
	}

	// Aliases if available
	if (domain.aliases && domain.aliases.length > 0) {
		output.push("ALIASES");
		output.push(`  ${domain.aliases.join(", ")}`);
		output.push("");
	}

	// Footer with reference to global help
	output.push(colorDim(`For global options, run: ${CLI_NAME} --help`));
	output.push("");

	return output;
}

/**
 * Format danger level indicator with appropriate coloring.
 */
function formatDangerLevel(level: string): string {
	switch (level.toLowerCase()) {
		case "high":
			return colorRed(`[DANGER: ${level.toUpperCase()}]`);
		case "medium":
			return colorYellow(`[CAUTION: ${level.toUpperCase()}]`);
		default:
			return colorDim(`[Risk: ${level}]`);
	}
}

/**
 * Format action-level help.
 * Shows action-specific usage within a domain context.
 * Uses the description resolver to source descriptions from upstream OpenAPI specs.
 */
export function formatActionHelp(domainName: string, action: string): string[] {
	const domain = getDomainInfo(domainName);
	const displayDomain = domain?.displayName ?? domainName;

	// Try to get operation description from upstream OpenAPI specs
	const opInfo = getOperationDescription(domainName, action);

	// Get action description from resolver (fallback)
	const actionDesc = descriptionResolver.getActionDescription(action);

	// Determine description - prefer upstream OpenAPI if available
	const desc = opInfo
		? opInfo.description || opInfo.summary || actionDesc.medium
		: actionDesc.medium;

	// Default usage patterns
	const usagePatterns: Record<string, string> = {
		list: `${CLI_NAME} ${domainName} list [--limit N] [--label key=value]`,
		get: `${CLI_NAME} ${domainName} get <name> [-o json|yaml|table]`,
		create: `${CLI_NAME} ${domainName} create <name> -f <file.yaml>`,
		delete: `${CLI_NAME} ${domainName} delete <name>`,
		replace: `${CLI_NAME} ${domainName} replace <name> -f <file.yaml>`,
		apply: `${CLI_NAME} ${domainName} apply -f <file.yaml>`,
		status: `${CLI_NAME} ${domainName} status <name>`,
		patch: `${CLI_NAME} ${domainName} patch <name> -f <patch.yaml>`,
		"add-labels": `${CLI_NAME} ${domainName} add-labels <name> key=value`,
		"remove-labels": `${CLI_NAME} ${domainName} remove-labels <name> key`,
	};

	const usage =
		usagePatterns[action] ??
		`${CLI_NAME} ${domainName} ${action} [options]`;

	// Build header with danger level warning if applicable
	const output = ["", colorBoldWhite(`${displayDomain} - ${action}`)];

	// Show danger level warning prominently in header
	if (opInfo?.dangerLevel && opInfo.dangerLevel !== "low") {
		output.push(`  ${formatDangerLevel(opInfo.dangerLevel)}`);
	}

	output.push("");
	output.push(`  ${desc}`);
	output.push("");

	// Add purpose from x-ves-operation-metadata if available
	if (opInfo?.purpose) {
		output.push("PURPOSE");
		output.push(`  ${opInfo.purpose}`);
		output.push("");
	}

	output.push("USAGE");
	output.push(`  ${usage}`);
	output.push("");

	// Add resource type information if from upstream
	if (opInfo?.resourceType) {
		output.push("RESOURCE TYPE");
		output.push(`  ${opInfo.resourceType}`);
		output.push("");
	}

	// Namespace scope information (from upstream enrichment)
	if (opInfo?.namespaceScope) {
		output.push("NAMESPACE SCOPE");
		const scopeDesc: Record<string, string> = {
			system: "This operation is only available in the 'system' namespace",
			shared: "This operation is available in 'shared' namespace",
			any: "This operation works in any namespace",
		};
		output.push(
			`  ${scopeDesc[opInfo.namespaceScope] ?? opInfo.namespaceScope}`,
		);
		output.push("");
	}

	// Required and optional fields (from upstream enrichment)
	if (opInfo?.requiredFields && opInfo.requiredFields.length > 0) {
		output.push("REQUIRED FIELDS");
		for (const field of opInfo.requiredFields) {
			output.push(`  - ${field}`);
		}
		output.push("");
	}

	if (opInfo?.optionalFields && opInfo.optionalFields.length > 0) {
		output.push("OPTIONAL FIELDS");
		for (const field of opInfo.optionalFields.slice(0, 10)) {
			output.push(colorDim(`  - ${field}`));
		}
		if (opInfo.optionalFields.length > 10) {
			output.push(
				colorDim(`  ... and ${opInfo.optionalFields.length - 10} more`),
			);
		}
		output.push("");
	}

	// Side effects warning (from upstream enrichment)
	if (opInfo?.sideEffects) {
		const { creates, updates, deletes } = opInfo.sideEffects;
		const hasSideEffects =
			(creates && creates.length > 0) ||
			(updates && updates.length > 0) ||
			(deletes && deletes.length > 0);

		if (hasSideEffects) {
			output.push("SIDE EFFECTS");
			if (creates && creates.length > 0) {
				output.push(`  Creates: ${creates.join(", ")}`);
			}
			if (updates && updates.length > 0) {
				output.push(`  Updates: ${updates.join(", ")}`);
			}
			if (deletes && deletes.length > 0) {
				output.push(colorYellow(`  Deletes: ${deletes.join(", ")}`));
			}
			output.push("");
		}
	}

	// Common errors with solutions (from upstream enrichment)
	if (opInfo?.commonErrors && opInfo.commonErrors.length > 0) {
		output.push("COMMON ERRORS");
		for (const err of opInfo.commonErrors.slice(0, 3)) {
			output.push(`  ${err.code}: ${err.message}`);
			if (err.solution) {
				output.push(colorDim(`    Solution: ${err.solution}`));
			}
		}
		output.push("");
	}

	output.push("OPTIONS");
	output.push("  -n, --name <name>     Resource name");
	output.push("  -ns, --namespace <ns> Target namespace");
	output.push(
		`  -o, --output <fmt>    Output format (${OUTPUT_FORMAT_HELP})`,
	);
	output.push("  -f, --file <path>     Configuration file");
	output.push("");

	// External documentation link (from upstream enrichment)
	if (opInfo?.externalDocs) {
		output.push("DOCUMENTATION");
		const docDesc = opInfo.externalDocs.description || "See documentation";
		output.push(`  ${docDesc}:`);
		output.push(colorDim(`  ${opInfo.externalDocs.url}`));
		output.push("");
	}

	output.push(
		colorDim(`For domain help, run: ${CLI_NAME} ${domainName} --help`),
	);
	output.push("");

	return output;
}

/**
 * Format help for a specific topic.
 */
export function formatTopicHelp(topic: string): string[] {
	const lowerTopic = topic.toLowerCase();

	// Check if it's a domain
	const domainInfo = getDomainInfo(lowerTopic);
	if (domainInfo) {
		return formatDomainHelp(domainInfo);
	}

	// Check for special topics
	switch (lowerTopic) {
		case "domains":
			return formatDomainsHelp();
		case "actions":
			return formatActionsHelp();
		case "navigation":
		case "nav":
			return formatNavigationHelp();
		case "env":
		case "environment":
			return ["", ...formatEnvironmentVariables(), ""];
		case "flags":
			return ["", ...formatGlobalFlags(), ""];
		default:
			return [
				"",
				`Unknown help topic: ${topic}`,
				"",
				"Available topics:",
				"  domains      List all available domains",
				"  actions      List available actions",
				"  navigation   Navigation commands",
				"  env          Environment variables",
				"  flags        Global flags",
				"  <domain>     Help for a specific domain (e.g., 'help dns')",
				"",
			];
	}
}

/**
 * Format domains list help.
 */
export function formatDomainsHelp(): string[] {
	const output: string[] = ["", colorBoldWhite("Available Domains"), ""];

	// Group domains by category
	const categories = new Map<string, DomainInfo[]>();

	for (const domain of domainRegistry.values()) {
		const category = domain.category ?? "Other";
		if (!categories.has(category)) {
			categories.set(category, []);
		}
		categories.get(category)?.push(domain);
	}

	// Sort categories and display
	const sortedCategories = Array.from(categories.keys()).sort();

	for (const category of sortedCategories) {
		const domains = categories.get(category) ?? [];
		output.push(colorBoldWhite(`  ${category}`));

		for (const domain of domains.sort((a, b) =>
			a.name.localeCompare(b.name),
		)) {
			const aliases =
				domain.aliases.length > 0
					? colorDim(` (${domain.aliases.join(", ")})`)
					: "";
			output.push(
				`    ${domain.name.padEnd(24)} ${domain.descriptionShort}`,
			);
			if (aliases) {
				output.push(`    ${"".padEnd(24)} Aliases:${aliases}`);
			}
		}
		output.push("");
	}

	return output;
}

/**
 * Format actions help.
 */
export function formatActionsHelp(): string[] {
	return [
		"",
		colorBoldWhite("Available Actions"),
		"",
		"  list              List all resources in the namespace",
		"  get               Get a specific resource by name",
		"  create            Create a new resource from a file",
		"  delete            Delete a resource by name",
		"  replace           Replace a resource configuration",
		"  apply             Apply configuration (create or update)",
		"  status            Get resource status",
		"  patch             Patch specific fields of a resource",
		"  add-labels        Add labels to a resource",
		"  remove-labels     Remove labels from a resource",
		"",
		"USAGE",
		`  ${CLI_NAME} <domain> <action> [options]`,
		"",
		"EXAMPLES",
		`  ${CLI_NAME} dns list`,
		`  ${CLI_NAME} lb get my-loadbalancer`,
		`  ${CLI_NAME} waf create my-policy -f policy.yaml`,
		"",
	];
}

/**
 * Format navigation help.
 */
export function formatNavigationHelp(): string[] {
	return [
		"",
		colorBoldWhite("Navigation Commands"),
		"",
		"  <domain>          Navigate into a domain context",
		"  /<domain>         Navigate directly to domain from anywhere",
		"  ..                Go up one level in context",
		"  /                 Return to root context",
		"  back              Go up one level (same as ..)",
		"  root              Return to root (same as /)",
		"",
		"CONTEXT DISPLAY",
		"  context           Show current navigation context",
		"  ctx               Alias for context",
		"",
		"EXAMPLES",
		"  xcsh> dns                    # Enter dns domain",
		"  dns> list                    # Execute list in dns context",
		"  dns> ..                      # Return to root",
		"  xcsh> /waf                   # Jump directly to waf",
		"  waf> /dns                    # Jump from waf to dns",
		"",
	];
}

/**
 * Format DOMAINS section for root help.
 * Displays all registered domains with their descriptions,
 * dynamically generated from the domain registry (no hardcoding).
 */
export function formatDomainsSection(): string[] {
	const output: string[] = ["DOMAINS"];

	// Get all domains sorted alphabetically
	const domains = Array.from(domainRegistry.values()).sort((a, b) =>
		a.name.localeCompare(b.name),
	);

	// Calculate max name length for column alignment
	const maxNameLen = Math.max(...domains.map((d) => d.name.length));

	for (const domain of domains) {
		const padding = " ".repeat(maxNameLen - domain.name.length + 2);
		output.push(`  ${domain.name}${padding}${domain.descriptionShort}`);
	}

	return output;
}

/**
 * Format help for a custom domain (login, cloudstatus, completion).
 * Mirrors formatDomainHelp() structure for consistency across all domains.
 */
export function formatCustomDomainHelp(domain: DomainDefinition): string[] {
	const output: string[] = ["", colorBoldWhite(domain.name), ""];

	// Description
	output.push("DESCRIPTION");
	output.push(...wrapText(domain.description, 80, 2));
	output.push("");

	// Usage
	output.push("USAGE");
	output.push(`  ${CLI_NAME} ${domain.name} <command> [options]`);
	output.push("");

	// Subcommands (if any)
	if (domain.subcommands.size > 0) {
		output.push("SUBCOMMANDS");
		for (const [name, group] of domain.subcommands) {
			output.push(`  ${name.padEnd(16)} ${group.descriptionShort}`);
		}
		output.push("");
	}

	// Commands (if any)
	if (domain.commands.size > 0) {
		output.push("COMMANDS");
		for (const [name, cmd] of domain.commands) {
			output.push(`  ${name.padEnd(16)} ${cmd.descriptionShort}`);
		}
		output.push("");
	}

	// Footer
	output.push(colorDim(`For global options, run: ${CLI_NAME} --help`));
	output.push("");

	return output;
}

/**
 * Format help for a subcommand group (e.g., login profile).
 * Mirrors formatDomainHelp() structure for consistency.
 */
export function formatSubcommandHelp(
	domainName: string,
	subcommand: SubcommandGroup,
): string[] {
	const output: string[] = [
		"",
		colorBoldWhite(`${domainName} ${subcommand.name}`),
		"",
	];

	// Description
	output.push("DESCRIPTION");
	output.push(...wrapText(subcommand.description, 80, 2));
	output.push("");

	// Usage
	output.push("USAGE");
	output.push(
		`  ${CLI_NAME} ${domainName} ${subcommand.name} <command> [options]`,
	);
	output.push("");

	// Commands
	if (subcommand.commands.size > 0) {
		output.push("COMMANDS");
		for (const [name, cmd] of subcommand.commands) {
			const usage = cmd.usage ? ` ${cmd.usage}` : "";
			output.push(
				`  ${name}${usage.padEnd(16 - name.length)} ${cmd.descriptionShort}`,
			);
		}
		output.push("");
	}

	// Footer
	output.push(
		colorDim(`For domain help, run: ${CLI_NAME} ${domainName} --help`),
	);
	output.push("");

	return output;
}
