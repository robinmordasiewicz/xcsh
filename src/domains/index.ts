/**
 * Domains index - Register all custom domains
 */

import { customDomains } from "./registry.js";
import { loginDomain } from "./login/index.js";
import { cloudstatusDomain, cloudstatusAliases } from "./cloudstatus/index.js";
import { completionDomain } from "./completion/index.js";
import { domainRegistry } from "../types/domains.js";
import {
	completionRegistry,
	fromCustomDomain,
	fromApiDomain,
} from "../completion/index.js";

// Register custom domains
// Only domains with no upstream API equivalent stay as custom domains
customDomains.register(loginDomain);
customDomains.register(cloudstatusDomain);
customDomains.register(completionDomain);

// Populate unified completion registry
// Custom domains first (higher priority)
for (const domain of customDomains.all()) {
	completionRegistry.registerDomain(fromCustomDomain(domain));
}

// API-generated domains
for (const [, info] of domainRegistry) {
	// Skip if already registered by custom domain (custom takes precedence)
	if (!completionRegistry.has(info.name)) {
		completionRegistry.registerDomain(fromApiDomain(info));
	}
}

// Domain alias mapping (alias -> canonical name)
const domainAliases = new Map<string, string>();

// Register cloudstatus aliases
for (const alias of cloudstatusAliases) {
	domainAliases.set(alias, "cloudstatus");
}

// Export registry and types
export { customDomains } from "./registry.js";
export type {
	DomainDefinition,
	SubcommandGroup,
	CommandDefinition,
	CommandHandler,
	CompletionHandler,
	DomainCommandResult,
} from "./registry.js";
export { successResult, errorResult } from "./registry.js";

// Export unified completion registry
export { completionRegistry } from "../completion/index.js";

// Export domain definitions for reference
export { loginDomain } from "./login/index.js";
export { cloudstatusDomain } from "./cloudstatus/index.js";
export { completionDomain } from "./completion/index.js";

/**
 * Resolve domain alias to canonical name
 */
export function resolveDomainAlias(name: string): string {
	return domainAliases.get(name) ?? name;
}

/**
 * Check if a domain name is a custom domain (including aliases)
 */
export function isCustomDomain(name: string): boolean {
	const canonical = resolveDomainAlias(name);
	return customDomains.has(canonical);
}

/**
 * Get list of all custom domain names
 */
export function getCustomDomainNames(): string[] {
	return customDomains.list();
}

/**
 * Get list of all domain aliases
 */
export function getDomainAliases(): Map<string, string> {
	return new Map(domainAliases);
}
