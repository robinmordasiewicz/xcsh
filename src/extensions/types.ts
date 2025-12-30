/**
 * Domain Extension Types
 *
 * Extensions augment dynamically generated API domains with xcsh-specific
 * CLI commands. Extensions provide unique functionality only - anything
 * useful to other downstream projects should be submitted to upstream.
 *
 * Design philosophy:
 * 1. Upstream First - prefer feature requests to upstream API repo
 * 2. xcsh-Specific Only - wizards, enhanced output, validation helpers
 * 3. Complement, Not Compete - unique names, no conflicts with API actions
 * 4. Single Source of Truth - upstream remains the authority
 */

import type {
	CommandDefinition,
	SubcommandGroup,
	DomainCommandResult,
} from "../domains/registry.js";
import type { DomainInfo } from "../types/domains.js";

/**
 * Extension definition for augmenting API domains with xcsh-specific commands
 */
export interface DomainExtension {
	/**
	 * Target domain to extend (canonical name from upstream specs)
	 * @example "sites", "http_loadbalancer", "virtual_k8s"
	 */
	targetDomain: string;

	/**
	 * Description of what this extension provides
	 */
	description: string;

	/**
	 * Whether extension works standalone if upstream domain doesn't exist yet
	 * When true, extension commands are available even before upstream adds the domain
	 */
	standalone: boolean;

	/**
	 * xcsh-specific commands to add to the domain
	 * These should have unique names that don't conflict with API actions
	 * (list, get, create, delete, replace, apply, status, patch, etc.)
	 */
	commands: Map<string, CommandDefinition>;

	/**
	 * Optional subcommand groups for hierarchical organization
	 */
	subcommands?: Map<string, SubcommandGroup>;
}

/**
 * Merged domain view combining upstream API domain + xcsh extension
 */
export interface MergedDomain {
	/** Canonical domain name */
	name: string;

	/** Display name for UI (from upstream or extension) */
	displayName: string;

	/** Combined description */
	description: string;

	/** Source of this merged domain */
	source: "generated" | "extension" | "merged";

	/** Whether upstream domain exists in specs */
	hasGeneratedDomain: boolean;

	/** Whether xcsh extension exists */
	hasExtension: boolean;

	/** All extension commands (xcsh-specific) */
	extensionCommands: Map<string, CommandDefinition>;

	/** All extension subcommands */
	extensionSubcommands: Map<string, SubcommandGroup>;

	/** Domain metadata from upstream specs (if exists) */
	metadata?: DomainInfo;
}

/**
 * Standard API actions that extensions should NOT override
 * These names are reserved for upstream API operations
 */
export const RESERVED_API_ACTIONS = new Set([
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
 * Check if a command name conflicts with reserved API actions
 */
export function isReservedAction(name: string): boolean {
	return RESERVED_API_ACTIONS.has(name.toLowerCase());
}

/**
 * Validate extension commands don't conflict with API actions
 * @throws Error if any command name conflicts
 */
export function validateExtension(extension: DomainExtension): void {
	const conflicts: string[] = [];

	for (const [name] of extension.commands) {
		if (isReservedAction(name)) {
			conflicts.push(name);
		}
	}

	if (conflicts.length > 0) {
		throw new Error(
			`Extension "${extension.targetDomain}" has commands that conflict with API actions: ${conflicts.join(", ")}. ` +
				`Use unique names or submit a feature request to upstream.`,
		);
	}
}

// Re-export types from domain registry for convenience
export type { DomainCommandResult, CommandDefinition, SubcommandGroup };
