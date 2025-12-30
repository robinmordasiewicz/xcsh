/**
 * Unified Completion Types
 *
 * Single source of truth for completion data across all domain types.
 * Both custom domains and API-generated domains use this structure.
 */

/**
 * Unified completion tree node
 * Works for domains, subcommands, actions, and flags
 */
export interface CompletionNode {
	/** Node name (domain, command, or action) */
	name: string;

	/** Short description for completion display (~60 chars) */
	description: string;

	/** Alternative names that resolve to this node */
	aliases?: string[];

	/** Child nodes (subcommands, actions, nested commands) */
	children?: Map<string, CompletionNode>;

	/** Available flags at this level */
	flags?: CompletionFlag[];

	/** Hide from completion suggestions */
	hidden?: boolean;

	/** Source type for debugging/prioritization */
	source?: "custom" | "api" | "extension";
}

/**
 * Flag definition for completion
 */
export interface CompletionFlag {
	/** Flag name (e.g., "--namespace", "-ns") */
	name: string;

	/** Short description for completion display */
	description: string;

	/** Short alias (e.g., "-n" for "--namespace") */
	shortName?: string;

	/** Flag expects a value argument */
	hasValue?: boolean;

	/** Predefined value completions */
	valueCompletions?: string[];

	/** Flag is required */
	required?: boolean;
}

/**
 * Complete completion tree (domain name → node)
 */
export type CompletionTree = Map<string, CompletionNode>;

/**
 * Alias registry (alias → canonical name)
 */
export type AliasMap = Map<string, string>;

/**
 * Shell types supported for completion generation
 */
export type ShellType = "bash" | "zsh" | "fish";

/**
 * Completion suggestion for REPL display
 */
export interface CompletionSuggestion {
	/** Text to insert */
	text: string;

	/** Description to display */
	description: string;

	/** Category for grouping */
	category: "domain" | "command" | "subcommand" | "action" | "flag" | "value";
}
