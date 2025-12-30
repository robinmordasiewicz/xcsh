/**
 * Common Completion Generator Utilities
 *
 * Shared escaping functions and tree-walking utilities for shell generators.
 */

import type { CompletionNode, ShellType } from "../types.js";

/**
 * Escape a string for safe use in Bash completion scripts
 * Bash uses double-quoted strings for completions
 */
export function escapeForBash(str: string): string {
	return str
		.replace(/\\/g, "\\\\") // Escape backslashes first
		.replace(/"/g, '\\"') // Escape double quotes
		.replace(/\$/g, "\\$") // Escape dollar signs
		.replace(/`/g, "\\`"); // Escape backticks
}

/**
 * Escape a string for safe use in Zsh completion scripts
 * Zsh uses single-quoted strings with special handling for colons
 */
export function escapeForZsh(str: string): string {
	return str
		.replace(/'/g, "'\\''") // Escape single quotes (end quote, escaped quote, start quote)
		.replace(/:/g, "\\:"); // Escape colons (special in zsh completion format)
}

/**
 * Escape a string for safe use in Fish completion scripts
 * Fish uses single-quoted strings
 */
export function escapeForFish(str: string): string {
	return str.replace(/'/g, "\\'"); // Escape single quotes
}

/**
 * Escape a description for the specified shell type
 */
export function escapeDescription(str: string, shell: ShellType): string {
	switch (shell) {
		case "bash":
			return escapeForBash(str);
		case "zsh":
			return escapeForZsh(str);
		case "fish":
			return escapeForFish(str);
		default:
			return str;
	}
}

/**
 * Format a name:description pair for shell completion
 */
export function formatCompletionItem(
	name: string,
	description: string,
	shell: ShellType,
): string {
	const escapedDesc = escapeDescription(description, shell);

	switch (shell) {
		case "bash":
			// Bash doesn't typically show descriptions in basic completions
			return name;
		case "zsh":
			return `'${name}:${escapedDesc}'`;
		case "fish":
			return `-a "${name}" -d '${escapedDesc}'`;
		default:
			return name;
	}
}

/**
 * Collect all nodes at a given depth in the tree
 * Returns flat array of nodes with their full path
 */
export interface NodeWithPath {
	node: CompletionNode;
	path: string[];
}

/**
 * Walk the completion tree and collect nodes
 */
export function walkTree(
	root: CompletionNode,
	maxDepth: number = Infinity,
): NodeWithPath[] {
	const results: NodeWithPath[] = [];

	function walk(node: CompletionNode, path: string[], depth: number): void {
		results.push({ node, path });

		if (depth >= maxDepth || !node.children) {
			return;
		}

		for (const child of node.children.values()) {
			if (!child.hidden) {
				walk(child, [...path, child.name], depth + 1);
			}
		}
	}

	walk(root, [root.name], 0);
	return results;
}

/**
 * Get all leaf nodes (nodes without children)
 */
export function getLeafNodes(root: CompletionNode): NodeWithPath[] {
	return walkTree(root).filter(
		({ node }) => !node.children || node.children.size === 0,
	);
}

/**
 * Get children of a node as formatted completion items
 */
export function getChildCompletions(
	node: CompletionNode,
	shell: ShellType,
): string[] {
	if (!node.children) {
		return [];
	}

	const items: string[] = [];
	for (const child of node.children.values()) {
		if (!child.hidden) {
			items.push(
				formatCompletionItem(child.name, child.description, shell),
			);
		}
	}

	return items.sort();
}

/**
 * Get child names only (for simple completions)
 */
export function getChildNames(node: CompletionNode): string[] {
	if (!node.children) {
		return [];
	}

	return Array.from(node.children.values())
		.filter((child) => !child.hidden)
		.map((child) => child.name)
		.sort();
}

/**
 * Check if a node has nested children (grandchildren)
 */
export function hasNestedChildren(node: CompletionNode): boolean {
	if (!node.children) {
		return false;
	}

	for (const child of node.children.values()) {
		if (child.children && child.children.size > 0) {
			return true;
		}
	}

	return false;
}

/**
 * Get the depth of a completion tree
 */
export function getTreeDepth(node: CompletionNode): number {
	if (!node.children || node.children.size === 0) {
		return 0;
	}

	let maxChildDepth = 0;
	for (const child of node.children.values()) {
		const childDepth = getTreeDepth(child);
		if (childDepth > maxChildDepth) {
			maxChildDepth = childDepth;
		}
	}

	return 1 + maxChildDepth;
}
