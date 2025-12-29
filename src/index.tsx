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
import { App } from "./repl/index.js";
import { CLI_NAME, CLI_VERSION } from "./branding/index.js";
import { executeCommand } from "./repl/executor.js";
import { REPLSession } from "./repl/session.js";
import { formatRootHelp } from "./repl/help.js";
import { isValidLogoMode, type LogoDisplayMode } from "./config/index.js";
import { renderBanner } from "./domains/login/banner/display.js";

const program = new Command();

// Custom help: override default help to show comprehensive root help
program.configureHelp({
	formatHelp: () => formatRootHelp().join("\n"),
});

program
	.name(CLI_NAME)
	.description("F5 Distributed Cloud Shell - Interactive CLI for F5 XC")
	.version(CLI_VERSION, "-v, --version", "Show version number")
	.option("-i, --interactive", "Force interactive mode")
	.option("--no-color", "Disable color output")
	.option("--logo <mode>", "Logo display mode: image, ascii, none")
	.option("-h, --help", "Show help") // Manual help option to prevent auto-exit
	.argument("[command...]", "Command to execute non-interactively")
	.allowUnknownOption(true) // Pass through unknown options to commands
	.helpOption(false) // Disable auto-help so domain --help works
	.action(
		async (
			commandArgs: string[],
			options: { interactive?: boolean; help?: boolean; logo?: string },
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

			// If no command args or --interactive flag, enter REPL mode
			if (commandArgs.length === 0 || options.interactive) {
				// Check if stdin is a TTY (interactive terminal)
				if (!process.stdin.isTTY && !options.interactive) {
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

				// Print banner to scrollback BEFORE Ink takes over
				// Use "startup" context for direct stdout with image support
				renderBanner(cliLogoMode, "startup");

				// Enter interactive REPL mode
				// WORKAROUND: Bun doesn't call process.stdin.resume() automatically,
				// which breaks Ink's useInput hook. This is a known Bun bug:
				// https://github.com/oven-sh/bun/issues/6862
				process.stdin.resume();
				render(<App />);
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
