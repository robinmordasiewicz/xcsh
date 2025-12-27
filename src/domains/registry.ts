/**
 * Domain Registry - Extensible custom command domain system
 * Mirrors Go implementation in xcsh-go/cmd/domains.go
 */

import type { REPLSession } from "../repl/session.js";

/**
 * Result from domain command execution
 * Compatible with ExecutionResult from executor
 */
export interface DomainCommandResult {
	/** Output lines to display */
	output: string[];
	/** Whether to exit the REPL */
	shouldExit: boolean;
	/** Whether to clear the screen */
	shouldClear: boolean;
	/** Whether the command modified context */
	contextChanged: boolean;
	/** Error message if command failed */
	error?: string;
}

/**
 * Command handler function signature
 */
export type CommandHandler = (
	args: string[],
	session: REPLSession,
) => Promise<DomainCommandResult>;

/**
 * Completion handler for a command
 */
export type CompletionHandler = (
	partial: string,
	args: string[],
	session: REPLSession,
) => Promise<string[]>;

/**
 * Definition of a single command within a domain
 */
export interface CommandDefinition {
	/** Command name (e.g., "list", "show", "create") */
	name: string;
	/** Short description for help/completion */
	description: string;
	/** Usage pattern (e.g., "<name> [options]") */
	usage?: string;
	/** Command execution handler */
	execute: CommandHandler;
	/** Optional completion handler */
	completion?: CompletionHandler;
	/** Aliases for this command */
	aliases?: string[];
}

/**
 * Definition of a subcommand group (e.g., "profile" under "login")
 */
export interface SubcommandGroup {
	/** Group name (e.g., "profile") */
	name: string;
	/** Short description */
	description: string;
	/** Commands within this group */
	commands: Map<string, CommandDefinition>;
}

/**
 * Definition of a custom domain
 */
export interface DomainDefinition {
	/** Domain name (e.g., "login") */
	name: string;
	/** Short description */
	description: string;
	/** Direct commands at domain level (e.g., "login" itself) */
	commands: Map<string, CommandDefinition>;
	/** Subcommand groups (e.g., "login profile") */
	subcommands: Map<string, SubcommandGroup>;
}

/**
 * Registry of custom domains
 * Custom domains take precedence over API-generated domains
 */
class DomainRegistry {
	private domains: Map<string, DomainDefinition> = new Map();

	/**
	 * Register a custom domain
	 */
	register(domain: DomainDefinition): void {
		this.domains.set(domain.name, domain);
	}

	/**
	 * Check if a domain is registered
	 */
	has(name: string): boolean {
		return this.domains.has(name);
	}

	/**
	 * Get a domain by name
	 */
	get(name: string): DomainDefinition | undefined {
		return this.domains.get(name);
	}

	/**
	 * Get all registered domain names
	 */
	list(): string[] {
		return Array.from(this.domains.keys());
	}

	/**
	 * Get all domains
	 */
	all(): DomainDefinition[] {
		return Array.from(this.domains.values());
	}

	/**
	 * Execute a command within a domain
	 * Handles routing through subcommand groups
	 *
	 * @param domainName - Name of the domain (e.g., "login")
	 * @param args - Command arguments (e.g., ["profile", "list"])
	 * @param session - REPL session
	 */
	async execute(
		domainName: string,
		args: string[],
		session: REPLSession,
	): Promise<DomainCommandResult> {
		const domain = this.domains.get(domainName);
		if (!domain) {
			return {
				output: [`Unknown domain: ${domainName}`],
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
				error: "Unknown domain",
			};
		}

		// No args - show domain help
		if (args.length === 0) {
			return this.showDomainHelp(domain);
		}

		const firstArg = args[0]?.toLowerCase() ?? "";
		const restArgs = args.slice(1);

		// Check for subcommand group first (e.g., "profile" in "login profile list")
		const subgroup = domain.subcommands.get(firstArg);
		if (subgroup) {
			// Need a command within the subgroup
			if (restArgs.length === 0) {
				return this.showSubcommandHelp(domain, subgroup);
			}

			const cmdName = restArgs[0]?.toLowerCase() ?? "";
			const cmdArgs = restArgs.slice(1);

			// Find command in subgroup
			const cmd = subgroup.commands.get(cmdName);
			if (cmd) {
				return cmd.execute(cmdArgs, session);
			}

			// Check aliases
			for (const [, command] of subgroup.commands) {
				if (command.aliases?.includes(cmdName)) {
					return command.execute(cmdArgs, session);
				}
			}

			return {
				output: [
					`Unknown command: ${domainName} ${firstArg} ${cmdName}`,
					``,
					`Run '${domainName} ${firstArg}' for available commands.`,
				],
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
				error: "Unknown command",
			};
		}

		// Check for direct command at domain level
		const cmd = domain.commands.get(firstArg);
		if (cmd) {
			return cmd.execute(restArgs, session);
		}

		// Check aliases
		for (const [, command] of domain.commands) {
			if (command.aliases?.includes(firstArg)) {
				return command.execute(restArgs, session);
			}
		}

		return {
			output: [
				`Unknown command: ${domainName} ${firstArg}`,
				``,
				`Run '${domainName}' for available commands.`,
			],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
			error: "Unknown command",
		};
	}

	/**
	 * Get completions for a domain command
	 */
	async getCompletions(
		domainName: string,
		args: string[],
		partial: string,
		session: REPLSession,
	): Promise<Array<{ text: string; description: string; category: string }>> {
		const domain = this.domains.get(domainName);
		if (!domain) {
			return [];
		}

		const suggestions: Array<{
			text: string;
			description: string;
			category: string;
		}> = [];

		// No args yet - suggest subcommands and commands
		if (args.length === 0) {
			// Add subcommand groups
			for (const [name, group] of domain.subcommands) {
				if (name.toLowerCase().startsWith(partial.toLowerCase())) {
					suggestions.push({
						text: name,
						description: group.description,
						category: "subcommand",
					});
				}
			}

			// Add direct commands
			for (const [name, cmd] of domain.commands) {
				if (name.toLowerCase().startsWith(partial.toLowerCase())) {
					suggestions.push({
						text: name,
						description: cmd.description,
						category: "command",
					});
				}
			}

			return suggestions;
		}

		// First arg is a subcommand group
		const firstArg = args[0]?.toLowerCase() ?? "";
		const subgroup = domain.subcommands.get(firstArg);
		if (subgroup && args.length === 1) {
			// Suggest commands within the subgroup
			for (const [name, cmd] of subgroup.commands) {
				if (name.toLowerCase().startsWith(partial.toLowerCase())) {
					suggestions.push({
						text: name,
						description: cmd.description,
						category: "command",
					});
				}
			}
			return suggestions;
		}

		// Delegate to command's completion handler if available
		if (subgroup && args.length >= 2) {
			const cmdName = args[1]?.toLowerCase() ?? "";
			const cmd = subgroup.commands.get(cmdName);
			if (cmd?.completion) {
				const completions = await cmd.completion(
					partial,
					args.slice(2),
					session,
				);
				return completions.map((text) => ({
					text,
					description: "",
					category: "argument",
				}));
			}
		}

		return suggestions;
	}

	/**
	 * Show help for a domain
	 */
	private showDomainHelp(domain: DomainDefinition): DomainCommandResult {
		const output: string[] = [`${domain.name} - ${domain.description}`, ``];

		if (domain.subcommands.size > 0) {
			output.push(`Subcommands:`);
			for (const [name, group] of domain.subcommands) {
				output.push(`  ${name}  ${group.description}`);
			}
			output.push(``);
		}

		if (domain.commands.size > 0) {
			output.push(`Commands:`);
			for (const [name, cmd] of domain.commands) {
				output.push(`  ${name}  ${cmd.description}`);
			}
		}

		return {
			output,
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	}

	/**
	 * Show help for a subcommand group
	 */
	private showSubcommandHelp(
		domain: DomainDefinition,
		subgroup: SubcommandGroup,
	): DomainCommandResult {
		const output: string[] = [
			`${domain.name} ${subgroup.name} - ${subgroup.description}`,
			``,
			`Commands:`,
		];

		for (const [name, cmd] of subgroup.commands) {
			const usage = cmd.usage ? ` ${cmd.usage}` : "";
			output.push(`  ${name}${usage}  ${cmd.description}`);
		}

		return {
			output,
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	}
}

// Singleton instance
export const customDomains = new DomainRegistry();

/**
 * Helper to create a success result
 */
export function successResult(
	output: string[],
	contextChanged: boolean = false,
): DomainCommandResult {
	return {
		output,
		shouldExit: false,
		shouldClear: false,
		contextChanged,
	};
}

/**
 * Helper to create an error result
 */
export function errorResult(message: string): DomainCommandResult {
	return {
		output: [message],
		shouldExit: false,
		shouldClear: false,
		contextChanged: false,
		error: message,
	};
}
