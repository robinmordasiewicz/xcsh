/**
 * Completer - Tab completion engine for the REPL
 * Provides context-aware suggestions for domains, actions, flags, and built-in commands
 */

import type { CompletionSuggestion, ParsedInput } from "./types.js";
import type { REPLSession } from "../session.js";
import { domainRegistry } from "../../types/domains.js";
import { customDomains, isCustomDomain } from "../../domains/index.js";
import { CompletionCache } from "./cache.js";

/**
 * Parse input text into args array, handling quoted strings
 */
export function parseInputArgs(text: string): string[] {
	const args: string[] = [];
	let current = "";
	let inQuote = false;
	let quoteChar = "";

	for (const char of text) {
		if (inQuote) {
			if (char === quoteChar) {
				inQuote = false;
			} else {
				current += char;
			}
		} else if (char === '"' || char === "'") {
			inQuote = true;
			quoteChar = char;
		} else if (char === " " || char === "\t") {
			if (current) {
				args.push(current);
				current = "";
			}
		} else {
			current += char;
		}
	}

	if (current) {
		args.push(current);
	}

	return args;
}

/**
 * Parse input for completion context
 */
export function parseInput(text: string): ParsedInput {
	const trimmed = text.trimStart();
	const args = parseInputArgs(trimmed);

	// Check for "/" prefix escape
	let isEscapedToRoot = false;
	if (args.length > 0 && args[0]?.startsWith("/")) {
		isEscapedToRoot = true;
		args[0] = args[0].slice(1);
		if (args[0] === "") {
			args.shift();
		}
	}

	// Check if input ends with whitespace - means we're starting a new word
	const endsWithSpace = trimmed.length > 0 && /\s$/.test(trimmed);

	// Get current word being typed
	// If input ends with space, we're starting a new word (empty currentWord)
	const currentWord = endsWithSpace ? "" : (args[args.length - 1] ?? "");

	// Check if completing a flag
	const isCompletingFlag = currentWord.startsWith("-");

	// Check if completing a flag value
	// Patterns: "--flag " or "--flag=value"
	let isCompletingFlagValue = false;
	let currentFlag: string | null = null;

	if (args.length >= 2 && !currentWord.startsWith("-")) {
		// Check if previous arg is a flag that expects a value
		const prevArg = args[args.length - 2];
		if (prevArg && prevArg.startsWith("-") && !prevArg.includes("=")) {
			// Check if it's a flag that expects a value (not a boolean flag)
			const valueFlagPatterns = [
				"--namespace",
				"-ns",
				"--output",
				"-o",
				"--name",
				"-n",
				"--file",
				"-f",
				"--limit",
				"--label",
			];
			if (valueFlagPatterns.some((f) => prevArg === f)) {
				isCompletingFlagValue = true;
				currentFlag = prevArg;
			}
		}
	}

	// Check for --flag=value pattern
	if (currentWord.includes("=")) {
		const eqIndex = currentWord.indexOf("=");
		const flagPart = currentWord.slice(0, eqIndex);
		if (flagPart.startsWith("-")) {
			isCompletingFlagValue = true;
			currentFlag = flagPart;
		}
	}

	return {
		args,
		currentWord,
		isEscapedToRoot,
		isCompletingFlag,
		isCompletingFlagValue,
		currentFlag,
	};
}

/**
 * Completer provides context-aware tab completion
 */
export class Completer {
	private session: REPLSession | null = null;
	private cache: CompletionCache;

	constructor() {
		this.cache = new CompletionCache();
	}

	/**
	 * Set the session for context-aware completions
	 */
	setSession(session: REPLSession): void {
		this.session = session;
	}

	/**
	 * Get suggestions for the given input text
	 */
	async complete(text: string): Promise<CompletionSuggestion[]> {
		const trimmed = text.trimStart();

		// Empty input - show contextual suggestions
		if (trimmed === "") {
			return await this.getContextualSuggestions();
		}

		const parsed = parseInput(trimmed);

		// "/" alone - show root suggestions
		if (parsed.isEscapedToRoot && parsed.args.length === 0) {
			return this.getRootContextSuggestions();
		}

		// Check if first arg is a custom domain - delegate to domain completion
		const firstArg = parsed.args[0]?.toLowerCase() ?? "";
		if (isCustomDomain(firstArg) && parsed.args.length >= 1) {
			return await this.getCustomDomainCompletions(
				firstArg,
				parsed.args.slice(1),
				parsed.currentWord,
			);
		}

		// Completing a flag value (e.g., "--output json" or "--output=json")
		if (parsed.isCompletingFlagValue && parsed.currentFlag) {
			return await this.getFlagValueCompletions(
				parsed.currentFlag,
				parsed.currentWord,
			);
		}

		// Completing a flag
		if (parsed.isCompletingFlag) {
			return this.getFlagCompletions(parsed.currentWord);
		}

		// Get base suggestions based on context
		let suggestions: CompletionSuggestion[];
		if (parsed.isEscapedToRoot) {
			suggestions = this.getRootContextSuggestions();
		} else {
			suggestions = await this.getContextualSuggestions();
		}

		// Filter by current word if present
		if (parsed.currentWord) {
			return this.filterSuggestions(suggestions, parsed.currentWord);
		}

		return suggestions;
	}

	/**
	 * Get completions for custom domain commands
	 */
	async getCustomDomainCompletions(
		domainName: string,
		args: string[],
		currentWord: string,
	): Promise<CompletionSuggestion[]> {
		const domain = customDomains.get(domainName);
		if (!domain) {
			return [];
		}

		const suggestions: CompletionSuggestion[] = [];

		// No args after domain - suggest subcommands and direct commands
		if (
			args.length === 0 ||
			(args.length === 1 && currentWord === args[0])
		) {
			// Add subcommand groups
			for (const [name, group] of domain.subcommands) {
				if (
					!currentWord ||
					name.toLowerCase().startsWith(currentWord.toLowerCase())
				) {
					suggestions.push({
						text: name,
						description: group.description,
						category: "subcommand",
					});
				}
			}

			// Add direct commands
			for (const [name, cmd] of domain.commands) {
				if (
					!currentWord ||
					name.toLowerCase().startsWith(currentWord.toLowerCase())
				) {
					suggestions.push({
						text: name,
						description: cmd.description,
						category: "command",
					});
				}
			}

			return suggestions;
		}

		// First arg is a subcommand group - suggest commands within
		const subgroupName = args[0]?.toLowerCase() ?? "";
		const subgroup = domain.subcommands.get(subgroupName);
		if (subgroup) {
			// Only one arg so far (the subgroup name) - suggest commands
			if (
				args.length === 1 ||
				(args.length === 2 && currentWord === args[1])
			) {
				const cmdPrefix = args.length === 2 ? currentWord : "";
				for (const [name, cmd] of subgroup.commands) {
					if (
						!cmdPrefix ||
						name.toLowerCase().startsWith(cmdPrefix.toLowerCase())
					) {
						suggestions.push({
							text: name,
							description: cmd.description,
							category: "command",
						});
					}
				}
				return suggestions;
			}

			// Deeper completion - delegate to command's completion handler
			if (args.length >= 2 && this.session) {
				const cmdName = args[1]?.toLowerCase() ?? "";
				const cmd = subgroup.commands.get(cmdName);
				if (cmd?.completion) {
					try {
						const completions = await cmd.completion(
							currentWord,
							args.slice(2),
							this.session,
						);
						return completions.map((text) => ({
							text,
							description: "",
							category: "argument" as const,
						}));
					} catch {
						// Ignore completion errors
					}
				}
			}
		}

		return suggestions;
	}

	/**
	 * Get suggestions based on current navigation context
	 */
	async getContextualSuggestions(): Promise<CompletionSuggestion[]> {
		if (!this.session) {
			return this.getRootContextSuggestions();
		}

		const ctx = this.session.getContextPath();

		if (ctx.isRoot()) {
			return this.getRootContextSuggestions();
		}

		if (ctx.isDomain()) {
			return this.getDomainContextSuggestions();
		}

		if (ctx.isAction()) {
			return await this.getActionContextSuggestions();
		}

		return this.getRootContextSuggestions();
	}

	/**
	 * Get suggestions for root context
	 */
	getRootContextSuggestions(): CompletionSuggestion[] {
		const suggestions: CompletionSuggestion[] = [];

		// Add domains from registry
		suggestions.push(...this.getDomainSuggestions());

		// Add built-in commands
		suggestions.push(
			{
				text: "quit",
				description: "Exit the shell",
				category: "builtin",
			},
			{
				text: "help",
				description: "Show help information",
				category: "builtin",
			},
			{
				text: "clear",
				description: "Clear the screen",
				category: "builtin",
			},
			{
				text: "history",
				description: "Show command history",
				category: "builtin",
			},
			{
				text: "namespace",
				description: "Set default namespace",
				category: "builtin",
			},
			{
				text: "ns",
				description: "Set default namespace (alias)",
				category: "builtin",
			},
			{
				text: "context",
				description: "Show current context",
				category: "builtin",
			},
			{
				text: "ctx",
				description: "Show current context (alias)",
				category: "builtin",
			},
		);

		return suggestions;
	}

	/**
	 * Get suggestions when in a domain context
	 */
	getDomainContextSuggestions(): CompletionSuggestion[] {
		const suggestions: CompletionSuggestion[] = [];

		// Add actions
		suggestions.push(...this.getActionSuggestions());

		// Add navigation commands
		suggestions.push(
			{
				text: "exit",
				description: "Go up to root context",
				category: "navigation",
			},
			{
				text: "back",
				description: "Go up to root context",
				category: "navigation",
			},
			{
				text: "..",
				description: "Go up to root context",
				category: "navigation",
			},
			{
				text: "help",
				description: "Show context help",
				category: "builtin",
			},
		);

		return suggestions;
	}

	/**
	 * Get suggestions when in an action context
	 */
	async getActionContextSuggestions(): Promise<CompletionSuggestion[]> {
		const suggestions: CompletionSuggestion[] = [];

		// Add flags for current action
		suggestions.push(...this.getActionFlagSuggestions());

		// Add navigation commands
		suggestions.push(
			{
				text: "exit",
				description: "Go up to domain context",
				category: "navigation",
			},
			{
				text: "back",
				description: "Go up to domain context",
				category: "navigation",
			},
			{
				text: "..",
				description: "Go up to domain context",
				category: "navigation",
			},
			{
				text: "root",
				description: "Go to root context",
				category: "navigation",
			},
			{
				text: "/",
				description: "Go to root context",
				category: "navigation",
			},
			{
				text: "help",
				description: "Show context help",
				category: "builtin",
			},
		);

		return suggestions;
	}

	/**
	 * Get domain suggestions from registry
	 */
	getDomainSuggestions(): CompletionSuggestion[] {
		const suggestions: CompletionSuggestion[] = [];

		// Add custom domains first (higher priority)
		for (const domain of customDomains.all()) {
			suggestions.push({
				text: domain.name,
				description: domain.description,
				category: "domain",
			});
		}

		// Add API-generated domains (skip if already added as custom)
		for (const [domain, meta] of domainRegistry) {
			if (isCustomDomain(domain)) continue;

			suggestions.push({
				text: domain,
				description: meta.description,
				category: "domain",
			});
		}

		return suggestions;
	}

	/**
	 * Get action suggestions
	 */
	getActionSuggestions(): CompletionSuggestion[] {
		return [
			{ text: "list", description: "List resources", category: "action" },
			{
				text: "get",
				description: "Get a specific resource",
				category: "action",
			},
			{
				text: "create",
				description: "Create a new resource",
				category: "action",
			},
			{
				text: "delete",
				description: "Delete a resource",
				category: "action",
			},
			{
				text: "replace",
				description: "Replace a resource",
				category: "action",
			},
			{
				text: "apply",
				description: "Apply configuration from file",
				category: "action",
			},
			{
				text: "status",
				description: "Get resource status",
				category: "action",
			},
			{
				text: "patch",
				description: "Patch a resource",
				category: "action",
			},
			{
				text: "add-labels",
				description: "Add labels to a resource",
				category: "action",
			},
			{
				text: "remove-labels",
				description: "Remove labels from a resource",
				category: "action",
			},
		];
	}

	/**
	 * Get flag suggestions for current action
	 */
	getActionFlagSuggestions(): CompletionSuggestion[] {
		// Common flags that apply to most actions
		const commonFlags: CompletionSuggestion[] = [
			{ text: "--name", description: "Resource name", category: "flag" },
			{
				text: "-n",
				description: "Resource name (short)",
				category: "flag",
			},
			{ text: "--namespace", description: "Namespace", category: "flag" },
			{
				text: "--output",
				description: "Output format (json, yaml, table)",
				category: "flag",
			},
			{
				text: "-o",
				description: "Output format (short)",
				category: "flag",
			},
		];

		if (!this.session) {
			return commonFlags;
		}

		const ctx = this.session.getContextPath();
		const action = ctx.action;

		// Add action-specific flags
		const actionFlags: CompletionSuggestion[] = [...commonFlags];

		switch (action) {
			case "list":
				actionFlags.push(
					{
						text: "--limit",
						description: "Maximum results to return",
						category: "flag",
					},
					{
						text: "--label",
						description: "Filter by label",
						category: "flag",
					},
				);
				break;
			case "get":
				actionFlags.push({
					text: "--show-labels",
					description: "Show resource labels",
					category: "flag",
				});
				break;
			case "create":
			case "apply":
				actionFlags.push(
					{
						text: "--file",
						description: "Configuration file path",
						category: "flag",
					},
					{
						text: "-f",
						description: "Configuration file (short)",
						category: "flag",
					},
				);
				break;
			case "delete":
				actionFlags.push(
					{
						text: "--force",
						description: "Force deletion",
						category: "flag",
					},
					{
						text: "--cascade",
						description: "Cascade delete",
						category: "flag",
					},
				);
				break;
		}

		return actionFlags;
	}

	/**
	 * Get flag completions filtered by prefix
	 */
	getFlagCompletions(prefix: string): CompletionSuggestion[] {
		const allFlags = this.getActionFlagSuggestions();
		return this.filterSuggestions(allFlags, prefix);
	}

	/**
	 * Get flag value completions based on the flag being completed
	 */
	async getFlagValueCompletions(
		flag: string,
		partial: string,
	): Promise<CompletionSuggestion[]> {
		// Extract partial after "=" if present
		let valuePartial = partial;
		if (partial.includes("=")) {
			valuePartial = partial.slice(partial.indexOf("=") + 1);
		}

		const lowerPartial = valuePartial.toLowerCase();

		switch (flag) {
			case "--output":
			case "-o":
				return [
					{
						text: "json",
						description: "JSON format",
						category: "value" as const,
					},
					{
						text: "yaml",
						description: "YAML format",
						category: "value" as const,
					},
					{
						text: "table",
						description: "Table format",
						category: "value" as const,
					},
				].filter((s) => s.text.toLowerCase().startsWith(lowerPartial));

			case "--namespace":
			case "-ns":
				return this.completeNamespace(valuePartial);

			case "--name":
			case "-n":
				// Resource name completion - requires domain context
				if (this.session) {
					const ctx = this.session.getContextPath();
					if (ctx.domain) {
						return this.completeResourceName(
							ctx.domain,
							ctx.domain,
							valuePartial,
						);
					}
				}
				return [];

			case "--limit":
				// Common limit values
				return [
					{
						text: "10",
						description: "10 results",
						category: "value" as const,
					},
					{
						text: "25",
						description: "25 results",
						category: "value" as const,
					},
					{
						text: "50",
						description: "50 results",
						category: "value" as const,
					},
					{
						text: "100",
						description: "100 results",
						category: "value" as const,
					},
				].filter((s) => s.text.startsWith(valuePartial));

			default:
				return [];
		}
	}

	/**
	 * Complete namespace names with caching
	 */
	async completeNamespace(partial: string): Promise<CompletionSuggestion[]> {
		const namespaces = await this.cache.getNamespaces(async () => {
			// TODO: Fetch from API client when available
			// For now, return common defaults
			return ["default", "system", "shared"];
		});

		return namespaces
			.filter((ns) => ns.toLowerCase().startsWith(partial.toLowerCase()))
			.map((ns) => ({
				text: ns,
				description: "Namespace",
				category: "value" as const,
			}));
	}

	/**
	 * Complete resource names with caching
	 */
	async completeResourceName(
		domain: string,
		resourceType: string,
		partial: string,
	): Promise<CompletionSuggestion[]> {
		const cacheKey = `${domain}:${resourceType}`;
		const names = await this.cache.getResourceNames(cacheKey, async () => {
			// TODO: Fetch from API client when available
			// For now, return empty - requires API integration
			return [];
		});

		return names
			.filter((name) =>
				name.toLowerCase().startsWith(partial.toLowerCase()),
			)
			.map((name) => ({
				text: name,
				description: "Resource",
				category: "argument" as const,
			}));
	}

	/**
	 * Get the completion cache (for testing/debugging)
	 */
	getCache(): CompletionCache {
		return this.cache;
	}

	/**
	 * Filter suggestions by prefix (case-insensitive)
	 */
	filterSuggestions(
		suggestions: CompletionSuggestion[],
		prefix: string,
	): CompletionSuggestion[] {
		if (!prefix) {
			return suggestions;
		}

		const lowerPrefix = prefix.toLowerCase();
		return suggestions.filter((s) =>
			s.text.toLowerCase().startsWith(lowerPrefix),
		);
	}
}

/**
 * Create and configure a completer for the session
 */
export function createCompleter(session?: REPLSession): Completer {
	const completer = new Completer();
	if (session) {
		completer.setSession(session);
	}
	return completer;
}

export default Completer;
