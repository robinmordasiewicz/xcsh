#!/usr/bin/env node
/**
 * F5 Distributed Cloud Shell - Interactive CLI for F5 XC
 * Main entry point for the xcsh CLI application.
 *
 * Usage:
 *   xcsh                    # Enter interactive REPL mode
 *   xcsh <command> [args]   # Execute command non-interactively
 *   xcsh --help             # Show help
 *   xcsh --version          # Show version
 */

import { render } from "ink";
import { Command } from "commander";
import { App, type AppProps } from "./repl/index.js";
import { CLI_NAME, CLI_VERSION, colors, ENV_PREFIX } from "./branding/index.js";
import { executeCommand } from "./repl/executor.js";
import { REPLSession } from "./repl/session.js";
import { formatRootHelp } from "./repl/help.js";
import { isValidLogoMode, type LogoDisplayMode } from "./config/index.js";
import { renderBanner } from "./domains/login/banner/display.js";
import { debugProtocol, emitSessionState } from "./debug/protocol.js";

const program = new Command();

// Custom help: override default help to show comprehensive root help
program.configureHelp({
	formatHelp: () => formatRootHelp().join("\n"),
});

program
	.name(CLI_NAME)
	.description("F5 Distributed Cloud Shell - Interactive CLI for F5 XC")
	.version(CLI_VERSION, "-v, --version", "Show version number")
	.option("--no-color", "Disable color output")
	.option("--logo <mode>", "Logo display mode: image, ascii, none")
	.option("-o, --output <format>", "Output format (json, yaml, table)")
	.option("--spec", "Output command specification as JSON (for AI)")
	.option("-h, --help", "Show help") // Manual help option to prevent auto-exit
	.argument("[command...]", "Command to execute non-interactively")
	.allowUnknownOption(true) // Pass through unknown options to commands
	.helpOption(false) // Disable auto-help so domain --help works
	.action(
		async (
			commandArgs: string[],
			options: {
				help?: boolean;
				logo?: string;
				output?: string;
				spec?: boolean;
			},
		) => {
			// Handle root-level help (xcsh --help or xcsh -h with no domain)
			if (options.help && commandArgs.length === 0) {
				formatRootHelp().forEach((line) => console.log(line));
				process.exit(0);
			}

			// If --help with a domain, re-inject --help into args for domain help
			// Commander consumes --help as an option, so we add it back
			if (options.help && commandArgs.length > 0) {
				commandArgs.push("--help");
			}

			// If --logo with a command, re-inject --logo into args for command handling
			// Commander consumes --logo as an option, so we add it back for non-interactive mode
			if (options.logo && commandArgs.length > 0) {
				commandArgs.push("--logo", options.logo);
			}

			// If --output with a command, re-inject --output into args for command handling
			// Commander consumes --output as an option, so we add it back for non-interactive mode
			if (options.output && commandArgs.length > 0) {
				commandArgs.push("--output", options.output);
			}

			// If --spec with a command, re-inject --spec into args for command handling
			// Commander consumes --spec as an option, so we add it back for non-interactive mode
			if (options.spec && commandArgs.length > 0) {
				commandArgs.push("--spec");
			}

			// If no command args, enter REPL mode
			if (commandArgs.length === 0) {
				// Check if stdin is a TTY (interactive terminal)
				if (!process.stdin.isTTY) {
					console.error(
						"Error: Interactive mode requires a terminal (TTY).",
					);
					console.error(
						"Use: xcsh <command> for non-interactive execution.",
					);
					process.exit(1);
				}

				// Parse logo mode from CLI option
				const cliLogoMode =
					options.logo && isValidLogoMode(options.logo)
						? (options.logo as LogoDisplayMode)
						: undefined;

				// Show initialization message first
				process.stdout.write("Initializing...");

				// Initialize session BEFORE Ink takes over
				const session = new REPLSession();
				await session.initialize();

				// Emit debug event for session state (helps AI/PTY debugging)
				debugProtocol.session("init", { mode: "repl" });
				emitSessionState(session);

				// Clear the "Initializing..." message
				process.stdout.write("\r\x1b[K");

				// Print banner to scrollback BEFORE Ink takes over
				// Use "startup" context for direct stdout with image support
				renderBanner(cliLogoMode, "startup");

				// Show warning if token validation failed
				if (
					session.isAuthenticated() &&
					!session.isTokenValidated() &&
					session.getValidationError()
				) {
					console.log("");
					console.log(
						`${colors.yellow}Warning: ${session.getValidationError()}${colors.reset}`,
					);
				}

				// Check if user needs guidance on connecting
				const profiles = await session.getProfileManager().list();
				const envConfigured =
					process.env[`${ENV_PREFIX}_API_URL`] &&
					process.env[`${ENV_PREFIX}_API_TOKEN`];

				if (profiles.length === 0 && !envConfigured) {
					console.log("");
					console.log(
						`${colors.yellow}No connection profiles found.${colors.reset}`,
					);
					console.log("");
					console.log(
						"Create a profile to connect to F5 Distributed Cloud:",
					);
					console.log("");
					console.log(
						`  ${colors.blue}login profile create${colors.reset} <name> --url <api-url> --token <api-token>`,
					);
					console.log("");
					console.log("Or set environment variables:");
					console.log("");
					console.log(
						`  ${colors.blue}export ${ENV_PREFIX}_API_URL${colors.reset}=https://tenant.console.ves.volterra.io`,
					);
					console.log(
						`  ${colors.blue}export ${ENV_PREFIX}_API_TOKEN${colors.reset}=<your-api-token>`,
					);
					console.log("");
				}

				// Enter interactive REPL mode with pre-initialized session
				// WORKAROUND: Bun doesn't call process.stdin.resume() automatically,
				// which breaks Ink's useInput hook. This is a known Bun bug:
				// https://github.com/oven-sh/bun/issues/6862
				process.stdin.resume();
				const appProps: AppProps = { initialSession: session };
				render(<App {...appProps} />);
				return;
			}

			// Non-interactive mode: execute command and exit
			await executeNonInteractive(commandArgs);
		},
	);

/**
 * Execute a command non-interactively
 */
async function executeNonInteractive(args: string[]): Promise<void> {
	const session = new REPLSession();
	await session.initialize();

	// Emit debug event for session state (helps AI/PTY debugging)
	debugProtocol.session("init", {
		mode: "non-interactive",
		command: args.join(" "),
	});
	emitSessionState(session);

	// Show warning if token validation failed
	if (
		session.isAuthenticated() &&
		!session.isTokenValidated() &&
		session.getValidationError()
	) {
		console.error(
			`${colors.yellow}Warning: ${session.getValidationError()}${colors.reset}`,
		);
	}

	// Join args into a command string
	const command = args.join(" ");

	// Execute the command
	const result = await executeCommand(command, session);

	// Output results
	result.output.forEach((line) => {
		// eslint-disable-next-line no-console
		console.log(line);
	});

	// Exit with appropriate code
	if (result.error) {
		process.exit(1);
	}

	process.exit(0);
}

// Parse command line arguments
program.parse();
