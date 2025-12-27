/**
 * Domains index - Register all custom domains
 */

import { customDomains } from "./registry.js";
import { loginDomain } from "./login/index.js";

// Register custom domains
customDomains.register(loginDomain);

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

/**
 * Check if a domain name is a custom domain
 */
export function isCustomDomain(name: string): boolean {
	return customDomains.has(name);
}

/**
 * Get list of all custom domain names
 */
export function getCustomDomainNames(): string[] {
	return customDomains.list();
}
