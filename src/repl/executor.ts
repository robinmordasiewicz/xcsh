/**
 * Executor - Command execution engine for the REPL
 * Handles built-in commands, domain navigation, and command routing
 */

import type { REPLSession } from "./session.js";
import type { ContextPath } from "./context.js";
import { allDomains, isValidDomain } from "../types/domains.js";
import { customDomains, isCustomDomain } from "../domains/index.js";

/**
 * Command execution result
 */
export interface ExecutionResult {
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
 * Built-in command names
 */
const BUILTIN_COMMANDS = new Set([
	"help",
	"clear",
	"quit",
	"exit",
	"back",
	"..",
	"root",
	"/",
	"context",
	"ctx",
	"history",
	"version",
	"domains",
]);

/**
 * Parse command input into components
 */
interface ParsedCommand {
	/** Original raw input */
	raw: string;
	/** Is this a "/" prefixed domain navigation? */
	isDirectNavigation: boolean;
	/** Target domain (if navigation) */
	targetDomain?: string | undefined;
	/** Target action (if specified) */
	targetAction?: string | undefined;
	/** Command arguments */
	args: string[];
	/** Is this a built-in command? */
	isBuiltin: boolean;
}

/**
 * Parse user input into structured command
 */
export function parseCommand(input: string): ParsedCommand {
	const trimmed = input.trim();

	// Handle empty input
	if (!trimmed) {
		return {
			raw: trimmed,
			isDirectNavigation: false,
			args: [],
			isBuiltin: false,
		};
	}

	// Handle "/" prefix for direct domain navigation
	if (trimmed.startsWith("/") && trimmed.length > 1) {
		const parts = trimmed.slice(1).split(/\s+/);
		const domainPart = parts[0] ?? "";

		// Check if it's a valid domain (custom or API-generated)
		if (isValidDomain(domainPart) || isCustomDomain(domainPart)) {
			return {
				raw: trimmed,
				isDirectNavigation: true,
				targetDomain: domainPart,
				targetAction: parts[1],
				args: parts.slice(2),
				isBuiltin: false,
			};
		}
	}

	// Handle just "/" for root navigation
	if (trimmed === "/") {
		return {
			raw: trimmed,
			isDirectNavigation: false,
			isBuiltin: true,
			args: [],
		};
	}

	// Check for built-in commands
	const firstWord = trimmed.split(/\s+/)[0]?.toLowerCase() ?? "";
	if (BUILTIN_COMMANDS.has(firstWord)) {
		return {
			raw: trimmed,
			isDirectNavigation: false,
			isBuiltin: true,
			args: trimmed.split(/\s+/).slice(1),
		};
	}

	// Regular command - split into parts
	const parts = trimmed.split(/\s+/);
	return {
		raw: trimmed,
		isDirectNavigation: false,
		isBuiltin: false,
		args: parts,
	};
}

/**
 * Execute a built-in command
 */
function executeBuiltin(
	cmd: ParsedCommand,
	session: REPLSession,
	ctx: ContextPath,
): ExecutionResult {
	const command = cmd.raw.toLowerCase();

	// Clear screen
	if (command === "clear") {
		return {
			output: [],
			shouldExit: false,
			shouldClear: true,
			contextChanged: false,
		};
	}

	// Exit/quit
	if (command === "quit") {
		return {
			output: ["Goodbye!"],
			shouldExit: true,
			shouldClear: false,
			contextChanged: false,
		};
	}

	// Navigate up (exit/back/..)
	if (command === "exit" || command === "back" || command === "..") {
		if (ctx.isRoot()) {
			return {
				output: ["Goodbye!"],
				shouldExit: true,
				shouldClear: false,
				contextChanged: false,
			};
		}
		ctx.navigateUp();
		const location = ctx.isRoot() ? "root" : ctx.domain;
		return {
			output: [`Navigated to ${location} context`],
			shouldExit: false,
			shouldClear: false,
			contextChanged: true,
		};
	}

	// Navigate to root
	if (command === "/" || command === "root") {
		if (!ctx.isRoot()) {
			ctx.reset();
			return {
				output: ["Navigated to root context"],
				shouldExit: false,
				shouldClear: false,
				contextChanged: true,
			};
		}
		return {
			output: ["Already at root context"],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	}

	// Show current context
	if (command === "context" || command === "ctx") {
		let contextStr: string;
		if (ctx.isRoot()) {
			contextStr = "root";
		} else if (ctx.isDomain()) {
			contextStr = ctx.domain;
		} else {
			contextStr = `${ctx.domain} > ${ctx.action}`;
		}
		return {
			output: [`Current context: ${contextStr}`],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	}

	// Show history
	if (command === "history") {
		const histMgr = session.getHistory();
		const hist = histMgr?.getHistory() ?? [];
		if (hist.length === 0) {
			return {
				output: ["No command history"],
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
			};
		}
		const lines = ["Command history:"];
		hist.slice(-20).forEach((histCmd, i) => {
			lines.push(`  ${i + 1}. ${histCmd}`);
		});
		return {
			output: lines,
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	}

	// Show version
	if (command === "version") {
		return {
			output: ["xcsh version 0.1.0 (Ink rewrite)"],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	}

	// List available domains
	if (command === "domains") {
		const lines = ["Available domains:"];
		allDomains().forEach((domain) => {
			lines.push(`  ${domain}`);
		});
		return {
			output: lines,
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	}

	// Help command
	if (command === "help" || command.startsWith("help ")) {
		return {
			output: [
				"xcsh - F5 Distributed Cloud Shell",
				"",
				"Navigation:",
				"  /domain         Navigate directly to domain",
				"  exit, back, ..  Navigate up one level (exits at root)",
				"  /, root         Navigate to root context",
				"  context, ctx    Show current context",
				"",
				"Built-in commands:",
				"  help            Show this help",
				"  clear           Clear the screen",
				"  quit            Exit the shell",
				"  history         Show command history",
				"  domains         List available domains",
				"  version         Show version info",
				"",
				"Keyboard shortcuts:",
				"  Tab             Trigger/cycle completions",
				"  Up/Down         Navigate history or suggestions",
				"  Ctrl+C twice    Exit (within 500ms)",
				"  Ctrl+D          Exit immediately",
				"  Escape          Cancel suggestions",
			],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	}

	// Unknown built-in
	return {
		output: [`Unknown command: ${cmd.raw}`],
		shouldExit: false,
		shouldClear: false,
		contextChanged: false,
		error: "Unknown command",
	};
}

/**
 * Handle direct domain navigation with "/" prefix
 */
async function handleDirectNavigation(
	cmd: ParsedCommand,
	ctx: ContextPath,
	session: REPLSession,
): Promise<ExecutionResult> {
	if (!cmd.targetDomain) {
		return {
			output: ["Invalid domain"],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
			error: "Invalid domain",
		};
	}

	// Check if it's a custom domain - execute directly with all args
	if (isCustomDomain(cmd.targetDomain)) {
		const allArgs = [cmd.targetAction, ...cmd.args].filter(
			(arg): arg is string => arg !== undefined,
		);
		return customDomains.execute(cmd.targetDomain, allArgs, session);
	}

	// Validate API-generated domain
	if (!isValidDomain(cmd.targetDomain)) {
		return {
			output: [`Unknown domain: ${cmd.targetDomain}`],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
			error: "Unknown domain",
		};
	}

	// Navigate to domain
	ctx.reset();
	ctx.setDomain(cmd.targetDomain);

	// If action was also specified, set it
	if (cmd.targetAction) {
		ctx.setAction(cmd.targetAction);
		return {
			output: [`Navigated to ${cmd.targetDomain} > ${cmd.targetAction}`],
			shouldExit: false,
			shouldClear: false,
			contextChanged: true,
		};
	}

	return {
		output: [`Navigated to ${cmd.targetDomain} context`],
		shouldExit: false,
		shouldClear: false,
		contextChanged: true,
	};
}

/**
 * Handle domain navigation (entering a domain from root)
 */
async function handleDomainNavigation(
	domain: string,
	args: string[],
	ctx: ContextPath,
	session: REPLSession,
): Promise<ExecutionResult> {
	// Check if it's a custom domain - execute directly
	if (isCustomDomain(domain)) {
		return customDomains.execute(domain, args, session);
	}

	// Check if it's an API-generated domain
	if (!isValidDomain(domain)) {
		return {
			output: [`Unknown domain: ${domain}`],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
			error: "Unknown domain",
		};
	}

	ctx.setDomain(domain);
	return {
		output: [`Entered ${domain} context`],
		shouldExit: false,
		shouldClear: false,
		contextChanged: true,
	};
}

/**
 * Execute a command in context
 */
export async function executeCommand(
	input: string,
	session: REPLSession,
): Promise<ExecutionResult> {
	const cmd = parseCommand(input);
	const ctx = session.getContextPath();

	// Handle empty input
	if (!cmd.raw) {
		return {
			output: [],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	}

	// Handle "/" prefix direct navigation
	if (cmd.isDirectNavigation) {
		return handleDirectNavigation(cmd, ctx, session);
	}

	// Handle built-in commands
	if (cmd.isBuiltin) {
		return executeBuiltin(cmd, session, ctx);
	}

	// At root context, check if first word is a domain (custom or API-generated)
	if (ctx.isRoot()) {
		const firstWord = cmd.args[0]?.toLowerCase() ?? "";
		if (isValidDomain(firstWord) || isCustomDomain(firstWord)) {
			// Pass remaining args for custom domain execution
			const domainArgs = cmd.args.slice(1);
			return handleDomainNavigation(firstWord, domainArgs, ctx, session);
		}
	}

	// In domain context, check if first word is an action
	if (ctx.isDomain() && !ctx.isAction()) {
		const firstWord = cmd.args[0] ?? "";
		// For now, just set the action context
		// TODO: Validate against domain's available actions
		if (firstWord && !firstWord.startsWith("-")) {
			ctx.setAction(firstWord);
			return {
				output: [`Entered ${ctx.domain} > ${firstWord} context`],
				shouldExit: false,
				shouldClear: false,
				contextChanged: true,
			};
		}
	}

	// Build full command with context prepending
	let fullCommand = cmd.raw;
	if (!ctx.isRoot()) {
		if (ctx.isAction()) {
			fullCommand = `${ctx.domain} ${ctx.action} ${cmd.raw}`;
		} else {
			fullCommand = `${ctx.domain} ${cmd.raw}`;
		}
	}

	// Add to history
	session.addToHistory(cmd.raw);

	// TODO: Execute actual F5 XC API commands
	// For now, return placeholder indicating the full command
	return {
		output: [
			`[Command execution placeholder]`,
			`Context: ${ctx.isRoot() ? "root" : ctx.isAction() ? `${ctx.domain}/${ctx.action}` : ctx.domain}`,
			`Full command: ${fullCommand}`,
		],
		shouldExit: false,
		shouldClear: false,
		contextChanged: false,
	};
}

/**
 * Get suggestions for current context and input
 */
export function getCommandSuggestions(
	input: string,
	session: REPLSession,
): Array<{ text: string; description: string; category: string }> {
	const ctx = session.getContextPath();
	const suggestions: Array<{
		text: string;
		description: string;
		category: string;
	}> = [];

	// At root, suggest domains and built-ins
	if (ctx.isRoot()) {
		// Add custom domain suggestions first (higher priority)
		for (const domain of customDomains.all()) {
			if (
				!input ||
				domain.name.toLowerCase().startsWith(input.toLowerCase())
			) {
				suggestions.push({
					text: domain.name,
					description: domain.description,
					category: "domain",
				});
			}
		}

		// Add API-generated domain suggestions
		allDomains().forEach((domain) => {
			// Skip if already added as custom domain
			if (isCustomDomain(domain)) return;

			if (
				!input ||
				domain.toLowerCase().startsWith(input.toLowerCase())
			) {
				suggestions.push({
					text: domain,
					description: `Navigate to ${domain} domain`,
					category: "domain",
				});
			}
		});

		// Add built-in commands
		BUILTIN_COMMANDS.forEach((cmd) => {
			if (!input || cmd.toLowerCase().startsWith(input.toLowerCase())) {
				suggestions.push({
					text: cmd,
					description: getBuiltinDescription(cmd),
					category: "builtin",
				});
			}
		});
	}

	// In domain context, suggest actions
	if (ctx.isDomain() && !ctx.isAction()) {
		// TODO: Get actual actions from domain registry
		// For now, suggest common actions
		const commonActions = ["list", "get", "create", "delete", "update"];
		commonActions.forEach((action) => {
			if (
				!input ||
				action.toLowerCase().startsWith(input.toLowerCase())
			) {
				suggestions.push({
					text: action,
					description: `${action} ${ctx.domain} resources`,
					category: "action",
				});
			}
		});
	}

	return suggestions;
}

/**
 * Get description for built-in command
 */
function getBuiltinDescription(cmd: string): string {
	const descriptions = new Map<string, string>([
		["help", "Show help information"],
		["clear", "Clear the screen"],
		["quit", "Exit the shell"],
		["exit", "Navigate up or exit"],
		["back", "Navigate up one level"],
		["..", "Navigate up one level"],
		["root", "Navigate to root"],
		["/", "Navigate to root"],
		["context", "Show current context"],
		["ctx", "Show current context"],
		["history", "Show command history"],
		["version", "Show version info"],
		["domains", "List available domains"],
	]);
	return descriptions.get(cmd) ?? "Built-in command";
}

export default {
	parseCommand,
	executeCommand,
	getCommandSuggestions,
};
