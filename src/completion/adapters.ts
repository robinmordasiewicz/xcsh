/**
 * Completion Adapters
 *
 * Convert existing domain types to unified CompletionNode structure.
 */

import type { CompletionNode } from "./types.js";
import type {
	DomainDefinition,
	CommandDefinition,
	SubcommandGroup,
} from "../domains/registry.js";
import type { DomainInfo } from "../types/domains.js";
import { validActions } from "../types/domains.js";

/**
 * Action descriptions for API domains
 * These descriptions appear in shell completions
 */
const actionDescriptions: Record<string, string> = {
	list: "List resources",
	get: "Get a specific resource",
	create: "Create a new resource",
	delete: "Delete a resource",
	replace: "Replace a resource",
	apply: "Apply configuration from file",
	status: "Get resource status",
	patch: "Patch a resource",
	"add-labels": "Add labels to a resource",
	"remove-labels": "Remove labels from a resource",
};

/**
 * Convert a CommandDefinition to CompletionNode
 */
function fromCommand(cmd: CommandDefinition): CompletionNode {
	const node: CompletionNode = {
		name: cmd.name,
		description: cmd.descriptionShort,
		source: "custom",
	};
	if (cmd.aliases && cmd.aliases.length > 0) {
		node.aliases = cmd.aliases;
	}
	return node;
}

/**
 * Convert a SubcommandGroup to CompletionNode with children
 */
function fromSubcommandGroup(group: SubcommandGroup): CompletionNode {
	const children = new Map<string, CompletionNode>();

	for (const [name, cmd] of group.commands) {
		children.set(name, fromCommand(cmd));
	}

	const node: CompletionNode = {
		name: group.name,
		description: group.descriptionShort,
		source: "custom",
	};
	if (children.size > 0) {
		node.children = children;
	}
	return node;
}

/**
 * Convert a custom DomainDefinition to CompletionNode
 *
 * Custom domains have hierarchical structure:
 * - Direct commands (e.g., `login banner`)
 * - Subcommand groups with nested commands (e.g., `login profile list`)
 */
export function fromCustomDomain(domain: DomainDefinition): CompletionNode {
	const children = new Map<string, CompletionNode>();

	// Add direct commands
	for (const [name, cmd] of domain.commands) {
		children.set(name, fromCommand(cmd));
	}

	// Add subcommand groups
	for (const [name, group] of domain.subcommands) {
		children.set(name, fromSubcommandGroup(group));
	}

	const node: CompletionNode = {
		name: domain.name,
		description: domain.descriptionShort,
		source: "custom",
	};
	if (children.size > 0) {
		node.children = children;
	}
	return node;
}

/**
 * Convert an API DomainInfo to CompletionNode
 *
 * API domains have flat structure:
 * - All use the same set of actions (list, get, create, etc.)
 * - Actions are children of the domain
 */
export function fromApiDomain(info: DomainInfo): CompletionNode {
	const children = new Map<string, CompletionNode>();

	// Add all valid actions as children
	for (const action of validActions) {
		children.set(action, {
			name: action,
			description: actionDescriptions[action] ?? action,
			source: "api",
		});
	}

	const node: CompletionNode = {
		name: info.name,
		description: info.descriptionShort,
		children,
		source: "api",
	};
	if (info.aliases.length > 0) {
		node.aliases = info.aliases;
	}
	return node;
}

/**
 * Create action nodes for API domains
 * Returns the standard action set as CompletionNodes
 */
export function getApiActionNodes(): Map<string, CompletionNode> {
	const actions = new Map<string, CompletionNode>();

	for (const action of validActions) {
		actions.set(action, {
			name: action,
			description: actionDescriptions[action] ?? action,
			source: "api",
		});
	}

	return actions;
}

/**
 * Get action descriptions map
 */
export function getActionDescriptions(): Record<string, string> {
	return { ...actionDescriptions };
}
