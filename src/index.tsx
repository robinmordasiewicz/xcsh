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

const program = new Command();

program
	.name(CLI_NAME)
	.description("F5 Distributed Cloud Shell - Interactive CLI for F5 XC")
	.version(CLI_VERSION, "-v, --version", "Show version number")
	.option("-i, --interactive", "Force interactive mode")
	.option("--no-color", "Disable color output")
	.argument("[command...]", "Command to execute non-interactively")
	.allowUnknownOption(true) // Pass through unknown options to commands
	.action(
		async (commandArgs: string[], options: { interactive?: boolean }) => {
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

				// Enter interactive REPL mode
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
