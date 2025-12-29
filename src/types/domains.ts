/**
 * Domain type definitions and registry interfaces.
 * These types mirror the Go types/domains.go structure.
 */

import {
	generatedDomains,
	SPEC_VERSION,
	DOMAIN_COUNT,
} from "./domains_generated.js";

/**
 * DomainInfo contains metadata about a resource domain
 */
export interface DomainInfo {
	name: string; // Canonical: "load_balancer"
	displayName: string; // Human: "Load Balancer"
	description: string; // Long description (~500 chars) for detailed help
	descriptionShort: string; // Short description (~60 chars) for completions, badges
	descriptionMedium: string; // Medium description (~150 chars) for tooltips, summaries
	aliases: string[]; // Short forms: ["lb"]

	// Fields from upstream specs
	complexity?: "simple" | "moderate" | "advanced";
	isPreview?: boolean;
	requiresTier?: string;
	category?: string;
	useCases?: string[];
	relatedDomains?: string[];
	cliMetadata?: Record<string, unknown>;
}

/**
 * Domain registry mapping canonical names to domain info
 */
export type DomainRegistry = Map<string, DomainInfo>;

/**
 * Alias registry mapping aliases/canonical names to canonical names
 */
export type AliasRegistry = Map<string, string>;

// Global registries - populated from generated data
export const domainRegistry: DomainRegistry = new Map(generatedDomains);
export const aliasRegistry: AliasRegistry = new Map();

// Export spec metadata
export { SPEC_VERSION, DOMAIN_COUNT };

/**
 * Initialize the alias registry from domain registry
 */
export function initializeAliasRegistry(): void {
	aliasRegistry.clear();

	for (const [canonical, info] of domainRegistry) {
		// Map canonical name to itself
		aliasRegistry.set(canonical, canonical);

		// Map all aliases to canonical name
		for (const alias of info.aliases) {
			aliasRegistry.set(alias, canonical);
		}
	}
}

/**
 * Resolve an alias or canonical name to the canonical domain name
 */
export function resolveDomain(nameOrAlias: string): string | undefined {
	return aliasRegistry.get(nameOrAlias);
}

/**
 * Get domain info by canonical name or alias
 */
export function getDomainInfo(nameOrAlias: string): DomainInfo | undefined {
	const canonical = resolveDomain(nameOrAlias);
	if (!canonical) {
		return undefined;
	}
	return domainRegistry.get(canonical);
}

/**
 * Get all canonical domain names
 */
export function allDomains(): string[] {
	return Array.from(domainRegistry.keys()).sort();
}

/**
 * Check if a name is a valid domain (canonical or alias)
 */
export function isValidDomain(name: string): boolean {
	return aliasRegistry.has(name);
}

/**
 * Valid action commands for domains
 */
export const validActions = new Set([
	"list",
	"get",
	"create",
	"delete",
	"replace",
	"apply",
	"status",
	"patch",
	"add-labels",
	"remove-labels",
]);

/**
 * Check if a name is a valid action command
 */
export function isValidAction(name: string): boolean {
	return validActions.has(name);
}

// Initialize alias registry at module load
initializeAliasRegistry();
