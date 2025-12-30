/**
 * Domain Extensions Module
 *
 * Provides xcsh-specific functionality to augment upstream API domains.
 *
 * Design philosophy:
 * 1. Upstream First - prefer feature requests to upstream API repo
 * 2. xcsh-Specific Only - wizards, enhanced output, validation helpers
 * 3. Complement, Not Compete - unique names, no conflicts with API actions
 * 4. Single Source of Truth - upstream remains the authority
 *
 * Usage:
 *   import { extensionRegistry } from './extensions/index.js';
 *
 *   // Check if domain has extension commands
 *   const merged = extensionRegistry.getMergedDomain('sites');
 *   if (merged?.hasExtension) {
 *     // Extension commands available
 *   }
 */

// Export types
export type { DomainExtension, MergedDomain } from "./types.js";
export {
	RESERVED_API_ACTIONS,
	isReservedAction,
	validateExtension,
} from "./types.js";

// Export registry
export { ExtensionRegistry, extensionRegistry } from "./registry.js";

// Import extensions for registration
// Extensions are registered when this module loads
// import { extensionRegistry } from "./registry.js";

/**
 * Initialize all extensions
 * Called automatically on module load
 */
function initializeExtensions(): void {
	// No extensions currently registered
	// Extensions can be added here when needed:
	// extensionRegistry.register(sitesExtension);
	// extensionRegistry.register(virtualExtension);
}

// Initialize on module load
initializeExtensions();

/**
 * Re-initialize extensions
 * Useful for testing or dynamic reloading
 */
export function reinitializeExtensions(): void {
	// extensionRegistry.clearCache();
	initializeExtensions();
}
