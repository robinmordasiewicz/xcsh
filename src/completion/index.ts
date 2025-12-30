/**
 * Unified Completion Module
 *
 * Single source of truth for completion data across all domain types.
 * Provides a unified tree structure for custom and API domains.
 */

// Types
export type {
	CompletionNode,
	CompletionFlag,
	CompletionTree,
	AliasMap,
	ShellType,
	CompletionSuggestion,
} from "./types.js";

// Registry
export { CompletionRegistry, completionRegistry } from "./registry.js";

// Adapters
export {
	fromCustomDomain,
	fromApiDomain,
	getApiActionNodes,
	getActionDescriptions,
} from "./adapters.js";

// Generator utilities
export {
	escapeForBash,
	escapeForZsh,
	escapeForFish,
	escapeDescription,
	formatCompletionItem,
	walkTree,
	getLeafNodes,
	getChildCompletions,
	getChildNames,
	hasNestedChildren,
	getTreeDepth,
	type NodeWithPath,
} from "./generators/common.js";
