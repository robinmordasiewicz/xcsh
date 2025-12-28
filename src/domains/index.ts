/**
 * Domains index - Register all custom domains
 *
 * Note: subscription is now an extension, not a custom domain.
 * Extensions augment API domains with xcsh-specific commands.
 * See src/extensions/ for extension implementations.
 */

import { customDomains } from "./registry.js";
import { loginDomain } from "./login/index.js";
import { cloudstatusDomain, cloudstatusAliases } from "./cloudstatus/index.js";
import { completionDomain } from "./completion/index.js";

// Register custom domains
// Only domains with no upstream API equivalent stay as custom domains
customDomains.register(loginDomain);
customDomains.register(cloudstatusDomain);
customDomains.register(completionDomain);

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

// Export domain definitions for reference
export { loginDomain } from "./login/index.js";
export { cloudstatusDomain } from "./cloudstatus/index.js";
export { completionDomain } from "./completion/index.js";
// Note: subscription is now an extension - see src/extensions/subscription/

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
