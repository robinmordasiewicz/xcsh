/**
 * Executor - Command execution engine for the REPL
 * Handles built-in commands, domain navigation, and command routing
 */

import type { REPLSession } from "./session.js";
import type { ContextPath } from "./context.js";
import {
	allDomains,
	isValidDomain,
	resolveDomain,
	validActions,
	aliasRegistry,
} from "../types/domains.js";
import {
	customDomains,
	isCustomDomain,
	resolveDomainAlias,
	getDomainAliases,
} from "../domains/index.js";
import { extensionRegistry } from "../extensions/index.js";
import { getWhoamiInfo, formatWhoami } from "../domains/login/whoami/index.js";
import { APIError } from "../api/index.js";
import {
	formatDomainOutput,
	formatAPIError,
	parseOutputFormat,
	getCommandSpec,
	formatSpec,
	ALL_OUTPUT_FORMATS,
	OUTPUT_FORMAT_HELP,
	type OutputFormat,
} from "../output/index.js";
import { CLI_NAME, CLI_VERSION } from "../branding/index.js";
import {
	formatRootHelp,
	formatDomainHelp,
	formatActionHelp,
	formatTopicHelp,
} from "./help.js";
import { getDomainInfo } from "../types/domains.js";
import {
	validateNamespaceScope,
	checkOperationSafety,
} from "../validation/index.js";

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
	/**
	 * Raw stdout content to write directly (bypassing Ink).
	 * When set, App.tsx will hide status bar first, write this content,
	 * then restore status bar. Used for commands that need cursor positioning
	 * like the image banner.
	 */
	rawStdout?: string;
}

/**
 * Built-in command names
 */
const BUILTIN_COMMANDS = new Set([
	"help",
	"--help",
	"-h",
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
	"whoami",
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

		// Check if it's a valid domain (custom, extension, or API-generated)
		const hasExtension = extensionRegistry.hasExtension(domainPart);
		if (
			isValidDomain(domainPart) ||
			isCustomDomain(domainPart) ||
			hasExtension
		) {
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
	// Normalize: strip leading "/" for builtin check (e.g., "/help" → "help")
	const normalizedFirst = firstWord.startsWith("/")
		? firstWord.slice(1)
		: firstWord;
	if (
		BUILTIN_COMMANDS.has(firstWord) ||
		BUILTIN_COMMANDS.has(normalizedFirst)
	) {
		// Map --help and -h to help for execution
		const effectiveCommand =
			normalizedFirst === "--help" || normalizedFirst === "-h"
				? "help"
				: normalizedFirst;
		return {
			raw: effectiveCommand,
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
 * Returns ExecutionResult or Promise<ExecutionResult> for async commands like whoami
 */
function executeBuiltin(
	cmd: ParsedCommand,
	session: REPLSession,
	ctx: ContextPath,
): ExecutionResult | Promise<ExecutionResult> {
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
			output: [`${CLI_NAME} version ${CLI_VERSION}`],
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

	// Show current user and connection info
	if (command === "whoami" || command.startsWith("whoami ")) {
		// Parse flags from command
		const parts = command.split(/\s+/).slice(1); // Skip "whoami"
		const options: {
			includeQuotas?: boolean;
			includeAddons?: boolean;
			verbose?: boolean;
			json?: boolean;
		} = {};

		for (const arg of parts) {
			const lowerArg = arg.toLowerCase();
			switch (lowerArg) {
				case "--quota":
				case "--quotas":
				case "-q":
					options.includeQuotas = true;
					break;
				case "--addons":
				case "--addon":
				case "-a":
					options.includeAddons = true;
					break;
				case "--verbose":
				case "-v":
					options.verbose = true;
					break;
				case "--json":
				case "-j":
					options.json = true;
					break;
			}
		}

		// getWhoamiInfo is async, but executeBuiltin is sync
		// Return a promise-based result that will be handled by the caller
		return getWhoamiInfo(session, options)
			.then((info) => ({
				output: formatWhoami(info, options),
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
			}))
			.catch((error: unknown) => ({
				output: [
					`Failed to get whoami info: ${error instanceof Error ? error.message : "Unknown error"}`,
				],
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
				error: "whoami failed",
			}));
	}

	// Help command - context-aware
	if (command === "help" || command.startsWith("help ")) {
		// Check if a specific topic was requested
		// Topic comes from cmd.args (parsed command arguments)
		const topic = cmd.args[0]?.trim() ?? "";

		if (topic) {
			// Help for a specific topic
			return {
				output: formatTopicHelp(topic),
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
			};
		}

		// Context-aware help
		if (ctx.isRoot()) {
			// Root context: show comprehensive help
			return {
				output: formatRootHelp(),
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
			};
		}

		if (ctx.isDomain()) {
			// Domain context: show domain-specific help (no global flags/env vars)
			const domain = ctx.domain ?? "";
			const domainInfo = getDomainInfo(domain);
			if (domainInfo) {
				return {
					output: formatDomainHelp(domainInfo),
					shouldExit: false,
					shouldClear: false,
					contextChanged: false,
				};
			}
			// Fallback for unknown domain
			return {
				output: [
					`Help for domain: ${domain}`,
					"",
					"Use 'help' at root for full help.",
				],
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
			};
		}

		if (ctx.isAction()) {
			// Action context: show action-specific help
			const domain = ctx.domain ?? "";
			const action = ctx.action ?? "";
			return {
				output: formatActionHelp(domain, action),
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
			};
		}

		// Fallback to root help
		return {
			output: formatRootHelp(),
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
 *
 * Resolution order:
 * 1. Custom domains (login, cloudstatus) - full custom implementation
 * 2. Extension commands - xcsh-specific augmentations
 * 3. API domains - generated from upstream specs
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

	// Check if it's a custom domain (including aliases) - execute directly with all args
	// Custom domains handle their own --help
	if (isCustomDomain(cmd.targetDomain)) {
		const canonicalDomain = resolveDomainAlias(cmd.targetDomain);
		const allArgs = [cmd.targetAction, ...cmd.args].filter(
			(arg): arg is string => arg !== undefined,
		);
		return customDomains.execute(canonicalDomain, allArgs, session);
	}

	// Check for --help or -h flag on API/extension domains - show domain-specific help
	if (
		cmd.targetAction === "--help" ||
		cmd.targetAction === "-h" ||
		cmd.targetAction === "help"
	) {
		const domainInfo = getDomainInfo(cmd.targetDomain);
		if (domainInfo) {
			return {
				output: formatDomainHelp(domainInfo),
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
			};
		}
		// For domains without domain info, show a basic message
		return {
			output: [
				`${cmd.targetDomain} - Run '${cmd.targetDomain}' for available commands.`,
				"",
				`For global options, run: help`,
			],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	}

	// Check for extension commands
	// Extensions augment API domains with xcsh-specific functionality
	const extensionDomain =
		aliasRegistry.get(cmd.targetDomain) ?? cmd.targetDomain;
	const merged = extensionRegistry.getMergedDomain(extensionDomain);

	if (merged?.hasExtension && cmd.targetAction) {
		// Check if action is an extension command (not an API action)
		const extCmd = extensionRegistry.getExtensionCommand(
			extensionDomain,
			cmd.targetAction,
		);
		if (extCmd) {
			// Execute extension command
			const result = await extCmd.execute(cmd.args, session);
			const execResult: ExecutionResult = {
				output: result.output,
				shouldExit: result.shouldExit ?? false,
				shouldClear: result.shouldClear ?? false,
				contextChanged: result.contextChanged ?? false,
			};
			if (result.error) {
				execResult.error = result.error;
			}
			return execResult;
		}
	}

	// Handle standalone extensions (extension exists but no API domain yet)
	if (merged?.hasExtension && !merged.hasGeneratedDomain) {
		// Standalone extension - if no action, show available commands
		if (!cmd.targetAction) {
			const lines = [`${merged.displayName} commands:`];
			lines.push("");
			for (const [name, cmdDef] of merged.extensionCommands) {
				const aliases = cmdDef.aliases
					? ` (${cmdDef.aliases.join(", ")})`
					: "";
				lines.push(`  ${name}${aliases} - ${cmdDef.description}`);
			}
			return {
				output: lines,
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
			};
		}

		// Action specified but not found in extension
		return {
			output: [
				`Unknown command: ${cmd.targetAction}`,
				"",
				`Available ${merged.displayName} commands:`,
				...Array.from(merged.extensionCommands.keys()).map(
					(c) => `  ${c}`,
				),
			],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
			error: `Unknown command: ${cmd.targetAction}`,
		};
	}

	// Validate API-generated domain
	if (!isValidDomain(cmd.targetDomain) && !merged?.hasExtension) {
		return {
			output: [`Unknown domain: ${cmd.targetDomain}`],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
			error: "Unknown domain",
		};
	}

	// Navigate to domain (API domain or merged domain with API support)
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
 *
 * Resolution order:
 * 1. Custom domains (login, cloudstatus) - full custom implementation
 * 2. Extension domains - check for extension commands
 * 3. API domains - generated from upstream specs
 */
async function handleDomainNavigation(
	domain: string,
	args: string[],
	ctx: ContextPath,
	session: REPLSession,
): Promise<ExecutionResult> {
	// Check if it's a custom domain (including aliases) - execute directly
	// Custom domains handle their own --help
	if (isCustomDomain(domain)) {
		const canonicalDomain = resolveDomainAlias(domain);
		return customDomains.execute(canonicalDomain, args, session);
	}

	// Check for --help or -h flag on API/extension domains - show domain-specific help
	const firstArg = args[0]?.toLowerCase() ?? "";
	if (firstArg === "--help" || firstArg === "-h" || firstArg === "help") {
		const domainInfo = getDomainInfo(domain);
		if (domainInfo) {
			return {
				output: formatDomainHelp(domainInfo),
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
			};
		}
		// For domains without domain info, show a basic message
		return {
			output: [
				`${domain} - Run '${domain}' for available commands.`,
				"",
				`For global options, run: help`,
			],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	}

	// Check for extension domain
	const merged = extensionRegistry.getMergedDomain(domain);

	if (merged?.hasExtension) {
		// If args provided, check if first arg is an extension command
		if (args.length > 0) {
			const action = args[0];
			const extCmd = extensionRegistry.getExtensionCommand(
				domain,
				action ?? "",
			);
			if (extCmd) {
				// Execute extension command
				const result = await extCmd.execute(args.slice(1), session);
				const execResult: ExecutionResult = {
					output: result.output,
					shouldExit: result.shouldExit ?? false,
					shouldClear: result.shouldClear ?? false,
					contextChanged: result.contextChanged ?? false,
				};
				if (result.error) {
					execResult.error = result.error;
				}
				return execResult;
			}
		}

		// Standalone extension with no API domain - show commands
		if (!merged.hasGeneratedDomain) {
			if (args.length === 0) {
				const lines = [`${merged.displayName} commands:`];
				lines.push("");
				for (const [name, cmdDef] of merged.extensionCommands) {
					const aliases = cmdDef.aliases
						? ` (${cmdDef.aliases.join(", ")})`
						: "";
					lines.push(`  ${name}${aliases} - ${cmdDef.description}`);
				}
				return {
					output: lines,
					shouldExit: false,
					shouldClear: false,
					contextChanged: false,
				};
			}

			// Unknown command for standalone extension
			return {
				output: [
					`Unknown command: ${args[0]}`,
					"",
					`Available ${merged.displayName} commands:`,
					...Array.from(merged.extensionCommands.keys()).map(
						(c) => `  ${c}`,
					),
				],
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
				error: `Unknown command: ${args[0]}`,
			};
		}
	}

	// Check if it's an API-generated domain
	if (!isValidDomain(domain) && !merged?.hasExtension) {
		return {
			output: [`Unknown domain: ${domain}`],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
			error: "Unknown domain",
		};
	}

	// Set domain context first
	ctx.setDomain(domain);

	// If there are remaining args, check if first is an action to execute
	if (args.length > 0) {
		const firstArg = args[0]?.toLowerCase() ?? "";
		const resourceActions = new Set(["list", "get", "delete", "status"]);

		if (resourceActions.has(firstArg)) {
			// Create a ParsedCommand from the args and execute API command
			const cmd: ParsedCommand = {
				raw: args.join(" "),
				args: args,
				isBuiltin: false,
				isDirectNavigation: false,
			};
			return await executeAPICommand(session, ctx, cmd);
		}
	}

	// No action args - just return context change
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

	// Add ALL non-empty commands to history (before execution)
	// This ensures direct navigation, builtins, and API commands are all recorded
	session.addToHistory(input);

	// Handle "/" prefix direct navigation
	if (cmd.isDirectNavigation) {
		return handleDirectNavigation(cmd, ctx, session);
	}

	// Handle built-in commands
	if (cmd.isBuiltin) {
		return executeBuiltin(cmd, session, ctx);
	}

	// At root context, check if first word is a domain (custom, extension, or API-generated)
	if (ctx.isRoot()) {
		const firstWord = cmd.args[0]?.toLowerCase() ?? "";
		const hasExtension = extensionRegistry.hasExtension(firstWord);
		if (
			isValidDomain(firstWord) ||
			isCustomDomain(firstWord) ||
			hasExtension
		) {
			// Pass remaining args for domain execution
			const domainArgs = cmd.args.slice(1);
			return handleDomainNavigation(firstWord, domainArgs, ctx, session);
		}
	}

	// In domain context, execute API command directly (no action sub-context)
	// Execute API command
	return await executeAPICommand(session, ctx, cmd);
}

/**
 * Convert domain name to API resource path
 */
function domainToResourcePath(domain: string): string {
	// Resolve alias to canonical name
	const canonical = resolveDomain(domain) ?? domain;

	// F5 XC API uses snake_case for resource paths (not kebab-case)
	// e.g., http_loadbalancer → http_loadbalancers (plural)
	const resourceName = canonical;

	// Add 's' for plural form (most F5 XC resources are plural in API)
	return resourceName.endsWith("s") ? resourceName : `${resourceName}s`;
}

/**
 * Parsed command arguments
 */
interface ParsedArgs {
	resourceType: string | undefined;
	name: string | undefined;
	namespace: string | undefined;
	outputFormat: OutputFormat | undefined;
	spec: boolean;
	noColor: boolean;
}

/**
 * Parse command arguments for resource type, name, namespace, output format, and other flags
 * @param args - Command arguments to parse
 * @param domainResourceTypes - Set of valid resource type names for the current domain
 */
function parseCommandArgs(
	args: string[],
	domainResourceTypes?: Set<string>,
): ParsedArgs {
	let resourceType: string | undefined;
	let name: string | undefined;
	let namespace: string | undefined;
	let outputFormat: OutputFormat | undefined;
	let spec = false;
	let noColor = false;
	let positionalIndex = 0;

	for (let i = 0; i < args.length; i++) {
		const arg = args[i] ?? "";

		if (arg.startsWith("--")) {
			const flagName = arg.slice(2).toLowerCase();
			const nextArg = args[i + 1];

			switch (flagName) {
				case "namespace":
				case "ns":
					namespace = nextArg;
					i++;
					break;
				case "name":
					name = nextArg;
					i++;
					break;
				case "output":
					if (nextArg) {
						outputFormat = parseOutputFormat(nextArg);
						i++;
					}
					break;
				case "spec":
					spec = true;
					break;
				case "no-color":
					noColor = true;
					break;
				default:
					// Skip other flags with values
					if (nextArg && !nextArg.startsWith("--")) {
						i++;
					}
			}
		} else if (arg.startsWith("-")) {
			const flagName = arg.slice(1);
			const nextArg = args[i + 1];

			switch (flagName) {
				case "n":
				case "ns":
					namespace = nextArg;
					i++;
					break;
				case "o":
					if (nextArg) {
						outputFormat = parseOutputFormat(nextArg);
						i++;
					}
					break;
				default:
					// Skip other flags with values
					if (nextArg && !nextArg.startsWith("-")) {
						i++;
					}
			}
		} else {
			// Positional argument
			if (positionalIndex === 0) {
				// First positional arg: could be resource type or resource name
				if (domainResourceTypes?.has(arg.toLowerCase())) {
					// It's a valid resource type for this domain
					resourceType = arg.toLowerCase();
				} else {
					// Not a resource type, treat as resource name
					name = arg;
				}
			} else if (positionalIndex === 1 && resourceType && !name) {
				// Second positional arg after resource type: resource name
				name = arg;
			}
			positionalIndex++;
		}
	}

	return { resourceType, name, namespace, outputFormat, spec, noColor };
}

/**
 * Execute an API command against the F5 XC API
 */
async function executeAPICommand(
	session: REPLSession,
	ctx: ContextPath,
	cmd: ParsedCommand,
): Promise<ExecutionResult> {
	const client = session.getAPIClient();

	// Check if connected
	if (!client) {
		return {
			output: ["Error: Not connected to F5 XC API"],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
			error: "Not connected. Use 'login profile use <profile>' to connect.",
		};
	}

	// Check if authenticated
	if (!client.isAuthenticated()) {
		return {
			output: ["Error: Not authenticated"],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
			error: "Not authenticated. Configure a profile with API token.",
		};
	}

	// Determine domain and action from context and command
	let domain: string;
	let action: string;
	let args: string[];

	if (ctx.isAction()) {
		// Already in action context, command is arguments
		domain = ctx.domain ?? "";
		action = ctx.action ?? "";
		args = cmd.args;
	} else if (ctx.domain) {
		// In domain context, first arg might be action
		domain = ctx.domain;
		const firstArg = cmd.args[0]?.toLowerCase() ?? "";
		if (validActions.has(firstArg)) {
			action = firstArg;
			args = cmd.args.slice(1);
		} else {
			// Default to list if no action specified
			action = "list";
			args = cmd.args;
		}
	} else {
		// At root, parse domain/action from command
		const parts = cmd.raw.split(/\s+/);
		domain = parts[0] ?? "";
		action = parts[1]?.toLowerCase() ?? "list";
		args = parts.slice(validActions.has(action) ? 2 : 1);

		// If second part is not a valid action, treat as args
		if (!validActions.has(action)) {
			action = "list";
		}
	}

	// Resolve domain alias
	const canonicalDomain = resolveDomain(domain) ?? domain;

	// Get valid resource types for this domain
	const domainInfo = getDomainInfo(canonicalDomain);
	const domainResourceTypes = new Set(
		domainInfo?.primaryResources?.map((r) => r.name) ?? [],
	);

	// Parse arguments (with resource type detection)
	const { resourceType, name, namespace, outputFormat, spec, noColor } =
		parseCommandArgs(args, domainResourceTypes);
	const effectiveNamespace = namespace ?? session.getNamespace();

	// Determine which resource to use for the API path
	// If a resource type was specified (e.g., "list http_loadbalancer"), use it
	// Otherwise, fall back to the domain name
	const effectiveResource = resourceType ?? canonicalDomain;

	// Validate namespace scope (from upstream enrichment)
	const nsValidation = validateNamespaceScope(
		canonicalDomain,
		action,
		effectiveNamespace,
	);
	if (!nsValidation.valid) {
		const errorMsg =
			nsValidation.message || "Invalid namespace for this operation";
		const suggestion = nsValidation.suggestion
			? `\nSuggestion: Use --namespace ${nsValidation.suggestion}`
			: "";
		return {
			output: [`Error: ${errorMsg}${suggestion}`],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
			error: errorMsg,
		};
	}

	// Check operation safety and show warnings for dangerous operations
	const safetyCheck = checkOperationSafety(canonicalDomain, action);
	const warningOutput: string[] = [];
	if (safetyCheck.warning) {
		warningOutput.push(safetyCheck.warning);
		warningOutput.push("");
	}

	// Handle --spec flag: return command specification for AI assistants
	if (spec) {
		const commandPath = `${canonicalDomain} ${action}`;
		const cmdSpec = getCommandSpec(commandPath);
		if (cmdSpec) {
			return {
				output: [formatSpec(cmdSpec)],
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
			};
		}
		// Build a basic spec for API commands
		const basicSpec = {
			command: `${CLI_NAME} ${canonicalDomain} ${action}`,
			description: `Execute ${action} on ${canonicalDomain} resources`,
			usage: `${CLI_NAME} ${canonicalDomain} ${action} [name] [options]`,
			flags: [
				{
					name: "--namespace",
					alias: "-ns",
					type: "string",
					description: "Target namespace",
				},
				{
					name: "--output",
					alias: "-o",
					type: "string",
					description: `Output format (${OUTPUT_FORMAT_HELP})`,
				},
				{
					name: "--name",
					type: "string",
					description: "Resource name",
				},
			],
			outputFormats: [...ALL_OUTPUT_FORMATS],
		};
		return {
			output: [JSON.stringify(basicSpec, null, 2)],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	}

	// Build API path using the effective resource (explicit resource type or domain)
	const resourcePath = domainToResourcePath(effectiveResource);
	let apiPath = `/api/config/namespaces/${effectiveNamespace}/${resourcePath}`;

	// Execute based on action
	try {
		let result: unknown;

		switch (action) {
			case "list": {
				const response = await client.get(apiPath);
				result = response.data;
				break;
			}

			case "get": {
				if (!name) {
					return {
						output: [
							"Error: Resource name required for 'get' action",
						],
						shouldExit: false,
						shouldClear: false,
						contextChanged: false,
						error: "Usage: get <name>",
					};
				}
				apiPath += `/${name}`;
				const response = await client.get(apiPath);
				result = response.data;
				break;
			}

			case "delete": {
				if (!name) {
					return {
						output: [
							"Error: Resource name required for 'delete' action",
						],
						shouldExit: false,
						shouldClear: false,
						contextChanged: false,
						error: "Usage: delete <name>",
					};
				}
				apiPath += `/${name}`;
				await client.delete(apiPath);
				result = { message: `Deleted ${canonicalDomain} '${name}'` };
				break;
			}

			case "create":
			case "replace":
			case "apply": {
				// For create/replace/apply, we need a request body
				// This would typically come from a file or stdin
				return {
					output: [
						`Action '${action}' requires a resource specification.`,
						"Use --file <path> to provide resource YAML/JSON.",
					],
					shouldExit: false,
					shouldClear: false,
					contextChanged: false,
				};
			}

			case "status": {
				if (!name) {
					return {
						output: [
							"Error: Resource name required for 'status' action",
						],
						shouldExit: false,
						shouldClear: false,
						contextChanged: false,
						error: "Usage: status <name>",
					};
				}
				// Status typically comes from a different endpoint
				apiPath += `/${name}/status`;
				const response = await client.get(apiPath);
				result = response.data;
				break;
			}

			default: {
				return {
					output: [`Unknown action: ${action}`],
					shouldExit: false,
					shouldClear: false,
					contextChanged: false,
					error: `Valid actions: ${Array.from(validActions).join(", ")}`,
				};
			}
		}

		// Format output
		// Priority: CLI flag (--output) > session format (from env var or default)
		const effectiveFormat = outputFormat ?? session.getOutputFormat();
		const formatted = formatDomainOutput(result, {
			format: effectiveFormat,
			noColor: noColor ?? false,
		});

		// Prepend safety warnings if any
		const finalOutput = [
			...warningOutput,
			...(formatted.length > 0 ? formatted : ["(no output)"]),
		];

		return {
			output: finalOutput,
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
		};
	} catch (error) {
		if (error instanceof APIError) {
			const formatted = formatAPIError(
				error.statusCode,
				error.response,
				`${action} ${canonicalDomain}`,
			);
			return {
				output: [formatted],
				shouldExit: false,
				shouldClear: false,
				contextChanged: false,
				error: error.message,
			};
		}

		// Handle other errors
		const message = error instanceof Error ? error.message : String(error);
		return {
			output: [`Error: ${message}`],
			shouldExit: false,
			shouldClear: false,
			contextChanged: false,
			error: message,
		};
	}
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

		// Add domain aliases
		for (const [alias, canonical] of getDomainAliases()) {
			if (!input || alias.toLowerCase().startsWith(input.toLowerCase())) {
				const domain = customDomains.get(canonical);
				suggestions.push({
					text: alias,
					description: domain
						? `${domain.description} (alias)`
						: `Alias for ${canonical}`,
					category: "domain",
				});
			}
		}

		// Add extension domain suggestions (standalone extensions not in API)
		for (const extDomain of extensionRegistry.getExtendedDomains()) {
			// Skip if already added as custom domain or API domain
			if (isCustomDomain(extDomain) || isValidDomain(extDomain)) continue;

			if (
				!input ||
				extDomain.toLowerCase().startsWith(input.toLowerCase())
			) {
				const merged = extensionRegistry.getMergedDomain(extDomain);
				suggestions.push({
					text: extDomain,
					description: merged?.description ?? `${extDomain} commands`,
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
				const merged = extensionRegistry.getMergedDomain(domain);
				const hasExt = merged?.hasExtension ? " (+ext)" : "";
				suggestions.push({
					text: domain,
					description: `Navigate to ${domain} domain${hasExt}`,
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
		const domain = ctx.domain ?? "";

		// Add extension commands first (higher priority for xcsh-specific)
		const extCmds = extensionRegistry.getExtensionCommandNames(domain);
		for (const cmd of extCmds) {
			if (!input || cmd.toLowerCase().startsWith(input.toLowerCase())) {
				const cmdDef = extensionRegistry.getExtensionCommand(
					domain,
					cmd,
				);
				suggestions.push({
					text: cmd,
					description: cmdDef?.description ?? `${cmd} command`,
					category: "extension",
				});
			}
		}

		// Add API actions
		const commonActions = ["list", "get", "create", "delete", "update"];
		commonActions.forEach((action) => {
			if (
				!input ||
				action.toLowerCase().startsWith(input.toLowerCase())
			) {
				suggestions.push({
					text: action,
					description: `${action} ${domain} resources`,
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
