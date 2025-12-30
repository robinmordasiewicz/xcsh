/**
 * Completion Registry
 *
 * Single registry for all completion data.
 * Populated by custom domains, API domains, and extensions.
 */

import type {
	CompletionNode,
	CompletionTree,
	AliasMap,
	CompletionSuggestion,
} from "./types.js";

/**
 * Centralized completion registry
 * Provides unified access to completion data from all domain sources
 */
export class CompletionRegistry {
	private tree: CompletionTree = new Map();
	private aliases: AliasMap = new Map();

	/**
	 * Register a domain with its completion subtree
	 * Later registrations for same name take precedence (custom > api)
	 */
	registerDomain(node: CompletionNode): void {
		this.tree.set(node.name, node);

		// Register aliases
		if (node.aliases) {
			for (const alias of node.aliases) {
				this.aliases.set(alias, node.name);
			}
		}
	}

	/**
	 * Add or merge children into an existing domain
	 * Used for extensions that augment API domains
	 */
	mergeChildren(
		domainName: string,
		children: Map<string, CompletionNode>,
	): void {
		const existing = this.tree.get(domainName);
		if (!existing) {
			return;
		}

		if (!existing.children) {
			existing.children = new Map();
		}

		for (const [name, node] of children) {
			existing.children.set(name, node);
		}
	}

	/**
	 * Get the complete completion tree
	 */
	getTree(): CompletionTree {
		return this.tree;
	}

	/**
	 * Get all registered domains as array (sorted by name)
	 */
	getDomains(): CompletionNode[] {
		return Array.from(this.tree.values()).sort((a, b) =>
			a.name.localeCompare(b.name),
		);
	}

	/**
	 * Resolve alias to canonical name
	 */
	resolveAlias(nameOrAlias: string): string {
		return this.aliases.get(nameOrAlias) ?? nameOrAlias;
	}

	/**
	 * Get a domain by name or alias
	 */
	get(nameOrAlias: string): CompletionNode | undefined {
		const canonical = this.resolveAlias(nameOrAlias);
		return this.tree.get(canonical);
	}

	/**
	 * Check if a domain exists (by name or alias)
	 */
	has(nameOrAlias: string): boolean {
		const canonical = this.resolveAlias(nameOrAlias);
		return this.tree.has(canonical);
	}

	/**
	 * Get all aliases
	 */
	getAliases(): AliasMap {
		return new Map(this.aliases);
	}

	/**
	 * Get domain suggestions for completion
	 */
	getDomainSuggestions(prefix = ""): CompletionSuggestion[] {
		const suggestions: CompletionSuggestion[] = [];
		const lowerPrefix = prefix.toLowerCase();

		for (const node of this.tree.values()) {
			if (node.hidden) continue;

			if (!prefix || node.name.toLowerCase().startsWith(lowerPrefix)) {
				suggestions.push({
					text: node.name,
					description: node.description,
					category: "domain",
				});
			}

			// Also include aliases
			if (node.aliases) {
				for (const alias of node.aliases) {
					if (
						!prefix ||
						alias.toLowerCase().startsWith(lowerPrefix)
					) {
						suggestions.push({
							text: alias,
							description: `Alias for ${node.name}`,
							category: "domain",
						});
					}
				}
			}
		}

		return suggestions.sort((a, b) => a.text.localeCompare(b.text));
	}

	/**
	 * Get child suggestions for a domain
	 */
	getChildSuggestions(
		domainName: string,
		prefix = "",
	): CompletionSuggestion[] {
		const node = this.get(domainName);
		if (!node?.children) {
			return [];
		}

		const suggestions: CompletionSuggestion[] = [];
		const lowerPrefix = prefix.toLowerCase();

		for (const child of node.children.values()) {
			if (child.hidden) continue;

			if (!prefix || child.name.toLowerCase().startsWith(lowerPrefix)) {
				suggestions.push({
					text: child.name,
					description: child.description,
					category: child.children ? "subcommand" : "action",
				});
			}
		}

		return suggestions.sort((a, b) => a.text.localeCompare(b.text));
	}

	/**
	 * Get nested child suggestions (e.g., login profile <TAB>)
	 */
	getNestedChildSuggestions(
		domainName: string,
		path: string[],
		prefix = "",
	): CompletionSuggestion[] {
		let node = this.get(domainName);
		if (!node) {
			return [];
		}

		// Navigate down the path
		for (const segment of path) {
			if (!node.children) {
				return [];
			}
			const child = node.children.get(segment);
			if (!child) {
				return [];
			}
			node = child;
		}

		if (!node.children) {
			return [];
		}

		const suggestions: CompletionSuggestion[] = [];
		const lowerPrefix = prefix.toLowerCase();

		for (const child of node.children.values()) {
			if (child.hidden) continue;

			if (!prefix || child.name.toLowerCase().startsWith(lowerPrefix)) {
				suggestions.push({
					text: child.name,
					description: child.description,
					category: "command",
				});
			}
		}

		return suggestions.sort((a, b) => a.text.localeCompare(b.text));
	}

	/**
	 * Clear the registry (useful for testing)
	 */
	clear(): void {
		this.tree.clear();
		this.aliases.clear();
	}
}

/**
 * Singleton completion registry instance
 */
export const completionRegistry = new CompletionRegistry();
